from fastapi import APIRouter, Depends

from src.models import (
    BatchInteractionsRequest,
    BatchInteractionsResponse,
    GetRecommendationsRequest,
    GetRecommendationsResponse,
    GetSimilarRequest,
    GetSimilarResponse,
    GetTrendingRequest,
    GetTrendingResponse,
    PersonalizeFeedRequest,
    PersonalizeFeedResponse,
    TrackInteractionResponse,
    UserInteraction,
)
from src.service import RecommendationService

router = APIRouter(prefix="/v1/recommendation")


async def get_service() -> RecommendationService:
    from src.main import app

    return app.state.service


@router.post("/get", response_model=GetRecommendationsResponse)
async def get_recommendations(
    request: GetRecommendationsRequest,
    service: RecommendationService = Depends(get_service),
) -> GetRecommendationsResponse:
    return await service.get_recommendations(request)


@router.post("/similar", response_model=GetSimilarResponse)
async def get_similar(
    request: GetSimilarRequest,
    service: RecommendationService = Depends(get_service),
) -> GetSimilarResponse:
    return await service.get_similar(request)


@router.post("/trending", response_model=GetTrendingResponse)
async def get_trending(
    request: GetTrendingRequest,
    service: RecommendationService = Depends(get_service),
) -> GetTrendingResponse:
    return await service.get_trending(request)


@router.post("/feed", response_model=PersonalizeFeedResponse)
async def personalize_feed(
    request: PersonalizeFeedRequest,
    service: RecommendationService = Depends(get_service),
) -> PersonalizeFeedResponse:
    return await service.personalize_feed(request)


@router.post("/interaction", response_model=TrackInteractionResponse)
async def track_interaction(
    interaction: UserInteraction,
    service: RecommendationService = Depends(get_service),
) -> TrackInteractionResponse:
    return await service.track_interaction(interaction)


@router.post("/interactions/batch", response_model=BatchInteractionsResponse)
async def batch_track(
    request: BatchInteractionsRequest,
    service: RecommendationService = Depends(get_service),
) -> BatchInteractionsResponse:
    return await service.batch_track(request.interactions)
