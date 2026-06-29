from fastapi import APIRouter

from src.models import (
    BatchFraudRequest,
    BatchFraudResponse,
    FraudAnalysisResponse,
    FraudRule,
    FraudScoreRequest,
    FraudScoreResponse,
    TransactionData,
)
from src.service import FraudService


def create_router(service: FraudService) -> APIRouter:
    router = APIRouter(prefix="/v1/fraud")

    @router.post("/score", response_model=FraudScoreResponse)
    async def score_transaction(body: FraudScoreRequest):
        return await service.analyze(body)

    @router.post("/batch", response_model=BatchFraudResponse)
    async def batch_score(body: BatchFraudRequest):
        results = await service.batch_analyze(body.requests)
        return BatchFraudResponse(results=results)

    @router.post("/rules/check", response_model=FraudAnalysisResponse)
    async def check_rules(body: TransactionData):
        return await service.analyze_rules(body)

    @router.get("/rules", response_model=list[FraudRule])
    async def list_rules():
        return await service.get_rules()

    return router
