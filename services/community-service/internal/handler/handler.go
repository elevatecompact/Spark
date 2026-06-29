package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/community-service/internal/domain"
	"github.com/elevatecompact/spark/services/community-service/internal/service"
)

type Handler struct {
	svc service.CommunityService
}

func New(svc service.CommunityService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r chi.Router) {
	r.Route("/v1/communities", func(r chi.Router) {
		r.Post("/", h.create)
		r.Get("/", h.list)
		r.Get("/{id}", h.getByID)
		r.Patch("/{id}", h.update)
		r.Delete("/{id}", h.delete)

		r.Post("/{id}/join", h.join)
		r.Post("/{id}/leave", h.leave)
		r.Get("/{id}/members", h.listMembers)
		r.Patch("/{id}/members/{userId}/role", h.updateMemberRole)

		r.Post("/{id}/posts", h.createPost)
		r.Get("/{id}/posts", h.listPosts)
	})
	r.Route("/v1/posts", func(r chi.Router) {
		r.Get("/{id}", h.getPost)
		r.Put("/{id}", h.updatePost)
		r.Delete("/{id}", h.deletePost)
		r.Post("/{id}/pin", h.pinPost)
		r.Post("/{id}/reactions", h.reactToPost)
		r.Post("/{id}/comments", h.createComment)
		r.Get("/{id}/comments", h.listComments)
	})
	r.Route("/v1/comments", func(r chi.Router) {
		r.Delete("/{id}", h.deleteComment)
		r.Post("/{id}/reactions", h.reactToComment)
	})
	r.Route("/v1/admin/communities", func(r chi.Router) {
		r.Post("/{id}/feature", h.feature)
		r.Post("/{id}/suspend", h.suspend)
	})
	r.Get("/v1/admin/stats", h.getAdminStats)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var c domain.Community
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
		writeError(w, http.StatusNotFound, "community not found")
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	var c domain.Community
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

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	communities, err := h.svc.List(r.Context(), category, page, size)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, communities)
}

func (h *Handler) join(w http.ResponseWriter, r *http.Request) {
	communityID := uuid.MustParse(chi.URLParam(r, "id"))
	var req struct{ UserID uuid.UUID `json:"userId"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.Join(r.Context(), communityID, req.UserID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "joined"})
}

func (h *Handler) leave(w http.ResponseWriter, r *http.Request) {
	communityID := uuid.MustParse(chi.URLParam(r, "id"))
	var req struct{ UserID uuid.UUID `json:"userId"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.Leave(r.Context(), communityID, req.UserID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "left"})
}

func (h *Handler) listMembers(w http.ResponseWriter, r *http.Request) {
	communityID := uuid.MustParse(chi.URLParam(r, "id"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	members, err := h.svc.ListMembers(r.Context(), communityID, page, size)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, members)
}

func (h *Handler) updateMemberRole(w http.ResponseWriter, r *http.Request) {
	communityID := uuid.MustParse(chi.URLParam(r, "id"))
	userID := uuid.MustParse(chi.URLParam(r, "userId"))
	var req struct{ Role domain.MemberRole `json:"role"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.UpdateMemberRole(r.Context(), communityID, userID, req.Role); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "role updated"})
}

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	communityID := uuid.MustParse(chi.URLParam(r, "id"))
	var p domain.CommunityPost
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	p.CommunityID = communityID
	result, err := h.svc.CreatePost(r.Context(), &p)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) getPost(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	p, err := h.svc.GetPost(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "post not found")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *Handler) updatePost(w http.ResponseWriter, r *http.Request) {
	var p domain.CommunityPost
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	p.ID = uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.UpdatePost(r.Context(), &p); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *Handler) deletePost(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.DeletePost(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *Handler) listPosts(w http.ResponseWriter, r *http.Request) {
	communityID := uuid.MustParse(chi.URLParam(r, "id"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	posts, err := h.svc.ListPosts(r.Context(), communityID, page, size)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, posts)
}

func (h *Handler) pinPost(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.PinPost(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "pinned"})
}

func (h *Handler) reactToPost(w http.ResponseWriter, r *http.Request) {
	postID := uuid.MustParse(chi.URLParam(r, "id"))
	var req struct {
		UserID uuid.UUID `json:"userId"`
		Emoji  string    `json:"emoji"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.ReactToPost(r.Context(), postID, req.UserID, req.Emoji); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "reacted"})
}

func (h *Handler) createComment(w http.ResponseWriter, r *http.Request) {
	postID := uuid.MustParse(chi.URLParam(r, "id"))
	var c domain.PostComment
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	c.PostID = postID
	result, err := h.svc.CreateComment(r.Context(), &c)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) listComments(w http.ResponseWriter, r *http.Request) {
	postID := uuid.MustParse(chi.URLParam(r, "id"))
	comments, err := h.svc.ListComments(r.Context(), postID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, comments)
}

func (h *Handler) deleteComment(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.DeleteComment(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *Handler) reactToComment(w http.ResponseWriter, r *http.Request) {
	commentID := uuid.MustParse(chi.URLParam(r, "id"))
	var req struct {
		UserID uuid.UUID `json:"userId"`
		Emoji  string    `json:"emoji"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.ReactToComment(r.Context(), commentID, req.UserID, req.Emoji); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "reacted"})
}

func (h *Handler) feature(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.FeatureCommunity(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "featured"})
}

func (h *Handler) suspend(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.SuspendCommunity(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "suspended"})
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
