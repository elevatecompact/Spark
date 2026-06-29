package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pion/interceptor"
	pion "github.com/pion/webrtc/v3"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/stream-service/internal/config"
	"github.com/elevatecompact/spark/services/stream-service/internal/domain"
)

type PeerConnectionEntry struct {
	ID           string
	PeerConn     *pion.PeerConnection
	ViewerID     uuid.UUID
	StreamID     uuid.UUID
	ConnectedAt  time.Time
	Quality      string
	AudioLevel   float64
	VideoEnabled bool
	AudioEnabled bool
	mu           sync.Mutex
}

type WebRTCService struct {
	cfg          *config.Config
	streamSvc    *StreamService
	mu           sync.RWMutex
	streamPeers  map[uuid.UUID]map[string]*PeerConnectionEntry
	settings     pion.SettingEngine
	interceptor  *interceptor.Interceptor
}

func NewWebRTCService(cfg *config.Config, streamSvc *StreamService) *WebRTCService {
	s := &WebRTCService{
		cfg:         cfg,
		streamSvc:   streamSvc,
		streamPeers: make(map[uuid.UUID]map[string]*PeerConnectionEntry),
	}

	return s
}

func (s *WebRTCService) CreatePeerConnection(ctx context.Context, streamID uuid.UUID, viewerID uuid.UUID) (*PeerConnectionEntry, error) {
	if s.GetViewerCount(ctx, streamID) >= s.cfg.WebRTC.MaxViewers {
		return nil, domain.ErrViewerLimitReached
	}

	stream, err := s.streamSvc.GetStream(ctx, streamID)
	if err != nil {
		return nil, err
	}
	if !stream.IsLive() {
		return nil, domain.ErrStreamNotLive
	}

	m := &pion.MediaEngine{}
	if err := m.RegisterDefaultCodecs(); err != nil {
		return nil, fmt.Errorf("register codecs: %w", err)
	}

	i := &interceptor.Registry{}
	if err := pion.RegisterDefaultInterceptors(m, i); err != nil {
		return nil, fmt.Errorf("register interceptors: %w", err)
	}

	se := pion.SettingEngine{}

	iceServers := make([]pion.ICEServer, 0, len(s.cfg.WebRTC.ICEServers))
	for _, is := range s.cfg.WebRTC.ICEServers {
		iceServers = append(iceServers, pion.ICEServer{
			URLs:           is.URLs,
			Username:       is.Username,
			Credential:     is.Credential,
		})
	}

	if len(iceServers) == 0 {
		for _, stun := range s.cfg.WebRTC.STUNServers {
			iceServers = append(iceServers, pion.ICEServer{
				URLs: []string{stun},
			})
		}
	}

	webrtcCfg := pion.Configuration{
		ICEServers: iceServers,
	}

	api := pion.NewAPI(pion.WithMediaEngine(m), pion.WithInterceptorRegistry(i), pion.WithSettingEngine(se))
	pc, err := api.NewPeerConnection(webrtcCfg)
	if err != nil {
		return nil, fmt.Errorf("create peer connection: %w", err)
	}

	entry := &PeerConnectionEntry{
		ID:           uuid.New().String(),
		PeerConn:     pc,
		ViewerID:     viewerID,
		StreamID:     streamID,
		ConnectedAt:  time.Now(),
		VideoEnabled: true,
		AudioEnabled: true,
		Quality:      "source",
	}

	pc.OnConnectionStateChange(func(state pion.PeerConnectionState) {
		log.Debug().
			Str("stream_id", streamID.String()).
			Str("viewer_id", viewerID.String()).
			Str("state", state.String()).
			Msg("WebRTC connection state changed")

		switch state {
		case pion.PeerConnectionStateFailed:
			s.RemoveViewer(ctx, streamID, entry.ID)
		case pion.PeerConnectionStateDisconnected:
			s.RemoveViewer(ctx, streamID, entry.ID)
		case pion.PeerConnectionStateClosed:
			s.RemoveViewer(ctx, streamID, entry.ID)
		}
	})

	s.addPeer(streamID, entry)

	return entry, nil
}

func (s *WebRTCService) HandleOffer(ctx context.Context, streamID uuid.UUID, offer domain.WebRTCOffer, viewerID uuid.UUID) (*domain.WebRTCAnswer, *PeerConnectionEntry, error) {
	entry, err := s.CreatePeerConnection(ctx, streamID, viewerID)
	if err != nil {
		return nil, nil, err
	}

	if err := entry.PeerConn.SetRemoteDescription(pion.SessionDescription{
		Type: pion.SDPTypeOffer,
		SDP:  offer.SDP,
	}); err != nil {
		s.RemoveViewer(ctx, streamID, entry.ID)
		return nil, nil, fmt.Errorf("set remote description: %w", err)
	}

	answer, err := entry.PeerConn.CreateAnswer(nil)
	if err != nil {
		s.RemoveViewer(ctx, streamID, entry.ID)
		return nil, nil, fmt.Errorf("create answer: %w", err)
	}

	gatherComplete := pion.GatheringCompletePromise(entry.PeerConn)
	if err := entry.PeerConn.SetLocalDescription(answer); err != nil {
		s.RemoveViewer(ctx, streamID, entry.ID)
		return nil, nil, fmt.Errorf("set local description: %w", err)
	}

	<-gatherComplete

	answer = *entry.PeerConn.LocalDescription()

	return &domain.WebRTCAnswer{
		SDP:  answer.SDP,
		Type: string(answer.Type),
	}, entry, nil
}

func (s *WebRTCService) HandleAnswer(ctx context.Context, streamID uuid.UUID, pcID string, answer domain.WebRTCAnswer) error {
	entry, err := s.findPeer(streamID, pcID)
	if err != nil {
		return err
	}

	return entry.PeerConn.SetRemoteDescription(pion.SessionDescription{
		Type: pion.SDPTypeAnswer,
		SDP:  answer.SDP,
	})
}

func (s *WebRTCService) HandleICECandidate(ctx context.Context, streamID uuid.UUID, pcID string, candidate domain.ICECandidate) error {
	entry, err := s.findPeer(streamID, pcID)
	if err != nil {
		return err
	}

	sdpMLineIndex := uint16(candidate.SDPMLineIndex)
	return entry.PeerConn.AddICECandidate(pion.ICECandidateInit{
		Candidate:     candidate.Candidate,
		SDPMid:        &candidate.SDPMid,
		SDPMLineIndex: &sdpMLineIndex,
	})
}

func (s *WebRTCService) AddViewer(ctx context.Context, streamID uuid.UUID, pc *pion.PeerConnection) error {
	return nil
}

func (s *WebRTCService) RemoveViewer(ctx context.Context, streamID uuid.UUID, pcID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	peers, ok := s.streamPeers[streamID]
	if !ok {
		return domain.ErrPeerConnectionNotFound
	}

	entry, ok := peers[pcID]
	if !ok {
		return domain.ErrPeerConnectionNotFound
	}

	if entry.PeerConn != nil {
		if err := entry.PeerConn.Close(); err != nil {
			log.Warn().Err(err).Str("pc_id", pcID).Msg("Error closing peer connection")
		}
	}

	delete(peers, pcID)
	if len(peers) == 0 {
		delete(s.streamPeers, streamID)
	}

	log.Debug().
		Str("stream_id", streamID.String()).
		Str("viewer_id", entry.ViewerID.String()).
		Int("remaining", len(peers)).
		Msg("Viewer removed from WebRTC stream")

	return nil
}

func (s *WebRTCService) GetViewerCount(ctx context.Context, streamID uuid.UUID) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	peers, ok := s.streamPeers[streamID]
	if !ok {
		return 0
	}
	return len(peers)
}

func (s *WebRTCService) GetViewerIDs(ctx context.Context, streamID uuid.UUID) []uuid.UUID {
	s.mu.RLock()
	defer s.mu.RUnlock()

	peers, ok := s.streamPeers[streamID]
	if !ok {
		return nil
	}

	ids := make([]uuid.UUID, 0, len(peers))
	for _, entry := range peers {
		ids = append(ids, entry.ViewerID)
	}
	return ids
}

func (s *WebRTCService) RemoveAllStreamViewers(ctx context.Context, streamID uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	peers, ok := s.streamPeers[streamID]
	if !ok {
		return
	}

	for id, entry := range peers {
		if entry.PeerConn != nil {
			_ = entry.PeerConn.Close()
		}
		delete(peers, id)
	}
	delete(s.streamPeers, streamID)
}

func (s *WebRTCService) addPeer(streamID uuid.UUID, entry *PeerConnectionEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.streamPeers[streamID]; !ok {
		s.streamPeers[streamID] = make(map[string]*PeerConnectionEntry)
	}
	s.streamPeers[streamID][entry.ID] = entry
}

func (s *WebRTCService) findPeer(streamID uuid.UUID, pcID string) (*PeerConnectionEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	peers, ok := s.streamPeers[streamID]
	if !ok {
		return nil, domain.ErrPeerConnectionNotFound
	}

	entry, ok := peers[pcID]
	if !ok {
		return nil, domain.ErrPeerConnectionNotFound
	}
	return entry, nil
}
