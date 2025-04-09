package bypass

type Hoster int

const (
	// Hoster
	_ Hoster = iota
	StreamTape
	Voe
	DoodStream
)

type Bypasser func(url string) (videoTitle string, videoURL string, _ error)
