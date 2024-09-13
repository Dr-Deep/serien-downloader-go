package bypass

import (
	"errors"
	"net/http"
	sdl "serien-downloader/internal"
)

type Hoster int

const (
	// Hoster
	_ Hoster = iota
	StreamTape
	Voe
	DoodStream
)

func Bypass(hoster Hoster, url string) (sdl.Element, error) {
	switch hoster {
	case StreamTape:
		title, link, err := GetStreamTapeVideo(url)
		if err != nil {
			return sdl.Element{}, err
		}

		return sdl.Element{
			Title:       title,
			Description: "",
			URLS:        []string{link},
		}, nil

	case Voe:
		//

	case DoodStream:
		//

	}

	return sdl.Element{}, errors.New("bypasser not found")
}

func BypassRedirect(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return url, err
	}

	if resp.StatusCode == http.StatusFound { // 302:Found
		_url, err := resp.Location()
		if err != nil {
			return url, err
		}

		return _url.String(), nil
	}

	return url, nil
}
