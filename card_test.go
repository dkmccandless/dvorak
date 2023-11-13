package dvorak

import (
	"reflect"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestParseWikitext(t *testing.T) {
	for _, tt := range []struct{ value, want string }{
		{
			"''Italics'' '''Bold''' '''''Both'''''",
			"<i>Italics</i> <b>Bold</b> <b><i>Both</i></b>",
		},
	} {
		if got := parseWikitext(tt.value); dump(got) != tt.want {
			t.Errorf("parseWikitext(%q): got %v, want %v",
				tt.value, got, tt.want,
			)
		}
	}
}

func dump(frag []*html.Node) string {
	var b strings.Builder
	for _, n := range frag {
		err := html.Render(&b, n)
		if err != nil {
			panic(err)
		}
	}
	return b.String()
}

func TestPopulateCard(t *testing.T) {
	for _, test := range []struct {
		params map[string]string
		c      Card
	}{
		{nil, Card{}},
		{map[string]string{"": ""}, Card{}},
		{map[string]string{"": "ABC"}, Card{}},
		{map[string]string{"ABC": ""}, Card{}},
		{map[string]string{"title": "ABC"}, Card{Title: text("ABC")}},
		{map[string]string{"longtitle": "y"}, Card{LongTitle: true}},
		{map[string]string{"text": "ABC"}, Card{Text: text("ABC")}},
		{map[string]string{"longtext": "y"}, Card{LongText: true}},
		{map[string]string{"type": "Action"}, Card{Type: text("Action")}},
		{map[string]string{"bgcolor": "600"}, Card{BGColor: "600"}},
		{map[string]string{"bgcolor": "006"}, Card{BGColor: "006"}},
		{map[string]string{"bgcolor": "666"}, Card{BGColor: "666"}},
		{map[string]string{"bgcolor": "000"}, Card{BGColor: "000"}},
		{map[string]string{"bgcolor": "FFF"}, Card{BGColor: "FFF"}},
		{map[string]string{"cornervalue": "4"}, Card{CornerValue: text("4")}},
		{map[string]string{"image": "ABC.png"}, Card{Image: "ABC.png"}},
		{map[string]string{"imgback": "FFD700"}, Card{ImgBack: "FFD700"}},
		{map[string]string{"flavortext": "ABC"}, Card{FlavorText: text("ABC")}},
		{map[string]string{"creator": "ABC"}, Card{Creator: text("ABC")}},
		{map[string]string{"minicard": "y"}, Card{MiniCard: true}},
		{
			map[string]string{"title": "A", "type": "Action", "bgcolor": "006"},
			Card{Title: text("A"), Type: text("Action"), BGColor: "006"},
		},
		{
			map[string]string{"title": "B", "type": "Thing", "bgcolor": "600"},
			Card{Title: text("B"), Type: text("Thing"), BGColor: "600"},
		},
		{
			map[string]string{
				"title":   "Fishing Rod",
				"type":    "Action",
				"bgcolor": "369",
				"text":    "Gain control of a fish.",
				"creator": "Binarius",
			},
			Card{
				Title:   text("Fishing Rod"),
				Type:    text("Action"),
				Text:    text("Gain control of a fish."),
				BGColor: "369",
				Creator: text("Binarius"),
			},
		},
	} {
		c := populateCard(test.params)
		if !reflect.DeepEqual(c, test.c) {
			t.Errorf("populateCard(%v): got %v, want %v",
				test.params, c, test.c,
			)
		}
	}
}

func TestWithDefaultColor(t *testing.T) {
	for _, tt := range []struct{ typ, bgcolor, want string }{
		{"", "", otherGray},
		{"Special", "", otherGray},
		{"Action", "", actionRed},
		{"Thing", "", thingBlue},
		{"Action - Song", "", otherGray},
		{"Thing - Moon", "", otherGray},
		{"Action", thingBlue, thingBlue},
		{"Thing", actionRed, actionRed},
		{"Void", "000", "000"},
	} {
		got := withDefaultColor(tt.typ, tt.bgcolor)
		if got != tt.want {
			t.Errorf("withDefaultColor(%v, %v): got %v, want %v",
				tt.typ, tt.bgcolor, got, tt.want,
			)
		}
	}
}

func text(s string) []*html.Node {
	return []*html.Node{{
		Type: html.TextNode,
		Data: s,
	}}
}
