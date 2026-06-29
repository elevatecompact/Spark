from fastapi import APIRouter

from src.models import (
    SpeechToTextRequest,
    SpeechToTextResponse,
    TextToSpeechRequest,
    TextToSpeechResponse,
    VoicesResponse,
)
from src.service import VoiceService


def create_router(service: VoiceService) -> APIRouter:
    router = APIRouter(prefix="/v1/voice")

    @router.post("/stt", response_model=SpeechToTextResponse)
    async def speech_to_text(body: SpeechToTextRequest):
        return await service.speech_to_text(body)

    @router.post("/tts", response_model=TextToSpeechResponse)
    async def text_to_speech(body: TextToSpeechRequest):
        return await service.text_to_speech(body)

    @router.get("/voices", response_model=VoicesResponse)
    async def get_voices():
        return await service.get_voices()

    return router
