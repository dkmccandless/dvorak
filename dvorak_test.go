package dvorak

import (
	"image/color"
	"reflect"
	"testing"
)

func TestParsePage(t *testing.T) {
	for _, test := range []struct {
		s string
		p *Page
	}{
		{"", &Page{}},
		{"{{Subpage}}", &Page{}},
		{"{{card", &Page{}},
		{"{{card}}", &Page{Cards: []*Card{{BGColor: otherGray}}}},
		{
			"{{Subpage|page=Cards 1-100}}",
			&Page{Subpages: []*Subpage{{Page: "Cards 1-100"}}},
		},
		{
			`{{card|title=A|text={{card|title=B|type=Thing}}card|type=Action}}
			{{card|title=C}}`,
			&Page{
				Cards: []*Card{
					{Title: "B", Type: "Thing", BGColor: thingBlue},
					{Title: "C", BGColor: otherGray},
				},
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
			&Page{
				Cards: []*Card{
					{Title: "A", Type: "Action", BGColor: actionRed},
					{Title: "C", Type: "Letter", BGColor: otherGray},
					{Title: "E", BGColor: otherGray},
				},
			},
		},
		{
			`
				{{Subpage || hide = true | page=Cards 1-100 }}
				{{Subpage || hide = true | page=Cards 101-200 }}
				{{Subpage || hide = true }}
				{{card|title=A|type=Action|bgcolor=900}}
				<!-- {{card|title=|type=Action|text=|creator=|bgcolor=600}} -->
				<!-- {{card|title=|type=Thing|text=|creator=|bgcolor=006}} -->
				{{card|title=B|type=Thing|bgcolor=090}}
			`,
			&Page{
				Subpages: []*Subpage{
					{Page: "Cards 1-100"},
					{Page: "Cards 101-200"},
				},
				Cards: []*Card{
					{
						Title:   "A",
						Type:    "Action",
						BGColor: &color.RGBA{153, 0, 0, 255},
					},
					{
						Title:   "B",
						Type:    "Thing",
						BGColor: &color.RGBA{0, 153, 0, 255},
					},
				},
			},
		},
	} {
		p := ParsePage(test.s)
		if !reflect.DeepEqual(p, test.p) {
			t.Errorf("ParsePage(%q): got %v; want %v",
				test.s, p, test.p,
			)
		}
	}
}
