from fastapi import APIRouter, Depends

from src.models import (
    AnalyzeAudienceRequest,
    AnalyzeAudienceResponse,
    CaptionRequest,
    CaptionResponse,
    ContentIdeaRequest,
    ContentIdeaResponse,
    HashtagRequest,
    HashtagResponse,
    OptimizeContentRequest,
    OptimizeContentResponse,
    ScheduleRequest,
    ScheduleResponse,
)
from src.service import CreatorService

router = APIRouter(prefix="/v1/creator")


async def get_service() -> CreatorService:
    from src.main import app
    return app.state.service


@router.post("/ideas", response_model=ContentIdeaResponse)
async def generate_ideas(
    request: ContentIdeaRequest,
    service: CreatorService = Depends(get_service),
) -> ContentIdeaResponse:
    return await service.generate_ideas(request)


@router.post("/optimize", response_model=OptimizeContentResponse)
async def optimize_content(
    request: OptimizeContentRequest,
    service: CreatorService = Depends(get_service),
) -> OptimizeContentResponse:
    return await service.optimize_content(request)


@router.post("/audience", response_model=AnalyzeAudienceResponse)
async def analyze_audience(
    request: AnalyzeAudienceRequest,
    service: CreatorService = Depends(get_service),
) -> AnalyzeAudienceResponse:
    return await service.analyze_audience(request)


@router.post("/schedule", response_model=ScheduleResponse)
async def suggest_schedule(
    request: ScheduleRequest,
    service: CreatorService = Depends(get_service),
) -> ScheduleResponse:
    return await service.suggest_schedule(request)


@router.post("/hashtags", response_model=HashtagResponse)
async def suggest_hashtags(
    request: HashtagRequest,
    service: CreatorService = Depends(get_service),
) -> HashtagResponse:
    return await service.suggest_hashtags(request)


@router.post("/captions", response_model=CaptionResponse)
async def generate_captions(
    request: CaptionRequest,
    service: CreatorService = Depends(get_service),
) -> CaptionResponse:
    return await service.generate_captions(request)
