from fastapi import APIRouter

from src.models import (
    BatchRankRequest,
    BatchRankResponse,
    RankingRequest,
    RankingResponse,
    TrainingData,
    TrainingResponse,
)
from src.service import RankingService


def create_router(service: RankingService) -> APIRouter:
    router = APIRouter(prefix="/v1/ranking")

    @router.post("/rank", response_model=RankingResponse)
    async def rank_items(body: RankingRequest):
        return await service.rank_items(body)

    @router.post("/batch", response_model=BatchRankResponse)
    async def batch_rank(body: BatchRankRequest):
        results = await service.batch_rank(body.requests)
        return BatchRankResponse(results=results)

    @router.post("/train", response_model=TrainingResponse)
    async def train_model(body: list[TrainingData]):
        return await service.train_model(body)

    @router.get("/model")
    async def get_model_info():
        return await service.get_model_info()

    return router
