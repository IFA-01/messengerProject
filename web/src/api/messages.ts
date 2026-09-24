import { request } from './client';
import type { Message } from '../types/message';

type CreateMessagePayload = {
  chat_id: number;
  content: string;
};

type MessagesResponse = {
  messages: Message[];
};

export function listMessages(chatId: number) {
  return request<MessagesResponse>(`/v1/messages/${chatId}`);
}

export function sendMessage(payload: CreateMessagePayload) {
  return request<Message>('/v1/messages', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}
