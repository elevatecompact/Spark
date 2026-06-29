from pydantic import BaseModel


class ContentIdea(BaseModel):
    title: str
    description: str
    format: str
    difficulty: str
    estimated_engagement: str
    tags: list[str]


class ContentIdeaRequest(BaseModel):
    niche: str
    count: int = 5
    audience: str | None = None
    format: str | None = None


class ContentIdeaResponse(BaseModel):
    ideas: list[ContentIdea]
    count: int
    model_version: str = "noop-v1"


class OptimizeContentRequest(BaseModel):
    content: str
    content_type: str
    platform: str = "spark"
    goals: list[str] | None = None


class OptimizeContentResponse(BaseModel):
    optimized_content: str
    suggestions: list[str]
    score: float
    improvements: list[str]
    model_version: str


class AnalyzeAudienceRequest(BaseModel):
    creator_id: str
    timeframe: str = "30d"


class AnalyzeAudienceResponse(BaseModel):
    total_followers: int
    growth_rate: float
    demographics: dict
    top_content_types: list[dict]
    best_posting_times: list[str]
    engagement_rate: float
    model_version: str


class ScheduleRequest(BaseModel):
    content: str
    preferred_times: list[str] | None = None
    content_type: str | None = None


class ScheduleResponse(BaseModel):
    suggested_time: str
    alternatives: list[str]
    reasoning: str
    model_version: str


class HashtagRequest(BaseModel):
    content: str
    count: int = 10
    category: str | None = None


class HashtagSuggestion(BaseModel):
    hashtag: str
    popularity: float
    relevance: float
    category: str | None = None


class HashtagResponse(BaseModel):
    hashtags: list[HashtagSuggestion]
    model_version: str


class CaptionRequest(BaseModel):
    image_description: str
    tone: str = "casual"
    length: str = "medium"


class CaptionResponse(BaseModel):
    captions: list[str]
    model_version: str


class HealthResponse(BaseModel):
    status: str = "ok"
    service: str = "creator-ai"
    version: str = "0.1.0"
