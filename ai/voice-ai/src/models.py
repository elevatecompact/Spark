from pydantic import BaseModel


class TranscriptSegment(BaseModel):
    start_time: float
    end_time: float
    text: str
    speaker: str | None
    confidence: float


class SpeechToTextRequest(BaseModel):
    audio_url: str
    language: str | None = None
    format: str = "wav"
    enable_diarization: bool = False


class SpeechToTextResponse(BaseModel):
    text: str
    confidence: float
    language: str
    segments: list[TranscriptSegment] | None
    duration_seconds: float | None
    model_version: str = "noop-v1"


class TextToSpeechRequest(BaseModel):
    text: str
    voice: str = "default"
    language: str = "en-US"
    speed: float = 1.0
    pitch: float = 1.0


class TextToSpeechResponse(BaseModel):
    audio_url: str
    duration_seconds: float
    format: str = "mp3"
    voice: str
    model_version: str = "noop-v1"


class VoiceProfile(BaseModel):
    profile_id: str
    name: str
    gender: str | None
    language: str
    style: str | None


class VoicesResponse(BaseModel):
    voices: list[VoiceProfile]
    model_version: str


class HealthResponse(BaseModel):
    status: str = "ok"
    service: str = "voice-ai"
    version: str = "0.1.0"
