package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/competition-service/internal/domain"
	"github.com/elevatecompact/spark/services/competition-service/internal/service"
)

type Handler struct {
	svc service.CompetitionService
}

func New(svc service.CompetitionService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r chi.Router) {
	r.Route("/v1/competitions", func(r chi.Router) {
		r.Post("/", h.create)
		r.Get("/", h.list)
		r.Get("/{id}", h.getByID)
		r.Patch("/{id}", h.update)
		r.Post("/{id}/start", h.start)
		r.Post("/{id}/end", h.end)
		r.Post("/{id}/register", h.register)
		r.Post("/{id}/withdraw", h.withdraw)
		r.Get("/{id}/participants", h.listParticipants)
		r.Get("/{id}/bracket", h.getBracket)
		r.Get("/{id}/leaderboard", h.getLeaderboard)
		r.Get("/{id}/results", h.getResults)
		r.Get("/{id}/prizes", h.getPrizes)
		r.Post("/{id}/prizes/distribute", h.distributePrizes)
		r.Post("/{id}/judges", h.assignJudge)
	})
	r.Route("/v1/matches", func(r chi.Router) {
		r.Post("/{id}/score", h.submitScore)
		r.Post("/{id}/confirm", h.confirmMatch)
		r.Post("/{id}/dispute", h.disputeMatch)
	})
	r.Route("/v1/submissions", func(r chi.Router) {
		r.Get("/", h.listSubmissions)
		r.Post("/{id}/score", h.scoreSubmission)
	})
	r.Route("/v1/admin/competitions", func(r chi.Router) {
		r.Post("/{id}/cancel", h.cancel)
	})
	r.Route("/v1/admin/matches", func(r chi.Router) {
		r.Post("/{id}/override", h.overrideMatch)
	})
	r.Get("/v1/admin/stats", h.getAdminStats)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var c domain.Competition
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	result, err := h.svc.Create(r.Context(), &c)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	c, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "competition not found")
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	var c domain.Competition
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	c.ID = uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.Update(r.Context(), &c); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	status := domain.CompetitionStatus(r.URL.Query().Get("status"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	comps, err := h.svc.List(r.Context(), status, page, size)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, comps)
}

func (h *Handler) start(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.Start(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "started"})
}

func (h *Handler) end(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.End(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ended"})
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	compID := uuid.MustParse(chi.URLParam(r, "id"))
	var req struct{ UserID uuid.UUID `json:"userId"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.Register(r.Context(), compID, req.UserID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "registered"})
}

func (h *Handler) withdraw(w http.ResponseWriter, r *http.Request) {
	compID := uuid.MustParse(chi.URLParam(r, "id"))
	var req struct{ UserID uuid.UUID `json:"userId"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.Withdraw(r.Context(), compID, req.UserID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "withdrawn"})
}

func (h *Handler) listParticipants(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	parts, err := h.svc.ListParticipants(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, parts)
}

func (h *Handler) getBracket(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	matches, err := h.svc.GetBracket(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, matches)
}

func (h *Handler) submitScore(w http.ResponseWriter, r *http.Request) {
	matchID := uuid.MustParse(chi.URLParam(r, "id"))
	var req struct {
		WinnerID uuid.UUID              `json:"winnerId"`
		Scores   map[string]interface{} `json:"scores"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.SubmitScore(r.Context(), matchID, req.WinnerID, req.Scores); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "scored"})
}

func (h *Handler) confirmMatch(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.ConfirmMatch(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "confirmed"})
}

func (h *Handler) disputeMatch(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.DisputeMatch(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "disputed"})
}

func (h *Handler) assignJudge(w http.ResponseWriter, r *http.Request) {
	compID := uuid.MustParse(chi.URLParam(r, "id"))
	var req struct{ JudgeID uuid.UUID `json:"judgeId"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.AssignJudge(r.Context(), compID, req.JudgeID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "assigned"})
}

func (h *Handler) listSubmissions(w http.ResponseWriter, r *http.Request) {
	compID := uuid.MustParse(r.URL.Query().Get("competitionId"))
	subs, err := h.svc.ListSubmissions(r.Context(), compID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, subs)
}

func (h *Handler) scoreSubmission(w http.ResponseWriter, r *http.Request) {
	subID := uuid.MustParse(chi.URLParam(r, "id"))
	var req struct {
		JudgeID uuid.UUID `json:"judgeId"`
		Score   float64   `json:"score"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.ScoreSubmission(r.Context(), subID, req.JudgeID, req.Score); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "scored"})
}

func (h *Handler) getLeaderboard(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	entries, err := h.svc.GetLeaderboard(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, entries)
}

func (h *Handler) getResults(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	results, err := h.svc.GetResults(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, results)
}

func (h *Handler) getPrizes(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	prizes, err := h.svc.GetPrizes(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, prizes)
}

func (h *Handler) distributePrizes(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.DistributePrizes(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "distributed"})
}

func (h *Handler) cancel(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.Cancel(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "cancelled"})
}

func (h *Handler) overrideMatch(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	var req struct{ WinnerID uuid.UUID `json:"winnerId"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.OverrideMatch(r.Context(), id, req.WinnerID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "overridden"})
}

func (h *Handler) getAdminStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.svc.GetAdminStats(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
