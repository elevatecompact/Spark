from fastapi import APIRouter, Depends, Request

from src.models import (
    BatchThumbnailRequest,
    BatchThumbnailResponse,
    ExtractFramesRequest,
    ExtractFramesResponse,
    GenerateThumbnailRequest,
    GenerateThumbnailResponse,
    SelectBestThumbnailRequest,
    SelectBestThumbnailResponse,
)
from src.service import ThumbnailService

router = APIRouter(prefix="/v1/thumbnail")


def get_service(request: Request) -> ThumbnailService:
    return request.app.state.service  # type: ignore[no-any-return]


@router.post("/generate", response_model=GenerateThumbnailResponse)
async def generate_thumbnail(
    request: GenerateThumbnailRequest,
    service: ThumbnailService = Depends(get_service),
) -> GenerateThumbnailResponse:
    return await service.generate_thumbnail(request)


@router.post("/batch", response_model=BatchThumbnailResponse)
async def batch_generate(
    request: BatchThumbnailRequest,
    service: ThumbnailService = Depends(get_service),
) -> BatchThumbnailResponse:
    results = await service.batch_generate(request.requests)
    return BatchThumbnailResponse(results=results)


@router.post("/extract", response_model=ExtractFramesResponse)
async def extract_frames(
    request: ExtractFramesRequest,
    service: ThumbnailService = Depends(get_service),
) -> ExtractFramesResponse:
    return await service.extract_frames(request)


@router.post("/select", response_model=SelectBestThumbnailResponse)
async def select_best(
    request: SelectBestThumbnailRequest,
    service: ThumbnailService = Depends(get_service),
) -> SelectBestThumbnailResponse:
    return await service.select_best(request)
