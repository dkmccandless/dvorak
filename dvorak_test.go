package dvorak

import (
	"image/color"
	"reflect"
	"testing"
)

func TestNextDelimiter(t *testing.T) {
	for _, test := range []struct {
		s string
		i int
	}{
		{"", -1},
		{"[[", -1},
		{"]]", -1},
		{"|", 0},
		{"}}", 0},
		{" |", 1},
		{" }}", 1},
		{"Action}}", 6},
		{"Action|longtext=true}}", 6},
		{"[[Dvorak]]", -1},
		{"[[Dvorak]]}}", 10},
		{"[[Dvorak]]|longtext=true}}", 10},
		{"[[User:ABC|ABC]]", -1},
		{"[[User:ABC|ABC}}", 10},
		{"[[User:ABC|ABC]]}}", 16},
		{"[[User:ABC|ABC]]|longtext=true}}", 16},
		{"[[User:ABC|ABC]], [[User:DEF|DEF]], and others", -1},
		{"[[|]], [[|]]}}", 12},
		{"[[|]], [[|]}}", 9},
		{"Action|creator=[[User:ABC|ABC]]}}", 6},
	} {
		if i := nextDelimiter(test.s); i != test.i {
			t.Errorf("nextDelimiter(%q): got %v, want %v", test.s, i, test.i)
		}
	}
}

func TestParseParameter(t *testing.T) {
	for _, test := range []struct{ s, name, value string }{
		{"", "", ""},
		{"title", "", ""},
		{"title=", "title", ""},
		{"title=ABC", "title", "ABC"},
		{"title = ABC", "title", "ABC"},
		{" title = ABC ", "title", "ABC"},
	} {
		name, value := parseParameter(test.s)
		if name != test.name || value != test.value {
			t.Errorf("parseParameter(%q): got %q, %q; want %q, %q",
				test.s, name, value, test.name, test.value,
			)
		}
	}
}

func TestParseRGB(t *testing.T) {
	for _, test := range []struct {
		s    string
		rgba *color.RGBA
	}{
		{"", nil},
		{"1", nil},
		{"xyz", nil},
		{"#123", nil},
		{"#808080", nil},

		{"000", &color.RGBA{0, 0, 0, 255}},
		{"123", &color.RGBA{17, 34, 51, 255}},
		{"ABC", &color.RGBA{170, 187, 204, 255}},
		{"abc", &color.RGBA{170, 187, 204, 255}},
		{"123456", &color.RGBA{18, 52, 86, 255}},
		{"AABBCC", &color.RGBA{170, 187, 204, 255}},
		{"b0a9e4", &color.RGBA{176, 169, 228, 255}},
	} {
		if rgba := parseRGB(test.s); !reflect.DeepEqual(rgba, test.rgba) {
			t.Errorf("parseRGB(%q): got %v, want %v", test.s, rgba, test.rgba)
		}
	}
}

func TestParseCard(t *testing.T) {
	var blank = &Card{BGColor: otherGray}

	for _, test := range []struct {
		s string
		c *Card
	}{
		{"", nil},
		{"{{}}", nil},
		{"{{card  ", nil},
		{"{{card||", nil},
		{"{{CARD}}", nil},
		{"{{card|title=ABC", nil},
		{"{{card}}|title=ABC", nil},
		{"{{card|title=ABC}}{{card}}", nil},
		{"{{card}}{{card|title=ABC}}", nil},
		{"{{card|title=ABC}}{{card|title=DEF}}", nil},
		{"{{card|title=ABC|text={{card|title=DEF}}}}", nil},

		{"{{card}}", blank},
		{"{{Card}}", blank},
		{"{{card }}", blank},
		{"{{ card}}", blank},
		{"{{\ncard\n}}", blank},
		{"{{template:card}}", blank},
		{"{{template:Card}}", blank},
		{"{{Template:card}}", blank},
		{"{{Template:Card}}", blank},
		{"{{card|}}", blank},
		{"{{card| }}", blank},
		{"{{card||}}", blank},
		{"{{card|ABC}}", blank},
		{"{{card|title=ABC}}", &Card{Title: "ABC", BGColor: otherGray}},
		{"{{ card | title = ABC }}", &Card{Title: "ABC", BGColor: otherGray}},
		{"{{card|longtitle=y}}", &Card{LongTitle: "y", BGColor: otherGray}},
		{"{{card|text=ABC}}", &Card{Text: "ABC", BGColor: otherGray}},
		{"{{card|longtext=y}}", &Card{LongText: "y", BGColor: otherGray}},
		{"{{card|type=Action}}", &Card{Type: "Action", BGColor: actionRed}},
		{"{{card|type=Thing}}", &Card{Type: "Thing", BGColor: thingBlue}},
		{"{{card|type=Action - Song}}", &Card{Type: "Action - Song", BGColor: otherGray}},
		{"{{card|type=Thing - Moon}}", &Card{Type: "Thing - Moon", BGColor: otherGray}},
		{"{{card|bgcolor=600}}", &Card{BGColor: actionRed}},
		{"{{card|bgcolor=006}}", &Card{BGColor: thingBlue}},
		{"{{card|bgcolor=666}}", &Card{BGColor: otherGray}},
		{"{{card|bgcolor=000}}", &Card{BGColor: &color.RGBA{A: 255}}},
		{"{{card|bgcolor=FFF}}", &Card{BGColor: &color.RGBA{255, 255, 255, 255}}},
		{"{{card|cornervalue=4}}", &Card{CornerValue: "4", BGColor: otherGray}},
		{"{{card|image=ABC.png}}", &Card{Image: "ABC.png", BGColor: otherGray}},
		{"{{card|imgback=FFD700}}", &Card{ImgBack: &color.RGBA{255, 215, 0, 255}, BGColor: otherGray}},
		{"{{card|flavortext=ABC}}", &Card{FlavorText: "ABC", BGColor: otherGray}},
		{"{{card|creator=ABC}}", &Card{Creator: "ABC", BGColor: otherGray}},
		{"{{card|minicard=y}}", &Card{MiniCard: "y", BGColor: otherGray}},
		{
			"{{card|title=A|type=Action|bgcolor=006}}",
			&Card{
				Title:   "A",
				Type:    "Action",
				BGColor: &color.RGBA{0, 0, 102, 255},
			},
		},
		{
			"{{card|title=B|type=Thing|bgcolor=600}}",
			&Card{
				Title:   "B",
				Type:    "Thing",
				BGColor: &color.RGBA{102, 0, 0, 255},
			},
		},
		{
			`{{card
			|title=Fishing Rod
			|type=Action
			|bgcolor=369
			|text=Gain control of a fish.
			|creator=Binarius
			}}`,
			&Card{
				Title:   "Fishing Rod",
				Type:    "Action",
				Text:    "Gain control of a fish.",
				BGColor: &color.RGBA{51, 102, 153, 255},
				Creator: "Binarius",
			},
		},
	} {
		c, err := ParseCard(test.s)
		if !reflect.DeepEqual(c, test.c) || (err != nil) != (c == nil) {
			t.Errorf("ParseCard(%q): got %v, %v; want %v",
				test.s, c, err, test.c,
			)
		}
	}
}

func TestParseSubpage(t *testing.T) {
	var sp1 = &subpage{page: "Cards 1-100"}

	for _, test := range []struct {
		s  string
		sp *subpage
	}{
		{"", nil},
		{"{{}}", nil},
		{"{{subpage  ", nil},
		{"{{subpage||", nil},
		{"{{SUBPAGE}}", nil},
		{"{{subpage|ABC}}", nil},
		{"{{subpage|page=ABC", nil},
		{"{{subpage|page=ABC}}|page=DEF", nil},
		{"{{subpage|page=ABC}}{{subpage|page=DEF}}", nil},

		{"{{subpage|page=Cards 1-100}}", sp1},
		{"{{Subpage|page=Cards 1-100}}", sp1},
		{"{{subpage |page=Cards 1-100}}", sp1},
		{"{{ subpage|page=Cards 1-100}}", sp1},
		{"{{\nsubpage\n|page=Cards 1-100}}", sp1},
		{"{{template:subpage|page=Cards 1-100}}", sp1},
		{"{{template:Subpage|page=Cards 1-100}}", sp1},
		{"{{Template:subpage|page=Cards 1-100}}", sp1},
		{"{{Template:Subpage|page=Cards 1-100}}", sp1},
		{"{{Subpage || hide = true | page=Cards 1-100 }}", sp1},
	} {
		sp, err := parseSubpage(test.s)
		if !reflect.DeepEqual(sp, test.sp) || (err != nil) != (sp == nil) {
			t.Errorf("parseSubpage(%q): got %v, %v; want %v",
				test.s, sp, err, test.sp,
			)
		}
	}
}

func TestParseDeck(t *testing.T) {
	for _, test := range []struct {
		s     string
		cards []*Card
	}{
		{"", nil},
		{"{{card", nil},
		{"{{card}}", []*Card{{BGColor: otherGray}}},
		{
			`{{card|title=A|text={{card|title=B|type=Thing}}card|type=Action}}
			{{card|title=C}}`,
			[]*Card{
				{Title: "B", Type: "Thing", BGColor: thingBlue},
				{Title: "C", BGColor: otherGray},
			},
		},
		{
			`
				{{card| title = A | type = Action }}
				card|title=B|type=Thing}}
				{{card|title=C|type=Letter}}
				{{card|title=D
				{{card|title=E}}
			`,
			[]*Card{
				{Title: "A", Type: "Action", BGColor: actionRed},
				{Title: "C", Type: "Letter", BGColor: otherGray},
				{Title: "E", BGColor: otherGray},
			},
		},
		{
			`
				{{card|title=A|type=Action|bgcolor=900}}
				<!-- {{card|title=|type=Action|text=|creator=|bgcolor=600}} -->
				<!-- {{card|title=|type=Thing|text=|creator=|bgcolor=006}} -->
				{{card|title=B|type=Thing|bgcolor=090}}
			`,
			[]*Card{
				{Title: "A", Type: "Action", BGColor: &color.RGBA{153, 0, 0, 255}},
				{Title: "B", Type: "Thing", BGColor: &color.RGBA{0, 153, 0, 255}},
			},
		},
	} {
		cards := ParseDeck(test.s)
		if !reflect.DeepEqual(cards, test.cards) {
			t.Errorf("ParseDeck(%q): got %v, want %v",
				test.s, cards, test.cards,
			)
		}
	}
}
