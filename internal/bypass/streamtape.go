package bypass

import (
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

func GetStreamtapeVideo(_url string) (vidTitle, vidURL string, _ error) {
	var (
		url   = strings.Replace(_url, "/e/", "/v/", 0) // make link downloadable
		base  string
		token string
		c     = colly.NewCollector(
			colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64; rv:129.0) Gecko/20100101 Firefox/129.0"),
		)
		lastCriticalErr error
	)

	/*
	 * Find Video Title
	 */
	c.OnHTML("meta", func(h *colly.HTMLElement) {
		if h.Attr("name") == "og:title" {
			vidTitle = h.Attr("content")
		}
	})

	/*
	 * Find Video URL
	 */
	c.OnHTML("script", func(h *colly.HTMLElement) {

		// grab js link expression
		var (
			exp   = regexp.MustCompile(`document.getElementById\('ideoolink'\).innerHTML = (.+);`)
			match string
		)
		if matches := exp.FindStringSubmatch(h.Text); len(matches) >= 2 {
			match = matches[1]
		} else {
			return
		}

		// grab token
		var (
			exp2 = regexp.MustCompile(`token=([^&']+)`)
		)

		token = exp2.FindString(match)
	})

	// grab base URL
	c.OnHTML("#ideoolink", func(h *colly.HTMLElement) {
		if h.Text != "" {
			base = h.Text
		}
	})

	if err := c.Visit(url); err != nil {
		return "", "", err
	}

	if lastCriticalErr != nil {
		return "", "", lastCriticalErr
	}

	// zusammenbauen
	base = strings.Split(base, "token=")[0]
	vidURL = "https:/" + base + token //?  stream=1?

	return vidTitle, vidURL, nil
}

/*
	if resp.StatusCode == http.StatusFound { // 302:Found
		_url, err := resp.Location()
		if err != nil {
			return url, err
		}

		return _url.String(), nil
	}
*/
