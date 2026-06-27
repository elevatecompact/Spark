package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/elevatecompact/spark/services/identity-service/internal/domain"
	"github.com/elevatecompact/spark/services/identity-service/internal/service"
)

type OAuthHandler struct {
	oauthSvc service.OAuthService
}

func NewOAuthHandler(oauthSvc service.OAuthService) *OAuthHandler {
	return &OAuthHandler{oauthSvc: oauthSvc}
}

type AuthorizeRequest struct {
	ClientID     string `json:"client_id"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
	State        string `json:"state"`
	ResponseType string `json:"response_type"`
}

func (h *OAuthHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromContext(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	var req AuthorizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.ResponseType != "code" {
		WriteError(w, http.StatusBadRequest, "unsupported response_type")
		return
	}

	authCode, err := h.oauthSvc.Authorize(r.Context(), req.ClientID, req.RedirectURI, req.Scope, req.State, user.ID)
	if err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	redirectURL, err := url.Parse(req.RedirectURI)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid redirect_uri")
		return
	}

	params := url.Values{}
	params.Set("code", authCode.Code)
	if req.State != "" {
		params.Set("state", req.State)
	}
	redirectURL.RawQuery = params.Encode()

	WriteJSON(w, http.StatusOK, map[string]string{
		"redirect_url": redirectURL.String(),
		"code":         authCode.Code,
	})
}

func (h *OAuthHandler) AuthorizeForm(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromContext(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login?redirect="+url.QueryEscape(r.RequestURI), http.StatusFound)
		return
	}

	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	scope := r.URL.Query().Get("scope")
	state := r.URL.Query().Get("state")

	authCode, err := h.oauthSvc.Authorize(r.Context(), clientID, redirectURI, scope, state, user.ID)
	if err != nil {
		WriteError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}

	redirectURL, err := url.Parse(redirectURI)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid redirect_uri")
		return
	}

	params := url.Values{}
	params.Set("code", authCode.Code)
	if state != "" {
		params.Set("state", state)
	}
	redirectURL.RawQuery = params.Encode()

	http.Redirect(w, r, redirectURL.String(), http.StatusFound)
}

type TokenRequest struct {
	GrantType    string `json:"grant_type"`
	Code         string `json:"code,omitempty"`
	RedirectURI  string `json:"redirect_uri,omitempty"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scope        string `json:"scope,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

func (h *OAuthHandler) Token(w http.ResponseWriter, r *http.Request) {
	var req TokenRequest

	contentType := r.Header.Get("Content-Type")
	if contentType == "application/x-www-form-urlencoded" {
		if err := r.ParseForm(); err != nil {
			WriteError(w, http.StatusBadRequest, "invalid form data")
			return
		}
		req.GrantType = r.FormValue("grant_type")
		req.Code = r.FormValue("code")
		req.RedirectURI = r.FormValue("redirect_uri")
		req.ClientID = r.FormValue("client_id")
		req.ClientSecret = r.FormValue("client_secret")
		req.Scope = r.FormValue("scope")
		req.RefreshToken = r.FormValue("refresh_token")
		req.Username = r.FormValue("username")
		req.Password = r.FormValue("password")
	} else {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, http.StatusBadRequest, "invalid request body")
			return
		}
	}

	var oauthToken *domain.OAuthToken
	var err error

	switch req.GrantType {
	case "authorization_code":
		oauthToken, err = h.oauthSvc.ExchangeAuthorizationCode(r.Context(), req.Code, req.ClientID, req.ClientSecret, req.RedirectURI)
	case "password":
		oauthToken, err = h.oauthSvc.ExchangePasswordCredentials(r.Context(), req.ClientID, req.ClientSecret, req.Username, req.Password, req.Scope)
	case "refresh_token":
		oauthToken, err = h.oauthSvc.ExchangeRefreshToken(r.Context(), req.ClientID, req.RefreshToken)
	default:
		WriteError(w, http.StatusBadRequest, "unsupported grant_type")
		return
	}

	if err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	expiresIn := int(time.Until(oauthToken.ExpiresAt).Seconds())
	if expiresIn < 0 {
		expiresIn = 0
	}

	resp := TokenResponse{
		AccessToken:  oauthToken.AccessToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		RefreshToken: oauthToken.RefreshToken,
		Scope:        oauthToken.Scope,
	}

	WriteJSON(w, http.StatusOK, resp)
}

type IntrospectRequest struct {
	Token string `json:"token"`
}

func (h *OAuthHandler) Introspect(w http.ResponseWriter, r *http.Request) {
	var req IntrospectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	introspection, err := h.oauthSvc.IntrospectToken(r.Context(), req.Token)
	if err != nil {
		WriteJSON(w, http.StatusOK, map[string]bool{"active": false})
		return
	}

	WriteJSON(w, http.StatusOK, introspection)
}

func (h *OAuthHandler) UserInfo(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		WriteError(w, http.StatusUnauthorized, "missing authorization header")
		return
	}

	user, err := h.oauthSvc.ValidateAccessToken(r.Context(), token)
	if err != nil {
		WriteJSON(w, http.StatusOK, map[string]string{"error": "invalid_token"})
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"sub":                user.ID.String(),
		"name":               user.DisplayName,
		"preferred_username": user.Username,
		"email":              user.Email,
		"email_verified":     user.Verified,
		"picture":            user.AvatarURL,
	})
}

func (h *OAuthHandler) OpenIDDiscovery(w http.ResponseWriter, r *http.Request) {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	baseURL := scheme + "://" + r.Host

	config := h.oauthSvc.GetDiscoveryConfig(r.Context(), baseURL)
	WriteJSON(w, http.StatusOK, config)
}
