package dvorak

import (
	"errors"
	"fmt"
	"strings"
)

// Subpage is a subpage of a Dvorak deck.
type Subpage struct {
	// http://dvorakgame.co.uk/index.php/Template:Subpage

	// Page is the subpage's relative URL.
	Page string
}

// ParseSubpage parses a single instance of Template:Subpage source code.
func ParseSubpage(s string) (*Subpage, error) {
	var errInvalid = errors.New("invalid syntax")

	if !strings.HasPrefix(s, "{{") {
		return nil, errInvalid
	}
	s = s[2:]

	next := nextDelimiter(s)
	if next == -1 {
		return nil, errInvalid
	}

	name := strings.TrimSpace(s[:next])
	if strings.HasPrefix(name, "Template:") ||
		strings.HasPrefix(name, "template:") {
		name = name[9:]
	}
	if name != "Subpage" && name != "subpage" {
		return nil, errInvalid
	}
	s = s[next:]

	sp := &Subpage{}
	for len(s) > 2 {
		// Disallow nested templates by only accepting "|".
		if !strings.HasPrefix(s, "|") {
			return nil, errInvalid
		}
		s = s[1:]

		next := nextDelimiter(s)
		if next == -1 {
			return nil, errInvalid
		}

		switch name, value := parseParameter(s[:next]); name {
		case "page":
			sp.Page = value
		}

		s = s[next:]
	}
	if s != "}}" {
		return nil, errInvalid
	}

	if sp.Page == "" {
		return nil, fmt.Errorf("empty page value")
	}
	return sp, nil
}
