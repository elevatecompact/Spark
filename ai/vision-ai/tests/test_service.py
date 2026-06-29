import pytest

from src.events import NoopProducer
from src.models import ClassifyImageRequest
from src.repository import VisionRepository
from src.service import VisionService


@pytest.fixture
def service() -> VisionService:
    return VisionService(
        repository=VisionRepository(),
        event_producer=NoopProducer(),
    )


@pytest.mark.asyncio
async def test_classify_returns_predictions_list(service: VisionService) -> None:
    request = ClassifyImageRequest(image_url="https://example.com/image.jpg")
    response = await service.classify_image(request)
    assert len(response.predictions) > 0
    assert response.dominant_category == "people"
    assert response.model_version == "noop-v1"
