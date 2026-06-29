package dev.spark;

import com.google.gson.Gson;
import com.google.gson.reflect.TypeToken;
import okhttp3.*;

import java.io.IOException;
import java.lang.reflect.Type;
import java.util.Map;
import java.util.concurrent.TimeUnit;

public class Spark {
    private static final String DEFAULT_BASE_URL = "https://api.spark.dev/api/v1";

    private final OkHttpClient client;
    private final Gson gson;
    private final String baseUrl;
    private String accessToken;

    public Spark() {
        this(DEFAULT_BASE_URL, null);
    }

    public Spark(String baseUrl, String accessToken) {
        this.baseUrl = baseUrl;
        this.accessToken = accessToken;
        this.gson = new Gson();
        this.client = new OkHttpClient.Builder()
                .connectTimeout(30, TimeUnit.SECONDS)
                .readTimeout(30, TimeUnit.SECONDS)
                .build();
    }

    public void setAccessToken(String token) {
        this.accessToken = token;
    }

    private Request.Builder auth() {
        Request.Builder builder = new Request.Builder();
        if (accessToken != null) {
            builder.header("Authorization", "Bearer " + accessToken);
        }
        builder.header("Content-Type", "application/json");
        return builder;
    }

    private <T> T request(String method, String path, Object body, Type type) throws IOException {
        String url = baseUrl + path;
        Request.Builder builder = auth().url(url);

        if (body != null) {
            String json = gson.toJson(body);
            builder.method(method, RequestBody.create(json, MediaType.get("application/json")));
        } else {
            builder.method(method, null);
        }

        try (Response response = client.newCall(builder.build()).execute()) {
            if (!response.isSuccessful()) {
                throw new IOException("HTTP " + response.code() + ": " + response.body().string());
            }
            String responseBody = response.body().string();
            return gson.fromJson(responseBody, type);
        }
    }

    public AuthResponse register(String email, String username, String password) throws IOException {
        Type type = new TypeToken<AuthResponse>() {}.getType();
        return request("POST", "/auth/register",
                Map.of("email", email, "username", username, "password", password), type);
    }

    public AuthResponse login(String email, String password) throws IOException {
        Type type = new TypeToken<AuthResponse>() {}.getType();
        AuthResponse result = request("POST", "/auth/login",
                Map.of("email", email, "password", password), type);
        if (result.accessToken != null) {
            this.accessToken = result.accessToken;
        }
        return result;
    }

    public User me() throws IOException {
        Type type = new TypeToken<User>() {}.getType();
        return request("GET", "/users/me", null, type);
    }

    public static class User {
        public String id;
        public String email;
        public String username;
        public String displayName;
        public String avatarUrl;
        public boolean isCreator;
        public boolean isVerified;
        public String createdAt;
    }

    public static class AuthResponse {
        public User user;
        public String accessToken;
        public String refreshToken;
        public int expiresIn;
    }
}
