import logging
import uuid

from src.events import EventProducer
from src.models import (
    SpeechToTextRequest,
    SpeechToTextResponse,
    TextToSpeechRequest,
    TextToSpeechResponse,
    TranscriptSegment,
    VoicesResponse,
)
from src.repository import VoiceRepository

logger = logging.getLogger(__name__)

_MOCK_TRANSCRIPTIONS: dict[str, str] = {
    "en": "This is a simulated transcription of the provided audio.",
    "es": "Esta es una transcripción simulada del audio proporcionado.",
    "fr": "Ceci est une transcription simulée de l'audio fourni.",
    "de": "Dies ist eine simulierte Transkription des bereitgestellten Audios.",
    "default": "This is a simulated transcription of the provided audio.",
}


class VoiceService:

    def __init__(self, repo: VoiceRepository, event_pub: EventProducer):
        self._repo = repo
        self._event_pub = event_pub

    async def speech_to_text(self, request: SpeechToTextRequest) -> SpeechToTextResponse:
        lang = request.language or "en"
        text = _MOCK_TRANSCRIPTIONS.get(lang, _MOCK_TRANSCRIPTIONS["default"])

        segments = [
            TranscriptSegment(
                start_time=0.0,
                end_time=2.5,
                text=text[: len(text) // 2],
                speaker="speaker_0" if request.enable_diarization else None,
                confidence=0.92,
            ),
            TranscriptSegment(
                start_time=2.5,
                end_time=5.0,
                text=text[len(text) // 2 :],
                speaker="speaker_0" if request.enable_diarization else None,
                confidence=0.95,
            ),
        ]

        response = SpeechToTextResponse(
            text=text,
            confidence=0.93,
            language=lang,
            segments=segments,
            duration_seconds=5.0,
            model_version="noop-v1",
        )

        await self._event_pub.publish(
            "voice.stt.completed",
            {
                "language": lang,
                "format": request.format,
                "enable_diarization": request.enable_diarization,
                "confidence": response.confidence,
            },
        )

        await self._repo.log_stt_request(audio_length=5.0)
        return response

    async def text_to_speech(self, request: TextToSpeechRequest) -> TextToSpeechResponse:
        duration = len(request.text) * 0.05
        audio_url = f"https://cdn.spark.dev/audio/tts/{uuid.uuid4()}.mp3"

        response = TextToSpeechResponse(
            audio_url=audio_url,
            duration_seconds=duration,
            format="mp3",
            voice=request.voice,
            model_version="noop-v1",
        )

        await self._event_pub.publish(
            "voice.tts.completed",
            {
                "voice": request.voice,
                "language": request.language,
                "text_length": len(request.text),
                "duration": duration,
            },
        )

        await self._repo.log_tts_request(text_length=len(request.text), voice=request.voice)
        return response

    async def get_voices(self) -> VoicesResponse:
        profiles = await self._repo.get_voices()
        return VoicesResponse(voices=profiles, model_version="noop-v1")
