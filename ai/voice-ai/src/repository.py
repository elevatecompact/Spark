import logging

from src.models import VoiceProfile

logger = logging.getLogger(__name__)


class VoiceRepository:

    async def get_voices(self) -> list[VoiceProfile]:
        return [
            VoiceProfile(
                profile_id="voice-en-f1",
                name="Alice",
                gender="female",
                language="en-US",
                style="neutral",
            ),
            VoiceProfile(
                profile_id="voice-en-m1",
                name="Bob",
                gender="male",
                language="en-US",
                style="neutral",
            ),
            VoiceProfile(
                profile_id="voice-en-f2",
                name="Clara",
                gender="female",
                language="en-US",
                style="cheerful",
            ),
            VoiceProfile(
                profile_id="voice-es-f1",
                name="Lucia",
                gender="female",
                language="es-ES",
                style="neutral",
            ),
            VoiceProfile(
                profile_id="voice-fr-m1",
                name="Pierre",
                gender="male",
                language="fr-FR",
                style="neutral",
            ),
            VoiceProfile(
                profile_id="voice-de-f1",
                name="Greta",
                gender="female",
                language="de-DE",
                style="neutral",
            ),
        ]

    async def log_tts_request(self, text_length: int, voice: str) -> None:
        logger.info("TTS request logged: text_length=%d, voice=%s", text_length, voice)

    async def log_stt_request(self, audio_length: float) -> None:
        logger.info("STT request logged: audio_length=%.2f", audio_length)
