package dvorak

import (
	"reflect"
	"testing"
)

func TestPopulateSubpage(t *testing.T) {
	for _, test := range []struct {
		params map[string]string
		sp     subpage
	}{
		{map[string]string{}, subpage{}},
		{map[string]string{"page": ""}, subpage{}},
		{
			map[string]string{"page": "Cards 1-100"},
			subpage{page: "Cards 1-100"},
		},
		{
			map[string]string{"": "", "hide": "true", "page": "Cards 1-100"},
			subpage{page: "Cards 1-100"},
		},
	} {
		sp, err := populateSubpage(test.params)
		if !reflect.DeepEqual(sp, test.sp) || (err != nil) != (sp == subpage{}) {
			t.Errorf("populateSubpage(%v): got %v, %v; want %v",
				test.params, sp, err, test.sp,
			)
		}
	}
}
