package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/stream-service/internal/config"
	"github.com/elevatecompact/spark/services/stream-service/internal/domain"
)

type HLSService struct {
	cfg       *config.Config
	streamSvc *StreamService
	mu        sync.RWMutex
	segments  map[string]map[int]bool
}

func NewHLSService(cfg *config.Config, streamSvc *StreamService) *HLSService {
	return &HLSService{
		cfg:       cfg,
		streamSvc: streamSvc,
		segments:  make(map[string]map[int]bool),
	}
}

func (s *HLSService) GenerateMasterPlaylist(ctx context.Context, streamID uuid.UUID) (string, error) {
	stream, err := s.streamSvc.GetStream(ctx, streamID)
	if err != nil {
		return "", err
	}

	qualities := stream.AvailableQualities
	if len(qualities) == 0 {
		cfgQualities := s.cfg.HLS.Qualities
		for _, q := range cfgQualities {
			qualities = append(qualities, q.Name)
		}
	}

	if len(qualities) == 0 {
		qualities = []string{"source"}
	}

	var sb strings.Builder
	sb.WriteString("#EXTM3U\n")
	sb.WriteString("#EXT-X-VERSION:6\n")
	sb.WriteString(fmt.Sprintf("# Created by Spark Stream Service - Stream: %s\n", streamID.String()))

	for _, qual := range qualities {
		var qCfg *config.QualityConfig
		for _, q := range s.cfg.HLS.Qualities {
			if q.Name == qual {
				qCfg = &q
				break
			}
		}

		if qCfg != nil {
			sb.WriteString(fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%dx%d,CODECS=\"avc1.640028,mp4a.40.2\"\n",
				qCfg.Bitrate*1000, qCfg.Width, qCfg.Height))
		} else {
			sb.WriteString(fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d\n", 8000000))
		}
		sb.WriteString(fmt.Sprintf("%s/index.m3u8\n", qual))
	}

	return sb.String(), nil
}

func (s *HLSService) GetManifest(ctx context.Context, streamID uuid.UUID, quality string) (string, error) {
	segmentDuration := s.cfg.HLS.SegmentDuration
	playlistLength := s.cfg.HLS.PlaylistLength

	if quality == "" {
		quality = "source"
	}

	key := s.segmentKey(streamID, quality)

	s.mu.RLock()
	segs, exists := s.segments[key]
	s.mu.RUnlock()

	var latestSegment int
	if exists {
		for segNum := range segs {
			if segNum > latestSegment {
				latestSegment = segNum
			}
		}
	}

	if latestSegment == 0 {
		return "", domain.ErrPlaylistNotFound
	}

	startSegment := latestSegment - playlistLength + 1
	if startSegment < 1 {
		startSegment = 1
	}

	var sb strings.Builder
	sb.WriteString("#EXTM3U\n")
	sb.WriteString("#EXT-X-VERSION:6\n")
	sb.WriteString(fmt.Sprintf("#EXT-X-TARGETDURATION:%d\n", segmentDuration))
	sb.WriteString("#EXT-X-MEDIA-SEQUENCE:" + strconv.Itoa(startSegment) + "\n")
	sb.WriteString("#EXT-X-DISCONTINUITY-SEQUENCE:0\n")

	for segNum := startSegment; segNum <= latestSegment; segNum++ {
		if _, ok := segs[segNum]; ok {
			duration := s.getSegmentDuration(segNum, segmentDuration)
			sb.WriteString(fmt.Sprintf("#EXTINF:%.3f,\n", duration))
			sb.WriteString(fmt.Sprintf("segment-%d.ts\n", segNum))
		}
	}

	sb.WriteString("#EXT-X-ENDLIST\n")

	return sb.String(), nil
}

func (s *HLSService) GetSegment(ctx context.Context, streamID uuid.UUID, quality string, segmentNumber int) ([]byte, error) {
	if quality == "" {
		quality = "source"
	}

	segmentPath := filepath.Join(s.cfg.HLS.OutputDir, streamID.String(), quality, fmt.Sprintf("segment-%d.ts", segmentNumber))

	data, err := os.ReadFile(segmentPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, domain.ErrSegmentNotFound
		}
		return nil, fmt.Errorf("read segment file: %w", err)
	}

	return data, nil
}

func (s *HLSService) IsSegmentReady(ctx context.Context, streamID uuid.UUID, quality string, segmentNumber int) bool {
	key := s.segmentKey(streamID, quality)

	s.mu.RLock()
	segs, ok := s.segments[key]
	s.mu.RUnlock()

	if ok && segs[segmentNumber] {
		return true
	}

	segmentPath := filepath.Join(s.cfg.HLS.OutputDir, streamID.String(), quality, fmt.Sprintf("segment-%d.ts", segmentNumber))
	if _, err := os.Stat(segmentPath); err == nil {
		s.mu.Lock()
		if s.segments[key] == nil {
			s.segments[key] = make(map[int]bool)
		}
		s.segments[key][segmentNumber] = true
		s.mu.Unlock()
		return true
	}

	return false
}

func (s *HLSService) GetAvailableQualities(ctx context.Context, streamID uuid.UUID) ([]string, error) {
	stream, err := s.streamSvc.GetStream(ctx, streamID)
	if err != nil {
		return nil, err
	}

	if len(stream.AvailableQualities) > 0 {
		return stream.AvailableQualities, nil
	}

	qualities := make([]string, len(s.cfg.HLS.Qualities))
	for i, q := range s.cfg.HLS.Qualities {
		qualities[i] = q.Name
	}
	return qualities, nil
}

func (s *HLSService) ReportSegment(ctx context.Context, streamID uuid.UUID, quality string, segmentNumber int) {
	key := s.segmentKey(streamID, quality)

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.segments[key] == nil {
		s.segments[key] = make(map[int]bool)
	}
	s.segments[key][segmentNumber] = true
}

func (s *HLSService) CleanupStream(ctx context.Context, streamID uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for key := range s.segments {
		if strings.HasPrefix(key, streamID.String()+"-") {
			delete(s.segments, key)
		}
	}

	outputDir := filepath.Join(s.cfg.HLS.OutputDir, streamID.String())
	if err := os.RemoveAll(outputDir); err != nil {
		log.Warn().Err(err).Str("stream_id", streamID.String()).Msg("Failed to clean up HLS output")
	}
}

func (s *HLSService) segmentKey(streamID uuid.UUID, quality string) string {
	return streamID.String() + "-" + quality
}

func (s *HLSService) getSegmentDuration(segmentNumber int, defaultDuration int) float64 {
	return float64(defaultDuration)
}
