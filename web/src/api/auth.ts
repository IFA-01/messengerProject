import { request } from './client';
import type { LoginResponse, User } from '../types/user';

type RegisterPayload = {
  username: string;
  nickname: string;
  email: string;
  password: string;
};

type LoginPayload = {
  email: string;
  password: string;
};

export function register(payload: RegisterPayload) {
  return request<User>('/v1/users', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export function login(payload: LoginPayload) {
  return request<LoginResponse>('/v1/login', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export function getMe() {
  return request<User>('/v1/users/me');
}
