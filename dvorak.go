// Package dvorak parses source code templates used by the Dvorak game wiki.
package dvorak

// Card is a Dvorak card.
type Card struct {
	// Title is the card's title.
	Title string

	// Text is the card's rule text.
	Text string

	// Type is the card's type, usually "Action" or "Thing".
	Type string

	// BGColor is the color of the card header background, as a hex value.
	// If omitted, the default values according to card type are
	// 600 for "Action", 006 for "Thing", or else 666.
	BGColor string

	// CornerValue is an optional value to print in the card's top right corner.
	CornerValue string

	// Image is the filename of an image for the card.
	Image string

	// ImgBack is the color to be shown behind the card image, as a hex value.
	ImgBack string

	// FlavorText is the card's flavor text. If not empty, this is shown under
	// the rule text, separated by a horizontal line.
	FlavorText string

	// Creator is the player who created the card. If not empty, this value is
	// displayed at the bottom of the card.
	Creator string

	// LongTitle, LongText, and MiniCard are optional string flags that modify
	// the layout of the card if they are not the empty string.

	// LongTitle indicates the title is too long to fit the standard header.
	// If not empty, it expands the header to fit the title and type.
	LongTitle string

	// LongText indicates the text is too long to fit the standard text area.
	// If not empty, it reduces the text and flavor text font size.
	LongText string

	// MiniCard indicates that the card is smaller than standard size,
	// e.g. for display in example texts.
	// If not empty, it reduces the card's size.
	MiniCard string
}
