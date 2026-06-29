import pytest

from src.models import GenerateThumbnailRequest
from src.service import ThumbnailService


@pytest.mark.asyncio
async def test_generate_returns_thumbnail_url() -> None:
    service = ThumbnailService()
    request = GenerateThumbnailRequest(video_url="https://example.com/video.mp4")
    response = await service.generate_thumbnail(request)
    assert response.thumbnail_url.startswith("https://cdn.spark.dev/media/thumbnails/")
    assert response.thumbnail_url.endswith(".jpg")
    assert response.time_seconds == 5.0
    assert 10000 <= response.file_size_bytes <= 50000


@pytest.mark.asyncio
async def test_generate_with_custom_time() -> None:
    service = ThumbnailService()
    request = GenerateThumbnailRequest(video_url="https://example.com/video.mp4", time_seconds=30.0)
    response = await service.generate_thumbnail(request)
    assert response.time_seconds == 30.0


@pytest.mark.asyncio
async def test_batch_generate() -> None:
    service = ThumbnailService()
    requests = [
        GenerateThumbnailRequest(video_url="https://example.com/1.mp4"),
        GenerateThumbnailRequest(video_url="https://example.com/2.mp4"),
    ]
    results = await service.batch_generate(requests)
    assert len(results) == 2
    assert all(r.thumbnail_url.startswith("https://cdn.spark.dev/media/thumbnails/") for r in results)
