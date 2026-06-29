from pydantic import BaseModel


class GenerateThumbnailRequest(BaseModel):
    video_url: str
    time_seconds: float | None = None
    width: int = 1280
    height: int = 720
    format: str = "jpg"
    strategy: str = "auto"


class GenerateThumbnailResponse(BaseModel):
    thumbnail_url: str
    time_seconds: float
    width: int
    height: int
    format: str
    file_size_bytes: int
    confidence: float
    model_version: str = "noop-v1"


class BatchThumbnailRequest(BaseModel):
    requests: list[GenerateThumbnailRequest]


class BatchThumbnailResponse(BaseModel):
    results: list[GenerateThumbnailResponse]


class FrameInfo(BaseModel):
    url: str
    time_seconds: float
    width: int
    height: int
    file_size_bytes: int
    confidence: float
    aesthetic_score: float | None = None


class ExtractFramesRequest(BaseModel):
    video_url: str
    interval_seconds: float = 10.0
    max_frames: int = 50
    format: str = "jpg"


class ExtractFramesResponse(BaseModel):
    frames: list[FrameInfo]
    total_frames: int
    model_version: str = "noop-v1"


class SelectBestThumbnailRequest(BaseModel):
    frame_urls: list[str]
    video_url: str | None = None
    count: int = 1


class SelectBestThumbnailResponse(BaseModel):
    selections: list[FrameInfo]
    model_version: str = "noop-v1"


class HealthResponse(BaseModel):
    status: str = "ok"
    service: str = "thumbnail-ai"
    version: str = "0.1.0"
