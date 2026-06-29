from pydantic import BaseModel


class RecommendedItem(BaseModel):
    item_id: str
    score: float
    reason: str | None = None
    position: int
    metadata: dict | None = None


class GetRecommendationsRequest(BaseModel):
    user_id: str
    context: str = "home"
    count: int = 20
    offset: int = 0
    filters: dict | None = None


class GetRecommendationsResponse(BaseModel):
    items: list[RecommendedItem]
    total: int
    context: str
    model_version: str = "noop-v1"


class GetSimilarRequest(BaseModel):
    item_id: str
    count: int = 10
    user_id: str | None = None


class GetSimilarResponse(BaseModel):
    items: list[RecommendedItem]
    source_item_id: str
    model_version: str = "noop-v1"


class GetTrendingRequest(BaseModel):
    timeframe: str = "day"
    category: str | None = None
    count: int = 20


class GetTrendingResponse(BaseModel):
    items: list[RecommendedItem]
    timeframe: str
    category: str | None
    model_version: str = "noop-v1"


class PersonalizeFeedRequest(BaseModel):
    user_id: str
    feed_type: str = "home"
    count: int = 30
    offset: int = 0


class PersonalizeFeedResponse(BaseModel):
    items: list[RecommendedItem]
    total: int
    feed_type: str
    personalization_score: float
    model_version: str = "noop-v1"


class UserInteraction(BaseModel):
    user_id: str
    item_id: str
    interaction_type: str
    weight: float = 1.0
    timestamp: str | None = None


class TrackInteractionResponse(BaseModel):
    status: str
    recorded: bool
    model_version: str = "noop-v1"


class BatchInteractionsRequest(BaseModel):
    interactions: list[UserInteraction]


class BatchInteractionsResponse(BaseModel):
    recorded_count: int
    model_version: str = "noop-v1"


class HealthResponse(BaseModel):
    status: str = "ok"
    service: str = "recommendation-ai"
    version: str = "0.1.0"
