package dl

import (
	"os"
	"os/exec"
	"sync"
)

/*
Options:
   * ytdlp / standalone bzw buildin
   * paralel?
   * concurrent?

   MP4 und HLS download mit "--concurrent-fragments 2"
*/

type DownloadManager struct {
	sync.WaitGroup
}

func NewDlMgr() DownloadManager {
	return DownloadManager{sync.WaitGroup{}}
}

/*
   mit dlmgr.Download() element der dl queue hinzufügen
   mit dlmgr.Wait() auf alle dl's warten und status anzeigen
*/

func (dlmgr *DownloadManager) Download(_filename string, url string) {
	cmd := exec.Command(
		"yt-dlp",
		//"-q",
		"--continue",
		"-f mp4",
		url,
		"-o", _filename,
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func (dlmgr *DownloadManager) Wait() {
	dlmgr.Wait()
}

// "ffmpeg -i input.mp4/stream -c:v libx264 -crf 0 output.mp4"
// ffmpeg -hwaccel vaapi -vaapi_device /dev/dri/renderD128 -i input.mp4
// -c:v h264_vaapi -b:v 4M -c:a copy output.mp4
// ffmpeg -hwaccel vaapi -vaapi_device /dev/dri/renderX -i X.mp4 -:c:v h264_vaapi -crf 0 X.mp4

/*
## downloads mit ffmpeg
ffmpeg -i link title.mp4
format: mp4 hohe qualität kleine größe
-loglevel (error,16)
concurrent downloads?
*/

// info/status von ffmpeg über pipe?
// und dann anzeigen
// return codes als error zeichen und stderr in logger senden

// kann ffmpeg auf HD hochskalieren?
