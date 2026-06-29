import pytest

from src.models import DetectClipsRequest
from src.service import ClipService


@pytest.mark.asyncio
async def test_detect_clips_returns_clips_list() -> None:
    service = ClipService()
    request = DetectClipsRequest(video_url="https://example.com/video.mp4")
    response = await service.detect_clips(request)
    assert response.clip_count > 0
    assert len(response.clips) == response.clip_count
    assert all(c.confidence > 0.0 for c in response.clips)
    assert response.model_version == "noop-v1"
