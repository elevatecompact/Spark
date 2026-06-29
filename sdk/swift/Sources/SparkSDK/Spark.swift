import Foundation

public struct User: Codable, Sendable {
    public let id: String
    public let email: String
    public let username: String
    public let displayName: String?
    public let avatarUrl: String?
    public let isCreator: Bool
    public let isVerified: Bool
    public let createdAt: String

    enum CodingKeys: String, CodingKey {
        case id, email, username
        case displayName = "display_name"
        case avatarUrl = "avatar_url"
        case isCreator = "is_creator"
        case isVerified = "is_verified"
        case createdAt = "created_at"
    }
}

public struct AuthResponse: Codable, Sendable {
    public let user: User
    public let accessToken: String
    public let refreshToken: String
    public let expiresIn: Int

    enum CodingKeys: String, CodingKey {
        case user
        case accessToken = "access_token"
        case refreshToken = "refresh_token"
        case expiresIn = "expires_in"
    }
}

public struct Stream: Codable, Sendable {
    public let id: String
    public let creatorId: String
    public let title: String
    public let description: String?
    public let thumbnailUrl: String?
    public let status: String
    public let isLive: Bool
    public let viewerCount: Int
    public let createdAt: String

    enum CodingKeys: String, CodingKey {
        case id, title, description, status
        case creatorId = "creator_id"
        case thumbnailUrl = "thumbnail_url"
        case isLive = "is_live"
        case viewerCount = "viewer_count"
        case createdAt = "created_at"
    }
}

public actor Spark {
    private let baseURL: String
    private var accessToken: String?
    private let session: URLSession

    public init(baseURL: String = "https://api.spark.dev/api/v1", accessToken: String? = nil) {
        self.baseURL = baseURL
        self.accessToken = accessToken
        let config = URLSessionConfiguration.default
        config.timeoutIntervalForRequest = 30
        self.session = URLSession(configuration: config)
    }

    public func setAccessToken(_ token: String) {
        self.accessToken = token
    }

    private func request<T: Decodable>(_ method: String, _ path: String, body: Encodable? = nil) async throws -> T {
        guard let url = URL(string: "\(baseURL)\(path)") else {
            throw SparkError.invalidURL
        }

        var req = URLRequest(url: url)
        req.httpMethod = method
        req.setValue("application/json", forHTTPHeaderField: "Content-Type")
        if let token = accessToken {
            req.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        if let body = body {
            req.httpBody = try JSONEncoder().encode(AnyEncodable(body))
        }

        let (data, response) = try await session.data(for: req)
        guard let httpResp = response as? HTTPURLResponse, (200...299).contains(httpResp.statusCode) else {
            throw SparkError.httpError(response)
        }

        return try JSONDecoder().decode(T.self, from: data)
    }

    public func register(email: String, username: String, password: String) async throws -> AuthResponse {
        let body = ["email": email, "username": username, "password": password]
        var result: AuthResponse = try await request("POST", "/auth/register", body: body)
        if !result.accessToken.isEmpty {
            self.accessToken = result.accessToken
        }
        return result
    }

    public func login(email: String, password: String) async throws -> AuthResponse {
        let body = ["email": email, "password": password]
        var result: AuthResponse = try await request("POST", "/auth/login", body: body)
        if !result.accessToken.isEmpty {
            self.accessToken = result.accessToken
        }
        return result
    }

    public func me() async throws -> User {
        try await request("GET", "/users/me")
    }

    public func getStream(id: String) async throws -> Stream {
        try await request("GET", "/streams/\(id)")
    }
}

enum SparkError: Error {
    case invalidURL
    case httpError(URLResponse)
}

struct AnyEncodable: Encodable {
    private let value: Encodable
    init(_ value: Encodable) { self.value = value }
    func encode(to encoder: Encoder) throws {
        try value.encode(to: encoder)
    }
}
