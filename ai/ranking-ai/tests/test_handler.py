import pytest
from httpx import ASGITransport, AsyncClient

from src.main import create_app


@pytest.fixture
def app():
    return create_app()


@pytest.mark.asyncio
async def test_post_rank_returns_200(app):
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        payload = {
            "items": [
                {"item_id": "a", "features": {"f1": 1.0, "f2": 2.0}},
                {"item_id": "b", "features": {"f1": 5.0}},
            ],
            "context": {"user_id": "u1"},
        }
        response = await client.post("/v1/ranking/rank", json=payload)

    assert response.status_code == 200
    data = response.json()
    assert "ranked_items" in data
    assert len(data["ranked_items"]) == 2
    assert data["ranked_items"][0]["item_id"] == "b"
    assert data["ranked_items"][1]["item_id"] == "a"
    assert data["model_version"] == "noop-v1"


@pytest.mark.asyncio
async def test_health_returns_200(app):
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        response = await client.get("/health")

    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "ok"
    assert data["service"] == "ranking-ai"
    assert data["version"] == "0.1.0"
