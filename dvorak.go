// Package dvorak parses source code templates used by the Dvorak game wiki.
package dvorak

import (
	"strings"
)

// Page contains the template information of a Dvorak wiki page.
type Page struct {
	// Subpages lists the page's Subpages, which may contain other Cards.
	Subpages []*Subpage

	// Cards lists the page's Cards.
	Cards []*Card
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
	for {
		cl := strings.Index(s, "}}")
		if cl == -1 {
			break
		}
		op := strings.LastIndex(s[:cl], "{{")
		if op == -1 {
			s = s[cl+2:]
			continue
		}

		t := s[op : cl+2]
		if c, err := ParseCard(t); err == nil {
			p.Cards = append(p.Cards, c)
		} else if sp, err := ParseSubpage(t); err == nil {
			p.Subpages = append(p.Subpages, sp)
		}
		s = s[cl+2:]
	}
	return p
}
