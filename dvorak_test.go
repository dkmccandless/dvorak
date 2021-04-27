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
		{"}}", -1},
		{"|", 0},
		{" |", 1},
		{"Action|longtext=true}}", 6},
		{"Action|creator=[[User:ABC|ABC]]", 6},
		{"[[User:ABC|ABC", 10},
		{"[[User:ABC|ABC]]", -1},
		{"[[User:ABC|ABC]]|longtext=true", 16},
		{"[[User:ABC|ABC]], [[User:DEF|DEF]], and others", -1},
	} {
		if i := nextDelimiter(test.s); i != test.i {
			t.Errorf("nextDelimiter(%q): got %v, want %v", test.s, i, test.i)
		}
	}
}

func TestParseParameter(t *testing.T) {
	for _, test := range []struct{ s, name, value string }{
		{"", "", ""},
		{"card", "", "card"},
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

func TestParseTemplate(t *testing.T) {
	for _, test := range []struct {
		s      string
		name   string
		params map[string]string
		isErr  bool
	}{
		{"{{card|title=ABC", "", nil, true},
		{"card|title=ABC}}", "", nil, true},
		{"{{card|title=ABC}}{{card|title=DEF}}", "", nil, true},
		{"{{}}", "", map[string]string{}, false},
		{"{{card}}", "card", map[string]string{}, false},
		{"{{Card}}", "Card", map[string]string{}, false},
		{"{{template:card}}", "card", map[string]string{}, false},
		{"{{template:Card}}", "Card", map[string]string{}, false},
		{"{{Template:card}}", "card", map[string]string{}, false},
		{"{{Template:Card}}", "Card", map[string]string{}, false},
		{
			"{{card|title=ABC|text=DEF}}",
			"card",
			map[string]string{"title": "ABC", "text": "DEF"},
			false,
		},
		{
			"{{ card | title = ABC | text = DEF }}",
			"card",
			map[string]string{"title": "ABC", "text": "DEF"},
			false,
		},
		{
			`{{card
			|title=ABC
			|text=DEF
			}}`,
			"card",
			map[string]string{"title": "ABC", "text": "DEF"},
			false,
		},
		{
			"{{card|creator=[[User:ABC|ABC]]|title=DEF}}",
			"card",
			map[string]string{"creator": "[[User:ABC|ABC]]", "title": "DEF"},
			false,
		},
	} {
		name, params, err := parseTemplate(test.s)
		if isErr := err != nil; isErr != test.isErr {
			t.Errorf("parseTemplate(%q): error=%v, want %v",
				test.s, isErr, test.isErr,
			)
		}
		if name != test.name || !reflect.DeepEqual(params, test.params) {
			t.Errorf("parseTemplate(%q): got %v, %v, want %v, %v",
				test.s, name, params, test.name, test.params,
			)
		}
	}
}

func TestParsePage(t *testing.T) {
	for _, test := range []struct {
		s string
		p page
	}{
		{"", page{}},
		{"{{Subpage}}", page{}},
		{"{{card", page{}},
		{"{{card}}", page{cards: []Card{{BGColor: otherGray}}}},
		{
			"{{Subpage|page=Cards 1-100}}",
			page{subpages: []subpage{{page: "Cards 1-100"}}},
		},
		{
			`{{card|title=A|text={{card|title=B|type=Thing}}card|type=Action}}
			{{card|title=C}}`,
			page{
				cards: []Card{
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
			page{
				cards: []Card{
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
			page{
				subpages: []subpage{
					{page: "Cards 1-100"},
					{page: "Cards 101-200"},
				},
				cards: []Card{
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
		p := parsePage(test.s)
		if !reflect.DeepEqual(p, test.p) {
			t.Errorf("parsePage(%q): got %v; want %v",
				test.s, p, test.p,
			)
		}
	}
}
