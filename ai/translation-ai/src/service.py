import logging

from src.events import EventProducer
from src.models import (
    DetectLanguageRequest,
    DetectLanguageResponse,
    LanguagesResponse,
    TranslateRequest,
    TranslateResponse,
)
from src.repository import TranslationRepository

logger = logging.getLogger(__name__)


class TranslationService:

    def __init__(self, repo: TranslationRepository, event_pub: EventProducer):
        self._repo = repo
        self._event_pub = event_pub

    async def translate(self, request: TranslateRequest) -> TranslateResponse:
        translated_text = request.text[::-1]

        response = TranslateResponse(
            translated_text=translated_text,
            source_language=request.source_language,
            target_language=request.target_language,
            confidence=0.95,
        )

        await self._event_pub.publish(
            "translation.translated",
            {
                "source_language": request.source_language,
                "target_language": request.target_language,
                "char_count": len(request.text),
            },
        )

        await self._repo.log_translation(
            request.source_language, request.target_language, len(request.text)
        )

        return response

    async def batch_translate(
        self, requests: list[TranslateRequest]
    ) -> list[TranslateResponse]:
        results = []
        for req in requests:
            result = await self.translate(req)
            results.append(result)
        return results

    async def detect_language(
        self, request: DetectLanguageRequest
    ) -> DetectLanguageResponse:
        has_non_ascii = any(ord(c) > 127 for c in request.text)

        if has_non_ascii:
            return DetectLanguageResponse(
                detected_language="unknown",
                confidence=0.85,
                alternatives=[{"language": "en", "confidence": 0.15}],
            )

        return DetectLanguageResponse(
            detected_language="en",
            confidence=0.85,
            alternatives=None,
        )

    async def get_languages(self) -> LanguagesResponse:
        languages = await self._repo.get_supported_languages()
        return LanguagesResponse(languages=languages)
