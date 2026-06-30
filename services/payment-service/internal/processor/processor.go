package processor

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

	"github.com/elevatecompact/spark/services/payment-service/internal/domain"
)

// HTTPClient is the small interface used by external payment processors so
// tests can plug in a fake.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

func defaultClient() HTTPClient { return &http.Client{Timeout: 30 * time.Second} }

// PaymentProcessor is the contract every gateway implementation satisfies.
type PaymentProcessor interface {
	Name() domain.PaymentProcessor
	CreateIntent(ctx context.Context, intent *domain.PaymentIntent) error
	ConfirmIntent(ctx context.Context, intent *domain.PaymentIntent, paymentMethodID string) error
	CancelIntent(ctx context.Context, intent *domain.PaymentIntent) error
	Refund(ctx context.Context, intent *domain.PaymentIntent, amountCents *int64) (string, error)
	CreatePayout(ctx context.Context, payout *domain.Payout) error
}

// -----------------------------------------------------------------------------
// Noop implementation (used when credentials are missing or in dev/test).
// -----------------------------------------------------------------------------

type noopProcessor struct {
	name domain.PaymentProcessor
}

func NewNoopProcessor(name domain.PaymentProcessor) PaymentProcessor {
	return &noopProcessor{name: name}
}

func (p *noopProcessor) Name() domain.PaymentProcessor { return p.name }

func (p *noopProcessor) CreateIntent(ctx context.Context, intent *domain.PaymentIntent) error {
	intent.ExternalID = string(p.name) + "_pi_noop_" + intent.ID.String()
	if intent.Status == "" {
		intent.Status = domain.IntentRequiresPaymentMethod
	}
	return nil
}

func (p *noopProcessor) ConfirmIntent(ctx context.Context, intent *domain.PaymentIntent, paymentMethodID string) error {
	intent.Status = domain.IntentSucceeded
	return nil
}

func (p *noopProcessor) CancelIntent(ctx context.Context, intent *domain.PaymentIntent) error {
	intent.Status = domain.IntentCanceled
	return nil
}

func (p *noopProcessor) Refund(ctx context.Context, intent *domain.PaymentIntent, amountCents *int64) (string, error) {
	return string(p.name) + "_rf_noop_" + intent.ID.String(), nil
}

func (p *noopProcessor) CreatePayout(ctx context.Context, payout *domain.Payout) error {
	payout.ExternalID = string(p.name) + "_po_noop_" + payout.ID.String()
	payout.Status = "completed"
	return nil
}

// -----------------------------------------------------------------------------
// Stripe processor.
// https://docs.stripe.com/api
// -----------------------------------------------------------------------------

type stripeProcessor struct {
	secretKey      string
	webhookSecret  string
	publishableKey string
	client         HTTPClient
	apiBase        string
}

func NewStripeProcessor(secretKey, webhookSecret, publishableKey string) PaymentProcessor {
	if strings.TrimSpace(secretKey) == "" || strings.HasPrefix(secretKey, "sk_test_noop") {
		return NewNoopProcessor(domain.ProcessorStripe)
	}
	return &stripeProcessor{
		secretKey:      secretKey,
		webhookSecret:  webhookSecret,
		publishableKey: publishableKey,
		client:         defaultClient(),
		apiBase:        "https://api.stripe.com/v1",
	}
}

func NewStripeProcessorWithClient(secretKey, webhookSecret, publishableKey string, client HTTPClient) PaymentProcessor {
	if strings.TrimSpace(secretKey) == "" {
		return NewNoopProcessor(domain.ProcessorStripe)
	}
	return &stripeProcessor{secretKey: secretKey, webhookSecret: webhookSecret, publishableKey: publishableKey, client: client, apiBase: "https://api.stripe.com/v1"}
}

func (p *stripeProcessor) Name() domain.PaymentProcessor { return domain.ProcessorStripe }

type stripeIntent struct {
	ID            string `json:"id"`
	ClientSecret  string `json:"client_secret"`
	Status        string `json:"status"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	LatestCharge  string `json:"latest_charge"`
}

type stripeError struct {
	Err struct {
		Type    string `json:"type"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (p *stripeProcessor) CreateIntent(ctx context.Context, intent *domain.PaymentIntent) error {
	form := fmt.Sprintf("amount=%d&currency=%s&metadata[user_id]=%s&metadata[intent_id]=%s",
		intent.AmountCents, strings.ToLower(intent.Currency), intent.UserID, intent.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.apiBase+"/payment_intents", strings.NewReader(form))
	if err != nil {
		return fmt.Errorf("stripe: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.secretKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("stripe: create intent: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("stripe: status %d: %s", resp.StatusCode, string(body))
	}
	var si stripeIntent
	if err := json.Unmarshal(body, &si); err != nil {
		return fmt.Errorf("stripe: decode intent: %w", err)
	}
	intent.ExternalID = si.ID
	intent.Status = mapStripeStatus(si.Status)
	return nil
}

func (p *stripeProcessor) ConfirmIntent(ctx context.Context, intent *domain.PaymentIntent, paymentMethodID string) error {
	form := "payment_method=" + paymentMethodID
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.apiBase+"/payment_intents/"+intent.ExternalID+"/confirm", strings.NewReader(form))
	if err != nil {
		return fmt.Errorf("stripe: build confirm request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.secretKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("stripe: confirm intent: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("stripe: confirm status %d: %s", resp.StatusCode, string(body))
	}
	var si stripeIntent
	if err := json.Unmarshal(body, &si); err != nil {
		return fmt.Errorf("stripe: decode confirm: %w", err)
	}
	intent.Status = mapStripeStatus(si.Status)
	return nil
}

func (p *stripeProcessor) CancelIntent(ctx context.Context, intent *domain.PaymentIntent) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.apiBase+"/payment_intents/"+intent.ExternalID+"/cancel", nil)
	if err != nil {
		return fmt.Errorf("stripe: build cancel request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.secretKey)
	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("stripe: cancel intent: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("stripe: cancel status %d: %s", resp.StatusCode, string(body))
	}
	var si stripeIntent
	if err := json.Unmarshal(body, &si); err != nil {
		return fmt.Errorf("stripe: decode cancel: %w", err)
	}
	intent.Status = mapStripeStatus(si.Status)
	return nil
}

func (p *stripeProcessor) Refund(ctx context.Context, intent *domain.PaymentIntent, amountCents *int64) (string, error) {
	form := "payment_intent=" + intent.ExternalID
	if amountCents != nil {
		form += fmt.Sprintf("&amount=%d", *amountCents)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.apiBase+"/refunds", strings.NewReader(form))
	if err != nil {
		return "", fmt.Errorf("stripe: build refund: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.secretKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("stripe: refund: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("stripe: refund status %d: %s", resp.StatusCode, string(body))
	}
	var refund struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &refund); err != nil {
		return "", fmt.Errorf("stripe: decode refund: %w", err)
	}
	log.Debug().Str("refund_id", refund.ID).Msg("stripe refund created")
	return refund.ID, nil
}

func (p *stripeProcessor) CreatePayout(ctx context.Context, payout *domain.Payout) error {
	form := fmt.Sprintf("amount=%d&currency=%s&metadata[user_id]=%s&metadata[payout_id]=%s",
		payout.AmountCents, strings.ToLower(payout.Currency), payout.UserID, payout.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.apiBase+"/payouts", strings.NewReader(form))
	if err != nil {
		return fmt.Errorf("stripe: build payout: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.secretKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("stripe: payout: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("stripe: payout status %d: %s", resp.StatusCode, string(body))
	}
	var pr struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(body, &pr); err != nil {
		return fmt.Errorf("stripe: decode payout: %w", err)
	}
	payout.ExternalID = pr.ID
	payout.Status = pr.Status
	return nil
}

func mapStripeStatus(s string) domain.IntentStatus {
	switch s {
	case "succeeded":
		return domain.IntentSucceeded
	case "requires_payment_method", "requires_confirmation":
		return domain.IntentRequiresPaymentMethod
	case "processing":
		return domain.IntentProcessing
	case "canceled":
		return domain.IntentCanceled
	case "requires_action":
		return domain.IntentProcessing
	default:
		return domain.IntentFailed
	}
}

// -----------------------------------------------------------------------------
// PayPal processor.
// https://developer.paypal.com/api/rest/
// -----------------------------------------------------------------------------

type paypalProcessor struct {
	clientID     string
	clientSecret string
	webhookID    string
	client       HTTPClient
	apiBase      string
}

func NewPayPalProcessor(clientID, clientSecret, webhookID string) PaymentProcessor {
	if strings.TrimSpace(clientID) == "" || strings.TrimSpace(clientSecret) == "" {
		return NewNoopProcessor(domain.ProcessorPayPal)
	}
	return &paypalProcessor{
		clientID:     clientID,
		clientSecret: clientSecret,
		webhookID:    webhookID,
		client:       defaultClient(),
		apiBase:      "https://api-m.paypal.com",
	}
}

func NewPayPalProcessorWithClient(clientID, clientSecret, webhookID string, client HTTPClient) PaymentProcessor {
	if strings.TrimSpace(clientID) == "" || strings.TrimSpace(clientSecret) == "" {
		return NewNoopProcessor(domain.ProcessorPayPal)
	}
	return &paypalProcessor{clientID: clientID, clientSecret: clientSecret, webhookID: webhookID, client: client, apiBase: "https://api-m.paypal.com"}
}

func (p *paypalProcessor) Name() domain.PaymentProcessor { return domain.ProcessorPayPal }

type paypalTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func (p *paypalProcessor) getAccessToken(ctx context.Context) (string, error) {
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
	var tr paypalTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", fmt.Errorf("paypal: decode token: %w", err)
	}
	return tr.AccessToken, nil
}

func (p *paypalProcessor) do(ctx context.Context, method, path string, body interface{}, out interface{}) error {
	token, err := p.getAccessToken(ctx)
	if err != nil {
		return err
	}
	var bodyReader io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("paypal: marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(buf)
	}
	req, err := http.NewRequestWithContext(ctx, method, p.apiBase+path, bodyReader)
	if err != nil {
		return fmt.Errorf("paypal: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("paypal: %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("paypal: %s %s status %d: %s", method, path, resp.StatusCode, string(respBody))
	}
	if out != nil {
		if err := json.Unmarshal(respBody, out); err != nil {
			return fmt.Errorf("paypal: decode %s %s: %w", method, path, err)
		}
	}
	return nil
}

type paypalOrder struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type paypalAmount struct {
	CurrencyCode string `json:"currency_code"`
	Value        string `json:"value"`
}

type paypalPurchaseUnit struct {
	ReferenceID string        `json:"reference_id"`
	Amount      paypalAmount  `json:"amount"`
}

type paypalCreateOrderRequest struct {
	Intent        string              `json:"intent"`
	PurchaseUnits []paypalPurchaseUnit `json:"purchase_units"`
}

type paypalRefund struct {
	ID string `json:"id"`
}

type paypalRefundRequest struct {
	Amount paypalAmount `json:"amount"`
}

type paypalPayout struct {
	BatchHeader struct {
		PayoutBatchID string `json:"payout_batch_id"`
		Status        string `json:"batch_status"`
	} `json:"batch_header"`
}

func (p *paypalProcessor) CreateIntent(ctx context.Context, intent *domain.PaymentIntent) error {
	body := paypalCreateOrderRequest{
		Intent: "CAPTURE",
		PurchaseUnits: []paypalPurchaseUnit{{
			ReferenceID: intent.ID.String(),
			Amount: paypalAmount{
				CurrencyCode: strings.ToUpper(intent.Currency),
				Value:        formatCents(intent.AmountCents),
			},
		}},
	}
	var order paypalOrder
	if err := p.do(ctx, http.MethodPost, "/v2/checkout/orders", body, &order); err != nil {
		return err
	}
	intent.ExternalID = order.ID
	intent.Status = mapPayPalStatus(order.Status)
	return nil
}

func (p *paypalProcessor) ConfirmIntent(ctx context.Context, intent *domain.PaymentIntent, paymentMethodID string) error {
	var order paypalOrder
	if err := p.do(ctx, http.MethodPost, "/v2/checkout/orders/"+intent.ExternalID+"/capture", nil, &order); err != nil {
		return err
	}
	intent.Status = mapPayPalStatus(order.Status)
	return nil
}

func (p *paypalProcessor) CancelIntent(ctx context.Context, intent *domain.PaymentIntent) error {
	intent.Status = domain.IntentCanceled
	return nil
}

func (p *paypalProcessor) Refund(ctx context.Context, intent *domain.PaymentIntent, amountCents *int64) (string, error) {
	amount := intent.AmountCents
	if amountCents != nil {
		amount = *amountCents
	}
	body := paypalRefundRequest{Amount: paypalAmount{
		CurrencyCode: strings.ToUpper(intent.Currency),
		Value:        formatCents(amount),
	}}
	var refund paypalRefund
	if err := p.do(ctx, http.MethodPost, "/v2/payments/captures/"+intent.ExternalID+"/refund", body, &refund); err != nil {
		return "", err
	}
	return refund.ID, nil
}

func (p *paypalProcessor) CreatePayout(ctx context.Context, payout *domain.Payout) error {
	body := map[string]interface{}{
		"sender_batch_header": map[string]interface{}{
			"sender_batch_id": payout.ID.String(),
		},
		"items": []map[string]interface{}{{
			"recipient_type": "EMAIL",
			"amount": map[string]interface{}{
				"value":        formatCents(payout.AmountCents),
				"currency":     strings.ToUpper(payout.Currency),
			},
			"sender_item_id": payout.ID.String(),
		}},
	}
	var pr paypalPayout
	if err := p.do(ctx, http.MethodPost, "/v1/payments/payouts", body, &pr); err != nil {
		return err
	}
	payout.ExternalID = pr.BatchHeader.PayoutBatchID
	payout.Status = pr.BatchHeader.Status
	return nil
}

func mapPayPalStatus(s string) domain.IntentStatus {
	switch strings.ToUpper(s) {
	case "COMPLETED":
		return domain.IntentSucceeded
	case "CREATED", "SAVED":
		return domain.IntentRequiresPaymentMethod
	case "APPROVED", "PAYER_ACTION_REQUIRED":
		return domain.IntentProcessing
	case "VOIDED":
		return domain.IntentCanceled
	default:
		return domain.IntentFailed
	}
}

func formatCents(cents int64) string {
	whole := cents / 100
	frac := cents % 100
	if frac < 0 {
		frac = -frac
	}
	return fmt.Sprintf("%d.%02d", whole, frac)
}
