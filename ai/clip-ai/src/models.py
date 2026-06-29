from pydantic import BaseModel
from typing import Optional


class ClipSegment(BaseModel):
    start_time: float
    end_time: float
    duration: float
    confidence: float
    label: str | None = None
    score: float = 0.0


class DetectClipsRequest(BaseModel):
    video_url: str
    max_clips: int = 10
    min_duration_seconds: float = 5.0
    max_duration_seconds: float = 120.0
    detection_strategy: str = "auto"


class DetectClipsResponse(BaseModel):
    clips: list[ClipSegment]
    total_duration_seconds: float
    clip_count: int
    model_version: str = "noop-v1"


class GenerateClipRequest(BaseModel):
    video_url: str
    start_time: float
    end_time: float
    title: str | None = None
    description: str | None = None
    include_audio: bool = True


class GenerateClipResponse(BaseModel):
    clip_url: str
    duration_seconds: float
    format: str = "mp4"
    file_size_bytes: int | None = None
    title: str
    model_version: str = "noop-v1"


class HighlightRequest(BaseModel):
    video_url: str
    duration_seconds: float = 60.0
    style: str = "auto"


class HighlightResponse(BaseModel):
    clips: list[ClipSegment]
    highlight_url: str | None = None
    total_duration_seconds: float
    model_version: str = "noop-v1"


class BatchClipRequest(BaseModel):
    requests: list[DetectClipsRequest]


class BatchClipResponse(BaseModel):
    results: list[DetectClipsResponse]


class HealthResponse(BaseModel):
    status: str = "ok"
    service: str = "clip-ai"
    version: str = "0.1.0"
