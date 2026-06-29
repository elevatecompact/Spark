using System.Net.Http.Headers;
using System.Text;
using System.Text.Json;
using System.Text.Json.Serialization;

namespace Spark.SDK;

public class SparkClient
{
    private static readonly JsonSerializerOptions JsonOptions = new()
    {
        PropertyNamingPolicy = JsonNamingPolicy.SnakeCaseLower,
        DefaultIgnoreCondition = JsonIgnoreCondition.WhenWritingNull,
    };

    private readonly HttpClient _http;
    private readonly string _baseUrl;

    public SparkClient(string baseUrl = "https://api.spark.dev/api/v1", string? accessToken = null)
    {
        _baseUrl = baseUrl.TrimEnd('/');
        _http = new HttpClient();
        if (accessToken != null)
        {
            _http.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", accessToken);
        }
        _http.DefaultRequestHeaders.Accept.Add(new MediaTypeWithQualityHeaderValue("application/json"));
    }

    public void SetAccessToken(string token)
    {
        _http.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", token);
    }

    private async Task<T> RequestAsync<T>(HttpMethod method, string path, object? body = null)
    {
        var request = new HttpRequestMessage(method, $"{_baseUrl}{path}");
        if (body != null)
        {
            var json = JsonSerializer.Serialize(body, JsonOptions);
            request.Content = new StringContent(json, Encoding.UTF8, "application/json");
        }

        var response = await _http.SendAsync(request);
        response.EnsureSuccessStatusCode();

        var content = await response.Content.ReadAsStringAsync();
        return JsonSerializer.Deserialize<T>(content, JsonOptions)!;
    }

    public Task<AuthResponse> RegisterAsync(string email, string username, string password)
    {
        return RequestAsync<AuthResponse>(HttpMethod.Post, "/auth/register",
            new { email, username, password });
    }

    public async Task<AuthResponse> LoginAsync(string email, string password)
    {
        var result = await RequestAsync<AuthResponse>(HttpMethod.Post, "/auth/login",
            new { email, password });
        if (!string.IsNullOrEmpty(result.AccessToken))
        {
            SetAccessToken(result.AccessToken);
        }
        return result;
    }

    public Task<User> MeAsync()
    {
        return RequestAsync<User>(HttpMethod.Get, "/users/me");
    }
}

public class AuthResponse
{
    [JsonPropertyName("access_token")]
    public string AccessToken { get; set; } = "";

    [JsonPropertyName("refresh_token")]
    public string RefreshToken { get; set; } = "";

    [JsonPropertyName("expires_in")]
    public int ExpiresIn { get; set; }

    public User User { get; set; } = new();
}

public class User
{
    public string Id { get; set; } = "";
    public string Email { get; set; } = "";
    public string Username { get; set; } = "";
    public string? DisplayName { get; set; }
    public string? AvatarUrl { get; set; }
    public bool IsCreator { get; set; }
    public bool IsVerified { get; set; }
    public string CreatedAt { get; set; } = "";
}
