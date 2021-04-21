// Package dvorak parses source code templates used by the Dvorak game wiki.
package dvorak

import (
	"encoding/hex"
	"errors"
	"image/color"
	"strings"
)

// Card is a Dvorak card.
type Card struct {
	// Title is the card's title.
	Title string

	// LongTitle indicates that Title is too long to fit the standard header.
	// If not empty, it expands the header to fit Title and Type.
	LongTitle string

	// Text is the card's rule text.
	Text string

	// LongText indicates that Text is too long to fit the standard text area.
	// If not empty, it reduces the Text and FlavorText font size.
	LongText string

	// Type is the card's type, usually "Action" or "Thing".
	Type string

	// BGColor is the color of the card header background.
	// If omitted, ParseCard sets a default value according to the card type.
	BGColor *color.RGBA

	// CornerValue is an optional value to print in the card's top right corner.
	CornerValue string

	// Image is the filename of an image for the card.
	Image string

	// ImgBack is the optional color to be shown behind the card image.
	ImgBack *color.RGBA

	// FlavorText is the card's flavor text. If not empty, this is displayed
	// under the rule text, separated by a horizontal line.
	FlavorText string

	// Creator is the player who created the card. If not empty, this is
	// displayed at the bottom of the card.
	Creator string

	// MiniCard indicates that the card is smaller than standard size,
	// e.g. for display in example texts.
	// If not empty, it reduces the card's size.
	MiniCard string
}

// ParseDeck parses a list of Template:Card values.
func ParseDeck(s string) []*Card {
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

	var cards []*Card
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

		c, err := ParseCard(s[op : cl+2])
		if err == nil {
			cards = append(cards, c)
		}
		s = s[cl+2:]
	}
	return cards
}

// ParseCard parses a single instance of Template:Card source code.
func ParseCard(s string) (*Card, error) {
	// http://dvorakgame.co.uk/index.php/Template:Card
	// https://meta.wikimedia.org/wiki/Help:Template

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
	if name != "Card" && name != "card" {
		return nil, errInvalid
	}
	s = s[next:]

	c := &Card{}
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
		case "title":
			c.Title = value
		case "longtitle":
			c.LongTitle = value
		case "text":
			c.Text = value
		case "longtext":
			c.LongText = value
		case "type":
			c.Type = value
		case "bgcolor":
			c.BGColor = parseRGB(value)
		case "cornervalue":
			c.CornerValue = value
		case "image":
			c.Image = value
		case "imgback":
			c.ImgBack = parseRGB(value)
		case "flavortext":
			c.FlavorText = value
		case "creator":
			c.Creator = value
		case "minicard":
			c.MiniCard = value
		}

		s = s[next:]
	}
	if s != "}}" {
		return nil, errInvalid
	}

	if c.BGColor == nil {
		c.BGColor = defaultBGColor(c.Type)
	}
	return c, nil
}

// parseParameter parses a named template parameter.
// Whitespace is trimmed from the returned strings.
// If s does not contain "=", parseParameter returns "", "" instead.
func parseParameter(s string) (name, value string) {
	eq := strings.Index(s, "=")
	if eq == -1 {
		return
	}
	return strings.TrimSpace(s[:eq]), strings.TrimSpace(s[eq+1:])
}

// parseRGB parses a hex string and returns the corresponding color.
// The string must be length 3 or 6 and contain only hexadecimal characters.
// Otherwise, parseRGB returns nil.
func parseRGB(s string) *color.RGBA {
	if len(s) == 3 {
		s = string([]byte{s[0], s[0], s[1], s[1], s[2], s[2]})
	}
	if len(s) != 6 {
		return nil
	}
	rgb, err := hex.DecodeString(s)
	if err != nil {
		return nil
	}
	return &color.RGBA{rgb[0], rgb[1], rgb[2], 255}
}

// Default header background colors
var (
	actionRed = parseRGB("600")
	thingBlue = parseRGB("006")
	otherGray = parseRGB("666")
)

// defaultBGColor returns the default header background color for typ.
func defaultBGColor(typ string) *color.RGBA {
	switch typ {
	case "Action":
		return actionRed
	case "Thing":
		return thingBlue
	default:
		return otherGray
	}
}

// nextDelimiter returns the index of the first delimiter character at the end
// of the first field in s. This is the first instance of "|" not contained
// within a matching pair of double brackets, or else the first "}}".
func nextDelimiter(s string) int {
	// minExists returns its smallest non-negative argument,
	// or -1 if both arguments are negative.
	minExists := func(x, y int) int {
		switch {
		case x < 0:
			if y < 0 {
				return -1
			}
			return y
		case y < 0, x < y:
			return x
		default:
			return y
		}
	}

	lbr := strings.Index(s, "[[")
	pipe := strings.Index(s, "|")
	cl := strings.Index(s, "}}")

	if minExists(lbr, pipe) == pipe || minExists(lbr, cl) == cl {
		return minExists(pipe, cl)
	}

	// Left double bracket occurs first; find the next right double bracket.
	rbroffset := strings.Index(s[lbr:], "]]")
	if rbroffset == -1 {
		return minExists(pipe, cl)
	}

	lenbr := lbr + rbroffset + 2
	next := nextDelimiter(s[lenbr:])
	if next == -1 {
		return -1
	}
	return lenbr + next
}
