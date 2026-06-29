from fastapi import APIRouter

from src.models import (
    ChatRequest,
    ChatResponse,
    GetHistoryResponse,
    ModifyContentRequest,
    ModifyContentResponse,
    SuggestRequest,
    SuggestResponse,
    SummarizeRequest,
    SummarizeResponse,
)
from src.service import AssistantService

router = APIRouter(prefix="/v1/assistant")


def create_handler(service: AssistantService) -> APIRouter:
    @router.post("/chat", response_model=ChatResponse)
    async def chat(request: ChatRequest) -> ChatResponse:
        return await service.chat(request)

    @router.get("/conversations/{conversation_id}", response_model=GetHistoryResponse)
    async def get_history(conversation_id: str) -> GetHistoryResponse:
        return await service.get_history(conversation_id)

    @router.post("/suggest", response_model=SuggestResponse)
    async def suggest(request: SuggestRequest) -> SuggestResponse:
        return await service.suggest(request)

    @router.post("/summarize", response_model=SummarizeResponse)
    async def summarize(request: SummarizeRequest) -> SummarizeResponse:
        return await service.summarize(request)

    @router.post("/modify", response_model=ModifyContentResponse)
    async def modify_content(request: ModifyContentRequest) -> ModifyContentResponse:
        return await service.modify_content(request)

    return router
