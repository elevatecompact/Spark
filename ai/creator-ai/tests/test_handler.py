import pytest
from httpx import ASGITransport, AsyncClient

from src.main import app
from src.service import CreatorService


@pytest.fixture(autouse=True)
def setup_service():
    app.state.service = CreatorService()
    yield
    app.state.service = None


@pytest.mark.asyncio
async def test_post_ideas_returns_200():
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        response = await client.post(
            "/v1/creator/ideas",
            json={"niche": "gaming", "count": 2},
        )
    assert response.status_code == 200
    data = response.json()
    assert data["count"] == 2
    assert len(data["ideas"]) == 2


@pytest.mark.asyncio
async def test_get_health_returns_200():
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        response = await client.get("/health")
    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "ok"
    assert data["service"] == "creator-ai"
