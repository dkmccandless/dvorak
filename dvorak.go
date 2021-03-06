// Package dvorak parses source code templates used by the Dvorak game wiki.
package dvorak

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// page contains the template information of a Dvorak wiki page.
type page struct {
	// subpages lists any subpages of the page.
	subpages []subpage

	// cards lists the page's cards.
	cards []Card
}

// Get returns the source code of a Dvorak deck,
// beginning with its subpages in order, if any.
func Get(rawURL string) ([]byte, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	if host := u.Hostname(); strings.ToLower(host) != "dvorakgame.co.uk" {
		return nil, fmt.Errorf("invalid host %q", host)
	}
	u.RawQuery = "action=raw"
	path := u.EscapedPath()

	main, err := readPage(u.String())
	if err != nil {
		return nil, err
	}
	var b []byte
	for _, sp := range parsePage(main).subpages {
		log.Print(sp.page)
		u.Path = path + "/" + sp.page
		sb, err := readPage(u.String())
		if err != nil {
			return nil, err
		}
		b = append(b, sb...)
	}

	return append(b, main...), nil
}

// readPage returns the body of the page at url.
// It returns an error if url cannot be accessed or read from.
func readPage(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%v: status %v", url, r.StatusCode)
	}
	defer r.Body.Close()

	return io.ReadAll(r.Body)
}

// Parse returns the Cards in b.
func Parse(b []byte) []Card {
	return parsePage(b).cards
}

// parsePage parses a page of wiki source code.
func parsePage(b []byte) page {
	s := string(b)

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

	p := page{}
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
			p.cards = append(p.cards, withDefaultColor(populateCard(params)))
		case "Subpage", "subpage":
			sp, err := populateSubpage(params)
			if err != nil {
				continue
			}
			p.subpages = append(p.subpages, sp)
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
