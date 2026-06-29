from pydantic import BaseModel


class TranslateRequest(BaseModel):
    text: str
    source_language: str
    target_language: str
    content_type: str | None = None


class TranslateResponse(BaseModel):
    translated_text: str
    source_language: str
    target_language: str
    confidence: float
    provider: str = "noop"
    model_version: str = "noop-v1"


class DetectLanguageRequest(BaseModel):
    text: str


class DetectLanguageResponse(BaseModel):
    detected_language: str
    confidence: float
    alternatives: list[dict] | None


class BatchTranslateRequest(BaseModel):
    requests: list[TranslateRequest]


class BatchTranslateResponse(BaseModel):
    results: list[TranslateResponse]


class SupportedLanguage(BaseModel):
    code: str
    name: str
    native_name: str
    is_rtl: bool = False


class LanguagesResponse(BaseModel):
    languages: list[SupportedLanguage]


class HealthResponse(BaseModel):
    status: str = "ok"
    service: str = "translation-ai"
    version: str = "0.1.0"
