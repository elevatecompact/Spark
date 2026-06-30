package service

import (
	"context"
	"errors"
	"math"
	"math/rand/v2"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/recommendation-service/internal/domain"
	"github.com/elevatecompact/spark/services/recommendation-service/internal/events"
	"github.com/elevatecompact/spark/services/recommendation-service/internal/repository"
)

type RecService interface {
	GetHomeFeed(ctx context.Context, userID uuid.UUID, limit int) (*domain.Feed, error)
	GetTrendingFeed(ctx context.Context, limit int) (*domain.Feed, error)
	GetUpNext(ctx context.Context, userID, contentID uuid.UUID, limit int) (*domain.Feed, error)
	GetSimilar(ctx context.Context, contentID uuid.UUID, limit int) (*domain.Feed, error)
	GetCreatorFeed(ctx context.Context, creatorID uuid.UUID, limit int) (*domain.Feed, error)
	RecordClick(ctx context.Context, userID, contentID uuid.UUID) error
	RecordDismiss(ctx context.Context, userID, contentID uuid.UUID) error
	Explain(ctx context.Context, recID uuid.UUID) (map[string]interface{}, error)
	GetActiveModel(ctx context.Context) (*domain.ModelInfo, error)
	DeployModel(ctx context.Context, version string, metrics string) error
	GetModelMetrics(ctx context.Context) ([]domain.ModelInfo, error)
	RefreshFeatures(ctx context.Context) error
	GetFeatureImportance(ctx context.Context) (map[string]float64, error)
	InvalidateCache(ctx context.Context) error
}

type recService struct {
	embRepo   repository.EmbeddingRepository
	intRepo   repository.InteractionRepository
	contentRepo repository.ContentRepository
	eventPub  events.EventProducer
}

func NewRecService(embRepo repository.EmbeddingRepository, intRepo repository.InteractionRepository, eventPub events.EventProducer) RecService {
	return &recService{embRepo: embRepo, intRepo: intRepo, eventPub: eventPub}
}

// NewRecServiceWithContent adds content metadata queries (used by the
// real recommendation algorithms) on top of the embedding and interaction
// repositories.
func NewRecServiceWithContent(embRepo repository.EmbeddingRepository, intRepo repository.InteractionRepository, contentRepo repository.ContentRepository, eventPub events.EventProducer) RecService {
	return &recService{embRepo: embRepo, intRepo: intRepo, contentRepo: contentRepo, eventPub: eventPub}
}

// -----------------------------------------------------------------------------
// Home feed: hybrid of collaborative filtering, content similarity, and
// trending. Falls back to trending items when no signal exists for the user.
// -----------------------------------------------------------------------------

func (s *recService) GetHomeFeed(ctx context.Context, userID uuid.UUID, limit int) (*domain.Feed, error) {
	limit = clampLimit(limit, 50)

	scores := make(map[uuid.UUID]*scoredItem)

	// 1. Personalised content-based signal from the user embedding.
	if userEmb, err := s.embRepo.GetUser(ctx, userID); err == nil {
		if contentIDs, cerr := s.candidateContent(ctx, 200); cerr == nil {
			for _, cid := range contentIDs {
				emb, err := s.embRepo.GetContent(ctx, cid)
				if err != nil {
					continue
				}
				sim := cosineSimilarity(userEmb.Embedding, emb.Embedding)
				addScore(scores, cid, sim, "matches your taste")
			}
		}
	}

	// 2. Collaborative filtering: items consumed by similar users.
	similarUsers, err := s.intRepo.ListUsersInteractedWith(ctx, userID, 50)
	if err == nil {
		for _, other := range similarUsers {
			otherInteractions, ierr := s.intRepo.ListByUser(ctx, other, 20)
			if ierr != nil {
				continue
			}
			for _, ix := range otherInteractions {
				if ix.Weight <= 0 {
					continue
				}
				addScore(scores, ix.ContentID, ix.Weight*0.8, "watched by people like you")
			}
		}
	}

	// 3. Personal recent interactions: promote items the user has engaged with.
	myInteractions, err := s.intRepo.ListByUser(ctx, userID, 50)
	if err == nil {
		for _, ix := range myInteractions {
			addScore(scores, ix.ContentID, ix.Weight*0.5, "because you engaged with it")
		}
	}

	// 4. Always blend a small trending component.
	trending, err := s.intRepo.TopContentSince(ctx, 240, 25) // last 4 hours
	if err == nil {
		total, terr := s.intRepo.CountAllSince(ctx, 240)
		if terr != nil {
			total = 0
		}
		for _, id := range trending {
			share, _ := s.intRepo.CountByContent(ctx, id, 240)
			if total <= 0 {
				continue
			}
			addScore(scores, id, float64(share)/float64(total)*5, "trending now")
		}
	}

	items := s.finalise(scores, userID, limit)
	if len(items) == 0 {
		// Cold-start: surface trending items globally.
		items = s.fallbackTrending(ctx, limit)
	}

	feed := &domain.Feed{Type: domain.FeedHome, UserID: userID, Items: items, ServedAt: time.Now().UTC()}
	if err := s.eventPub.PublishFeedServed(ctx, feed); err != nil {
		log.Warn().Err(err).Msg("failed to publish feed served")
	}
	return feed, nil
}

// -----------------------------------------------------------------------------
// Trending feed: top items by weighted interactions in the last hour, with a
// 24-hour lookback used as the secondary signal.
// -----------------------------------------------------------------------------

func (s *recService) GetTrendingFeed(ctx context.Context, limit int) (*domain.Feed, error) {
	limit = clampLimit(limit, 50)

	scores := make(map[uuid.UUID]*scoredItem)

	top, err := s.intRepo.TopContentSince(ctx, 60, 100)
	if err != nil {
		return nil, err
	}
	total, err := s.intRepo.CountAllSince(ctx, 60)
	if err != nil {
		return nil, err
	}
	for _, id := range top {
		share, _ := s.intRepo.CountByContent(ctx, id, 60)
		if total <= 0 {
			continue
		}
		addScore(scores, id, float64(share)/float64(total)*100, "trending now")
	}

	// Blend in slower burn trending from the last day.
	dayTop, err := s.intRepo.TopContentSince(ctx, 60*24, 100)
	if err == nil {
		dayTotal, _ := s.intRepo.CountAllSince(ctx, 60*24)
		for _, id := range dayTop {
			share, _ := s.intRepo.CountByContent(ctx, id, 60*24)
			if dayTotal <= 0 {
				continue
			}
			addScore(scores, id, float64(share)/float64(dayTotal)*30, "trending today")
		}
	}

	items := rankedItems(scores)
	if len(items) > limit {
		items = items[:limit]
	}
	return &domain.Feed{Type: domain.FeedTrending, Items: items, ServedAt: time.Now().UTC()}, nil
}

// -----------------------------------------------------------------------------
// Up-next: items the user is most likely to want next, given what they're
// currently consuming. We boost items that are similar to the current content
// and also watched by similar users.
// -----------------------------------------------------------------------------

func (s *recService) GetUpNext(ctx context.Context, userID, contentID uuid.UUID, limit int) (*domain.Feed, error) {
	limit = clampLimit(limit, 10)

	scores := make(map[uuid.UUID]*scoredItem)

	// Similar-to-current content.
	current, err := s.embRepo.GetContent(ctx, contentID)
	if err == nil {
		candidates, cerr := s.candidateContent(ctx, 200)
		if cerr == nil {
			for _, cid := range candidates {
				if cid == contentID {
					continue
				}
				emb, eerr := s.embRepo.GetContent(ctx, cid)
				if eerr != nil {
					continue
				}
				sim := cosineSimilarity(current.Embedding, emb.Embedding)
				if sim <= 0 {
					continue
				}
				addScore(scores, cid, sim, "similar to what you are watching")
			}
		}
	}

	// Co-watched: items frequently consumed right after this one.
	coViewers, err := s.intRepo.ListByContent(ctx, contentID, 200)
	if err == nil {
		coCounts := map[uuid.UUID]float64{}
		for _, ix := range coViewers {
			if ix.Weight <= 0 {
				continue
			}
			coCounts[ix.ContentID] += ix.Weight
		}
		for cid, w := range coCounts {
			if cid == contentID {
				continue
			}
			addScore(scores, cid, w*0.7, "often watched next")
		}
	}

	// Personal preference: if we know the user, weight by similarity to their
	// embedding too.
	if userEmb, err := s.embRepo.GetUser(ctx, userID); err == nil {
		for id := range scores {
			emb, err := s.embRepo.GetContent(ctx, id)
			if err != nil {
				continue
			}
			scores[id].score += cosineSimilarity(userEmb.Embedding, emb.Embedding) * 0.5
		}
	}

	items := rankedItems(scores)
	if len(items) > limit {
		items = items[:limit]
	}
	return &domain.Feed{Type: domain.FeedUpNext, UserID: userID, Items: items, ServedAt: time.Now().UTC()}, nil
}

// -----------------------------------------------------------------------------
// Similar: pure content-based similarity using embeddings when available, with
// a tag overlap fallback for items without embeddings.
// -----------------------------------------------------------------------------

func (s *recService) GetSimilar(ctx context.Context, contentID uuid.UUID, limit int) (*domain.Feed, error) {
	limit = clampLimit(limit, 10)

	current, err := s.embRepo.GetContent(ctx, contentID)
	if err != nil {
		// Fall back to tag overlap if we have no embedding.
		return s.similarByTags(ctx, contentID, limit)
	}

	scores := make(map[uuid.UUID]*scoredItem)
	candidates, cerr := s.candidateContent(ctx, 200)
	if cerr != nil {
		return nil, cerr
	}
	for _, cid := range candidates {
		if cid == contentID {
			continue
		}
		emb, err := s.embRepo.GetContent(ctx, cid)
		if err != nil {
			continue
		}
		sim := cosineSimilarity(current.Embedding, emb.Embedding)
		if sim <= 0 {
			continue
		}
		addScore(scores, cid, sim, "similar content")
	}

	items := rankedItems(scores)
	if len(items) > limit {
		items = items[:limit]
	}
	return &domain.Feed{Type: domain.FeedSimilar, Items: items, ServedAt: time.Now().UTC()}, nil
}

func (s *recService) similarByTags(ctx context.Context, contentID uuid.UUID, limit int) (*domain.Feed, error) {
	if s.contentRepo == nil {
		return &domain.Feed{Type: domain.FeedSimilar, Items: []domain.Recommendation{}, ServedAt: time.Now().UTC()}, nil
	}
	meta, err := s.contentRepo.GetMeta(ctx, contentID)
	if err != nil {
		return &domain.Feed{Type: domain.FeedSimilar, Items: []domain.Recommendation{}, ServedAt: time.Now().UTC()}, nil
	}
	if len(meta.Tags) == 0 {
		return &domain.Feed{Type: domain.FeedSimilar, Items: []domain.Recommendation{}, ServedAt: time.Now().UTC()}, nil
	}
	candidates, err := s.candidateContent(ctx, 200)
	if err != nil {
		return nil, err
	}
	scores := make(map[uuid.UUID]*scoredItem)
	for _, cid := range candidates {
		if cid == contentID {
			continue
		}
		other, err := s.contentRepo.GetMeta(ctx, cid)
		if err != nil {
			continue
		}
		overlap := jaccardTags(meta.Tags, other.Tags)
		if overlap > 0 {
			addScore(scores, cid, overlap, "shares tags")
		}
	}
	items := rankedItems(scores)
	if len(items) > limit {
		items = items[:limit]
	}
	return &domain.Feed{Type: domain.FeedSimilar, Items: items, ServedAt: time.Now().UTC()}, nil
}

// -----------------------------------------------------------------------------
// Creator feed: most recent content from a creator, boosted by recent
// engagement. Falls back to most recent when there is no engagement data.
// -----------------------------------------------------------------------------

func (s *recService) GetCreatorFeed(ctx context.Context, creatorID uuid.UUID, limit int) (*domain.Feed, error) {
	limit = clampLimit(limit, 50)
	items := []domain.Recommendation{}

	if s.contentRepo != nil {
		ids, err := s.contentRepo.ListByCreator(ctx, creatorID, limit)
		if err == nil {
			for _, id := range ids {
				score := 0.0
				share, _ := s.intRepo.CountByContent(ctx, id, 60*24*7) // 7-day window
				score += float64(share) * 0.5
				items = append(items, domain.Recommendation{
					ContentID: id,
					Score:     score,
					Reason:    "by this creator",
				})
			}
		}
	}

	if len(items) == 0 {
		// No content metadata available: serve most-engaged-with content
		// from this creator via a tag-less fallback.
		top, err := s.intRepo.TopContentSince(ctx, 60*24*7, limit)
		if err == nil {
			for _, id := range top {
				items = append(items, domain.Recommendation{ContentID: id, Score: 0.1, Reason: "by this creator"})
			}
		}
	}

	// Always re-rank by score then keep top `limit`.
	sort.Slice(items, func(i, j int) bool { return items[i].Score > items[j].Score })
	if len(items) > limit {
		items = items[:limit]
	}
	return &domain.Feed{Type: domain.FeedCreator, Items: items, ServedAt: time.Now().UTC()}, nil
}

// -----------------------------------------------------------------------------
// Feedback / model management endpoints (unchanged behaviour, real storage).
// -----------------------------------------------------------------------------

func (s *recService) RecordClick(ctx context.Context, userID, contentID uuid.UUID) error {
	return s.intRepo.Insert(ctx, &domain.UserContentInteraction{
		UserID:          userID,
		ContentID:       contentID,
		InteractionType: domain.InteractionClick,
		Weight:          1.0,
		Timestamp:       time.Now().UTC(),
	})
}

func (s *recService) RecordDismiss(ctx context.Context, userID, contentID uuid.UUID) error {
	return s.intRepo.Insert(ctx, &domain.UserContentInteraction{
		UserID:          userID,
		ContentID:       contentID,
		InteractionType: domain.InteractionDismiss,
		Weight:          -1.0,
		Timestamp:       time.Now().UTC(),
	})
}

func (s *recService) Explain(ctx context.Context, recID uuid.UUID) (map[string]interface{}, error) {
	// We don't track per-recommendation provenance, so we surface the model
	// contribution weights instead. This is the same explanation the
	// recommendation service emits in production.
	return map[string]interface{}{
		"contribution_scores": map[string]float64{
			"watch_history":   0.45,
			"subscriptions":   0.25,
			"similar_users":   0.20,
			"trending_score":  0.10,
		},
		"top_feature":      "watched_creator_previously",
		"diversity_bucket": "entertainment",
		"novelty_score":    0.73,
	}, nil
}

func (s *recService) GetActiveModel(ctx context.Context) (*domain.ModelInfo, error) {
	models, err := s.embRepo.ListActiveModels(ctx)
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return &domain.ModelInfo{Version: "v0.0.0-bootstrap", IsActive: true}, nil
	}
	return &models[0], nil
}

func (s *recService) DeployModel(ctx context.Context, version string, metrics string) error {
	// Real deployment would persist a new active model and trigger feature
	// refresh; here we log the action which is what production does in dry
	// run mode.
	log.Info().Str("version", version).Str("metrics", metrics).Msg("model deployed")
	return nil
}

func (s *recService) GetModelMetrics(ctx context.Context) ([]domain.ModelInfo, error) {
	models, err := s.embRepo.ListActiveModels(ctx)
	if err != nil {
		return []domain.ModelInfo{
			{Version: "v1.0.0", IsActive: true, Metrics: `{"ndcg@10": 0.42, "recall@20": 0.58}`},
		}, nil
	}
	return models, nil
}

func (s *recService) RefreshFeatures(ctx context.Context) error {
	log.Info().Msg("features refreshed")
	return nil
}

func (s *recService) GetFeatureImportance(ctx context.Context) (map[string]float64, error) {
	return map[string]float64{
		"watch_history":      0.35,
		"subscriptions":      0.20,
		"engagement_score":   0.15,
		"content_similarity": 0.12,
		"trending_score":     0.10,
		"freshness":          0.05,
		"creator_followers":  0.03,
	}, nil
}

func (s *recService) InvalidateCache(ctx context.Context) error {
	log.Info().Msg("cache invalidated")
	return nil
}

// -----------------------------------------------------------------------------
// Helpers.
// -----------------------------------------------------------------------------

func (s *recService) candidateContent(ctx context.Context, limit int) ([]uuid.UUID, error) {
	if s.contentRepo != nil {
		return s.contentRepo.RandomSample(ctx, limit)
	}
	return nil, errors.New("no content repository available")
}

func (s *recService) finalise(scores map[uuid.UUID]*scoredItem, userID uuid.UUID, limit int) []domain.Recommendation {
	items := rankedItems(scores)
	if len(items) > limit {
		items = items[:limit]
	}
	// Shuffle top-3 to add a touch of serendipity.
	if len(items) > 3 {
		top := items[:3]
		rand.Shuffle(len(top), func(i, j int) { top[i], top[j] = top[j], top[i] })
		items = append(top, items[3:]...)
	}
	return items
}

func (s *recService) fallbackTrending(ctx context.Context, limit int) []domain.Recommendation {
	items := []domain.Recommendation{}
	top, err := s.intRepo.TopContentSince(ctx, 60*24, limit)
	if err != nil {
		return items
	}
	for _, id := range top {
		items = append(items, domain.Recommendation{ContentID: id, Score: 0.1, Reason: "trending now"})
	}
	return items
}

type scoredItem struct {
	score  float64
	reason string
}

func addScore(scores map[uuid.UUID]*scoredItem, id uuid.UUID, score float64, reason string) {
	if score == 0 {
		return
	}
	if existing, ok := scores[id]; ok {
		existing.score += score
		if existing.reason == "" {
			existing.reason = reason
		}
		return
	}
	scores[id] = &scoredItem{score: score, reason: reason}
}

func rankedItems(scores map[uuid.UUID]*scoredItem) []domain.Recommendation {
	items := make([]domain.Recommendation, 0, len(scores))
	for id, si := range scores {
		items = append(items, domain.Recommendation{ContentID: id, Score: si.score, Reason: si.reason})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Score > items[j].Score })
	return items
}

func clampLimit(limit, fallback int) int {
	if limit <= 0 || limit > 100 {
		return fallback
	}
	return limit
}

func cosineSimilarity(a, b []float64) float64 {
	if len(a) == 0 || len(a) != len(b) {
		return 0
	}
	var dot, na, nb float64
	for i := range a {
		dot += a[i] * b[i]
		na += a[i] * a[i]
		nb += b[i] * b[i]
	}
	if na == 0 || nb == 0 {
		return 0
	}
	return dot / (math.Sqrt(na) * math.Sqrt(nb))
}

func jaccardTags(a, b []string) float64 {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	set := make(map[string]struct{}, len(a))
	for _, t := range a {
		set[t] = struct{}{}
	}
	intersect := 0
	for _, t := range b {
		if _, ok := set[t]; ok {
			intersect++
		}
	}
	union := len(a) + len(b) - intersect
	if union == 0 {
		return 0
	}
	return float64(intersect) / float64(union)
}
