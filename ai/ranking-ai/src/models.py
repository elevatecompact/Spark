from pydantic import BaseModel


class RankItem(BaseModel):
    item_id: str
    features: dict[str, float]


class RankContext(BaseModel):
    user_id: str | None = None
    session_id: str | None = None
    page: str | None = None


class RankingRequest(BaseModel):
    items: list[RankItem]
    context: RankContext | None = None


class RankedItem(BaseModel):
    item_id: str
    score: float
    position: int
    reason: str | None = None


class RankingResponse(BaseModel):
    ranked_items: list[RankedItem]
    scores: list[float]
    model_version: str = "noop-v1"


class TrainingData(BaseModel):
    item_id: str
    features: dict[str, float]
    label: float
    weight: float = 1.0


class TrainingResponse(BaseModel):
    status: str
    samples_trained: int
    model_version: str


class BatchRankRequest(BaseModel):
    requests: list[RankingRequest]


class BatchRankResponse(BaseModel):
    results: list[RankingResponse]


class HealthResponse(BaseModel):
    status: str = "ok"
    service: str = "ranking-ai"
    version: str = "0.1.0"
