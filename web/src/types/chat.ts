export type Chat = {
  id: number;
  name: string;
  is_group: boolean;
  created_at?: string;
  updated_at?: string;
};

export type CreateChatRequest = {
  username: string;
  name: string;
};