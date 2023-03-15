package html

import "testing"

func TestRemoveHTMLTags(t *testing.T) {
	for _, test := range []struct{ s, want string }{
		{"", ""},
		{"abc", "abc"},
		{"<b>", "<b>"},
		{"<b>abc</b>", "<b>abc</b>"},
		{"<br>", "<br>"},
		{"<br/>", "<br/>"},
		{"<opponent>", "&lt;opponent&gt;"},
		{"<card title>", "&lt;card title&gt;"},
		{"<card  title   >", "&lt;card  title   &gt;"},
		{"<target opponent>", "&lt;target opponent&gt;"},
		{
			"Destroy target Thing.<br>Draw a card.",
			"Destroy target Thing.<br>Draw a card.",
		},
		{
			"<font color=FFD700>Golden Text</font>",
			"<font color=\"FFD700\">Golden Text</font>",
		},
		{
			"<font color='FFD700'>Golden Text</font>",
			"<font color=\"FFD700\">Golden Text</font>",
		},
		{
			"<font color=\"FFD700\">Golden Text</font>",
			"<font color=\"FFD700\">Golden Text</font>",
		},
		{
			"Replace <metal> with the type of metal",
			"Replace &lt;metal&gt; with the type of metal",
		},
	} {
		if got := RemoveHTMLTags(test.s); got != test.want {
			t.Errorf("RemoveHTMLTags(%v): got %v, want %v",
				test.s, got, test.want,
			)
		}
	}
}

func TestRemoveHTMLComments(t *testing.T) {
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
		if got := removeHTMLComments(tt.s); got != tt.want {
			t.Errorf("removeHTMLComments(%q): got %q, want %q", tt.s, got, tt.want)
		}
	}
}
