# wallet-service — Database Schema
## wallets: id UUID PK, user_id UUID UNIQUE FK, balance_cents BIGINT CHECK>=0, currency VARCHAR, status(active,frozen,closed), version INT (optimistic lock)
## 	ransactions (append-only ledger): id UUID PK, idempotency_key VARCHAR UNIQUE, from_wallet_id FK, to_wallet_id FK, amount_cents, currency, type, status(pending,settled,failed), failure_reason, created_at, settled_at. SERIALIZABLE isolation.
## payouts: id UUID PK, wallet_id FK, amount_cents, method(paypal,bank,crypto), status(requested,processing,completed,failed), external_ref
