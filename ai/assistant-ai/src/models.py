from pydantic import BaseModel
from typing import Optional


class ChatRequest(BaseModel):
    message: str
    conversation_id: Optional[str] = None
    user_id: Optional[str] = None
    context: Optional[dict] = None
    stream: bool = False


class ChatMessage(BaseModel):
    role: str
    content: str
    timestamp: Optional[str] = None


class ChatResponse(BaseModel):
    reply: str
    conversation_id: str
    tokens_used: int
    model: str = "noop-gpt"
    model_version: str = "noop-v1"


class ConversationHistory(BaseModel):
    conversation_id: str
    messages: list[ChatMessage]
    created_at: str
    updated_at: str


class GetHistoryResponse(BaseModel):
    conversation: ConversationHistory
    message_count: int


class SuggestRequest(BaseModel):
    prompt: str
    context_type: str
    max_suggestions: int = 3


class SuggestResponse(BaseModel):
    suggestions: list[str]
    model_version: str


class SummarizeRequest(BaseModel):
    text: str
    max_length: int = 200
    format: str = "paragraph"


class SummarizeResponse(BaseModel):
    summary: str
    original_length: int
    summary_length: int
    compression_ratio: float
    model_version: str


class ModifyContentRequest(BaseModel):
    text: str
    instruction: str
    tone: Optional[str] = None
    language: Optional[str] = None


class ModifyContentResponse(BaseModel):
    modified_text: str
    changes_made: list[str]
    model_version: str


class HealthResponse(BaseModel):
    status: str = "ok"
    service: str = "assistant-ai"
    version: str = "0.1.0"
