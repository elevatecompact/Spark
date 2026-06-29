from datetime import datetime, timezone


class RankingRepository:

    async def get_model_metadata(self) -> dict:
        return {
            "version": "noop-v1",
            "last_trained": datetime.now(timezone.utc).isoformat(),
        }

    async def save_training_log(self, entries: list) -> int:
        return len(entries)
