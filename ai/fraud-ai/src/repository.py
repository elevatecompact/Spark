import logging

from src.models import FraudRule

logger = logging.getLogger(__name__)


class FraudRepository:

    async def get_rules(self) -> list[FraudRule]:
        return [
            FraudRule(
                rule_id="FR001",
                name="high_amount",
                description="Transaction amount exceeds threshold",
                severity="high",
                category="amount",
            ),
            FraudRule(
                rule_id="FR002",
                name="new_account",
                description="Account age less than 7 days",
                severity="medium",
                category="account_age",
            ),
            FraudRule(
                rule_id="FR003",
                name="high_velocity",
                description="Multiple transactions in short period",
                severity="high",
                category="velocity",
            ),
            FraudRule(
                rule_id="FR004",
                name="unusual_location",
                description="Transaction from unusual geographic location",
                severity="medium",
                category="location",
            ),
            FraudRule(
                rule_id="FR005",
                name="device_mismatch",
                description="Device fingerprint does not match known devices",
                severity="low",
                category="device",
            ),
        ]

    async def log_fraud_check(self, transaction_id: str, score: float, decision: str):
        logger.info("Fraud check logged: %s score=%.3f decision=%s", transaction_id, score, decision)
