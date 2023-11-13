package dvorak

import (
	"strings"

	"github.com/dkmccandless/dvorak/sanitize"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Default header background colors
const (
	actionRed = "600"
	thingBlue = "006"
	otherGray = "666"
)

// Card is a Dvorak card.
type Card struct {
	// http://dvorakgame.co.uk/index.php/Template:Card

	// Title is the card's title.
	Title []*html.Node

	// LongTitle indicates that Title is too long to fit the standard header.
	LongTitle bool

	// Text is the card's rule text.
	Text []*html.Node

	// LongText indicates that Text is too long to fit the standard text area.
	LongText bool

	// Type is the card's type, usually "Action" or "Thing".
	Type []*html.Node

	// BGColor is the color of the card header background,
	// as a three- or six-digit hex triplet.
	// If omitted, Parse sets a default value according to the card type.
	BGColor string

	// CornerValue is an optional value to print in the card's top right corner.
	CornerValue []*html.Node

	// Image is the filename of an image for the card.
	Image string

	// ImgBack is the optional color to be shown behind the card image,
	// as a three- or six-digit hex triplet.
	ImgBack string

	// FlavorText is the card's flavor text. If not empty, this is displayed
	// under the rule text, separated by a horizontal line.
	FlavorText []*html.Node

	// Creator is the player who created the card. If not empty, this is
	// displayed at the bottom of the card.
	Creator []*html.Node

	// MiniCard indicates that the card is smaller than standard size,
	// e.g. for display in example texts.
	MiniCard bool

	// ID is the card's position within the deck.
	ID int
}

// populateCard returns a Card populated with params.
func populateCard(params map[string]string) Card {
	return Card{
		Title:       parseWikitext(params["title"]),
		LongTitle:   params["longtitle"] != "",
		Text:        parseWikitext(params["text"]),
		LongText:    params["longtext"] != "",
		Type:        parseWikitext(params["type"]),
		BGColor:     params["bgcolor"],
		CornerValue: parseWikitext(params["cornervalue"]),
		Image:       params["image"],
		ImgBack:     params["imgback"],
		FlavorText:  parseWikitext(params["flavortext"]),
		Creator:     parseWikitext(params["creator"]),
		MiniCard:    params["minicard"] != "",
	}
}

// parseWikitext parses wikitext and wiki markup as HTML.
func parseWikitext(s string) []*html.Node {
	s = sanitize.RemoveHTMLTags(s)

	s = replacePair(s, "'''''", "<b><i>", "</i></b>")
	s = replacePair(s, "'''", "<b>", "</b>")
	s = replacePair(s, "''", "<i>", "</i>")

	r := strings.NewReader(s)
	frag, err := html.ParseFragmentWithOptions(r,
		&html.Node{
			Type:     html.ElementNode,
			DataAtom: atom.Div,
			Data:     atom.Div.String(),
		},
		html.ParseOptionEnableScripting(false),
	)
	if err != nil {
		return []*html.Node{{
			Type: html.ErrorNode,
			Data: err.Error(),
		}}
	}
	return frag
}

// withDefaultColor returns a default color according to typ if bgcolor is empty.
func withDefaultColor(typ, bgcolor string) string {
	if bgcolor == "" {
		switch strings.ToLower(typ) {
		case "action":
			return actionRed
		case "thing":
			return thingBlue
		default:
			return otherGray
		}
	}
	return bgcolor
}
