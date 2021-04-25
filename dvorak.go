// Package dvorak parses source code templates used by the Dvorak game wiki.
package dvorak

import (
	"fmt"
	"strings"
)

// Page contains the template information of a Dvorak wiki page.
type Page struct {
	// Subpages lists any subpages of the main page.
	Subpages []Subpage

	// Cards lists the main page's cards.
	Cards []Card
}

// ParsePage parses a page of wiki source code.
func ParsePage(s string) *Page {
	// Elide wiki hidden text
	for {
		op := strings.Index(s, "<!--")
		if op == -1 {
			break
		}
		cl := strings.Index(s[op:], "-->")
		if cl == -1 {
			break
		}
		s = s[:op] + s[op+cl+3:]
	}

	p := &Page{}
	for _, s := range strings.SplitAfter(s, "}}") {
		op := strings.LastIndex(s, "{{")
		if op == -1 {
			continue
		}

		name, params, err := parseTemplate(s[op:])
		if err != nil {
			continue
		}
		switch name {
		case "Card", "card":
			p.Cards = append(p.Cards, PopulateCard(params))
		case "Subpage", "subpage":
			sp, err := PopulateSubpage(params)
			if err != nil {
				continue
			}
			p.Subpages = append(p.Subpages, sp)
		}
	}
	return p
}

// parseTemplate parses a template and returns its name and parameters.
// Whitespace is trimmed from all returned strings.
// If s is not a single well-formed template or has any nested subtemplates,
// parseTemplate returns an error instead.
func parseTemplate(s string) (name string, params map[string]string, err error) {
	// https://meta.wikimedia.org/wiki/Help:Template

	var errInvalid = fmt.Errorf("invalid template syntax")

	if !strings.HasPrefix(s, "{{") || !strings.HasSuffix(s, "}}") {
		return "", nil, errInvalid
	}
	s = strings.TrimPrefix(s, "{{")
	s = strings.TrimSuffix(s, "}}")

	if strings.Contains(s, "{{") || strings.Contains(s, "}}") {
		return "", nil, errInvalid
	}

	var fields []string

	// Fast path for template values with no internal links
	if !strings.Contains(s, "[[") || !strings.Contains(s, "]]") {
		fields = strings.Split(s, "|")
	} else {
		for {
			next := nextDelimiter(s)
			if next == -1 {
				break
			}
			fields = append(fields, s[:next])
			s = s[next+1:]
		}
		fields = append(fields, s)
	}

	name = strings.TrimSpace(fields[0])
	if strings.HasPrefix(name, "Template:") ||
		strings.HasPrefix(name, "template:") {
		name = name[9:]
	}

	params = make(map[string]string)
	for _, f := range fields[1:] {
		key, value := parseParameter(f)
		params[key] = value
	}
	return
}

// parseParameter parses a named template parameter.
// Whitespace is trimmed from the returned strings.
// If s does not contain "=", name is the empty string.
func parseParameter(s string) (name, value string) {
	eq := strings.Index(s, "=")
	return strings.TrimSpace(strings.TrimSuffix(s[:eq+1], "=")),
		strings.TrimSpace(s[eq+1:])
}

// nextDelimiter returns the index of the first "|" in s
// that is not enclosed within matching double brackets,
// or -1 if no unenclosed "|" is present in s.
func nextDelimiter(s string) int {
	lbr := strings.Index(s, "[[")
	pipe := strings.Index(s, "|")

	if lbr == -1 || pipe != -1 && pipe < lbr {
		return pipe
	}

	// Left double bracket occurs first; find the next right double bracket.
	rbroffset := strings.Index(s[lbr:], "]]")
	if rbroffset == -1 {
		return pipe
	}

	endbr := lbr + rbroffset + 2
	next := nextDelimiter(s[endbr:])
	if next == -1 {
		return -1
	}
	return endbr + next
}
