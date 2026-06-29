import { api } from "./client";

export interface AuthResponse {
  token: string;
  user: {
    id: string;
    email: string;
    username: string;
  };
}

export const authApi = {
  login: (email: string, password: string) =>
    api.post<AuthResponse>("/auth/login", { email, password }),

  register: (email: string, username: string, password: string) =>
    api.post<AuthResponse>("/auth/register", { email, username, password }),

  me: () => api.get<AuthResponse["user"]>("/auth/me"),
};
