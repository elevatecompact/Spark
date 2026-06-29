import datetime

from src.models import UserInteraction


class RecommendationRepository:
    def __init__(self) -> None:
        self._interactions: list[dict] = []

    def get_user_profile(self, user_id: str) -> dict:
        return {
            "user_id": user_id,
            "preferred_categories": ["technology", "science", "arts"],
            "recently_viewed": [],
        }

    def get_item_metadata(self, item_id: str) -> dict:
        return {
            "item_id": item_id,
            "title": f"Item {item_id}",
            "category": "general",
            "tags": [],
        }

    def store_interaction(self, interaction: UserInteraction) -> None:
        self._interactions.append(
            {
                "user_id": interaction.user_id,
                "item_id": interaction.item_id,
                "interaction_type": interaction.interaction_type,
                "weight": interaction.weight,
                "timestamp": interaction.timestamp or datetime.datetime.now(datetime.timezone.utc).isoformat(),
            }
        )

    def get_recent_interactions(self, user_id: str, limit: int = 100) -> list[dict]:
        user_interactions = [i for i in self._interactions if i["user_id"] == user_id]
        return user_interactions[-limit:] if len(user_interactions) > limit else user_interactions
