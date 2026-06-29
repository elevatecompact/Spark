import pytest
from httpx import ASGITransport, AsyncClient

from src.main import create_app


@pytest.fixture
def app():
    return create_app()


@pytest.mark.asyncio
async def test_post_text_returns_200(app):
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        payload = {
            "text": "Hello world",
            "content_id": "c1",
        }
        response = await client.post("/v1/moderation/text", json=payload)

    assert response.status_code == 200
    data = response.json()
    assert data["overall_safe"] is True
    assert data["text_result"]["label"] == "clean"
    assert data["model_version"] == "noop-v1"


@pytest.mark.asyncio
async def test_post_image_returns_200(app):
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        payload = {
            "image_url": "https://example.com/img.jpg",
            "content_id": "c2",
        }
        response = await client.post("/v1/moderation/image", json=payload)

    assert response.status_code == 200
    data = response.json()
    assert data["overall_safe"] is True
    assert data["image_result"]["label"] == "clean"


@pytest.mark.asyncio
async def test_post_batch_returns_200(app):
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        payload = {
            "requests": [
                {"text": "Hello", "content_id": "c1"},
                {"image_url": "https://example.com/img.jpg", "content_id": "c2"},
            ]
        }
        response = await client.post("/v1/moderation/batch", json=payload)

    assert response.status_code == 200
    data = response.json()
    assert len(data["results"]) == 2


@pytest.mark.asyncio
async def test_get_categories_returns_200(app):
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        response = await client.get("/v1/moderation/categories")

    assert response.status_code == 200
    data = response.json()
    assert "categories" in data
    assert len(data["categories"]) == 7


@pytest.mark.asyncio
async def test_health_returns_200(app):
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        response = await client.get("/health")

    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "ok"
    assert data["service"] == "moderation-ai"
    assert data["version"] == "0.1.0"
