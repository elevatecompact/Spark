import pytest

from src.models import GetRecommendationsRequest
from src.repository import RecommendationRepository
from src.service import RecommendationService

pytestmark = pytest.mark.asyncio


async def test_get_recommendations_returns_items() -> None:
    repo = RecommendationRepository()
    service = RecommendationService(repository=repo)
    request = GetRecommendationsRequest(user_id="test-user", count=5)
    response = await service.get_recommendations(request)
    assert len(response.items) == 5
    assert response.total == 15
    assert response.context == "home"
    assert response.model_version == "noop-v1"
    for item in response.items:
        assert item.item_id.startswith("rec-")
        assert 0 <= item.score <= 1
        assert item.reason is not None


async def test_get_recommendations_sorted_by_score() -> None:
    service = RecommendationService()
    request = GetRecommendationsRequest(user_id="test-user", count=20)
    response = await service.get_recommendations(request)
    scores = [item.score for item in response.items]
    assert scores == sorted(scores, reverse=True)
    for i, item in enumerate(response.items):
        assert item.position == i
