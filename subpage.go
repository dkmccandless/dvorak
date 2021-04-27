package dvorak

import (
	"fmt"
)

// subpage is a subpage of a Dvorak deck.
// If a deck has subpages, parse their cards in order first,
// followed by the cards on the main page.
type subpage struct {
	// http://dvorakgame.co.uk/index.php/Template:Subpage

	// page is the subpage's URL relative to the main page.
	page string
}

// populateSubpage returns a Subpage populated with params.
func populateSubpage(params map[string]string) (subpage, error) {
	var sp subpage
	sp.page = params["page"]
	if sp.page == "" {
		return subpage{}, fmt.Errorf("empty page value")
	}
	return sp, nil
}
