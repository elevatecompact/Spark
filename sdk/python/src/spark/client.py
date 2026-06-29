from __future__ import annotations

from typing import Any

import httpx


class Spark:
    def __init__(
        self,
        base_url: str = "https://api.spark.dev/api/v1",
        access_token: str | None = None,
    ):
        self.base_url = base_url.rstrip("/")
        headers: dict[str, str] = {}
        if access_token:
            headers["Authorization"] = f"Bearer {access_token}"
        self.client = httpx.AsyncClient(base_url=self.base_url, headers=headers)

    async def close(self) -> None:
        await self.client.aclose()

    def _headers(self) -> dict[str, str]:
        return dict(self.client.headers)

    async def _request(
        self,
        method: str,
        path: str,
        **kwargs: Any,
    ) -> Any:
        url = f"{self.base_url}{path}"
        response = await self.client.request(method, url, **kwargs)
        response.raise_for_status()
        return response.json()

    # Auth
    async def register(
        self,
        email: str,
        username: str,
        password: str,
    ) -> dict[str, Any]:
        return await self._request(
            "POST",
            "/auth/register",
            json={"email": email, "username": username, "password": password},
        )

    async def login(
        self,
        email: str,
        password: str,
    ) -> dict[str, Any]:
        result = await self._request(
            "POST",
            "/auth/login",
            json={"email": email, "password": password},
        )
        if "access_token" in result:
            self.client.headers["Authorization"] = f"Bearer {result['access_token']}"
        return result

    async def set_access_token(self, token: str) -> None:
        self.client.headers["Authorization"] = f"Bearer {token}"

    # Users
    async def me(self) -> dict[str, Any]:
        return await self._request("GET", "/users/me")

    async def get_user(self, user_id: str) -> dict[str, Any]:
        return await self._request("GET", f"/users/{user_id}")

    # Streams
    async def list_streams(
        self,
        page: int = 1,
        page_size: int = 20,
        is_live: bool | None = None,
        category: str | None = None,
    ) -> dict[str, Any]:
        params: dict[str, Any] = {"page": page, "page_size": page_size}
        if is_live is not None:
            params["is_live"] = str(is_live).lower()
        if category:
            params["category"] = category
        return await self._request("GET", "/streams", params=params)

    async def get_stream(self, stream_id: str) -> dict[str, Any]:
        return await self._request("GET", f"/streams/{stream_id}")

    async def create_stream(
        self,
        title: str,
        description: str | None = None,
        category: str | None = None,
        tags: list[str] | None = None,
    ) -> dict[str, Any]:
        body: dict[str, Any] = {"title": title}
        if description:
            body["description"] = description
        if category:
            body["category"] = category
        if tags:
            body["tags"] = tags
        return await self._request("POST", "/streams", json=body)

    # Wallet
    async def get_balance(self) -> dict[str, Any]:
        return await self._request("GET", "/wallet/balance")

    async def get_transactions(
        self,
        page: int = 1,
        page_size: int = 20,
        type: str | None = None,
        status: str | None = None,
    ) -> dict[str, Any]:
        params: dict[str, Any] = {"page": page, "page_size": page_size}
        if type:
            params["type"] = type
        if status:
            params["status"] = status
        return await self._request("GET", "/wallet/transactions", params=params)

    # Notifications
    async def get_notifications(
        self,
        page: int = 1,
        page_size: int = 20,
        unread_only: bool | None = None,
    ) -> dict[str, Any]:
        params: dict[str, Any] = {"page": page, "page_size": page_size}
        if unread_only is not None:
            params["unread_only"] = str(unread_only).lower()
        return await self._request("GET", "/notifications", params=params)

    # Search
    async def search(
        self,
        query: str,
        type: str | None = None,
        page: int = 1,
        page_size: int = 20,
    ) -> dict[str, Any]:
        params: dict[str, Any] = {"q": query, "page": page, "page_size": page_size}
        if type:
            params["type"] = type
        return await self._request("GET", "/search", params=params)
