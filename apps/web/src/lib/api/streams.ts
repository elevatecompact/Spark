import { api } from "./client";

export interface Stream {
  id: string;
  title: string;
  creatorId: string;
  creatorUsername: string;
  thumbnailUrl: string;
  viewerCount: number;
  isLive: boolean;
  startedAt: string;
}

export interface StreamResponse {
  streams: Stream[];
  total: number;
}

export const streamApi = {
  list: (params?: Record<string, string>) => {
    const query = params ? `?${new URLSearchParams(params)}` : "";
    return api.get<StreamResponse>(`/streams${query}`);
  },

  get: (id: string) => api.get<Stream>(`/streams/${id}`),
};
