package sdl

import (
	"path"
	"serien-downloader/internal/bypass"
	"serien-downloader/internal/dl"
	"strings"
)

type SerienDownloader struct {
	BypassModules []bypass.Bypasser
	SiteModules   []Site
	Dlmgr         dl.DownloadManager
}

func (sdl SerienDownloader) Search(query string) ([]Element, error) {
	var (
		elems           []Element
		lastCriticalErr error
	)

	for _, site := range sdl.SiteModules {
		r, err := site.Search(query)
		if len(r) == 0 && err != nil {
			lastCriticalErr = err
		}

		elems = append(elems, r...)
	}

	return elems, lastCriticalErr
}

// Get(url string) ([]sdl.Element, error)
func (sdl SerienDownloader) Get(url string) ([]Element, error) {
	var (
		elems           []Element
		lastCriticalErr error
	)

	for _, site := range sdl.SiteModules {
		e, err := site.Get(url)
		if len(e) == 0 && err != nil {
			lastCriticalErr = err
		}

		elems = append(elems, e...)
	}

	return elems, lastCriticalErr
}

func (sdl SerienDownloader) Download(e *Element) error {

	var (
		filename = strings.ReplaceAll(path.Clean(e.Title), " ", ".") + ".mp4"
	)

	// time

	sdl.Dlmgr.Download(filename, e.URLS[0])

	return nil

	// progress?

	// filepath:
	// './{title}.mp4'; wenn '.mp4' nicht schon da und wirklich mp4 file
	// invalide zeichen entfernen (utf8 check, filesys check)

	// mit ffmpeg recoden?
	// go-routine

	/*
	   go dlmgr.Download() -> ffmpeg
	   dlmgr managed parallele downloads
	   dlmgr.Wait()

	*/
}

func (sdl SerienDownloader) Bypass(hoster bypass.Hoster, url string) (Element, error) {
	switch hoster {
	case bypass.StreamTape:
		title, link, err := bypass.GetStreamtapeVideo(url)
		if err != nil {
			return Element{}, err
		}

		return Element{
			Title:       title,
			Description: "",
			URLS:        []string{link},
		}, nil

	case bypass.Voe:
		title, link, err := bypass.GetVoeVideo(url)
		if err != nil {
			return Element{}, err
		}

		return Element{
			Title:       title,
			Description: "",
			URLS:        []string{link},
		}, nil

	case bypass.DoodStream:
		//

	}

	return Element{}, ErrBypasserNotFound
}
