package dvorak

import (
	"reflect"
	"testing"
)

func TestPopulateSubpage(t *testing.T) {
	for _, test := range []struct {
		params map[string]string
		sp     Subpage
	}{
		{map[string]string{}, Subpage{}},
		{map[string]string{"page": ""}, Subpage{}},
		{
			map[string]string{"page": "Cards 1-100"},
			Subpage{Page: "Cards 1-100"},
		},
		{
			map[string]string{"": "", "hide": "true", "page": "Cards 1-100"},
			Subpage{Page: "Cards 1-100"},
		},
	} {
		sp, err := PopulateSubpage(test.params)
		if !reflect.DeepEqual(sp, test.sp) || (err != nil) != (sp == Subpage{}) {
			t.Errorf("PopulateSubpage(%v): got %v, %v; want %v",
				test.params, sp, err, test.sp,
			)
		}
	}
}
