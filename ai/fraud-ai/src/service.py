import hashlib
import logging
import random

from src.events import EventProducer
from src.models import (
    FraudAnalysisResponse,
    FraudRule,
    FraudScoreRequest,
    FraudScoreResponse,
    TransactionData,
)
from src.repository import FraudRepository

logger = logging.getLogger(__name__)

FLAGS = ["large_amount", "new_account", "high_velocity", "unusual_location", "device_mismatch"]


class FraudService:

    def __init__(self, repo: FraudRepository, event_pub: EventProducer):
        self._repo = repo
        self._event_pub = event_pub

    async def analyze(self, request: FraudScoreRequest) -> FraudScoreResponse:
        raw = f"{request.transaction.amount}:{request.transaction.currency}:{request.transaction.payment_method}"
        digest = hashlib.md5(raw.encode()).hexdigest()
        score = int(digest[:8], 16) / 0xFFFFFFFF

        risk_level = self._risk_level(score)
        decision = self._decision(score)
        flags = random.sample(FLAGS, random.randint(0, 3))

        response = FraudScoreResponse(
            score=round(score, 4),
            risk_level=risk_level,
            flags=flags,
            decision=decision,
            model_version="noop-v1",
        )

        await self._event_pub.publish(
            "fraud.analysis.completed",
            {
                "transaction": request.transaction.model_dump(),
                "score": score,
                "decision": decision,
                "flags": flags,
            },
        )

        return response

    async def batch_analyze(self, requests: list[FraudScoreRequest]) -> list[FraudScoreResponse]:
        results = []
        for req in requests:
            result = await self.analyze(req)
            results.append(result)
        return results

    async def get_rules(self) -> list[FraudRule]:
        return await self._repo.get_rules()

    async def analyze_rules(self, transaction: TransactionData) -> FraudAnalysisResponse:
        rules = await self._repo.get_rules()
        triggered = []
        triggered.extend(
            rule for rule in rules
            if rule.rule_id == "FR001" and transaction.amount > 10000
        )
        triggered.extend(
            rule for rule in rules
            if rule.rule_id == "FR005" and transaction.device_fingerprint is None
        )
        triggered.extend(
            rule for rule in rules
            if rule.rule_id == "FR004" and transaction.ip_address is None
        )
        if len(triggered) > 0:
            triggered = triggered[:2]

        overall_score = min(1.0, len(triggered) * 0.25)

        return FraudAnalysisResponse(
            rules_triggered=triggered,
            overall_score=overall_score,
            risk_level=self._risk_level(overall_score),
            decision=self._decision(overall_score),
        )

    @staticmethod
    def _risk_level(score: float) -> str:
        if score < 0.3:
            return "low"
        if score < 0.6:
            return "medium"
        if score < 0.8:
            return "high"
        return "critical"

    @staticmethod
    def _decision(score: float) -> str:
        if score < 0.6:
            return "allow"
        if score < 0.8:
            return "review"
        return "block"
