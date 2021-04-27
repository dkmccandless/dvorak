package dvorak

import (
	"encoding/hex"
	"image/color"
)

// Card is a Dvorak card.
type Card struct {
	// http://dvorakgame.co.uk/index.php/Template:Card

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

// PopulateCard returns a Card populated with params.
func PopulateCard(params map[string]string) Card {
	c := Card{
		Title:       params["title"],
		LongTitle:   params["longtitle"],
		Text:        params["text"],
		LongText:    params["longtext"],
		Type:        params["type"],
		BGColor:     parseRGB(params["bgcolor"]),
		CornerValue: params["cornervalue"],
		Image:       params["image"],
		ImgBack:     parseRGB(params["imgback"]),
		FlavorText:  params["flavortext"],
		Creator:     params["creator"],
		MiniCard:    params["minicard"],
	}
	return c
}

// withDefaultColor adds a default background color to c if it has none.
func withDefaultColor(c Card) Card {
	if c.BGColor == nil {
		switch c.Type {
		case "Action":
			c.BGColor = actionRed
		case "Thing":
			c.BGColor = thingBlue
		default:
			c.BGColor = otherGray
		}
	}
	return c
}

// Default header background colors
var (
	actionRed = parseRGB("600")
	thingBlue = parseRGB("006")
	otherGray = parseRGB("666")
)

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
