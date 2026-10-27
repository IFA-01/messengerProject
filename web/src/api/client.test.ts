import { afterEach, describe, expect, it, vi } from 'vitest';

import { ApiError, request } from './client';

function mockFetch(response: Response) {
  const fetchMock = vi.fn().mockResolvedValue(response);
  vi.stubGlobal('fetch', fetchMock);
  return fetchMock;
}

describe('request', () => {
  afterEach(() => {
    vi.unstubAllGlobals();
    localStorage.clear();
  });

  it('attaches the bearer token when one is stored', async () => {
    localStorage.setItem('token', 'abc123');
    const fetchMock = mockFetch(new Response(JSON.stringify({ ok: true }), { status: 200 }));

    await request('/v1/users/me');

    const [, options] = fetchMock.mock.calls[0] as [string, RequestInit];
    const headers = options.headers as Headers;
    expect(headers.get('Authorization')).toBe('Bearer abc123');
    expect(headers.get('Content-Type')).toBe('application/json');
  });

  it('does not send Authorization without a token', async () => {
    const fetchMock = mockFetch(new Response('{}', { status: 200 }));

    await request('/v1/login');

    const [, options] = fetchMock.mock.calls[0] as [string, RequestInit];
    expect((options.headers as Headers).has('Authorization')).toBe(false);
  });

  it('returns parsed JSON on success', async () => {
    mockFetch(new Response(JSON.stringify({ chats: [] }), { status: 200 }));

    await expect(request('/v1/chats')).resolves.toEqual({ chats: [] });
  });

  it('throws ApiError with the status code on failure', async () => {
    mockFetch(new Response('Unauthorized', { status: 401 }));

    const error = await request('/v1/chats').catch((e: unknown) => e);

    expect(error).toBeInstanceOf(ApiError);
    expect((error as ApiError).status).toBe(401);
    expect((error as ApiError).message).toBe('Unauthorized');
  });
});
