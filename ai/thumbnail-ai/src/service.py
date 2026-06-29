import math
import uuid
from random import randint

from src.events import EventProducer, NoopProducer
from src.models import (
    ExtractFramesRequest,
    ExtractFramesResponse,
    FrameInfo,
    GenerateThumbnailRequest,
    GenerateThumbnailResponse,
    SelectBestThumbnailRequest,
    SelectBestThumbnailResponse,
)
from src.repository import ThumbnailRepository


class ThumbnailService:
    def __init__(
        self,
        repository: ThumbnailRepository | None = None,
        event_producer: EventProducer | None = None,
    ):
        self._repository = repository or ThumbnailRepository()
        self._event_producer = event_producer or NoopProducer()

    async def generate_thumbnail(self, request: GenerateThumbnailRequest) -> GenerateThumbnailResponse:
        time_sec = request.time_seconds if request.time_seconds is not None else 5.0
        thumbnail_id = str(uuid.uuid4())
        thumbnail_url = f"https://cdn.spark.dev/media/thumbnails/{thumbnail_id}.{request.format}"

        self._repository.save_thumbnail(request.video_url, thumbnail_url, time_sec)

        await self._event_producer.publish(
            "thumbnail.generate.completed",
            request.video_url,
            {"thumbnail_url": thumbnail_url, "time_seconds": time_sec},
        )

        return GenerateThumbnailResponse(
            thumbnail_url=thumbnail_url,
            time_seconds=time_sec,
            width=request.width,
            height=request.height,
            format=request.format,
            file_size_bytes=randint(10000, 50000),
            confidence=0.95,
        )

    async def batch_generate(self, requests: list[GenerateThumbnailRequest]) -> list[GenerateThumbnailResponse]:
        results = [await self.generate_thumbnail(req) for req in requests]

        await self._event_producer.publish(
            "thumbnail.batch.completed",
            "batch",
            {"count": len(results)},
        )

        return results

    async def extract_frames(self, request: ExtractFramesRequest) -> ExtractFramesResponse:
        video_duration = 600.0
        num_frames = min(request.max_frames, int(math.ceil(video_duration / request.interval_seconds)))
        frames: list[FrameInfo] = []

        for i in range(num_frames):
            time_sec = round(i * request.interval_seconds, 1)
            frame_id = str(uuid.uuid4())
            frames.append(
                FrameInfo(
                    url=f"https://cdn.spark.dev/media/frames/{frame_id}.{request.format}",
                    time_seconds=time_sec,
                    width=1280,
                    height=720,
                    file_size_bytes=randint(8000, 40000),
                    confidence=round(0.5 + (i / num_frames) * 0.4, 2),
                )
            )

        self._repository.log_extraction(request.video_url, num_frames)

        await self._event_producer.publish(
            "thumbnail.extract.completed",
            request.video_url,
            {"frame_count": num_frames},
        )

        return ExtractFramesResponse(frames=frames, total_frames=num_frames)

    async def select_best(self, request: SelectBestThumbnailRequest) -> SelectBestThumbnailResponse:
        count = min(request.count, len(request.frame_urls))
        selections: list[FrameInfo] = []

        for i in range(count):
            selections.append(
                FrameInfo(
                    url=request.frame_urls[i],
                    time_seconds=float(i * 10),
                    width=1280,
                    height=720,
                    file_size_bytes=randint(10000, 50000),
                    confidence=round(0.9 - i * 0.05, 2),
                    aesthetic_score=round(7.5 - i * 0.3, 1),
                )
            )

        await self._event_producer.publish(
            "thumbnail.select.completed",
            request.video_url or "unknown",
            {"count": count},
        )

        return SelectBestThumbnailResponse(selections=selections)
