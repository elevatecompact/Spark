from fastapi import APIRouter

from src.models import (
    BatchClipRequest,
    BatchClipResponse,
    DetectClipsRequest,
    DetectClipsResponse,
    GenerateClipRequest,
    GenerateClipResponse,
    HighlightRequest,
    HighlightResponse,
)
from src.service import ClipService

router = APIRouter(prefix="/v1/clip")
service = ClipService()


@router.post("/detect", response_model=DetectClipsResponse)
async def detect_clips(request: DetectClipsRequest) -> DetectClipsResponse:
    return await service.detect_clips(request)


@router.post("/generate", response_model=GenerateClipResponse)
async def generate_clip(request: GenerateClipRequest) -> GenerateClipResponse:
    return await service.generate_clip(request)


@router.post("/highlights", response_model=HighlightResponse)
async def generate_highlights(request: HighlightRequest) -> HighlightResponse:
    return await service.generate_highlights(request)


@router.post("/batch", response_model=BatchClipResponse)
async def batch_detect(request: BatchClipRequest) -> BatchClipResponse:
    results = await service.batch_detect(request.requests)
    return BatchClipResponse(results=results)
