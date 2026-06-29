import pytest

from src.models import (
    DetectLanguageRequest,
    TranslateRequest,
)
from src.repository import TranslationRepository
from src.events import NoopProducer
from src.service import TranslationService


@pytest.fixture
def service():
    repo = TranslationRepository()
    event_pub = NoopProducer()
    return TranslationService(repo=repo, event_pub=event_pub)


@pytest.mark.asyncio
async def test_translate_returns_correct_structure(service):
    request = TranslateRequest(
        text="hello",
        source_language="en",
        target_language="fr",
    )

    response = await service.translate(request)

    assert response.translated_text == "olleh"
    assert response.source_language == "en"
    assert response.target_language == "fr"
    assert response.confidence == 0.95
    assert response.provider == "noop"
    assert response.model_version == "noop-v1"


@pytest.mark.asyncio
async def test_detect_language_returns_a_language(service):
    request = DetectLanguageRequest(text="hello world")

    response = await service.detect_language(request)

    assert response.detected_language == "en"
    assert response.confidence == 0.85


@pytest.mark.asyncio
async def test_batch_translate_returns_all_results(service):
    requests = [
        TranslateRequest(text="hello", source_language="en", target_language="fr"),
        TranslateRequest(text="world", source_language="en", target_language="es"),
    ]

    results = await service.batch_translate(requests)

    assert len(results) == 2
    assert results[0].translated_text == "olleh"
    assert results[1].translated_text == "dlrow"


@pytest.mark.asyncio
async def test_get_languages_returns_list(service):
    response = await service.get_languages()

    assert len(response.languages) == 15
    assert response.languages[0].code == "en"
    assert response.languages[0].name == "English"
