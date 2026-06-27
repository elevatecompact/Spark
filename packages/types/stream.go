package types

type StreamStatus string

const (
	StreamOffline StreamStatus = "offline"
	StreamLive    StreamStatus = "live"
	StreamEnded   StreamStatus = "ended"
	StreamError   StreamStatus = "error"
)

type StreamQuality string

const (
	QualitySD  StreamQuality = "480p"
	QualityHD  StreamQuality = "720p"
	QualityFHD StreamQuality = "1080p"
	QualityQHD StreamQuality = "1440p"
	QualityUHD StreamQuality = "2160p"
)

type StreamType string

const (
	StreamVertical   StreamType = "vertical"
	StreamHorizontal StreamType = "horizontal"
)
