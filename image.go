package dvorak

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// imageInfo is the relevant part of the MediaWiki API's imageinfo query result.
type imageInfo struct {
	Query struct {
		Pages map[string]struct {
			// Title is the filename in namespace normalized form
			Title     string
			ImageInfo []struct {
				URL string
			}
		}
	}
}

// ImageURLs queries the Dvorak wiki API and returns a map of image filenames
// to their URLs.
//
// Filenames are in MediaWiki normalized form, with the first character
// capitalized and spaces instead of underscores, but without the "File:"
// namespace prefix.
func ImageURLs(cards []Card) (map[string]string, error) {
	// maxTitles is the number of titles MediaWiki allows in a query.
	const maxTitles = 50

	var images []string
	for _, c := range cards {
		if c.Image != "" {
			images = append(images, "File:"+c.Image)
		}
	}
	if len(images) == 0 {
		return nil, nil
	}

	m := make(map[string]string)
	for len(images) > 0 {
		n := len(images)
		if n > maxTitles {
			n = maxTitles
		}
		urls, err := queryImages(images[:n])
		if err != nil {
			return nil, err
		}
		for name, url := range urls {
			// The wiki's TLS certificate is for dvorakgame.co.uk
			m[name] = strings.Replace(url, "www.", "", 1)
		}
		images = images[n:]
	}
	return m, nil
}

// queryImages returns a map of normalized image filenames to their URLs.
func queryImages(images []string) (map[string]string, error) {
	// query is the query URL that the file names will be appended to.
	// MediaWiki etiquette prefers batching files in a single query if possible.
	// https://www.mediawiki.org/w/api.php?action=help&modules=query%2Bimageinfo
	const query = "https://dvorakgame.co.uk/api.php?action=query&prop=imageinfo&iiprop=url&iilimit=1&format=json&titles="

	titles := strings.ReplaceAll(strings.Join(images, "|"), " ", "_")
	resp, err := http.Get(query + titles)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var ii imageInfo
	json.Unmarshal(b, &ii)

	m := make(map[string]string)
	for _, p := range ii.Query.Pages {
		for _, s := range p.ImageInfo {
			title := strings.TrimPrefix(p.Title, "File:")
			m[title] = s.URL
		}
	}
	return m, nil
}
