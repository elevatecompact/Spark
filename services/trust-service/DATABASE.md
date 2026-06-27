# trust-service — Database Schema
## PostgreSQL — Reputation, signals, fraud cases
### eputation_scores: user_id UUID PK, overall_score INT(0-1000), trust_level ENUM(low,medium,high,verified), positive_signal_weight INT, negative_signal_weight INT, score_decay_factor FLOAT(0-1), model_version VARCHAR, calculated_at, next_recalculation_at
### 	rust_signals: id UUID PK, user_id FK, signal_type(positive,negative), category(identity,payment,content,community,account_age), weight INT(-100 to +100), description, source_entity_type, source_entity_id, expires_at nullable, recorded_at. Append-only, retained 2 years.
### isk_assessments: id UUID PK, user_id FK, action_type, context JSONB, risk_score FLOAT(0-1), risk_level(low,medium,high,critical), triggered_rules TEXT[], recommended_action(allow,challenge,block,review), assessed_at
### raud_cases: id UUID PK, user_id FK, case_type(payment_fraud,account_takeover,coordinated_behavior,spam), status(open,investigating,confirmed,false_positive,resolved), evidence JSONB, automated_decision, reviewed_by UUID nullable, resolved_at, created_at
### isk_rules: id UUID PK, name, category, conditions JSONB, risk_score_impact, action, is_active, priority, created_at
## Redis — Real-time risk counters (failed login attempts, rapid actions), fraud signal cache (IP reputation, device fingerprint), rate limit state for risk checks
