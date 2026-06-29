import pytest

from src.models import ChatRequest, SuggestRequest
from src.service import AssistantService


@pytest.fixture
def service() -> AssistantService:
    return AssistantService()


@pytest.mark.asyncio
async def test_chat_returns_reply(service: AssistantService) -> None:
    request = ChatRequest(message="Hello")
    response = await service.chat(request)

    assert "Hello" in response.reply
    assert response.conversation_id is not None
    assert response.tokens_used > 0
    assert response.model == "noop-gpt"


@pytest.mark.asyncio
async def test_suggest_returns_list(service: AssistantService) -> None:
    request = SuggestRequest(
        prompt="Write a post", context_type="post", max_suggestions=3
    )
    response = await service.suggest(request)

    assert len(response.suggestions) == 3
    assert "post" in response.suggestions[0]
    assert response.model_version == "noop-v1"
