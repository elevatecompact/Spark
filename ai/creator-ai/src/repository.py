import random


class CreatorRepository:

    def get_creator_profile(self, creator_id: str) -> dict:
        return {
            "creator_id": creator_id,
            "name": f"Creator {creator_id}",
            "followers": random.randint(1000, 500000),
            "niches": ["technology", "lifestyle", "entertainment"],
            "avg_engagement": round(random.uniform(1.0, 6.0), 2),
            "top_formats": ["video", "image", "text"],
        }

    def log_generation(self, creator_id: str | None, feature: str, content_summary: str) -> None:
        pass
