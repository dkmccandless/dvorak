// Package html implements the MediaWiki sanitizer's removeHTMLtags function.
package html

import (
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// wfURLProtocols lists valid URL protocols
// as a single pipe-delimited string for use in regular expressions.
// It emulates the output of MediaWiki's UrlUtils::validProtocols,
// which MediaWiki's wfUrlProtocols is deprecated in favor of.
// https://www.mediawiki.org/wiki/Release_notes/1.39
// https://doc.wikimedia.org/mediawiki-core/master/php/UrlUtils_8php_source.html
const wfURLProtocols = `bitcoin:|ftp:\/\/|ftps:\/\/|geo:|git:\/\/|gopher:\/\/|http:\/\/|` +
	`https:\/\/|irc:\/\/|ircs:\/\/|magnet:|mailto:|matrix:|mms:\/\/|` +
	`news:|nntp:\/\/|redis:\/\/|sftp:\/\/|sip:|sips:|sms:|` +
	`ssh:\/\/|svn:\/\/|tel:|telnet:\/\/|urn:|worldwind:\/\/|xmpp:|` +
	`\/\/`

var (
	// attrRegex matches HTML/XML attribute pairs within a tag.
	attrRegex = regexp.MustCompile(
		`((?s)(?:[^\s\/=]|=)[^\s\/=]*)` +
			`([\s]*=[\s]*` +
			`(?:` +
			// The attribute value, double- or single-quoted or unquoted
			`\"([^\"]*)(?:\"|\$)` +
			`|'([^']*)(?:'|\$)` +
			`|([^\s>]*)` +
			`)` +
			`)?`,
	)

	// attrNameRegex describes valid attribute names.
	attrNameRegex = regexp.MustCompile(
		`(?s:^([:_\p{L}\p{N}][:_\.\-\p{L}\p{N}]*)$)`,
	)

	// Alternative cases of MediaWiki's charRefsRegex.
	// charRefsRegex matches "various types of character references".
	// HTML5 allows some named entities to omit the trailing semicolon,
	// but wikitext syntax requires the semicolon; charRefsRegex reflects this.
	charRefsEntityRegex    = regexp.MustCompile(`&([A-Za-z0-9\x80-\xff]+;)`)
	charRefsDecRegex       = regexp.MustCompile(`&\#([0-9]+);`)
	charRefsHexRegex       = regexp.MustCompile(`&\#[xX]([0-9A-Fa-f]+);`)
	charRefsAmpersandRegex = regexp.MustCompile(`(&)`)

	// cssCommentRegex matches CSS comments.
	cssCommentRegex = regexp.MustCompile(`\s*/\*[^*\\/]*\*/\s*$`)

	// cssControlCharsRegex matches CSS control characters.
	cssControlCharsRegex = regexp.MustCompile(`[\000-\101\103\016-\037\177]`)

	// cssProblematicRegex matches "problematic" CSS keywords.
	cssProblematicRegex = regexp.MustCompile(
		`(?i:` +
			`expression` +
			`|accelerator\s*:` +
			`|-o-link\s*:` +
			`|-o-link-source\s*:` +
			`|-o-replace\s*:` +
			`|url\s*\(` +
			`|image\s*\(` +
			`|image-set\s*\(` +
			`|attr\s*\([^)]+[\s,]+url` +
			`)`,
	)

	// dataAttrRegex matches attribute names beginning with "data-".
	dataAttrRegex = regexp.MustCompile(`(?i:^data-[^:]*$)`)

	// dataReservedRegex matches attribute names beginning with "data-"
	// that are reserved for use by MediaWiki code.
	dataReservedRegex = regexp.MustCompile(`(?i:^data-(ooui|mw|parsoid))`)

	// decodeEscapeRegex matches CSS escape sequences.
	decodeEscapeRegex = regexp.MustCompile(
		`\\\\` + // backslash
			`(?:` +
			`((?:\\n|\\r\\n|\\r|\\f))` + // line continuation
			`|([0-9A-Fa-f]{1,6})\s?` + // character number
			`|(.)` + // backslash cancelling special meaning
			`|()` + // backslash at end of string
			`|)`,
	)

	// elementBitsRegex describes valid tag names.
	elementBitsRegex = regexp.MustCompile(`^(/?)([A-Za-z][^\s/>\0]*)([^>]*?)(/?>)([^<]*)$`)

	// evilURIpattern matches "evil uris like javascript".
	// "WARNING: DO NOT use this in any place that actually requires denying
	// certain URIs for security reasons. There are NUMEROUS[1] ways to bypass
	// pattern-based deny lists; the only way to be secure from javascript:
	// uri based xss vectors is to allow only things that you know are safe
	// and deny everything else.
	// [1]: http://ha.ckers.org/xss.html"
	evilURIPattern = regexp.MustCompile(`(?i:(^|\s|\*/\s*)(javascript|vbscript)([^\w]|$))`)

	// hrefExp matches valid href attribute values.
	hrefExp = regexp.MustCompile(`^(` + wfURLProtocols + `)[^\s]+$`)

	// spaceRegex matches consecutive whitespace.
	spaceRegex = regexp.MustCompile(`[\s]+`)

	// wfURLProtocolsRegex matches valid URL protocols.
	wfURLProtocolsRegex = regexp.MustCompile(`((?i)` + wfURLProtocols + `)`)
)

var (
	// htmlElements contains all valid HTML tags.
	htmlElements = makeSet(
		// htmlsingle
		"br", "wbr", "hr", "li", "dt", "dd", "meta", "link",

		// htmlpairsStatic
		"b", "bdi", "del", "i", "ins", "u", "font", "big", "small", "sub", "sup",
		"h1", "h2", "h3", "h4", "h5", "h6", "cite", "code", "em", "s",
		"strike", "strong", "tt", "var", "div", "center",
		"blockquote", "ol", "ul", "dl", "table", "caption", "pre",
		"ruby", "rb", "rp", "rt", "rtc", "p", "span", "abbr", "dfn",
		"kbd", "samp", "data", "time", "mark",

		// htmlnest
		"table", "tr", "td", "th", "div", "blockquote", "ol", "ul",
		"li", "dl", "dt", "dd", "font", "big", "small", "sub", "sup", "span",
		"var", "kbd", "samp", "em", "strong", "q", "ruby", "bdo",
	)

	// htmlSingle contains tags that can be self-closed.
	// For tags not also in htmlSingleOnly,
	// a self-closed tag will be emitted as an empty element.
	htmlSingle = makeSet(
		"br", "wbr", "hr", "li", "dt", "dd", "meta", "link",
	)

	// htmlSingleOnly contains elements that cannot have close tags.
	htmlSingleOnly = makeSet(
		"br", "wbr", "hr", "meta", "link",
	)
)

// mwEntityAliases contains character entity aliases accepted by MediaWiki
// that are not part of the HTML standard.
var mwEntityAliases = map[string]string{
	"רלמ;": "rlm;",
	"رلم;": "rlm;",
}

// RemoveHTMLTags cleans up HTML. It removes HTML comments
// and dangerous tags and attributes, and escapes invalid tags.
// It emulates the MediaWiki sanitizer's removeHTMLtags function with default arguments.
// https://github.com/wikimedia/mediawiki/blob/80d72fc07d509916224555c9a062892fc3690864/includes/parser/Sanitizer.php
func RemoveHTMLTags(s string) string {
	var b strings.Builder
	s = removeHTMLComments(s)
	subs := strings.Split(s, "<")
	b.WriteString(strings.ReplaceAll(subs[0], ">", "&gt;"))
	subs = subs[1:]

	for _, t := range subs {
		if !elementBitsRegex.MatchString(t) {
			b.WriteString("&lt;" + strings.ReplaceAll(t, ">", "&gt;"))
			continue
		}
		regs := elementBitsRegex.FindStringSubmatch(t)
		var (
			slash  = regs[1]
			tag    = regs[2]
			params = regs[3]
			brace  = regs[4]
			rest   = regs[5]
		)

		tag = strings.ToLower(tag)
		if !htmlElements[tag] {
			b.WriteString("&lt;" + strings.ReplaceAll(t, ">", "&gt;"))
			continue
		}

		if ok := validateTag(params, tag); !ok {
			b.WriteString("&lt;" + strings.ReplaceAll(t, ">", "&gt;"))
			continue
		}

		params = fixTagAttributes(params, tag)

		if brace == "/>" && !(htmlSingle[tag] || htmlSingleOnly[tag]) {
			// Remove the self-closing slash for consistency with HTML5 semantics.
			brace = ">"
		}
		if brace == "/>" && !htmlSingleOnly[tag] {
			// Interpret self-closing tags as empty tags,
			// even when HTML5 would interpret them as start tags.
			// This usage is commonly seen on Wikimedia wikis.
			brace = "></" + tag + ">"
		}

		rest = strings.ReplaceAll(rest, ">", "&gt;")
		b.WriteString("<" + slash + tag + params + brace + rest)
	}
	return b.String()
}

// validateTag reports whether a tag is allowed to be present.
// "This DOES NOT validate the attributes, nor does it validate the
// tags themselves. This method only handles the special circumstances
// where we may want to allow a tag within content but ONLY when it has
// specific attributes set."
func validateTag(params, element string) bool {
	if element != "meta" && element != "link" {
		return true
	}
	switch attrs := decodeTagAttributes(params); element {
	case "meta":
		return attrs["itemprop"] != "" && attrs["content"] != ""
	case "link":
		return attrs["itemprop"] != "" && attrs["href"] != ""
	default:
		panic("unreached")
	}
}

// decodeTagAttributes returns a map of attribute names and values
// from a partial tag string.
func decodeTagAttributes(text string) map[string]string {
	if strings.TrimSpace(text) == "" {
		return nil
	}

	sets := attrRegex.FindAllStringSubmatch(text, -1)
	if sets == nil {
		return nil
	}

	attrs := make(map[string]string)
	for _, set := range sets {
		attr := strings.ToLower(set[1])
		if !attrNameRegex.MatchString(attr) {
			continue
		}
		var value string
		switch {
		case set[5] != "":
			// unquoted
			value = set[5]
		case set[4] != "":
			// single-quoted
			value = set[4]
		case set[3] != "":
			// double-quoted
			value = set[3]
		case set[2] == "":
			// In XHTML, attributes must have a value.
			// https://www.w3.org/TR/html5/syntax.html#syntax-attribute-name
		default:
			// tag conditions not met
			continue
		}
		value = spaceRegex.ReplaceAllString(value, " ")
		attrs[attr] = decodeCharReferences(value)
	}
	return attrs
}

// fixTagAttributes normalizes an HTML element's attributes to well-formed XML,
// discarding unwanted attributes.
func fixTagAttributes(text, element string) string {
	if strings.TrimSpace(text) == "" {
		return ""
	}
	attrs := decodeTagAttributes(text)
	attrs = validateTagAttributes(attrs, element)
	return safeEncodeTagAttributes(attrs)
}

// validateTagAttributes normalizes attribute values
// and discards illegal and unsafe values for element.
func validateTagAttributes(attribs map[string]string, element string) map[string]string {
	return validateAttributes(attribs, attributesAllowed[element])
}

func validateAttributes(attrs map[string]string, allowed map[string]bool) map[string]string {
	out := make(map[string]string)
	for attr, value := range attrs {
		// Allow any attribute beginning with "data-"
		// except those reserved for use by MediaWiki code.
		// Ensure that attr is not namespaced by banning colons.
		if !dataAttrRegex.MatchString(attr) && !allowed[attr] ||
			dataReservedRegex.MatchString(attr) {
			continue
		}

		switch attr {
		case "style":
			// Strip JavaScript "expression" from stylesheets.
			// https://msdn.microsoft.com/en-us/library/ms537634.aspx
			value = checkCSS(value)
		case "id":
			value = escapeIDForAttribute(value)
		case "aria-describedby", "aria-flowto", "aria-labelledby", "aria-owns":
			value = escapeIDReferenceList(value)
		case "rel", "rev",
			// RDFa
			"about", "property", "resource", "datatype", "typeof",
			// HTML5 microdata
			"itemid", "itemprop", "itemref", "itemscope", "itemtype":
			// "Paranoia. Allow 'simple' values but suppress javascript"
			if evilURIPattern.MatchString(value) {
				continue
			}
		case "href", "src", "poster":
			// "NOTE: even though elements using href/src are not allowed directly, supply
			// validation code that can be used by tag hook handlers, etc"
			if !hrefExp.MatchString(value) {
				// "drop any href or src attributes not using an allowed protocol.
				// NOTE: this also drops all relative URLs"
				continue
			}
		case "tabindex":
			// "Only allow tabindex of 0, which is useful for accessibility."
			if value != "0" {
				continue
			}
		}
		out[attr] = value
	}

	// "itemtype, itemid, itemref don't make sense without itemscope"
	if out["itemscope"] == "" {
		delete(out, "itemtype")
		delete(out, "itemid")
		delete(out, "itemref")
	}
	// "TODO: Strip itemprop if we aren't descendants of an itemscope
	// or pointed to by an itemref."

	return out
}

// encodeAttribute encodes an attribute value for HTML output.
func encodeAttribute(s string) string {
	return strings.NewReplacer(
		`&`, "&amp;",
		`<`, "&lt;",
		`>`, "&gt;",
		`"`, "&#34;",
		`'`, "&#39;",
		"\n", `&#10;`,
		"\r", `&#13;`,
		"\t", `&#9;`,
	).Replace(s)
}

// safeEncodeTagAttributes builds a partial tag string
// from a map of attribute names and values as returned by decodeTagAttributes.
func safeEncodeTagAttributes(m map[string]string) string {
	var b strings.Builder
	rep := strings.NewReplacer(
		// htmlspecialchars with encoding ENT_COMPAT excludes single quotes
		`&`, "&amp;",
		`<`, "&lt;",
		`>`, "&gt;",
		`"`, "&#34;",
	)
	for attr, value := range m {
		s := fmt.Sprintf(" %v=\"%v\"",
			rep.Replace(attr),
			safeEncodeAttribute(value),
		)
		b.WriteString(s)
	}
	return b.String()
}

// safeEncodeAttribute encodes an attribute value for HTML tags,
// "with extra armoring against further wiki processing."
func safeEncodeAttribute(s string) string {
	s = encodeAttribute(s)

	// "Templates and links may be expanded in later parsing,
	// creating invalid or dangerous output. Suppress this."
	s = strings.NewReplacer(
		`<`, `&lt;`, // "This should never happen,"
		`>`, `&gt;`, // "we've received invalid input"
		`"`, `&quot;`, // "which should have been escaped."
		`{`, `&#123;`,
		`}`, `&#125;`, // "prevent unpaired language conversion syntax"
		`[`, `&#91;`,
		`]`, `&#93;`,
		`''`, `&#39;&#39;`,
		`ISBN`, `&#73;SBN`,
		`RFC`, `&#82;FC`,
		`PMID`, `&#80;MID`,
		`|`, `&#124;`,
		`__`, `&#95;_`,
	).Replace(s)

	// "Stupid hack"
	replaceColons := func(s string) string {
		return strings.ReplaceAll(s, ":", `&#58;`)
	}
	return wfURLProtocolsRegex.ReplaceAllStringFunc(s, replaceColons)
}

// normalizeCSS decodes character references and escape sequences,
// and strips comments unless s is a single valid comment.
func normalizeCSS(s string) string {
	s = decodeCharReferences(s)

	// "Decode escape sequences and line continuation
	// See the grammar in the CSS 2 spec, appendix D.
	// This has to be done AFTER decoding character references.
	// This means it isn't possible for this function to return
	// unsanitized escape sequences. It is possible to manufacture
	// input that contains character references that decode to
	// escape sequences that decode to character references, but
	// it's OK for the return value to contain character references
	// because the caller is supposed to escape those anyway."
	switch matches := decodeEscapeRegex.FindStringSubmatch(s); {
	case matches[1] != "":
		// line continuation
		return ""
	case matches[2] != "":
		n, err := strconv.ParseInt(matches[2], 16, 64)
		if err != nil {
			return ""
		}
		s = string(rune(n))
	case matches[3] != "":
		s = matches[3]
	default:
		s = `\\`
	}
	if s == `\n` || s == `"` || s == "`" || s == `\\` {
		// "These characters need to be escaped in strings
		// Clean up the escape sequence to avoid parsing errors by clients"
		s = fmt.Sprintf(`\\`+"%x ", s[0])
	}

	// "Let the value through if it's nothing but a single comment,
	// to allow other functions which may reject it
	// to pass some error message through."
	if cssCommentRegex.MatchString(s) {
		return s
	}

	// "Remove any comments; IE gets token splitting wrong
	// This must be done AFTER decoding character references and
	// escape sequences, because those steps can introduce comments
	// This step cannot introduce character references or escape
	// sequences, because it replaces comments with spaces rather
	// than removing them completely."
	for {
		before, rest, found := strings.Cut(s, "/*")
		if !found {
			break
		}
		_, after, found := strings.Cut(rest, "*/")
		if !found {
			break
		}
		s = before + " " + after
	}

	// "Remove anything after a comment-start token,
	// to guard against incorrect client implementations."
	s, _, _ = strings.Cut(s, "/*")
	return s
}

// checkCSS normalizes and sanitizes s,
// removing forbidden or unsafe structures.
// The returned string may still contain cleverly encoded character references
// and must be escaped before it is embedded in HTML.
//
// If s cannot be sanitized, checkCSS returns a comment string describing why.
func checkCSS(s string) string {
	s = normalizeCSS(s)
	if cssControlCharsRegex.MatchString(s) ||
		strings.Contains(s, string(utf8.RuneError)) {
		return "/* invalid control char */"
	}
	if cssProblematicRegex.MatchString(s) {
		return "/* insecure input */"
	}
	return s
}

// decodeCharReferences decodes numeric or named entities in text.
func decodeCharReferences(text string) string {
	text = charRefsEntityRegex.ReplaceAllStringFunc(text, decodeEntity)
	text = charRefsDecRegex.ReplaceAllStringFunc(text, decodeDec)
	text = charRefsHexRegex.ReplaceAllStringFunc(text, decodeHex)
	return charRefsAmpersandRegex.ReplaceAllString(text, "&")
}

// decodeEntity returns the UTF-8 encoding of a defined named entity,
// or else the pseudo-entity source.
func decodeEntity(name string) string {
	if s := mwEntityAliases[name]; s != "" {
		return html.UnescapeString(s)
	}
	return html.UnescapeString(name)
}

// decodeDec parses s as a decimal code point
// and returns the corresponding string.
func decodeDec(s string) string {
	n, err := strconv.Atoi(s)
	if err != nil {
		return string(utf8.RuneError)
	}
	return decodeChar(rune(n))
}

// decodeHex parses s as a hexadecimal code point
// and returns the corresponding string.
func decodeHex(s string) string {
	n, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return string(utf8.RuneError)
	}
	return decodeChar(rune(n))
}

// decodeChar returns r's string representation if r is a valid code point,
// or else utf8.RuneError.
func decodeChar(r rune) string {
	if ok := validateCodePoint(r); !ok {
		return string(utf8.RuneError)
	}
	return string(r)
}

// validateCodePoint reports whether r is a valid code point in both HTML5 and XML.
func validateCodePoint(r rune) bool {
	// U+000C is valid in HTML5 but not allowed in XML.
	// U+000D is valid in XML but not allowed in HTML5.
	// U+007F-U+009F (control characters) are disallowed in HTML5.
	return r == 0x09 ||
		r == 0x0a ||
		r >= 0x20 && r <= 0x7e ||
		r >= 0xa0 && r <= 0xd7ff ||
		r >= 0xe000 && r <= 0xfffd ||
		r >= 0x10000 && r <= 0x10ffff
}

// escapeIDForAttribute escapes id to be a valid HTML ID attribute.
// "WARNING: The output of this function is not guaranteed to be HTML safe,
// so be sure to use proper escaping."
func escapeIDForAttribute(id string) string {
	// "Truncate overly-long IDs.
	// This isn't an HTML limit, it's just griefer protection."
	if len(id) >= 1024 {
		var last int
		for i := range id {
			if i >= 1024 {
				id = id[:last]
				break
			}
			last = i
		}
	}

	// "html5 spec says ids must not have any of the following:
	// U+0009 TAB, U+000A LF, U+000C FF, U+000D CR, or U+0020 SPACE
	// In practice, in wikitext, only tab, LF, CR (and SPACE) are
	// possible using either Lua or html entities."
	escapeSpace := func(r rune) rune {
		if unicode.IsSpace(r) {
			return '_'
		}
		return r
	}
	return strings.Map(escapeSpace, id)
}

// escapeIDReferenceList interprets ids as a space-delimited list of IDs
// and escapes each.
func escapeIDReferenceList(ids string) string {
	fields := strings.Fields(ids)
	for i, s := range fields {
		fields[i] = escapeIDForAttribute(s)
	}
	return strings.Join(fields, " ")
}

// removeHTMLComments removes HTML comments.
// If a comment is preceded and followed by a newline (ignoring spaces),
// removeHTMLComments removes the spaces and one of the newlines as well.
func removeHTMLComments(s string) string {
	for {
		op := strings.Index(s, "<!--")
		if op == -1 {
			break
		}
		cl := strings.Index(s[op:], "-->")
		if cl == -1 {
			break
		}
		cl = op + cl + 3
		lead := op
		for lead > 0 && s[lead-1] == ' ' {
			lead--
		}
		trail := cl
		for trail < len(s) && s[trail] == ' ' {
			trail++
		}
		if lead > 0 && s[lead-1] == '\n' &&
			trail < len(s) && s[trail] == '\n' {
			s = s[:lead-1] + "\n" + s[trail+1:]
		} else {
			s = s[:op] + s[cl:]
		}
	}
	return s
}

// makeSet returns a map of each element of ts to the boolean value true.
func makeSet[T comparable](ts ...T) map[T]bool {
	m := make(map[T]bool)
	for _, t := range ts {
		m[t] = true
	}
	return m
}

// appendSet returns a map containing the elements of set and values.
func appendSet[T comparable](set map[T]bool, values ...T) map[T]bool {
	m := make(map[T]bool)
	for t := range set {
		m[t] = true
	}
	for _, t := range values {
		m[t] = true
	}
	return m
}
