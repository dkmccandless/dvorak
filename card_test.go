package dvorak

import (
	"reflect"
	"testing"
)

func TestPopulateCard(t *testing.T) {
	for _, test := range []struct {
		params map[string]string
		c      Card
	}{
		{nil, Card{}},
		{map[string]string{"": ""}, Card{}},
		{map[string]string{"": "ABC"}, Card{}},
		{map[string]string{"ABC": ""}, Card{}},
		{map[string]string{"title": "ABC"}, Card{Title: "ABC"}},
		{map[string]string{"longtitle": "y"}, Card{LongTitle: "y"}},
		{map[string]string{"text": "ABC"}, Card{Text: "ABC"}},
		{map[string]string{"longtext": "y"}, Card{LongText: "y"}},
		{map[string]string{"type": "Action"}, Card{Type: "Action"}},
		{map[string]string{"bgcolor": "600"}, Card{BGColor: "600"}},
		{map[string]string{"bgcolor": "006"}, Card{BGColor: "006"}},
		{map[string]string{"bgcolor": "666"}, Card{BGColor: "666"}},
		{map[string]string{"bgcolor": "000"}, Card{BGColor: "000"}},
		{map[string]string{"bgcolor": "FFF"}, Card{BGColor: "FFF"}},
		{map[string]string{"cornervalue": "4"}, Card{CornerValue: "4"}},
		{map[string]string{"image": "ABC.png"}, Card{Image: "ABC.png"}},
		{map[string]string{"imgback": "FFD700"}, Card{ImgBack: "FFD700"}},
		{map[string]string{"flavortext": "ABC"}, Card{FlavorText: "ABC"}},
		{map[string]string{"creator": "ABC"}, Card{Creator: "ABC"}},
		{map[string]string{"minicard": "y"}, Card{MiniCard: "y"}},
		{
			map[string]string{"title": "A", "type": "Action", "bgcolor": "006"},
			Card{Title: "A", Type: "Action", BGColor: "006"},
		},
		{
			map[string]string{"title": "B", "type": "Thing", "bgcolor": "600"},
			Card{Title: "B", Type: "Thing", BGColor: "600"},
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
				BGColor: "369",
				Creator: "Binarius",
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
	for _, test := range []struct{ in, c Card }{
		{Card{}, Card{BGColor: otherGray}},
		{Card{Type: "Special"}, Card{Type: "Special", BGColor: otherGray}},
		{Card{Type: "Action"}, Card{Type: "Action", BGColor: actionRed}},
		{Card{Type: "Thing"}, Card{Type: "Thing", BGColor: thingBlue}},
		{
			Card{Type: "Action - Song"},
			Card{Type: "Action - Song", BGColor: otherGray},
		},
		{
			Card{Type: "Thing - Moon"},
			Card{Type: "Thing - Moon", BGColor: otherGray},
		},
		{
			Card{Type: "Action", BGColor: thingBlue},
			Card{Type: "Action", BGColor: thingBlue},
		},
		{
			Card{Type: "Thing", BGColor: actionRed},
			Card{Type: "Thing", BGColor: actionRed},
		},
		{
			Card{Type: "Void", BGColor: "000"},
			Card{Type: "Void", BGColor: "000"},
		},
	} {
		c := withDefaultColor(test.in)
		if !reflect.DeepEqual(c, test.c) {
			t.Errorf("withDefaultColor(%v): got %v, want %v",
				test.in, c, test.c,
			)
		}
	}
}
