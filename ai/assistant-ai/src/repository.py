from src.models import ConversationHistory


class AssistantRepository:
    def __init__(self) -> None:
        self._conversations: dict[str, ConversationHistory] = {}

    def save_conversation(self, conversation: ConversationHistory) -> None:
        self._conversations[conversation.conversation_id] = conversation

    def get_conversation(self, conversation_id: str) -> ConversationHistory | None:
        return self._conversations.get(conversation_id)

    def list_conversations(
        self, user_id: str, limit: int = 20
    ) -> list[ConversationHistory]:
        return list(self._conversations.values())[:limit]

    def delete_conversation(self, conversation_id: str) -> None:
        self._conversations.pop(conversation_id, None)
