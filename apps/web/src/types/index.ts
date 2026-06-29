export interface User {
  id: string;
  email: string;
  username: string;
  displayName: string;
  avatarUrl: string;
  bio: string;
  isCreator: boolean;
  createdAt: string;
}

export interface Stream {
  id: string;
  title: string;
  description: string;
  creatorId: string;
  creatorUsername: string;
  thumbnailUrl: string;
  viewerCount: number;
  isLive: boolean;
  category: string;
  tags: string[];
  startedAt: string;
}

export interface ChatMessage {
  id: string;
  streamId: string;
  userId: string;
  username: string;
  content: string;
  isModerator: boolean;
  isSubscriber: boolean;
  createdAt: string;
}

export interface Creator {
  id: string;
  userId: string;
  username: string;
  displayName: string;
  avatarUrl: string;
  bannerUrl: string;
  bio: string;
  followerCount: number;
  subscriberCount: number;
  isVerified: boolean;
}

export interface Community {
  id: string;
  name: string;
  description: string;
  avatarUrl: string;
  memberCount: number;
  isPrivate: boolean;
}

export interface Event {
  id: string;
  title: string;
  description: string;
  type: "virtual" | "inperson" | "hybrid";
  startDate: string;
  endDate: string;
  creatorId: string;
  attendeeCount: number;
}
