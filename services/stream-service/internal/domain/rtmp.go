package domain

import (
	"time"

	"github.com/google/uuid"
)

type RTMPPublishRequest struct {
	StreamKey string
	AppName   string
	FlashVer  string
	IP        string
}

type RTMPUnpublishRequest struct {
	StreamKey string
	Reason    string
}

type RTMPStreamInfo struct {
	StreamID    uuid.UUID
	StreamKey   string
	CreatorID   uuid.UUID
	Connected   bool
	ConnectedAt time.Time
	BytesIn     int64
	BytesOut    int64
	Bitrate     int
	AudioCodec  string
	VideoCodec  string
	Width       int
	Height      int
	FrameRate   float64
}
