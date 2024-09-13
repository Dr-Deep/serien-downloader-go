package dl

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

const (
	parallelConnections = 4
)

var (
	ErrUnabelToGetFile = errors.New("Unable to fetch file")
)

type DownloadManager struct {
	sync.WaitGroup
}

func NewDlMgr() *DownloadManager {
	return &DownloadManager{sync.WaitGroup{}}
}

func (dlmgr *DownloadManager) Download(url string, filePath string) error {

	// Get file size from the server
	resp, err := http.Head(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrUnabelToGetFile
	}

	var (
		fileSize = resp.ContentLength
		partSize = fileSize / parallelConnections
	)

	// Output-File
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Download
	for i := 0; i < parallelConnections; i++ {
		dlmgr.Add(1)
		offset := int64(i) * partSize

		// For the last part, make sure to download any remaining bytes
		if i == parallelConnections-1 {
			partSize = fileSize - offset
		}

		go dlmgr.downloadPart(url, i+1, file, offset, partSize)
	}

	// Wait for downloads
	dlmgr.Wait()

	return nil
}

func (dlmgr *DownloadManager) downloadPart(url string, partNum int, file *os.File, offset int64, chunkSize int64) {
	defer dlmgr.Done()

	// Create an HTTP request with a "Range" header to download a specific part
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set(
		"Range",
		fmt.Sprintf("bytes=%d-%d", offset, offset+chunkSize-1),
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Write the part to the file at the correct offset
	file.Seek(offset, 0)
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		panic(err)
	}
}
