package dvorak

import (
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
		s        string
		rawLinks bool
		name     string
		params   map[string]string
		isErr    bool
	}{
		{"{{card|title=ABC", false, "", nil, true},
		{"card|title=ABC}}", false, "", nil, true},
		{"{{card|title=ABC}}{{card|title=DEF}}", false, "", nil, true},
		{"{{}}", false, "", map[string]string{}, false},
		{"{{card}}", false, "card", map[string]string{}, false},
		{"{{Card}}", false, "Card", map[string]string{}, false},
		{"{{template:card}}", false, "card", map[string]string{}, false},
		{"{{template:Card}}", false, "Card", map[string]string{}, false},
		{"{{Template:card}}", false, "card", map[string]string{}, false},
		{"{{Template:Card}}", false, "Card", map[string]string{}, false},
		{
			"{{card|title=ABC|text=DEF}}",
			false,
			"card",
			map[string]string{"title": "ABC", "text": "DEF"},
			false,
		},
		{
			"{{ card | title = ABC | text = DEF }}",
			false,
			"card",
			map[string]string{"title": "ABC", "text": "DEF"},
			false,
		},
		{
			`{{card
			|title=ABC
			|text=DEF
			}}`,
			false,
			"card",
			map[string]string{"title": "ABC", "text": "DEF"},
			false,
		},
		{
			"{{card|creator=[[User:ABC|ABC]]|title=DEF}}",
			false,
			"card",
			map[string]string{"creator": "ABC", "title": "DEF"},
			false,
		},
		{
			"{{card|creator=[[User:ABC|ABC]]|title=DEF}}",
			true,
			"card",
			map[string]string{"creator": "[[User:ABC|ABC]]", "title": "DEF"},
			false,
		},
		{
			"{{card|text=[[File: ABC.jpg]]DEF}}",
			false,
			"card",
			map[string]string{"text": "DEF", "image": "ABC.jpg"},
			false,
		},
		{
			"{{card|creator=[[User:ABC|ABC]] ([[User talk:ABC|talk]])|title=DEF}}",
			false,
			"card",
			map[string]string{"creator": "ABC", "title": "DEF"},
			false,
		},
		{
			"{{card|creator=[[User:ABC|ABC]] ([[User talk:ABC|talk]]) 21:33, 25 July 2012 (UTC)|title=DEF}}",
			false,
			"card",
			map[string]string{"creator": "ABC", "title": "DEF"},
			false,
		},
		{
			"{{card|title=''ABC''|text='''DEF'''. '''''GHI!'''''}}",
			false,
			"card",
			map[string]string{"title": "<i>ABC</i>", "text": "<b>DEF</b>. <b><i>GHI!</i></b>"},
			false,
		},
	} {
		name, params, err := parseTemplate(test.s, test.rawLinks)
		if isErr := err != nil; isErr != test.isErr {
			t.Errorf("parseTemplate(%q): error=%v, want %v",
				test.s, isErr, test.isErr,
			)
		}
		if name != test.name || !reflect.DeepEqual(params, test.params) {
			t.Errorf("parseTemplate(%q, %v): got %v, %v, want %v, %v",
				test.s, test.rawLinks, name, params, test.name, test.params,
			)
		}
	}
}

func TestParsePage(t *testing.T) {
	for _, test := range []struct {
		s string
		p *page
	}{
		{"", &page{}},
		{"{{Subpage}}", &page{}},
		{"{{card", &page{}},
		{"{{card}}", &page{cards: []Card{{BGColor: otherGray, ID: 1}}}},
		{
			"{{Subpage|page=Cards 1-100}}",
			&page{subpages: []subpage{{page: "Cards 1-100"}}},
		},
		{
			`{{card|title=A|text={{card|title=B|type=Thing}}card|type=Action}}
			{{card|title=C}}`,
			&page{
				cards: []Card{
					{Title: "B", Type: "Thing", BGColor: thingBlue, ID: 1},
					{Title: "C", BGColor: otherGray, ID: 2},
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
			&page{
				cards: []Card{
					{Title: "A", Type: "Action", BGColor: actionRed, ID: 1},
					{Title: "C", Type: "Letter", BGColor: otherGray, ID: 2},
					{Title: "E", BGColor: otherGray, ID: 3},
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
			&page{
				subpages: []subpage{
					{page: "Cards 1-100"},
					{page: "Cards 101-200"},
				},
				cards: []Card{
					{Title: "A", Type: "Action", BGColor: "900", ID: 1},
					{Title: "B", Type: "Thing", BGColor: "090", ID: 2},
				},
			},
		},
	} {
		b := []byte(test.s)
		p := parsePage(b)
		if !reflect.DeepEqual(p, test.p) {
			t.Errorf("parsePage(%v): got %v; want %v",
				b, p, test.p,
			)
		}
	}
}

func TestParseLinkText(t *testing.T) {
	for _, tt := range []struct {
		s, want string
	}{
		{"[[abc]]", "abc"},
		{"[[abc def]]", "abc def"},
		{"[[abc|def]]", "def"},
		{"[[User:abc|abc]]", "abc"},
		{"[[User talk:abc|talk]]", "talk"},
		{"[[file: abc.jpg]]", "abc.jpg"},
		{"[[file: abc.jpg|center|frameless]]", "abc.jpg"},
		{"[[file:abc.exe|harmless]]", ""},
	} {
		if got := parseLinkText(tt.s); got != tt.want {
			t.Errorf("parseLinkText(%v): got %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestReplacePair(t *testing.T) {
	for _, tt := range []struct {
		s, old, new1, new2, want string
	}{
		{"", "a", "b", "c", ""},
		{"onomatopoeia", "u", "a", "o", "onomatopoeia"},
		{"onomatopoeia", "o", "a", "u", "anumatapueia"},
		{"''abc", "''", "<i>", "</i>", "<i>abc"},
		{
			"'''Action:''' Draw a card. '''Action:''' Destroy target Thing.",
			"'''", "<b>", "</b>",
			"<b>Action:</b> Draw a card. <b>Action:</b> Destroy target Thing.",
		},
	} {
		if got := replacePair(tt.s, tt.old, tt.new1, tt.new2); got != tt.want {
			t.Errorf("replacePair(%v, %v, %v, %v): got %v, want %v",
				tt.s, tt.old, tt.new1, tt.new2, got, tt.want,
			)
		}
	}
}

func TestRemoveComments(t *testing.T) {
	for _, tt := range []struct {
		s, want string
	}{
		{"", ""},
		{"a", "a"},
		{"<!--", "<!--"},
		{"-->", "-->"},
		{"<!---->", ""},
		{"<!--comment", "<!--comment"},
		{"<!--comment-->", ""},
		{"abc<!--comment-->", "abc"},
		{"abc <!--comment--> def", "abc  def"},
		{"abc\n<!--comment-->", "abc\n"},
		{"abc\n<!--comment-->\n", "abc\n"},
		{"abc\n <!--comment--> ", "abc\n  "},
		{"abc \n   <!--comment-->   \n  def", "abc \n  def"},
	} {
		if got := removeComments(tt.s); got != tt.want {
			t.Errorf("removeComments(%q): got %q, want %q", tt.s, got, tt.want)
		}
	}
}
