# payment-service — Database Schema
## payment_intents: id UUID PK, external_id VARCHAR(255), processor, amount_cents, currency, status(requires_payment_method,processing,succeeded,failed,canceled), idempotency_key UNIQUE, metadata JSONB, user_id FK
## payment_methods: id UUID PK, user_id FK, external_id (processor token), processor, type(card,paypal,bank), fingerprint (dedup), last4, exp_month/year, is_default
## webhook_events: id UUID PK, processor, external_event_id UNIQUE, type, body JSONB, status(received,processed,failed). Only reference data stored — no PCI.
