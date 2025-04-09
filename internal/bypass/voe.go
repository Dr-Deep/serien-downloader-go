package bypass

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

func GetVoeVideo(url string) (vidTitle, vidURL string, _ error) {
	var (
		c = colly.NewCollector(
			colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64; rv:129.0) Gecko/20100101 Firefox/129.0"),
		)
		lastCriticalErr error

		// Regexp's
		redirExp = regexp.MustCompile("window.location.href = '[^']+'")
	)

	/*
	 * Bypass Javascript Redirects
	 */
	c.OnHTML("html", func(h *colly.HTMLElement) {
		if redirURL := redirExp.FindString(h.Text); redirURL != "" {
			redirURL = strings.ReplaceAll(redirURL, "window.location.href = ", "")
			redirURL = strings.ReplaceAll(redirURL, "'", "")

			// reload!
			if err := c.Visit(redirURL); err != nil {
				lastCriticalErr = err
			}
			return
		}
	})

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

		exp := regexp.MustCompile("var sources = {[^}]+}")
		sourcesElement := exp.FindString(h.Text)
		if sourcesElement == "" {
			return
		}

		sourcesElement = strings.ReplaceAll(sourcesElement, "var sources = ", "")
		sourcesElement = strings.ReplaceAll(sourcesElement, "'", "\"")
		i := strings.LastIndex(sourcesElement, ",")
		sourcesElement = sourcesElement[0:i] + sourcesElement[i+1:]

		if videoURL, err := findVideoURLinSource(sourcesElement); err != nil {
			lastCriticalErr = err
		} else {
			vidURL = videoURL
		}
	})

	if err := c.Visit(url); err != nil {
		return "", "", err
	}

	if lastCriticalErr != nil {
		return "", "", lastCriticalErr
	}

	return vidTitle, vidURL, lastCriticalErr
}

func findVideoURLinSource(sus2 string) (string, error) {
	var v map[string]any

	if err := json.Unmarshal([]byte(sus2), &v); err != nil {
		return "", err
	}

	var (
		_, ok1 = v["mp4"]
		_, ok2 = v["hls"]
		coded  string
	)

	switch {
	case ok1:
		coded = v["mp4"].(string)

	case ok2:
		coded = v["hls"].(string)

	default:
		return "", errors.New("couldnt find video elements")

	}

	// decode base64 link
	l, err := base64.StdEncoding.DecodeString(coded)
	if err != nil {
		return "", err
	}

	return string(l), nil
}
