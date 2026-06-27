package domain

import "github.com/google/uuid"

type WebRTCOffer struct {
	SDP  string `json:"sdp"`
	Type string `json:"type"`
}

type WebRTCAnswer struct {
	SDP  string `json:"sdp"`
	Type string `json:"type"`
}

type ICECandidate struct {
	Candidate     string `json:"candidate"`
	SDPMid        string `json:"sdpMid"`
	SDPMLineIndex int    `json:"sdpMLineIndex"`
}

type ViewerJoinRequest struct {
	StreamID uuid.UUID `json:"stream_id"`
	ViewerID uuid.UUID `json:"viewer_id"`
}

type ViewerLeaveRequest struct {
	StreamID   uuid.UUID `json:"stream_id"`
	ViewerID   uuid.UUID `json:"viewer_id"`
	PeerConnID string    `json:"peer_connection_id"`
}

type WebRTCStreamInfo struct {
	StreamID    uuid.UUID `json:"stream_id"`
	ViewerCount int       `json:"viewer_count"`
	IsLive      bool      `json:"is_live"`
	Bitrate     int       `json:"bitrate"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	FrameRate   float64   `json:"frame_rate"`
}
