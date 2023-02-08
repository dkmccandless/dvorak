package dvorak

import (
	"strings"

	"golang.org/x/net/html"
)

// escapeNonTags escapes special characters in tag-like substrings of s
// that do not begin with a recognized HTML tag or attribute string.
//
// escapeNonTags only considers strings lexically,
// without any syntactic context.
// It does not escape syntactically invalid strings
// that start with an HTML tag or attribute.
func escapeNonTags(s string) string {
	var b strings.Builder
	z := html.NewTokenizer(strings.NewReader(s))
	for {
		tt := z.Next()
		switch raw := z.Raw(); tt {
		case html.ErrorToken:
			return b.String()
		case html.StartTagToken, html.EndTagToken, html.SelfClosingTagToken:
			if z.Token().DataAtom == 0 {
				b.WriteString(html.EscapeString(string(raw)))
			} else {
				b.Write(raw)
			}
		default:
			b.Write(raw)
		}
	}
}
