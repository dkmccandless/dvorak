package dvorak

import (
	"image/color"
	"reflect"
	"testing"
)

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

func TestPopulateCard(t *testing.T) {
	var blank = Card{BGColor: otherGray}

	for _, test := range []struct {
		params map[string]string
		c      Card
	}{
		{nil, blank},
		{map[string]string{"": ""}, blank},
		{map[string]string{"": "ABC"}, blank},
		{map[string]string{"ABC": ""}, blank},
		{map[string]string{"title": "ABC"}, Card{Title: "ABC", BGColor: otherGray}},
		{map[string]string{"longtitle": "y"}, Card{LongTitle: "y", BGColor: otherGray}},
		{map[string]string{"text": "ABC"}, Card{Text: "ABC", BGColor: otherGray}},
		{map[string]string{"longtext": "y"}, Card{LongText: "y", BGColor: otherGray}},
		{map[string]string{"type": "Action"}, Card{Type: "Action", BGColor: actionRed}},
		{map[string]string{"type": "Thing"}, Card{Type: "Thing", BGColor: thingBlue}},
		{map[string]string{"type": "Action - Song"}, Card{Type: "Action - Song", BGColor: otherGray}},
		{map[string]string{"type": "Thing - Moon"}, Card{Type: "Thing - Moon", BGColor: otherGray}},
		{map[string]string{"bgcolor": "600"}, Card{BGColor: actionRed}},
		{map[string]string{"bgcolor": "006"}, Card{BGColor: thingBlue}},
		{map[string]string{"bgcolor": "666"}, Card{BGColor: otherGray}},
		{map[string]string{"bgcolor": "000"}, Card{BGColor: &color.RGBA{A: 255}}},
		{map[string]string{"bgcolor": "FFF"}, Card{BGColor: &color.RGBA{255, 255, 255, 255}}},
		{map[string]string{"cornervalue": "4"}, Card{CornerValue: "4", BGColor: otherGray}},
		{map[string]string{"image": "ABC.png"}, Card{Image: "ABC.png", BGColor: otherGray}},
		{map[string]string{"imgback": "FFD700"}, Card{ImgBack: &color.RGBA{255, 215, 0, 255}, BGColor: otherGray}},
		{map[string]string{"flavortext": "ABC"}, Card{FlavorText: "ABC", BGColor: otherGray}},
		{map[string]string{"creator": "ABC"}, Card{Creator: "ABC", BGColor: otherGray}},
		{map[string]string{"minicard": "y"}, Card{MiniCard: "y", BGColor: otherGray}},
		{
			map[string]string{
				"title":   "A",
				"type":    "Action",
				"bgcolor": "006",
			},
			Card{
				Title:   "A",
				Type:    "Action",
				BGColor: &color.RGBA{0, 0, 102, 255},
			},
		},
		{
			map[string]string{
				"title":   "B",
				"type":    "Thing",
				"bgcolor": "600",
			},
			Card{
				Title:   "B",
				Type:    "Thing",
				BGColor: &color.RGBA{102, 0, 0, 255},
			},
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
				Title:   "Fishing Rod",
				Type:    "Action",
				Text:    "Gain control of a fish.",
				BGColor: &color.RGBA{51, 102, 153, 255},
				Creator: "Binarius",
			},
		},
	} {
		c := PopulateCard(test.params)
		if !reflect.DeepEqual(c, test.c) {
			t.Errorf("PopulateCard(%v): got %v, want %v",
				test.params, c, test.c,
			)
		}
	}
}
