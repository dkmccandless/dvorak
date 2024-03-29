package dvorak

import (
	"reflect"
	"testing"

	"kr.dev/diff"
)

func TestParseParameter(t *testing.T) {
	for _, test := range []struct{ s, name, value string }{
		{"", "", ""},
		{"card", "", "card"},
		{"title=", "title", ""},
		{"title=ABC", "title", "ABC"},
		{"title = ABC", "title", "ABC"},
		{" title = ABC ", "title", "ABC"},
		{
			"text=<b>Action:</b>Destroy target Thing.<br>Draw a card.",
			"text", "<b>Action:</b>Destroy target Thing.<br>Draw a card.",
		},
		{
			"title=<font color=FFD700>Golden Title</font>",
			"title", "<font color=FFD700>Golden Title</font>",
		},
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
			map[string]string{"creator": "ABC", "title": "DEF"},
			false,
		},
		{
			"{{card|text=[[File: ABC.jpg]]DEF}}",
			"card",
			map[string]string{"text": "DEF", "image": "ABC.jpg"},
			false,
		},
		{
			"{{card|creator=[[User:ABC|ABC]] ([[User talk:ABC|talk]])|title=DEF}}",
			"card",
			map[string]string{"creator": "ABC", "title": "DEF"},
			false,
		},
		{
			"{{card|creator=[[User:ABC|ABC]] ([[User talk:ABC|talk]]) 21:33, 25 July 2012 (UTC)|title=DEF}}",
			"card",
			map[string]string{"creator": "ABC", "title": "DEF"},
			false,
		},
		{
			"{{card| title = A | type = Action }}",
			"card",
			map[string]string{"title": "A", "type": "Action"},
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
			`
				{{card| title = A | type = Action }}
				card|title=B|type=Thing}}
				{{card|title=C|type=Letter}}
				{{card|title=D
				{{card|title=E}}
			`,
			&page{
				cards: []Card{
					{Title: text("A"), Type: text("Action"), BGColor: actionRed, ID: 1},
					{Title: text("C"), Type: text("Letter"), BGColor: otherGray, ID: 2},
					{Title: text("E"), BGColor: otherGray, ID: 3},
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
					{Title: text("A"), Type: text("Action"), BGColor: "900", ID: 1},
					{Title: text("B"), Type: text("Thing"), BGColor: "090", ID: 2},
				},
			},
		},
	} {
		b := []byte(test.s)
		p := parsePage(b)
		diff.Test(t, t.Errorf, p, test.p)
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
