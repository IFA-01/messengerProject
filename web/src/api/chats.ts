import { request } from './client';
import type { Chat } from '../types/chat';

type CreateChatPayload = {
  username: string;
  name: string;
};

type ChatsResponse = {
  chats: Chat[];
}

export function createChat(payload: CreateChatPayload) {
  return request<Chat>('/v1/chats', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export function listChats() {
  return request<ChatsResponse>('/v1/chats',{
    method: 'GET',
  });
}