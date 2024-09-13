package sdl

import "serien-downloader/internal/dl"

type SerienDownloader struct {
	SiteModules []Site
	Dlmgr       *dl.DownloadManager
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
	// progress?

	// filepath:
	// './{title}.mp4'; wenn '.mp4' nicht schon da und wirklich mp4 file
	// invalide zeichen entfernen (utf8 check, filesys check)
	return nil
}
