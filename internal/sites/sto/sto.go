package sto

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	sdl "serien-downloader/internal"
)

// type STO_Site sdl.Site {}
//var STO_Site sdl.Site = struct{}{}

const (
	BaseURI = "https://s.to"
)

var (
	// Regexps
	searchRegexp  = regexp.MustCompile("/serie/stream/[a-zA-Z,-]+$")
	seasonRegexp  = regexp.MustCompile("/serie/stream/[a-zA-Z,-]+/staffel-\\d+$")
	episodeRegexp = regexp.MustCompile("/serie/stream/[a-zA-Z,-]+/staffel-\\d/episode-\\d+$")

	// Errors
	ErrNotFound    error = errors.New("not found")
	ErrNoResults   error = errors.New("No Search-Results found")
	ErrServerError error = errors.New("Server Error, please try again later")
)

// from  sdl.Site
type STO_Site struct{}

// TODO: header mit richtigen encoder, URL encode
func (site STO_Site) Search(query string) ([]sdl.Element, error) {

	test := url.Values{}
	test.Set("keyword", query)

	//
	resp, err := http.PostForm(
		fmt.Sprintf("%s/ajax/search", BaseURI),
		test,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, ErrServerError
	}

	//
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//
	if len(body) == 0 {
		return nil, ErrNoResults
	}

	type sus struct {
		Title       string
		Description string
		Link        string
	}

	var (
		rawResults    []sus
		searchResults []sdl.Element
	)

	if err := json.Unmarshal(body, &rawResults); err != nil {
		return nil, err
	}

	for _, r := range rawResults {
		if searchRegexp.MatchString(r.Link) {
			searchResults = append(
				searchResults,
				sdl.Element{
					Title:       removeHTMLFragmentsFromString(r.Title),
					Description: removeHTMLFragmentsFromString(r.Description),
					URLS:        []string{BaseURI + r.Link},
				},
			)
		}
	}

	return searchResults, nil

}

func (site STO_Site) Get(url string) ([]sdl.Element, error) {
	switch true {

	case searchRegexp.MatchString(url):
		// Staffeln liste
		elems, err := site.getSeasons(url)
		if err != nil {
			return nil, err
		}

		return elems, nil

	case seasonRegexp.MatchString(url):
		// Episoden Liste
		elems, err := site.getSeason(url)
		if err != nil {
			return nil, err
		}

		return elems, nil

	case episodeRegexp.MatchString(url):
		// bypass && episoden mp4 zurück geben
		return site.getEpisode(url)

	default:
		return nil, ErrNotFound
	}
}
