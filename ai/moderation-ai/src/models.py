from pydantic import BaseModel


class ModerateTextRequest(BaseModel):
    text: str
    content_id: str
    user_id: str | None = None
    language: str | None = None


class ModerateImageRequest(BaseModel):
    image_url: str
    content_id: str
    user_id: str | None = None


class ModerationResult(BaseModel):
    is_ok: bool
    category: str | None
    confidence: float
    label: str
    details: dict | None = None


class ModerationResponse(BaseModel):
    content_id: str
    text_result: ModerationResult | None = None
    image_result: ModerationResult | None = None
    overall_safe: bool
    requires_review: bool
    model_version: str = "noop-v1"


class BatchModerationRequest(BaseModel):
    requests: list[ModerateTextRequest | ModerateImageRequest]


class BatchModerationResponse(BaseModel):
    results: list[ModerationResponse]


class ModerationCategory(BaseModel):
    name: str
    description: str
    severity: str


class CategoriesResponse(BaseModel):
    categories: list[ModerationCategory]
    model_version: str


class HealthResponse(BaseModel):
    status: str = "ok"
    service: str = "moderation-ai"
    version: str = "0.1.0"
