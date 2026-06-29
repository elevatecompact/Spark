import pytest
from httpx import ASGITransport, AsyncClient

from src.main import create_app


@pytest.fixture
def app():
    return create_app()


@pytest.mark.asyncio
async def test_post_score_returns_200(app):
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        payload = {
            "transaction": {
                "amount": 150.00,
                "currency": "USD",
                "payment_method": "credit_card",
            },
        }
        response = await client.post("/v1/fraud/score", json=payload)

    assert response.status_code == 200
    data = response.json()
    assert "score" in data
    assert "risk_level" in data
    assert "decision" in data
    assert "flags" in data
    assert data["model_version"] == "noop-v1"


@pytest.mark.asyncio
async def test_health_returns_200(app):
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        response = await client.get("/health")

    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "ok"
    assert data["service"] == "fraud-ai"
    assert data["version"] == "0.1.0"
