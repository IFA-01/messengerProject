import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { getMe } from '../api/auth';
import { listChats } from '../api/chats';
import MessengerPage from './MessengerPage';

vi.mock('../api/auth', () => ({ getMe: vi.fn() }));
vi.mock('../api/chats', () => ({ listChats: vi.fn(), createChat: vi.fn() }));
vi.mock('../api/messages', () => ({
  listMessages: vi.fn().mockResolvedValue({ messages: [] }),
  sendMessage: vi.fn(),
}));

function renderMessenger() {
  render(
    <MemoryRouter initialEntries={['/']}>
      <Routes>
        <Route path="/" element={<MessengerPage />} />
        <Route path="/login" element={<div>Login page</div>} />
      </Routes>
    </MemoryRouter>,
  );
}

describe('MessengerPage', () => {
  beforeEach(() => {
    localStorage.setItem('token', 'jwt-token');
    vi.mocked(getMe).mockResolvedValue({
      id: 1,
      username: 'IFA',
      nickname: 'IFA',
      email: 'ifa@example.com',
    });
    vi.mocked(listChats).mockResolvedValue({
      chats: [
        { id: 10, name: 'IFA-alice', is_group: false },
        { id: 11, name: 'bob-IFA', is_group: false },
      ],
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
  });

  it('shows the peer name as the chat title', async () => {
    renderMessenger();

    expect(await screen.findByText('alice')).toBeInTheDocument();
    expect(screen.getByText('bob')).toBeInTheDocument();
    expect(screen.queryByText('IFA-alice')).not.toBeInTheDocument();

    await userEvent.click(screen.getByText('alice'));
    expect(screen.getAllByText('alice')).toHaveLength(2);
  });

  it('collapses and restores the chat list', async () => {
    renderMessenger();
    await screen.findByText('alice');
    const messenger = screen.getByTestId('messenger');

    await userEvent.click(screen.getByRole('button', { name: 'Скрыть список чатов' }));
    expect(messenger).toHaveClass('sidebar-collapsed');

    await userEvent.click(screen.getByRole('button', { name: 'Показать список чатов' }));
    expect(messenger).not.toHaveClass('sidebar-collapsed');
  });

  it('keeps the chat list open after selecting a chat on desktop', async () => {
    renderMessenger();

    await userEvent.click(await screen.findByText('alice'));

    expect(screen.getByTestId('messenger')).not.toHaveClass('sidebar-collapsed');
  });

  it('hides the chat list after selecting a chat on mobile', async () => {
    vi.stubGlobal(
      'matchMedia',
      vi.fn().mockImplementation((query: string) => ({
        matches: query === '(max-width: 720px)',
        addEventListener: vi.fn(),
        removeEventListener: vi.fn(),
      })),
    );
    renderMessenger();

    await userEvent.click(await screen.findByText('alice'));

    expect(screen.getByTestId('messenger')).toHaveClass('sidebar-collapsed');
    vi.unstubAllGlobals();
  });

  it('logs out by clearing the token and redirecting to login', async () => {
    renderMessenger();

    await userEvent.click(await screen.findByRole('button', { name: 'Выйти' }));

    expect(await screen.findByText('Login page')).toBeInTheDocument();
    expect(localStorage.getItem('token')).toBeNull();
  });
});
