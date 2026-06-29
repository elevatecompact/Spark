import pytest

from src.events import NoopProducer
from src.models import ModerateImageRequest, ModerateTextRequest
from src.repository import ModerationRepository
from src.service import ModerationService


@pytest.fixture
def service():
    repo = ModerationRepository()
    event_pub = NoopProducer()
    return ModerationService(repo=repo, event_pub=event_pub)


@pytest.mark.asyncio
async def test_moderate_text_returns_moderation_result(service):
    request = ModerateTextRequest(
        text="Hello world",
        content_id="c1",
        user_id="u1",
        language="en",
    )
    result = await service.moderate_text(request)

    assert result.is_ok is True
    assert result.confidence == 0.95
    assert result.label == "clean"
    assert result.category is None


@pytest.mark.asyncio
async def test_moderate_returns_overall_safe(service):
    request = ModerateTextRequest(
        text="Hello world",
        content_id="c1",
    )
    response = await service.moderate(request)

    assert response.overall_safe is True
    assert response.requires_review is False
    assert response.model_version == "noop-v1"
    assert response.text_result is not None
    assert response.image_result is None


@pytest.mark.asyncio
async def test_moderate_image_returns_moderation_result(service):
    request = ModerateImageRequest(
        image_url="https://example.com/img.jpg",
        content_id="c2",
    )
    result = await service.moderate_image(request)

    assert result.is_ok is True
    assert result.confidence == 0.95
    assert result.label == "clean"


@pytest.mark.asyncio
async def test_batch_moderate_returns_multiple_results(service):
    requests = [
        ModerateTextRequest(text="Hello", content_id="c1"),
        ModerateImageRequest(image_url="https://example.com/img.jpg", content_id="c2"),
    ]
    results = await service.batch_moderate(requests)

    assert len(results) == 2
    assert all(r.overall_safe is True for r in results)
