import pytest
from httpx import ASGITransport, AsyncClient

from src.main import create_app


@pytest.fixture
def app():
    return create_app()


@pytest.mark.asyncio
async def test_post_stt_returns_200(app):
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        payload = {
            "audio_url": "https://cdn.spark.dev/audio/sample.wav",
            "language": "en",
            "format": "wav",
        }
        response = await client.post("/v1/voice/stt", json=payload)

    assert response.status_code == 200
    data = response.json()
    assert "text" in data
    assert "confidence" in data
    assert "model_version" in data
    assert data["model_version"] == "noop-v1"


@pytest.mark.asyncio
async def test_health_returns_200(app):
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        response = await client.get("/health")

    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "ok"
    assert data["service"] == "voice-ai"
    assert data["version"] == "0.1.0"
