package dvorak

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

	// BGColor is the color of the card header background,
	// as a three- or six-digit hex triplet.
	// If omitted, Parse sets a default value according to the card type.
	BGColor string

	// CornerValue is an optional value to print in the card's top right corner.
	CornerValue string

	// Image is the filename of an image for the card.
	Image string

	// ImgBack is the optional color to be shown behind the card image,
	// as a three- or six-digit hex triplet.
	ImgBack string

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

// populateCard returns a Card populated with params.
func populateCard(params map[string]string) Card {
	return Card{
		Title:       params["title"],
		LongTitle:   params["longtitle"],
		Text:        params["text"],
		LongText:    params["longtext"],
		Type:        params["type"],
		BGColor:     params["bgcolor"],
		CornerValue: params["cornervalue"],
		Image:       params["image"],
		ImgBack:     params["imgback"],
		FlavorText:  params["flavortext"],
		Creator:     params["creator"],
		MiniCard:    params["minicard"],
	}
}

// withDefaultColor adds a default background color to c if it has none.
func withDefaultColor(c Card) Card {
	if c.BGColor == "" {
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
