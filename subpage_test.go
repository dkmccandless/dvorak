package dvorak

import (
	"reflect"
	"testing"
)

func TestParseSubpage(t *testing.T) {
	var sp1 = &Subpage{Page: "Cards 1-100"}

	for _, test := range []struct {
		s  string
		sp *Subpage
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
		sp, err := ParseSubpage(test.s)
		if !reflect.DeepEqual(sp, test.sp) || (err != nil) != (sp == nil) {
			t.Errorf("ParseSubpage(%q): got %v, %v; want %v",
				test.s, sp, err, test.sp,
			)
		}
	}
}
