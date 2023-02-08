package dvorak

import "testing"

func TestEscapeInvalidTags(t *testing.T) {
	for _, test := range []struct{ s, want string }{
		{"abc", "abc"},
		{"<b>", "<b>"},
		{"<b>abc</b>", "<b>abc</b>"},
		{"<br>", "<br>"},
		{"<br/>", "<br/>"},
		{"<opponent>", "&lt;opponent&gt;"},
		{"<card title>", "&lt;card title&gt;"},
		{"<card  title   >", "&lt;card  title   &gt;"},
		{
			"Destroy target Thing.<br>Draw a card.",
			"Destroy target Thing.<br>Draw a card.",
		},
		{
			"<font color=FFD700>Golden Text</font>",
			"<font color=FFD700>Golden Text</font>",
		},
		{
			"Replace <metal> with the type of metal",
			"Replace &lt;metal&gt; with the type of metal",
		},
	} {
		if got := escapeNonTags(test.s); got != test.want {
			t.Errorf("escapeInvalidTags(%v): got %v, want %v",
				test.s, got, test.want,
			)
		}
	}
}
