import pytest
from httpx import ASGITransport, AsyncClient

from src.events import NoopProducer
from src.handler import _init as init_handler
from src.handler import router
from src.repository import VisionRepository
from src.service import VisionService


@pytest.fixture
def app():
    from fastapi import FastAPI

    app = FastAPI()
    repository = VisionRepository()
    producer = NoopProducer()
    service = VisionService(repository=repository, event_producer=producer)
    init_handler(service)
    app.include_router(router)

    @app.get("/health")
    async def health():
        return {"status": "ok", "service": "vision-ai", "version": "0.1.0"}

    return app


@pytest.mark.asyncio
async def test_post_classify_returns_200(app) -> None:
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        response = await client.post(
            "/v1/vision/classify",
            json={"image_url": "https://example.com/img.jpg"},
        )
    assert response.status_code == 200
    data = response.json()
    assert "predictions" in data
    assert len(data["predictions"]) > 0


@pytest.mark.asyncio
async def test_get_health_returns_200(app) -> None:
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        response = await client.get("/health")
    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "ok"
    assert data["service"] == "vision-ai"
