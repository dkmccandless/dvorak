package dvorak

import (
	"fmt"
)

// Subpage is a subpage of a Dvorak deck.
// If a deck has subpages, parse their cards in order first,
// followed by the cards on the main page.
type Subpage struct {
	// http://dvorakgame.co.uk/index.php/Template:Subpage

	// Page is the subpage's URL relative to the main page.
	Page string
}

// PopulateSubpage returns a Subpage populated with params.
func PopulateSubpage(params map[string]string) (Subpage, error) {
	var sp Subpage
	sp.Page = params["page"]
	if sp.Page == "" {
		return Subpage{}, fmt.Errorf("empty page value")
	}
	return sp, nil
}
