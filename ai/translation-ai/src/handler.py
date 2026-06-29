from fastapi import APIRouter

from src.models import (
    BatchTranslateRequest,
    BatchTranslateResponse,
    DetectLanguageRequest,
    DetectLanguageResponse,
    LanguagesResponse,
    TranslateRequest,
    TranslateResponse,
)
from src.service import TranslationService


def create_router(service: TranslationService) -> APIRouter:
    router = APIRouter(prefix="/v1/translation")

    @router.post("/translate", response_model=TranslateResponse)
    async def translate(body: TranslateRequest):
        return await service.translate(body)

    @router.post("/batch", response_model=BatchTranslateResponse)
    async def batch_translate(body: BatchTranslateRequest):
        results = await service.batch_translate(body.requests)
        return BatchTranslateResponse(results=results)

    @router.post("/detect", response_model=DetectLanguageResponse)
    async def detect_language(body: DetectLanguageRequest):
        return await service.detect_language(body)

    @router.get("/languages", response_model=LanguagesResponse)
    async def get_languages():
        return await service.get_languages()

    return router
