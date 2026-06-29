import pytest

from src.models import SpeechToTextRequest, TextToSpeechRequest
from src.repository import VoiceRepository
from src.events import NoopProducer
from src.service import VoiceService


@pytest.fixture
def service():
    repo = VoiceRepository()
    event_pub = NoopProducer()
    return VoiceService(repo=repo, event_pub=event_pub)


@pytest.mark.asyncio
async def test_speech_to_text_returns_text(service):
    request = SpeechToTextRequest(
        audio_url="https://cdn.spark.dev/audio/sample.wav",
        language="en",
    )

    response = await service.speech_to_text(request)

    assert response.text
    assert response.confidence > 0
    assert response.language == "en"
    assert response.model_version == "noop-v1"
    assert len(response.segments) > 0


@pytest.mark.asyncio
async def test_text_to_speech_returns_audio_url(service):
    request = TextToSpeechRequest(
        text="Hello, this is a test of the voice service.",
        voice="default",
    )

    response = await service.text_to_speech(request)

    assert response.audio_url.startswith("https://cdn.spark.dev/audio/tts/")
    assert response.audio_url.endswith(".mp3")
    assert response.duration_seconds > 0
    assert response.voice == "default"
    assert response.model_version == "noop-v1"
