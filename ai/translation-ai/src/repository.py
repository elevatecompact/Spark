import logging

from src.models import SupportedLanguage

logger = logging.getLogger(__name__)


class TranslationRepository:

    async def get_supported_languages(self) -> list[SupportedLanguage]:
        return [
            SupportedLanguage(code="en", name="English", native_name="English", is_rtl=False),
            SupportedLanguage(code="es", name="Spanish", native_name="Español", is_rtl=False),
            SupportedLanguage(code="fr", name="French", native_name="Français", is_rtl=False),
            SupportedLanguage(code="de", name="German", native_name="Deutsch", is_rtl=False),
            SupportedLanguage(code="zh", name="Chinese", native_name="中文", is_rtl=False),
            SupportedLanguage(code="ja", name="Japanese", native_name="日本語", is_rtl=False),
            SupportedLanguage(code="ko", name="Korean", native_name="한국어", is_rtl=False),
            SupportedLanguage(code="ar", name="Arabic", native_name="العربية", is_rtl=True),
            SupportedLanguage(code="hi", name="Hindi", native_name="हिन्दी", is_rtl=False),
            SupportedLanguage(code="pt", name="Portuguese", native_name="Português", is_rtl=False),
            SupportedLanguage(code="ru", name="Russian", native_name="Русский", is_rtl=False),
            SupportedLanguage(code="it", name="Italian", native_name="Italiano", is_rtl=False),
            SupportedLanguage(code="nl", name="Dutch", native_name="Nederlands", is_rtl=False),
            SupportedLanguage(code="pl", name="Polish", native_name="Polski", is_rtl=False),
            SupportedLanguage(code="tr", name="Turkish", native_name="Türkçe", is_rtl=False),
        ]

    async def log_translation(self, source_lang: str, target_lang: str, char_count: int) -> None:
        logger.info(
            "Translation logged: %s -> %s (%d chars)", source_lang, target_lang, char_count
        )
