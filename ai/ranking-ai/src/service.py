import logging

from src.models import (
    RankingRequest,
    RankingResponse,
    RankedItem,
    TrainingData,
    TrainingResponse,
)
from src.repository import RankingRepository
from src.events import EventProducer

logger = logging.getLogger(__name__)


class RankingService:

    def __init__(self, repo: RankingRepository, event_pub: EventProducer):
        self._repo = repo
        self._event_pub = event_pub

    async def rank_items(self, request: RankingRequest) -> RankingResponse:
        scored = []
        for item in request.items:
            score = sum(item.features.values())
            scored.append((item.item_id, score))

        scored.sort(key=lambda x: x[1], reverse=True)

        ranked = []
        scores = []
        for position, (item_id, score) in enumerate(scored):
            ranked.append(
                RankedItem(
                    item_id=item_id,
                    score=score,
                    position=position,
                    reason="feature_sum",
                )
            )
            scores.append(score)

        response = RankingResponse(
            ranked_items=ranked,
            scores=scores,
            model_version="noop-v1",
        )

        await self._event_pub.publish(
            "ranking.recommendations.served",
            {
                "context": request.context.model_dump() if request.context else None,
                "item_count": len(ranked),
            },
        )

        return response

    async def batch_rank(
        self, requests: list[RankingRequest]
    ) -> list[RankingResponse]:
        results = []
        for req in requests:
            result = await self.rank_items(req)
            results.append(result)
        return results

    async def train_model(self, data: list[TrainingData]) -> TrainingResponse:
        count = await self._repo.save_training_log(
            [d.model_dump() for d in data]
        )

        await self._event_pub.publish(
            "ranking.model.trained",
            {"samples": count, "model_version": "noop-v1"},
        )

        return TrainingResponse(
            status="ok",
            samples_trained=count,
            model_version="noop-v1",
        )

    async def get_model_info(self) -> dict:
        return await self._repo.get_model_metadata()
