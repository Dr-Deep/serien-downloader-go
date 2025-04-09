package imdb

import sdl "serien-downloader/internal"

type IMDB_Site struct{}

// return sto search result link
func (site IMDB_Site) Search(query string) ([]sdl.Element, error) {
	return nil, nil
}

/*
GET /suggestion/x/mr%20robot.json HTTP/2
Host: v3.sg.media-imdb.com
Host: https://v3.sg.media-imdb.com/suggestion/x/mr%20robot.json

Accept: application/json
*/

func (site IMDB_Site) Get(url string) ([]sdl.Element, error) {
	return nil, nil
}
