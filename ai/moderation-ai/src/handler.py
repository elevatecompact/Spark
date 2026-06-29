from fastapi import APIRouter

from src.models import (
    BatchModerationRequest,
    BatchModerationResponse,
    CategoriesResponse,
    ModerateImageRequest,
    ModerateTextRequest,
    ModerationResponse,
)
from src.service import ModerationService


def create_router(service: ModerationService) -> APIRouter:
    router = APIRouter(prefix="/v1/moderation")

    @router.post("/text", response_model=ModerationResponse)
    async def moderate_text(body: ModerateTextRequest):
        return await service.moderate(body)

    @router.post("/image", response_model=ModerationResponse)
    async def moderate_image(body: ModerateImageRequest):
        return await service.moderate(body)

    @router.post("/batch", response_model=BatchModerationResponse)
    async def batch_moderate(body: BatchModerationRequest):
        results = await service.batch_moderate(body.requests)
        return BatchModerationResponse(results=results)

    @router.get("/categories", response_model=CategoriesResponse)
    async def get_categories():
        return await service.get_categories()

    return router
