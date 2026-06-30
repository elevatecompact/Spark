package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// PaymentResult communicates the outcome of an external payment call.
type PaymentResult struct {
	Success     bool
	ExternalRef string
	Error       string
}

// HTTPClient is a small abstraction over net/http that processors use.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

func defaultClient() HTTPClient { return &http.Client{Timeout: 30 * time.Second} }

// PaymentProcessor is the contract every external money-movement service
// satisfies. Implementations include the in-memory noop (used in dev/test)
// and the Stripe Connect + PayPal Payout + bank ACH backends below.
type PaymentProcessor interface {
	Deposit(ctx context.Context, userID string, amountCents int64, currency string) (*PaymentResult, error)
	Withdraw(ctx context.Context, userID string, amountCents int64, currency string, method string) (*PaymentResult, error)
	Payout(ctx context.Context, userID string, amountCents int64, currency string, method string) (*PaymentResult, error)
	Refund(ctx context.Context, externalRef string, amountCents int64) (*PaymentResult, error)
}

// -----------------------------------------------------------------------------
// Noop processor (default when credentials missing).
// -----------------------------------------------------------------------------

type noopPaymentProcessor struct{}

func NewNoopPaymentProcessor() PaymentProcessor { return &noopPaymentProcessor{} }

func (p *noopPaymentProcessor) Deposit(ctx context.Context, userID string, amountCents int64, currency string) (*PaymentResult, error) {
	log.Debug().Str("user_id", userID).Int64("amount_cents", amountCents).Msg("noop: deposit")
	return &PaymentResult{Success: true, ExternalRef: "noop-deposit-" + userID}, nil
}

func (p *noopPaymentProcessor) Withdraw(ctx context.Context, userID string, amountCents int64, currency string, method string) (*PaymentResult, error) {
	log.Debug().Str("user_id", userID).Int64("amount_cents", amountCents).Str("method", method).Msg("noop: withdraw")
	return &PaymentResult{Success: true, ExternalRef: "noop-withdraw-" + userID}, nil
}

func (p *noopPaymentProcessor) Payout(ctx context.Context, userID string, amountCents int64, currency string, method string) (*PaymentResult, error) {
	log.Debug().Str("user_id", userID).Int64("amount_cents", amountCents).Str("method", method).Msg("noop: payout")
	return &PaymentResult{Success: true, ExternalRef: "noop-payout-" + userID}, nil
}

func (p *noopPaymentProcessor) Refund(ctx context.Context, externalRef string, amountCents int64) (*PaymentResult, error) {
	log.Debug().Str("external_ref", externalRef).Int64("amount_cents", amountCents).Msg("noop: refund")
	return &PaymentResult{Success: true, ExternalRef: "refund-" + externalRef}, nil
}

// -----------------------------------------------------------------------------
// Stripe Connect processor.
// https://docs.stripe.com/connect
// -----------------------------------------------------------------------------

type stripeConnectProcessor struct {
	secretKey string
	client    HTTPClient
	apiBase   string
}

func NewStripeConnectProcessor(secretKey string) PaymentProcessor {
	if strings.TrimSpace(secretKey) == "" {
		return NewNoopPaymentProcessor()
	}
	return &stripeConnectProcessor{
		secretKey: secretKey,
		client:    defaultClient(),
		apiBase:   "https://api.stripe.com/v1",
	}
}

func NewStripeConnectProcessorWithClient(secretKey string, client HTTPClient) PaymentProcessor {
	if strings.TrimSpace(secretKey) == "" {
		return NewNoopPaymentProcessor()
	}
	return &stripeConnectProcessor{secretKey: secretKey, client: client, apiBase: "https://api.stripe.com/v1"}
}

func (p *stripeConnectProcessor) do(ctx context.Context, method, path, form string) (map[string]interface{}, error) {
	var body io.Reader
	if form != "" {
		body = strings.NewReader(form)
	}
	req, err := http.NewRequestWithContext(ctx, method, p.apiBase+path, body)
	if err != nil {
		return nil, fmt.Errorf("stripe: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.secretKey)
	if form != "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("stripe: %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("stripe: %s %s status %d: %s", method, path, resp.StatusCode, string(respBody))
	}
	var out map[string]interface{}
	if len(respBody) == 0 {
		return out, nil
	}
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("stripe: decode %s %s: %w", method, path, err)
	}
	return out, nil
}

func (p *stripeConnectProcessor) Deposit(ctx context.Context, userID string, amountCents int64, currency string) (*PaymentResult, error) {
	form := fmt.Sprintf("amount=%d&currency=%s&customer=%s", amountCents, strings.ToLower(currency), userID)
	out, err := p.do(ctx, http.MethodPost, "/payment_intents", form)
	if err != nil {
		return &PaymentResult{Success: false, Error: err.Error()}, err
	}
	id, _ := out["id"].(string)
	return &PaymentResult{Success: true, ExternalRef: id}, nil
}

func (p *stripeConnectProcessor) Withdraw(ctx context.Context, userID string, amountCents int64, currency string, method string) (*PaymentResult, error) {
	form := fmt.Sprintf("amount=%d&currency=%s&destination=%s&source_type=%s", amountCents, strings.ToLower(currency), userID, method)
	out, err := p.do(ctx, http.MethodPost, "/transfers", form)
	if err != nil {
		return &PaymentResult{Success: false, Error: err.Error()}, err
	}
	id, _ := out["id"].(string)
	return &PaymentResult{Success: true, ExternalRef: id}, nil
}

func (p *stripeConnectProcessor) Payout(ctx context.Context, userID string, amountCents int64, currency string, method string) (*PaymentResult, error) {
	form := fmt.Sprintf("amount=%d&currency=%s&method=%s&metadata[user_id]=%s", amountCents, strings.ToLower(currency), strings.ToLower(method), userID)
	out, err := p.do(ctx, http.MethodPost, "/payouts", form)
	if err != nil {
		return &PaymentResult{Success: false, Error: err.Error()}, err
	}
	id, _ := out["id"].(string)
	return &PaymentResult{Success: true, ExternalRef: id}, nil
}

func (p *stripeConnectProcessor) Refund(ctx context.Context, externalRef string, amountCents int64) (*PaymentResult, error) {
	form := fmt.Sprintf("payment_intent=%s&amount=%d", externalRef, amountCents)
	out, err := p.do(ctx, http.MethodPost, "/refunds", form)
	if err != nil {
		return &PaymentResult{Success: false, Error: err.Error()}, err
	}
	id, _ := out["id"].(string)
	return &PaymentResult{Success: true, ExternalRef: id}, nil
}

// -----------------------------------------------------------------------------
// PayPal Payouts processor.
// https://developer.paypal.com/api/rest/reference/payouts/
// -----------------------------------------------------------------------------

type paypalPayoutProcessor struct {
	clientID     string
	clientSecret string
	client       HTTPClient
	apiBase      string
}

func NewPayPalPayoutProcessor(clientID, clientSecret string) PaymentProcessor {
	if strings.TrimSpace(clientID) == "" || strings.TrimSpace(clientSecret) == "" {
		return NewNoopPaymentProcessor()
	}
	return &paypalPayoutProcessor{
		clientID:     clientID,
		clientSecret: clientSecret,
		client:       defaultClient(),
		apiBase:      "https://api-m.paypal.com",
	}
}

func NewPayPalPayoutProcessorWithClient(clientID, clientSecret string, client HTTPClient) PaymentProcessor {
	if strings.TrimSpace(clientID) == "" || strings.TrimSpace(clientSecret) == "" {
		return NewNoopPaymentProcessor()
	}
	return &paypalPayoutProcessor{clientID: clientID, clientSecret: clientSecret, client: client, apiBase: "https://api-m.paypal.com"}
}

func (p *paypalPayoutProcessor) getAccessToken(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.apiBase+"/v1/oauth2/token", strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return "", fmt.Errorf("paypal: build token request: %w", err)
	}
	req.SetBasicAuth(p.clientID, p.clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("paypal: get token: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("paypal: token status %d: %s", resp.StatusCode, string(body))
	}
	var tr struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", fmt.Errorf("paypal: decode token: %w", err)
	}
	return tr.AccessToken, nil
}

func (p *paypalPayoutProcessor) do(ctx context.Context, method, path string, body interface{}) (map[string]interface{}, error) {
	token, err := p.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	var bodyReader io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("paypal: marshal: %w", err)
		}
		bodyReader = bytes.NewReader(buf)
	}
	req, err := http.NewRequestWithContext(ctx, method, p.apiBase+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("paypal: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("paypal: %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("paypal: %s %s status %d: %s", method, path, resp.StatusCode, string(respBody))
	}
	out := map[string]interface{}{}
	if len(respBody) == 0 {
		return out, nil
	}
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("paypal: decode %s %s: %w", method, path, err)
	}
	return out, nil
}

func (p *paypalPayoutProcessor) Deposit(ctx context.Context, userID string, amountCents int64, currency string) (*PaymentResult, error) {
	return &PaymentResult{Success: false, Error: "PayPal does not support wallet deposits; use a card processor"}, nil
}

func (p *paypalPayoutProcessor) Withdraw(ctx context.Context, userID string, amountCents int64, currency string, method string) (*PaymentResult, error) {
	return &PaymentResult{Success: false, Error: "PayPal wallet withdrawals are not supported"}, nil
}

func (p *paypalPayoutProcessor) Payout(ctx context.Context, userID string, amountCents int64, currency string, method string) (*PaymentResult, error) {
	body := map[string]interface{}{
		"sender_batch_header": map[string]interface{}{
			"sender_batch_id": userID,
		},
		"items": []map[string]interface{}{{
			"recipient_type": "PAYPAL_ID",
			"receiver":       userID,
			"amount": map[string]interface{}{
				"value":        fmt.Sprintf("%.2f", float64(amountCents)/100.0),
				"currency":     strings.ToUpper(currency),
			},
			"sender_item_id": userID + "-item",
		}},
	}
	out, err := p.do(ctx, http.MethodPost, "/v1/payments/payouts", body)
	if err != nil {
		return &PaymentResult{Success: false, Error: err.Error()}, err
	}
	batchHeader, _ := out["batch_header"].(map[string]interface{})
	id, _ := batchHeader["payout_batch_id"].(string)
	return &PaymentResult{Success: true, ExternalRef: id}, nil
}

func (p *paypalPayoutProcessor) Refund(ctx context.Context, externalRef string, amountCents int64) (*PaymentResult, error) {
	return &PaymentResult{Success: false, Error: "PayPal refunds are handled through the payment service"}, nil
}
