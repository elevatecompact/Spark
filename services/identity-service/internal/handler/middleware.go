package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/identity-service/internal/domain"
	"github.com/elevatecompact/spark/services/identity-service/internal/service"
)

type contextKey string

const (
	ContextUserKey    contextKey = "user"
	ContextRequestID  contextKey = "request_id"
	ContextSessionKey contextKey = "session"
)

func AuthRequired(tokenSvc service.TokenService, authSvc service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				WriteError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				WriteError(w, http.StatusUnauthorized, "invalid authorization header format")
				return
			}

			tokenString := parts[1]

			user, err := authSvc.ValidateToken(r.Context(), tokenString)
			if err != nil {
				WriteError(w, domain.HTTPStatusFromError(err), err.Error())
				return
			}

			ctx := context.WithValue(r.Context(), ContextUserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func OptionalAuth(tokenSvc service.TokenService, authSvc service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				next.ServeHTTP(w, r)
				return
			}

			user, err := authSvc.ValidateToken(r.Context(), parts[1])
			if err == nil && user != nil {
				ctx := context.WithValue(r.Context(), ContextUserKey, user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func AdminRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(ContextUserKey).(*domain.User)
		if !ok || user == nil {
			WriteError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		if user.Role != domain.RoleAdmin {
			WriteError(w, http.StatusForbidden, "admin access required")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ModOrAdminRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(ContextUserKey).(*domain.User)
		if !ok || user == nil {
			WriteError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		if user.Role != domain.RoleAdmin && user.Role != domain.RoleMod {
			WriteError(w, http.StatusForbidden, "moderator or admin access required")
			return
		}

		next.ServeHTTP(w, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)
		duration := time.Since(start)

		log.Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Int("status", rw.statusCode).
			Int("size", rw.size).
			Dur("duration", duration).
			Str("remote_addr", r.RemoteAddr).
			Str("user_agent", r.UserAgent()).
			Msg("request")
	})
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		w.Header().Set("X-Request-Id", requestID)
		ctx := context.WithValue(r.Context(), ContextRequestID, requestID)

		logger := log.With().Str("request_id", requestID).Logger()
		ctx = logger.WithContext(ctx)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Error().
					Interface("panic", rec).
					Str("url", r.URL.String()).
					Msg("panic recovered")
				WriteError(w, http.StatusInternalServerError, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*rateLimiterEntry
	rate     int
	burst    int
}

type rateLimiterEntry struct {
	tokens    int
	lastCheck time.Time
}

func NewRateLimiter(rate, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*rateLimiterEntry),
		rate:     rate,
		burst:    burst,
	}

	go func() {
		ticker := time.NewTicker(time.Minute)
		for range ticker.C {
			rl.mu.Lock()
			rl.visitors = make(map[string]*rateLimiterEntry)
			rl.mu.Unlock()
		}
	}()

	return rl
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	entry, exists := rl.visitors[ip]
	now := time.Now()

	if !exists {
		rl.visitors[ip] = &rateLimiterEntry{tokens: rl.burst - 1, lastCheck: now}
		return true
	}

	elapsed := now.Sub(entry.lastCheck).Seconds()
	entry.tokens += int(elapsed * float64(rl.rate))
	if entry.tokens > rl.burst {
		entry.tokens = rl.burst
	}
	entry.lastCheck = now

	if entry.tokens > 0 {
		entry.tokens--
		return true
	}

	return false
}

func RateLimitMiddleware(rate, burst int) func(http.Handler) http.Handler {
	limiter := NewRateLimiter(rate, burst)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if idx := strings.LastIndex(ip, ":"); idx != -1 {
				ip = ip[:idx]
			}

			forwarded := r.Header.Get("X-Forwarded-For")
			if forwarded != "" {
				ips := strings.Split(forwarded, ",")
				ip = strings.TrimSpace(ips[0])
			}

			if !limiter.allow(ip) {
				w.Header().Set("Retry-After", "60")
				WriteError(w, http.StatusTooManyRequests, "rate limit exceeded")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(data)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}

func GetUserFromContext(r *http.Request) (*domain.User, error) {
	user, ok := r.Context().Value(ContextUserKey).(*domain.User)
	if !ok || user == nil {
		return nil, domain.ErrUnauthorized
	}
	return user, nil
}
