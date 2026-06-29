import { ofetch, type FetchOptions } from "ofetch";

export interface SparkConfig {
  baseURL?: string;
  accessToken?: string;
}

export interface User {
  id: string;
  email: string;
  username: string;
  displayName: string | null;
  avatarUrl: string | null;
  isCreator: boolean;
  isVerified: boolean;
  createdAt: string;
}

export interface AuthResponse {
  user: User;
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
}

export class Spark {
  private client: ReturnType<typeof ofetch>;

  constructor(config: SparkConfig = {}) {
    const baseURL = config.baseURL || "https://api.spark.dev/api/v1";

    this.client = ofetch.create({
      baseURL,
      headers: config.accessToken
        ? { Authorization: `Bearer ${config.accessToken}` }
        : {},
    });
  }

  setAccessToken(token: string) {
    this.client = ofetch.create({
      baseURL: this.client.defaults.baseURL,
      headers: { Authorization: `Bearer ${token}` },
    });
  }

  // Auth
  async register(
    email: string,
    username: string,
    password: string,
  ): Promise<AuthResponse> {
    return this.client("/auth/register", {
      method: "POST",
      body: { email, username, password },
    });
  }

  async login(email: string, password: string): Promise<AuthResponse> {
    return this.client("/auth/login", {
      method: "POST",
      body: { email, password },
    });
  }

  // Users
  async getMe(): Promise<User> {
    return this.client("/users/me");
  }

  async getUser(id: string): Promise<User> {
    return this.client(`/users/${id}`);
  }

  // Streams
  async getStreams(params?: {
    page?: number;
    pageSize?: number;
    isLive?: boolean;
    category?: string;
  }) {
    return this.client("/streams", { params });
  }

  async getStream(id: string) {
    return this.client(`/streams/${id}`);
  }

  async createStream(data: {
    title: string;
    description?: string;
    category?: string;
    tags?: string[];
  }) {
    return this.client("/streams", { method: "POST", body: data });
  }

  // Wallet
  async getBalance() {
    return this.client("/wallet/balance");
  }

  async getTransactions(params?: {
    page?: number;
    pageSize?: number;
    type?: string;
    status?: string;
  }) {
    return this.client("/wallet/transactions", { params });
  }

  // Notifications
  async getNotifications(params?: {
    page?: number;
    pageSize?: number;
    unreadOnly?: boolean;
  }) {
    return this.client("/notifications", { params });
  }

  // Search
  async search(query: string, params?: {
    type?: string;
    page?: number;
    pageSize?: number;
  }) {
    return this.client("/search", { params: { q: query, ...params } });
  }

  // Subscriptions
  async getPlans(creatorId: string) {
    return this.client(`/plans`, { params: { creator_id: creatorId } });
  }

  async subscribe(planId: string, paymentMethodId: string) {
    return this.client("/subscriptions", {
      method: "POST",
      body: { plan_id: planId, payment_method_id: paymentMethodId },
    });
  }
}

export default Spark;
