package handler

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/stream-service/internal/config"
	"github.com/elevatecompact/spark/services/stream-service/internal/domain"
)

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	StreamIDKey  contextKey = "stream_id"
	RoleKey      contextKey = "role"
	RequestIDKey contextKey = "request_id"
)

type Middleware struct {
	cfg *config.Config
}

func NewMiddleware(cfg *config.Config) *Middleware {
	return &Middleware{cfg: cfg}
}

func (m *Middleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("query", r.URL.RawQuery).
			Int("status", ww.Status()).
			Int("bytes", ww.BytesWritten()).
			Dur("duration", time.Since(start)).
			Str("ip", r.RemoteAddr).
			Str("user_agent", r.UserAgent()).
			Msg("request")
	})
}

func (m *Middleware) Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Error().
					Interface("recover", rec).
					Str("path", r.URL.Path).
					Msg("panic recovered")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, domain.ErrStreamInsufficientPermissions
			}
			return []byte(m.cfg.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, `{"error":"invalid token claims"}`, http.StatusUnauthorized)
			return
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			http.Error(w, `{"error":"invalid token subject"}`, http.StatusUnauthorized)
			return
		}

		userID, err := uuid.Parse(sub)
		if err != nil {
			http.Error(w, `{"error":"invalid user id in token"}`, http.StatusUnauthorized)
			return
		}

		role, _ := claims["role"].(string)

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, RoleKey, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) StreamAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		streamIDStr := r.PathValue("id")
		if streamIDStr == "" {
			streamIDStr = r.URL.Query().Get("stream_id")
		}

		if streamIDStr != "" {
			streamID, err := uuid.Parse(streamIDStr)
			if err == nil {
				ctx := context.WithValue(r.Context(), StreamIDKey, streamID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) CreatorOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
		if !ok {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		role, _ := r.Context().Value(RoleKey).(string)
		if role == "admin" || role == "moderator" {
			next.ServeHTTP(w, r)
			return
		}

		creatorIDStr := r.PathValue("creator_id")
		if creatorIDStr != "" {
			creatorID, err := uuid.Parse(creatorIDStr)
			if err == nil && creatorID == userID {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, `{"error":"forbidden: creator access required"}`, http.StatusForbidden)
	})
}

func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return id, ok
}

func GetStreamID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(StreamIDKey).(uuid.UUID)
	return id, ok
}
