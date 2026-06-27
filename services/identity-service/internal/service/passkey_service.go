package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/identity-service/internal/domain"
	"github.com/elevatecompact/spark/services/identity-service/internal/repository"
)

type CredentialPublicKey struct {
	KeyType   int    `json:"kty"`
	Algorithm int    `json:"alg"`
	Curve     int    `json:"crv,omitempty"`
	X         []byte `json:"x,omitempty"`
	Y         []byte `json:"y,omitempty"`
	N         []byte `json:"n,omitempty"`
	E         []byte `json:"e,omitempty"`
}

type PasskeyCredential struct {
	ID              string              `json:"id"`
	PublicKey       CredentialPublicKey `json:"publicKey"`
	AttestationType string              `json:"attestationType"`
	Transports      []string            `json:"transports"`
	AAGUID          string              `json:"aaguid"`
	SignCount       uint32              `json:"signCount"`
}

type CredentialCreationOptions struct {
	Challenge          string   `json:"challenge"`
	RP                 RPInfo   `json:"rp"`
	User               UserInfo `json:"user"`
	PubKeyCredParams   []struct {
		Type string `json:"type"`
		Alg  int    `json:"alg"`
	} `json:"pubKeyCredParams"`
	AuthenticatorSelection AuthenticatorSelection `json:"authenticatorSelection"`
	Attestation            string                 `json:"attestation"`
}

type CredentialAssertionOptions struct {
	Challenge        string   `json:"challenge"`
	RPID             string   `json:"rpId"`
	AllowCredentials []struct {
		Type       string   `json:"type"`
		ID         string   `json:"id"`
		Transports []string `json:"transports,omitempty"`
	} `json:"allowCredentials"`
	UserVerification string `json:"userVerification"`
}

type RPInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type AuthenticatorSelection struct {
	AuthenticatorAttachment string `json:"authenticatorAttachment"`
	ResidentKey             string `json:"residentKey"`
	UserVerification        string `json:"userVerification"`
}

type AuthenticatorAttestationResponse struct {
	ID              string `json:"id"`
	RawID           string `json:"rawId"`
	AttestationObj  string `json:"attestationObject"`
	ClientDataJSON  string `json:"clientDataJSON"`
	Transports      []string `json:"transports"`
}

type AuthenticatorAssertionResponse struct {
	ID              string `json:"id"`
	RawID           string `json:"rawId"`
	AuthenticatorData string `json:"authenticatorData"`
	ClientDataJSON  string `json:"clientDataJSON"`
	Signature       string `json:"signature"`
	UserHandle      string `json:"userHandle,omitempty"`
}

type storedChallenge struct {
	Challenge string    `json:"challenge"`
	UserID    uuid.UUID `json:"user_id"`
	Purpose   string    `json:"purpose"`
	ExpiresAt time.Time `json:"expires_at"`
}

type PasskeyService interface {
	BeginRegistration(ctx context.Context, user *domain.User) (*CredentialCreationOptions, error)
	FinishRegistration(ctx context.Context, user *domain.User, response AuthenticatorAttestationResponse) error
	BeginAuthentication(ctx context.Context, user *domain.User) (*CredentialAssertionOptions, error)
	FinishAuthentication(ctx context.Context, user *domain.User, response AuthenticatorAssertionResponse) (*domain.Session, error)
	GetPasskeys(ctx context.Context, userID uuid.UUID) ([]*domain.Passkey, error)
	DeletePasskey(ctx context.Context, passkeyID uuid.UUID, userID uuid.UUID) error
}

type passkeyStore interface {
	SaveCredential(ctx context.Context, passkey *domain.Passkey) error
	GetCredentialsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Passkey, error)
	GetCredentialByID(ctx context.Context, credentialID string) (*domain.Passkey, error)
	DeleteCredential(ctx context.Context, id uuid.UUID) error
	SaveChallenge(ctx context.Context, challenge storedChallenge) error
	GetChallenge(ctx context.Context, challengeID string) (*storedChallenge, error)
	DeleteChallenge(ctx context.Context, challengeID string) error
}

type inMemoryPasskeyStore struct {
	credentials    map[string]*domain.Passkey
	userCreds      map[string][]*domain.Passkey
	challenges     map[string]storedChallenge
}

func newInMemoryPasskeyStore() *inMemoryPasskeyStore {
	return &inMemoryPasskeyStore{
		credentials: make(map[string]*domain.Passkey),
		userCreds:   make(map[string][]*domain.Passkey),
		challenges:  make(map[string]storedChallenge),
	}
}

func (s *inMemoryPasskeyStore) SaveCredential(ctx context.Context, passkey *domain.Passkey) error {
	s.credentials[passkey.CredentialID] = passkey
	uid := passkey.UserID.String()
	s.userCreds[uid] = append(s.userCreds[uid], passkey)
	return nil
}

func (s *inMemoryPasskeyStore) GetCredentialsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Passkey, error) {
	creds := s.userCreds[userID.String()]
	if creds == nil {
		return []*domain.Passkey{}, nil
	}
	return creds, nil
}

func (s *inMemoryPasskeyStore) GetCredentialByID(ctx context.Context, credentialID string) (*domain.Passkey, error) {
	cred, ok := s.credentials[credentialID]
	if !ok {
		return nil, domain.ErrPasskeyNotFound
	}
	return cred, nil
}

func (s *inMemoryPasskeyStore) DeleteCredential(ctx context.Context, id uuid.UUID) error {
	for cid, cred := range s.credentials {
		if cred.ID == id {
			delete(s.credentials, cid)
			uid := cred.UserID.String()
			creds := s.userCreds[uid]
			for i, c := range creds {
				if c.ID == id {
					s.userCreds[uid] = append(creds[:i], creds[i+1:]...)
					break
				}
			}
			return nil
		}
	}
	return domain.ErrPasskeyNotFound
}

func (s *inMemoryPasskeyStore) SaveChallenge(ctx context.Context, challenge storedChallenge) error {
	h := sha256.Sum256([]byte(challenge.Challenge))
	id := base64.RawURLEncoding.EncodeToString(h[:])
	s.challenges[id] = challenge
	return nil
}

func (s *inMemoryPasskeyStore) GetChallenge(ctx context.Context, challengeID string) (*storedChallenge, error) {
	h := sha256.Sum256([]byte(challengeID))
	id := base64.RawURLEncoding.EncodeToString(h[:])
	ch, ok := s.challenges[id]
	if !ok {
		return nil, domain.ErrPasskeyNotFound
	}
	if time.Now().After(ch.ExpiresAt) {
		delete(s.challenges, id)
		return nil, domain.ErrPasskeyNotFound
	}
	return &ch, nil
}

func (s *inMemoryPasskeyStore) DeleteChallenge(ctx context.Context, challengeID string) error {
	h := sha256.Sum256([]byte(challengeID))
	id := base64.RawURLEncoding.EncodeToString(h[:])
	delete(s.challenges, id)
	return nil
}

type passkeyService struct {
	store      passkeyStore
	tokenSvc   TokenService
	sessionRepo repository.SessionRepository
	sessionTTL time.Duration
	cfg        PasskeyConfig
}

type PasskeyConfig struct {
	RPID   string
	RPName string
	Origin string
}

func NewPasskeyService(
	sessionRepo repository.SessionRepository,
	tokenSvc TokenService,
	cfg PasskeyConfig,
) PasskeyService {
	return &passkeyService{
		store:       newInMemoryPasskeyStore(),
		tokenSvc:    tokenSvc,
		sessionRepo: sessionRepo,
		sessionTTL:  24 * time.Hour,
		cfg:         cfg,
	}
}

func (s *passkeyService) BeginRegistration(ctx context.Context, user *domain.User) (*CredentialCreationOptions, error) {
	challenge := make([]byte, 32)
	if _, err := rand.Read(challenge); err != nil {
		return nil, fmt.Errorf("failed to generate challenge: %w", err)
	}
	challengeB64 := base64.RawURLEncoding.EncodeToString(challenge)

	userIDBytes := user.ID[:]

	sc := storedChallenge{
		Challenge: challengeB64,
		UserID:    user.ID,
		Purpose:   "registration",
		ExpiresAt: time.Now().UTC().Add(5 * time.Minute),
	}
	if err := s.store.SaveChallenge(ctx, sc); err != nil {
		return nil, fmt.Errorf("failed to save challenge: %w", err)
	}

	userIDForWebAuthn := base64.RawURLEncoding.EncodeToString(userIDBytes)

	options := &CredentialCreationOptions{
		Challenge: challengeB64,
		RP: RPInfo{
			ID:   s.cfg.RPID,
			Name: s.cfg.RPName,
		},
		User: UserInfo{
			ID:          userIDForWebAuthn,
			Name:        user.Username,
			DisplayName: user.DisplayName,
		},
		PubKeyCredParams: []struct {
			Type string `json:"type"`
			Alg  int    `json:"alg"`
		}{
			{Type: "public-key", Alg: -7},
			{Type: "public-key", Alg: -257},
		},
		AuthenticatorSelection: AuthenticatorSelection{
			AuthenticatorAttachment: "platform",
			ResidentKey:             "preferred",
			UserVerification:        "preferred",
		},
		Attestation: "none",
	}

	return options, nil
}

func (s *passkeyService) FinishRegistration(ctx context.Context, user *domain.User, response AuthenticatorAttestationResponse) error {
	clientDataJSON, err := base64.RawURLEncoding.DecodeString(response.ClientDataJSON)
	if err != nil {
		return fmt.Errorf("%w: invalid client data json encoding", domain.ErrPasskeyRegistration)
	}

	var clientData struct {
		Type      string `json:"type"`
		Challenge string `json:"challenge"`
		Origin    string `json:"origin"`
	}
	if err := json.Unmarshal(clientDataJSON, &clientData); err != nil {
		return fmt.Errorf("%w: invalid client data json", domain.ErrPasskeyRegistration)
	}

	if clientData.Type != "webauthn.create" {
		return fmt.Errorf("%w: invalid client data type", domain.ErrPasskeyRegistration)
	}

	if clientData.Origin != s.cfg.Origin {
		return fmt.Errorf("%w: invalid origin", domain.ErrPasskeyRegistration)
	}

	stored, err := s.store.GetChallenge(ctx, clientData.Challenge)
	if err != nil {
		return fmt.Errorf("%w: challenge not found or expired", domain.ErrPasskeyRegistration)
	}
	if stored.UserID != user.ID {
		return fmt.Errorf("%w: challenge user mismatch", domain.ErrPasskeyRegistration)
	}
	if stored.Purpose != "registration" {
		return fmt.Errorf("%w: challenge purpose mismatch", domain.ErrPasskeyRegistration)
	}

	s.store.DeleteChallenge(ctx, clientData.Challenge)

	attestationBytes, err := base64.RawURLEncoding.DecodeString(response.AttestationObj)
	if err != nil {
		return fmt.Errorf("%w: invalid attestation object", domain.ErrPasskeyRegistration)
	}

	cred := &domain.Passkey{
		ID:              uuid.New(),
		UserID:          user.ID,
		CredentialID:    response.ID,
		PublicKey:       attestationBytes,
		AttestationType: "none",
		Transports:      response.Transports,
		SignCount:       0,
		Name:            fmt.Sprintf("Device - %s", time.Now().Format("Jan 2, 2006")),
		DeviceType:      "platform",
		BackedUp:        false,
		CreatedAt:       time.Now().UTC(),
		LastUsedAt:      nil,
	}

	if err := s.store.SaveCredential(ctx, cred); err != nil {
		return fmt.Errorf("failed to save credential: %w", err)
	}

	return nil
}

func (s *passkeyService) BeginAuthentication(ctx context.Context, user *domain.User) (*CredentialAssertionOptions, error) {
	challenge := make([]byte, 32)
	if _, err := rand.Read(challenge); err != nil {
		return nil, fmt.Errorf("failed to generate challenge: %w", err)
	}
	challengeB64 := base64.RawURLEncoding.EncodeToString(challenge)

	sc := storedChallenge{
		Challenge: challengeB64,
		UserID:    user.ID,
		Purpose:   "authentication",
		ExpiresAt: time.Now().UTC().Add(5 * time.Minute),
	}
	if err := s.store.SaveChallenge(ctx, sc); err != nil {
		return nil, fmt.Errorf("failed to save challenge: %w", err)
	}

	credentials, err := s.store.GetCredentialsByUserID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	options := &CredentialAssertionOptions{
		Challenge:        challengeB64,
		RPID:             s.cfg.RPID,
		AllowCredentials: make([]struct {
			Type       string   `json:"type"`
			ID         string   `json:"id"`
			Transports []string `json:"transports,omitempty"`
		}, 0),
		UserVerification: "preferred",
	}

	for _, cred := range credentials {
		options.AllowCredentials = append(options.AllowCredentials, struct {
			Type       string   `json:"type"`
			ID         string   `json:"id"`
			Transports []string `json:"transports,omitempty"`
		}{
			Type:       "public-key",
			ID:         cred.CredentialID,
			Transports: cred.Transports,
		})
	}

	return options, nil
}

func (s *passkeyService) FinishAuthentication(ctx context.Context, user *domain.User, response AuthenticatorAssertionResponse) (*domain.Session, error) {
	clientDataJSON, err := base64.RawURLEncoding.DecodeString(response.ClientDataJSON)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid client data json encoding", domain.ErrPasskeyAuthentication)
	}

	var clientData struct {
		Type      string `json:"type"`
		Challenge string `json:"challenge"`
		Origin    string `json:"origin"`
	}
	if err := json.Unmarshal(clientDataJSON, &clientData); err != nil {
		return nil, fmt.Errorf("%w: invalid client data json", domain.ErrPasskeyAuthentication)
	}

	if clientData.Type != "webauthn.get" {
		return nil, fmt.Errorf("%w: invalid client data type", domain.ErrPasskeyAuthentication)
	}

	if clientData.Origin != s.cfg.Origin {
		return nil, fmt.Errorf("%w: invalid origin", domain.ErrPasskeyAuthentication)
	}

	stored, err := s.store.GetChallenge(ctx, clientData.Challenge)
	if err != nil {
		return nil, fmt.Errorf("%w: challenge not found or expired", domain.ErrPasskeyAuthentication)
	}
	if stored.UserID != user.ID {
		return nil, fmt.Errorf("%w: challenge user mismatch", domain.ErrPasskeyAuthentication)
	}
	if stored.Purpose != "authentication" {
		return nil, fmt.Errorf("%w: challenge purpose mismatch", domain.ErrPasskeyAuthentication)
	}

	s.store.DeleteChallenge(ctx, clientData.Challenge)

	cred, err := s.store.GetCredentialByID(ctx, response.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: credential not found", domain.ErrPasskeyAuthentication)
	}

	if cred.UserID != user.ID {
		return nil, fmt.Errorf("%w: credential user mismatch", domain.ErrPasskeyAuthentication)
	}

	now := time.Now().UTC()
	cred.LastUsedAt = &now

	accessToken, err := s.tokenSvc.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, _, err := s.tokenSvc.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	session := &domain.Session{
		ID:           uuid.New(),
		UserID:       user.ID,
		Token:        accessToken,
		RefreshToken: refreshToken,
		TokenHash:    s.tokenSvc.HashToken(accessToken),
		RefreshHash:  s.tokenSvc.HashToken(refreshToken),
		IPAddress:    "",
		UserAgent:    "passkey",
		ExpiresAt:    now.Add(s.sessionTTL),
		CreatedAt:    now,
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

func (s *passkeyService) GetPasskeys(ctx context.Context, userID uuid.UUID) ([]*domain.Passkey, error) {
	return s.store.GetCredentialsByUserID(ctx, userID)
}

func (s *passkeyService) DeletePasskey(ctx context.Context, passkeyID uuid.UUID, userID uuid.UUID) error {
	return s.store.DeleteCredential(ctx, passkeyID)
}


