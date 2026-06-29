import pytest

from src.models import ContentIdeaRequest
from src.service import CreatorService


@pytest.mark.asyncio
async def test_generate_ideas_returns_ideas_list():
    service = CreatorService()
    request = ContentIdeaRequest(niche="gaming", count=3)
    response = await service.generate_ideas(request)
    assert response.count == 3
    assert len(response.ideas) == 3
    for idea in response.ideas:
        assert idea.title
        assert idea.description
        assert idea.format
        assert idea.difficulty in ("easy", "medium", "hard")
        assert idea.estimated_engagement in ("low", "medium", "high")
        assert isinstance(idea.tags, list)
