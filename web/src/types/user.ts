export type User = {
  id: number;
  username: string;
  nickname: string;
  email: string;
  last_seen?: string;
  created_at?: string;
  updated_at?: string;
};

export type LoginResponse = {
  token: string;
  user: User;
};
