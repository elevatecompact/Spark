# subscription-service — Database Schema
## subscription_plans: id UUID PK, creator_id FK (nullable=platform plans), name, price_cents, currency, billing_period(monthly,yearly), benefits JSONB, is_active
## subscriptions: id UUID PK, user_id FK, plan_id FK, status(active,cancelled,expired,grace_period), current_period_start/end, cancelled_at, grace_period_end
## invoices: id UUID PK, subscription_id FK, amount_cents, currency, status(pending,paid,failed,refunded), paid_at, period_start/end
Indexes on user_id, plan_id, status, next_billing_date.
