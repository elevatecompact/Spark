from src.models import ModerationCategory


class ModerationRepository:
    async def log_moderation(self, content_id: str, result: dict) -> None:
        pass

    async def get_categories(self) -> list[ModerationCategory]:
        return [
            ModerationCategory(
                name="hate_speech",
                description="Hateful or discriminatory language",
                severity="high",
            ),
            ModerationCategory(
                name="harassment",
                description="Harassing or bullying content",
                severity="high",
            ),
            ModerationCategory(
                name="violence",
                description="Violent or gory content",
                severity="critical",
            ),
            ModerationCategory(
                name="self_harm",
                description="Content promoting self-harm or suicide",
                severity="critical",
            ),
            ModerationCategory(
                name="sexual",
                description="Sexually explicit content",
                severity="medium",
            ),
            ModerationCategory(
                name="spam",
                description="Unsolicited or repetitive content",
                severity="low",
            ),
            ModerationCategory(
                name="misinformation",
                description="False or misleading information",
                severity="medium",
            ),
        ]
