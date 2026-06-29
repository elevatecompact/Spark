import pytest
from httpx import AsyncClient, ASGITransport

from src.main import app
from src.models import DetectClipsRequest


@pytest.mark.asyncio
async def test_post_detect_returns_200() -> None:
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        payload = DetectClipsRequest(video_url="https://example.com/video.mp4").model_dump()
        response = await client.post("/v1/clip/detect", json=payload)
    assert response.status_code == 200
    data = response.json()
    assert "clips" in data
    assert data["clip_count"] > 0


@pytest.mark.asyncio
async def test_get_health_returns_200() -> None:
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        response = await client.get("/health")
    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "ok"
    assert data["service"] == "clip-ai"
