import pytest
from httpx import ASGITransport, AsyncClient

from src.events import NoopProducer
from src.main import app
from src.repository import RecommendationRepository
from src.service import RecommendationService

pytestmark = pytest.mark.asyncio


@pytest.fixture(autouse=True)
async def setup_app() -> None:
    repository = RecommendationRepository()
    producer = NoopProducer()
    await producer.start()
    service = RecommendationService(repository=repository, producer=producer)
    app.state.service = service
    yield
    del app.state.service


@pytest.fixture
async def client() -> AsyncClient:
    async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test") as ac:
        yield ac


async def test_get_recommendations_returns_200(client: AsyncClient) -> None:
    payload = {"user_id": "test-user", "count": 3}
    response = await client.post("/v1/recommendation/get", json=payload)
    assert response.status_code == 200
    data = response.json()
    assert len(data["items"]) == 3
    assert data["total"] == 9
    assert data["context"] == "home"
    assert "model_version" in data


async def test_health_returns_200(client: AsyncClient) -> None:
    response = await client.get("/health")
    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "ok"
    assert data["service"] == "recommendation-ai"
    assert data["version"] == "0.1.0"


async def test_get_similar_returns_200(client: AsyncClient) -> None:
    payload = {"item_id": "item-1", "count": 5}
    response = await client.post("/v1/recommendation/similar", json=payload)
    assert response.status_code == 200
    data = response.json()
    assert len(data["items"]) == 5
    assert data["source_item_id"] == "item-1"


async def test_get_trending_returns_200(client: AsyncClient) -> None:
    payload = {"timeframe": "week", "count": 10}
    response = await client.post("/v1/recommendation/trending", json=payload)
    assert response.status_code == 200
    data = response.json()
    assert len(data["items"]) == 10
    assert data["timeframe"] == "week"


async def test_personalize_feed_returns_200(client: AsyncClient) -> None:
    payload = {"user_id": "test-user", "feed_type": "home", "count": 5}
    response = await client.post("/v1/recommendation/feed", json=payload)
    assert response.status_code == 200
    data = response.json()
    assert len(data["items"]) == 5
    assert "personalization_score" in data


async def test_track_interaction_returns_200(client: AsyncClient) -> None:
    payload = {"user_id": "test-user", "item_id": "item-1", "interaction_type": "view"}
    response = await client.post("/v1/recommendation/interaction", json=payload)
    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "success"
    assert data["recorded"] is True


async def test_batch_track_returns_200(client: AsyncClient) -> None:
    payload = {
        "interactions": [
            {"user_id": "test-user", "item_id": "item-1", "interaction_type": "view"},
            {"user_id": "test-user", "item_id": "item-2", "interaction_type": "like"},
        ]
    }
    response = await client.post("/v1/recommendation/interactions/batch", json=payload)
    assert response.status_code == 200
    data = response.json()
    assert data["recorded_count"] == 2
