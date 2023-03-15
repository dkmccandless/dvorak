package html

// attributesAllowed maps HTML elements to their allowed attributes.
var attributesAllowed = func() map[string]map[string]bool {
	commonAttrs := []string{
		// HTML
		"id",
		"class",
		"style",
		"lang",
		"dir",
		"title",
		"tabindex",

		// WAI-ARIA
		"aria-describedby",
		"aria-flowto",
		"aria-hidden",
		"aria-label",
		"aria-labelledby",
		"aria-owns",
		"role",

		// RDFa
		// These attributes are specified in section 9 of
		// https://www.w3.org/TR/2008/REC-rdfa-syntax-20081014
		"about",
		"property",
		"resource",
		"datatype",
		"typeof",

		// Microdata. These are specified by
		// https://html.spec.whatwg.org/multipage/microdata.html#the-microdata-model
		"itemid",
		"itemprop",
		"itemref",
		"itemscope",
		"itemtype",
	}
	tableAlignAttrs := []string{"align", "valign"}
	tableCellAttrs := []string{
		"abbr",
		"axis",
		"headers",
		"scope",
		"rowspan",
		"colspan",
		"nowrap",  // deprecated
		"width",   // deprecated
		"height",  // deprecated
		"bgcolor", // deprecated
	}
	commonAlignAttrs := append(commonAttrs, tableAlignAttrs...)
	commonAlignCellAttrs := append(commonAlignAttrs, tableCellAttrs...)

	common := makeSet(commonAttrs...)
	block := appendSet(common, "align")
	commonAlign := makeSet(commonAlignAttrs...)
	commonAlignCell := makeSet(commonAlignCellAttrs...)

	// Numbers refer to sections in HTML 4.01 standard describing the element.
	// See: https://www.w3.org/TR/html4/
	return map[string]map[string]bool{
		// 7.5.4
		"div":    block,
		"center": common, // deprecated
		"span":   common,

		// 7.5.5
		"h1": block,
		"h2": block,
		"h3": block,
		"h4": block,
		"h5": block,
		"h6": block,

		// 7.5.6
		// address

		// 8.2.4
		"bdo": common,

		// 9.2.1
		"em":     common,
		"strong": common,
		"cite":   common,
		"dfn":    common,
		"code":   common,
		"samp":   common,
		"kbd":    common,
		"var":    common,
		"abbr":   common,
		// acronym

		// 9.2.2
		"blockquote": appendSet(common, "cite"),
		"q":          appendSet(common, "cite"),

		// 9.2.3
		"sub": common,
		"sup": common,

		// 9.3.1
		"p": block,

		// 9.3.2
		"br": appendSet(common, "clear"),

		// https://www.w3.org/TR/html5/text-level-semantics.html#the-wbr-element
		"wbr": common,

		// 9.3.4
		"pre": appendSet(common, "width"),

		// 9.4
		"ins": appendSet(common, "cite", "datetime"),
		"del": appendSet(common, "cite", "datetime"),

		// 10.2
		"ul": appendSet(common, "type"),
		"ol": appendSet(common, "type", "start", "reversed"),
		"li": appendSet(common, "type", "value"),

		// 10.3
		"dl": common,
		"dd": common,
		"dt": common,

		// 11.2.1
		"table": appendSet(common,
			"summary", "width", "border", "frame",
			"rules", "cellspacing", "cellpadding",
			"align", "bgcolor",
		),

		// 11.2.2
		"caption": block,

		// 11.2.3
		"thead": common,
		"tfoot": common,
		"tbody": common,

		// 11.2.4
		"colgroup": appendSet(common, "span"),
		"col":      appendSet(common, "span"),

		// 11.2.5
		"tr": appendSet(commonAlign, "bgcolor"),

		// 11.2.6
		"td": commonAlignCell,
		"th": commonAlignCell,

		// 12.2
		// NOTE: <a> is not allowed directly, but this list of allowed
		// attributes is used from the Parser object
		"a": appendSet(common, "href", "rel", "rev"), // rel/rev esp. for RDFa

		// 13.2
		// Not usually allowed, but may be used for extension-style hooks
		// such as <math> when it is rasterized, or if $wgAllowImageTag is
		// true
		"img": appendSet(common, "alt", "src", "width", "height", "srcset"),
		// Attributes for A/V tags added in T163583 / T133673
		"audio":  appendSet(common, "controls", "preload", "width", "height"),
		"video":  appendSet(common, "poster", "controls", "preload", "width", "height"),
		"source": appendSet(common, "type", "src"),
		"track":  appendSet(common, "type", "src", "srclang", "kind", "label"),

		// 15.2.1
		"tt":     common,
		"b":      common,
		"i":      common,
		"big":    common,
		"small":  common,
		"strike": common,
		"s":      common,
		"u":      common,

		// 15.2.2
		"font": appendSet(common, "size", "color", "face"),
		// basefont

		// 15.3
		"hr": appendSet(common, "width"),

		// HTML Ruby annotation text module, simple ruby only.
		// https://www.w3.org/TR/html5/text-level-semantics.html#the-ruby-element
		"ruby": common,
		// rbc
		"rb":  common,
		"rp":  common,
		"rt":  common, // $merge( $common, [ "rbspan" ] ),
		"rtc": common,

		// MathML root element, where used for extensions
		// 'title' may not be 100% valid here; it's XHTML
		// https://www.w3.org/TR/REC-MathML/
		"math": makeSet("class", "style", "id", "title"),

		// HTML 5 section 4.5
		"figure":     common,
		"figcaption": common,

		// HTML 5 section 4.6
		"bdi": common,

		// HTML5 elements, defined by:
		// https://html.spec.whatwg.org/multipage/semantics.html#the-data-element
		"data": appendSet(common, "value"),
		"time": appendSet(common, "datetime"),
		"mark": common,

		// meta and link are only permitted by internalRemoveHtmlTags when Microdata
		// is enabled so we don't bother adding a conditional to hide these
		// Also meta and link are only valid in WikiText as Microdata elements
		// (ie: validateTag rejects tags missing the attributes needed for Microdata)
		// So we don't bother including $common attributes that have no purpose.
		"meta": makeSet("itemprop", "content"),
		"link": makeSet("itemprop", "href", "title"),

		// HTML 5 section 4.3.5
		"aside": common,
	}
}()
