package dev.spark

import io.ktor.client.*
import io.ktor.client.call.*
import io.ktor.client.plugins.contentnegotiation.*
import io.ktor.client.request.*
import io.ktor.http.*
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.Json

@Serializable
data class User(
    val id: String,
    val email: String,
    val username: String,
    @kotlinx.serialization.SerialName("display_name") val displayName: String? = null,
    @kotlinx.serialization.SerialName("avatar_url") val avatarUrl: String? = null,
    @kotlinx.serialization.SerialName("is_creator") val isCreator: Boolean = false,
    @kotlinx.serialization.SerialName("is_verified") val isVerified: Boolean = false,
    @kotlinx.serialization.SerialName("created_at") val createdAt: String,
)

@Serializable
data class AuthResponse(
    val user: User,
    @kotlinx.serialization.SerialName("access_token") val accessToken: String,
    @kotlinx.serialization.SerialName("refresh_token") val refreshToken: String,
    @kotlinx.serialization.SerialName("expires_in") val expiresIn: Int,
)

class Spark(
    private val baseUrl: String = "https://api.spark.dev/api/v1",
    accessToken: String? = null,
) {
    private var token: String? = accessToken

    private val client = HttpClient {
        install(ContentNegotiation) {
            json(Json { ignoreUnknownKeys = true })
        }
    }

    fun setAccessToken(token: String) {
        this.token = token
    }

    private suspend fun <T> request(
        method: HttpMethod,
        path: String,
        body: Any? = null,
    ): T {
        val response = client.request("$baseUrl$path") {
            this.method = method
            token?.let { header(HttpHeaders.Authorization, "Bearer $it") }
            contentType(ContentType.Application.Json)
            if (body != null) {
                setBody(body)
            }
        }
        return response.body()
    }

    suspend fun register(email: String, username: String, password: String): AuthResponse {
        val result: AuthResponse = request(
            HttpMethod.Post,
            "/auth/register",
            mapOf("email" to email, "username" to username, "password" to password),
        )
        token = result.accessToken
        return result
    }

    suspend fun login(email: String, password: String): AuthResponse {
        val result: AuthResponse = request(
            HttpMethod.Post,
            "/auth/login",
            mapOf("email" to email, "password" to password),
        )
        token = result.accessToken
        return result
    }

    suspend fun me(): User {
        return request(HttpMethod.Get, "/users/me")
    }

    fun close() {
        client.close()
    }
}
