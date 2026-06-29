import pytest
from httpx import ASGITransport, AsyncClient

from src.main import create_app


@pytest.fixture
def app():
    return create_app()


@pytest.mark.asyncio
async def test_post_translate_returns_200(app):
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        payload = {
            "text": "hello",
            "source_language": "en",
            "target_language": "fr",
        }
        response = await client.post("/v1/translation/translate", json=payload)

    assert response.status_code == 200
    data = response.json()
    assert data["translated_text"] == "olleh"
    assert data["source_language"] == "en"
    assert data["target_language"] == "fr"
    assert data["confidence"] == 0.95
    assert data["provider"] == "noop"
    assert data["model_version"] == "noop-v1"


@pytest.mark.asyncio
async def test_post_batch_returns_200(app):
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        payload = {
            "requests": [
                {"text": "hello", "source_language": "en", "target_language": "fr"},
                {"text": "world", "source_language": "en", "target_language": "es"},
            ]
        }
        response = await client.post("/v1/translation/batch", json=payload)

    assert response.status_code == 200
    data = response.json()
    assert len(data["results"]) == 2


@pytest.mark.asyncio
async def test_post_detect_returns_200(app):
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        payload = {"text": "hello world"}
        response = await client.post("/v1/translation/detect", json=payload)

    assert response.status_code == 200
    data = response.json()
    assert data["detected_language"] == "en"
    assert data["confidence"] == 0.85


@pytest.mark.asyncio
async def test_get_languages_returns_200(app):
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        response = await client.get("/v1/translation/languages")

    assert response.status_code == 200
    data = response.json()
    assert len(data["languages"]) == 15
    assert data["languages"][0]["code"] == "en"


@pytest.mark.asyncio
async def test_health_returns_200(app):
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        response = await client.get("/health")

    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "ok"
    assert data["service"] == "translation-ai"
    assert data["version"] == "0.1.0"
