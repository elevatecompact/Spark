import logging

from src.events import EventProducer
from src.models import (
    CategoriesResponse,
    ModerateImageRequest,
    ModerateTextRequest,
    ModerationResponse,
    ModerationResult,
)
from src.repository import ModerationRepository

logger = logging.getLogger(__name__)


class ModerationService:
    def __init__(self, repo: ModerationRepository, event_pub: EventProducer):
        self._repo = repo
        self._event_pub = event_pub

    async def moderate_text(self, request: ModerateTextRequest) -> ModerationResult:
        result = ModerationResult(
            is_ok=True,
            category=None,
            confidence=0.95,
            label="clean",
            details=None,
        )
        await self._repo.log_moderation(request.content_id, result.model_dump())
        await self._event_pub.publish(
            "moderation.analysis.completed",
            {
                "content_id": request.content_id,
                "type": "text",
                "result": result.model_dump(),
            },
        )
        return result

    async def moderate_image(self, request: ModerateImageRequest) -> ModerationResult:
        result = ModerationResult(
            is_ok=True,
            category=None,
            confidence=0.95,
            label="clean",
            details=None,
        )
        await self._repo.log_moderation(request.content_id, result.model_dump())
        await self._event_pub.publish(
            "moderation.analysis.completed",
            {
                "content_id": request.content_id,
                "type": "image",
                "result": result.model_dump(),
            },
        )
        return result

    async def moderate(
        self, request: ModerateTextRequest | ModerateImageRequest
    ) -> ModerationResponse:
        if isinstance(request, ModerateTextRequest):
            text_result = await self.moderate_text(request)
            image_result = None
        else:
            text_result = None
            image_result = await self.moderate_image(request)

        return ModerationResponse(
            content_id=request.content_id,
            text_result=text_result,
            image_result=image_result,
            overall_safe=True,
            requires_review=False,
            model_version="noop-v1",
        )

    async def batch_moderate(
        self, requests: list[ModerateTextRequest | ModerateImageRequest]
    ) -> list[ModerationResponse]:
        results = []
        for req in requests:
            result = await self.moderate(req)
            results.append(result)
        await self._event_pub.publish(
            "moderation.batch.completed",
            {"batch_size": len(requests), "results": [r.model_dump() for r in results]},
        )
        return results

    async def get_categories(self) -> CategoriesResponse:
        categories = await self._repo.get_categories()
        return CategoriesResponse(categories=categories, model_version="noop-v1")
