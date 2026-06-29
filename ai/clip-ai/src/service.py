import math
import uuid

from src.events import EventProducer, NoopProducer
from src.models import (
    ClipSegment,
    DetectClipsRequest,
    DetectClipsResponse,
    GenerateClipRequest,
    GenerateClipResponse,
    HighlightRequest,
    HighlightResponse,
)
from src.repository import ClipRepository


class ClipService:

    def __init__(
        self,
        repository: ClipRepository | None = None,
        producer: EventProducer | None = None,
    ) -> None:
        self._repository = repository or ClipRepository()
        self._producer = producer or NoopProducer()

    async def detect_clips(self, request: DetectClipsRequest) -> DetectClipsResponse:
        mock_duration = 600.0
        count = min(request.max_clips, 5)
        count = max(count, 3)
        interval = mock_duration / (count + 1)
        clips: list[ClipSegment] = []
        for i in range(count):
            start = interval * (i + 1) - 10
            end = interval * (i + 1) + 10
            confidence = round(0.98 - i * 0.08, 2)
            clips.append(
                ClipSegment(
                    start_time=max(0.0, start),
                    end_time=min(mock_duration, end),
                    duration=end - start,
                    confidence=confidence,
                    label=f"segment_{i + 1}",
                    score=confidence,
                )
            )
        response = DetectClipsResponse(
            clips=clips,
            total_duration_seconds=mock_duration,
            clip_count=len(clips),
        )
        await self._producer.send(
            "clip.detect.completed",
            key=request.video_url,
            value=response.model_dump(),
        )
        return response

    async def generate_clip(self, request: GenerateClipRequest) -> GenerateClipResponse:
        duration = request.end_time - request.start_time
        clip_url = f"https://cdn.spark.dev/video/clips/{uuid.uuid4()}.mp4"
        title = request.title or f"clip_{request.start_time}_{request.end_time}"
        response = GenerateClipResponse(
            clip_url=clip_url,
            duration_seconds=duration,
            title=title,
        )
        await self._producer.send(
            "clip.generate.completed",
            key=request.video_url,
            value=response.model_dump(),
        )
        return response

    async def generate_highlights(self, request: HighlightRequest) -> HighlightResponse:
        clips = [
            ClipSegment(
                start_time=10.0,
                end_time=30.0,
                duration=20.0,
                confidence=0.97,
                label="highlight_1",
                score=0.97,
            ),
            ClipSegment(
                start_time=120.0,
                end_time=150.0,
                duration=30.0,
                confidence=0.94,
                label="highlight_2",
                score=0.94,
            ),
            ClipSegment(
                start_time=300.0,
                end_time=330.0,
                duration=30.0,
                confidence=0.91,
                label="highlight_3",
                score=0.91,
            ),
        ]
        total = sum(c.duration for c in clips)
        response = HighlightResponse(
            clips=clips,
            total_duration_seconds=total,
        )
        await self._producer.send(
            "clip.highlights.completed",
            key=request.video_url,
            value=response.model_dump(),
        )
        return response

    async def batch_detect(
        self, requests: list[DetectClipsRequest]
    ) -> list[DetectClipsResponse]:
        results: list[DetectClipsResponse] = []
        for req in requests:
            result = await self.detect_clips(req)
            results.append(result)
        return results
