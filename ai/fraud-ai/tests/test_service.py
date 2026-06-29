import pytest

from src.models import FraudScoreRequest, TransactionData
from src.repository import FraudRepository
from src.events import NoopProducer
from src.service import FraudService


@pytest.fixture
def service():
    repo = FraudRepository()
    event_pub = NoopProducer()
    return FraudService(repo=repo, event_pub=event_pub)


@pytest.mark.asyncio
async def test_analyze_returns_valid_score(service):
    request = FraudScoreRequest(
        transaction=TransactionData(
            amount=150.00,
            currency="USD",
            payment_method="credit_card",
        ),
    )

    response = await service.analyze(request)

    assert 0.0 <= response.score <= 1.0
    assert response.risk_level in ("low", "medium", "high", "critical")
    assert response.decision in ("allow", "review", "block")
    assert isinstance(response.flags, list)
    assert response.model_version == "noop-v1"


@pytest.mark.asyncio
async def test_get_rules_returns_list(service):
    rules = await service.get_rules()
    assert isinstance(rules, list)
    assert len(rules) > 0
    assert rules[0].rule_id is not None
