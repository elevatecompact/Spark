package processor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// HTTPClient is the interface used by the various processors to send requests.
// It is a small subset of *http.Client so tests can plug in fakes.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

func defaultClient() HTTPClient { return &http.Client{Timeout: 15 * time.Second} }

// PushProcessor delivers push notifications to a device token.
type PushProcessor interface {
	Send(ctx context.Context, deviceToken string, title, body string, data map[string]string) error
}

// EmailProcessor sends transactional email.
type EmailProcessor interface {
	Send(ctx context.Context, to, subject, body string) error
}

// SMSProcessor sends SMS messages.
type SMSProcessor interface {
	Send(ctx context.Context, to, body string) error
}

// -----------------------------------------------------------------------------
// Noop implementations (used when credentials are missing or in dev/test).
// -----------------------------------------------------------------------------

type noopPush struct{}

func NewNoopPush() PushProcessor { return &noopPush{} }

func (p *noopPush) Send(ctx context.Context, deviceToken string, title, body string, data map[string]string) error {
	safe := deviceToken
	if len(safe) > 8 {
		safe = safe[:8] + "..."
	}
	log.Debug().Str("token", safe).Str("title", title).Msg("noop push sent")
	return nil
}

type noopEmail struct{}

func NewNoopEmail() EmailProcessor { return &noopEmail{} }

func (e *noopEmail) Send(ctx context.Context, to, subject, body string) error {
	log.Debug().Str("to", to).Str("subject", subject).Msg("noop email sent")
	return nil
}

type noopSMS struct{}

func NewNoopSMS() SMSProcessor { return &noopSMS{} }

func (s *noopSMS) Send(ctx context.Context, to, body string) error {
	log.Debug().Str("to", to).Msg("noop sms sent")
	return nil
}

// -----------------------------------------------------------------------------
// Firebase Cloud Messaging push implementation.
// https://firebase.google.com/docs/cloud-messaging/send-message
// -----------------------------------------------------------------------------

type fcmPush struct {
	serverKey string
	endpoint  string
	client    HTTPClient
}

func NewFCMPush(serverKey string) PushProcessor {
	if strings.TrimSpace(serverKey) == "" {
		return NewNoopPush()
	}
	return &fcmPush{
		serverKey: serverKey,
		endpoint:  "https://fcm.googleapis.com/fcm/send",
		client:    defaultClient(),
	}
}

func NewFCMPushWithClient(serverKey string, client HTTPClient) PushProcessor {
	if strings.TrimSpace(serverKey) == "" {
		return NewNoopPush()
	}
	return &fcmPush{serverKey: serverKey, endpoint: "https://fcm.googleapis.com/fcm/send", client: client}
}

type fcmPayload struct {
	To           string            `json:"to"`
	Notification *fcmNotification  `json:"notification,omitempty"`
	Data         map[string]string `json:"data,omitempty"`
}

type fcmNotification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type fcmResponse struct {
	Success      int      `json:"success"`
	Failure      int      `json:"failure"`
	MulticastID  int64    `json:"multicast_id"`
	FailedTokens []string `json:"failed_tokens,omitempty"`
}

func (p *fcmPush) Send(ctx context.Context, deviceToken string, title, body string, data map[string]string) error {
	payload := fcmPayload{
		To:           deviceToken,
		Notification: &fcmNotification{Title: title, Body: body},
		Data:         data,
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("fcm: marshal payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("fcm: build request: %w", err)
	}
	req.Header.Set("Authorization", "key="+p.serverKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("fcm: send: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("fcm: status %d: %s", resp.StatusCode, string(respBody))
	}

	var fcmResp fcmResponse
	if err := json.Unmarshal(respBody, &fcmResp); err != nil {
		return fmt.Errorf("fcm: decode response: %w", err)
	}
	if fcmResp.Failure > 0 {
		return fmt.Errorf("fcm: %d tokens failed", fcmResp.Failure)
	}
	log.Debug().Str("endpoint", p.endpoint).Int("success", fcmResp.Success).Msg("fcm push delivered")
	return nil
}

// -----------------------------------------------------------------------------
// SendGrid email implementation.
// https://docs.sendgrid.com/api-reference/mail-send/mail-send
// -----------------------------------------------------------------------------

type sendgridEmail struct {
	apiKey  string
	from    string
	endpoint string
	client  HTTPClient
}

func NewSendGridEmail(apiKey string) EmailProcessor {
	from := os.Getenv("SENDGRID_FROM_EMAIL")
	if from == "" {
		from = "no-reply@spark.dev"
	}
	if strings.TrimSpace(apiKey) == "" {
		return NewNoopEmail()
	}
	return &sendgridEmail{
		apiKey:   apiKey,
		from:     from,
		endpoint: "https://api.sendgrid.com/v3/mail/send",
		client:   defaultClient(),
	}
}

func NewSendGridEmailWithClient(apiKey string, client HTTPClient) EmailProcessor {
	if strings.TrimSpace(apiKey) == "" {
		return NewNoopEmail()
	}
	from := os.Getenv("SENDGRID_FROM_EMAIL")
	if from == "" {
		from = "no-reply@spark.dev"
	}
	return &sendgridEmail{apiKey: apiKey, from: from, endpoint: "https://api.sendgrid.com/v3/mail/send", client: client}
}

type sendgridPayload struct {
	Personalizations []sendgridPersonalization `json:"personalizations"`
	From             sendgridAddress           `json:"from"`
	Subject          string                    `json:"subject"`
	Content          []sendgridContent         `json:"content"`
}

type sendgridPersonalization struct {
	To []sendgridAddress `json:"to"`
}

type sendgridAddress struct {
	Email string `json:"email"`
}

type sendgridContent struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (e *sendgridEmail) Send(ctx context.Context, to, subject, body string) error {
	payload := sendgridPayload{
		Personalizations: []sendgridPersonalization{{To: []sendgridAddress{{Email: to}}}},
		From:             sendgridAddress{Email: e.from},
		Subject:          subject,
		Content:          []sendgridContent{{Type: "text/plain", Value: body}},
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("sendgrid: marshal payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("sendgrid: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+e.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return fmt.Errorf("sendgrid: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("sendgrid: status %d: %s", resp.StatusCode, string(respBody))
	}
	log.Debug().Str("to", to).Str("subject", subject).Msg("sendgrid email delivered")
	return nil
}

// -----------------------------------------------------------------------------
// Twilio SMS implementation.
// https://www.twilio.com/docs/sms/api
// -----------------------------------------------------------------------------

type twilioSMS struct {
	accountSID string
	authToken  string
	from       string
	endpoint   string
	client     HTTPClient
}

func NewTwilioSMS(sid, token, from string) SMSProcessor {
	if strings.TrimSpace(sid) == "" || strings.TrimSpace(token) == "" || strings.TrimSpace(from) == "" {
		return NewNoopSMS()
	}
	return &twilioSMS{
		accountSID: sid,
		authToken:  token,
		from:       from,
		endpoint:   "https://api.twilio.com/2010-04-01/Accounts/" + sid + "/Messages.json",
		client:     defaultClient(),
	}
}

func NewTwilioSMSWithClient(sid, token, from string, client HTTPClient) SMSProcessor {
	if strings.TrimSpace(sid) == "" || strings.TrimSpace(token) == "" || strings.TrimSpace(from) == "" {
		return NewNoopSMS()
	}
	return &twilioSMS{
		accountSID: sid,
		authToken:  token,
		from:       from,
		endpoint:   "https://api.twilio.com/2010-04-01/Accounts/" + sid + "/Messages.json",
		client:     client,
	}
}

func (s *twilioSMS) Send(ctx context.Context, to, body string) error {
	form := fmt.Sprintf("To=%s&From=%s&Body=%s", urlEncode(to), urlEncode(s.from), urlEncode(body))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.endpoint, strings.NewReader(form))
	if err != nil {
		return fmt.Errorf("twilio: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(s.accountSID, s.authToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("twilio: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("twilio: status %d: %s", resp.StatusCode, string(respBody))
	}
	log.Debug().Str("to", to).Msg("twilio sms delivered")
	return nil
}

// urlEncode is a tiny URL-form encoder that does not pull in net/url so the
// dependency surface stays light.
func urlEncode(s string) string {
	enc := ""
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == ' ':
			enc += "+"
		case (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' || c == '.' || c == '~':
			enc += string(c)
		default:
			enc += fmt.Sprintf("%%%02X", c)
		}
	}
	return enc
}
