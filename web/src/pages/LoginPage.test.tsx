import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { afterEach, describe, expect, it, vi } from 'vitest';

import { login } from '../api/auth';
import { ApiError } from '../api/client';
import LoginPage from './LoginPage';

vi.mock('../api/auth', () => ({
  login: vi.fn(),
}));

const mockedLogin = vi.mocked(login);

function renderLogin() {
  render(
    <MemoryRouter initialEntries={['/login']}>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/" element={<div>Messenger home</div>} />
      </Routes>
    </MemoryRouter>,
  );
}

async function fillAndSubmit() {
  const user = userEvent.setup();
  await user.type(screen.getByLabelText('Email'), 'ifa@example.com');
  await user.type(screen.getByLabelText('Пароль'), 'secret123');
  await user.click(screen.getByRole('button', { name: 'Войти' }));
}

describe('LoginPage', () => {
  afterEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
  });

  it('stores the token and redirects home on success', async () => {
    mockedLogin.mockResolvedValue({
      token: 'jwt-token',
      user: { id: 1, username: 'ifa', nickname: 'IFA', email: 'ifa@example.com' },
    });

    renderLogin();
    await fillAndSubmit();

    expect(mockedLogin).toHaveBeenCalledWith({
      email: 'ifa@example.com',
      password: 'secret123',
    });
    expect(await screen.findByText('Messenger home')).toBeInTheDocument();
    expect(localStorage.getItem('token')).toBe('jwt-token');
  });

  it('shows the API error message on failure', async () => {
    mockedLogin.mockRejectedValue(new ApiError('Invalid password', 400));

    renderLogin();
    await fillAndSubmit();

    expect(await screen.findByText('Invalid password')).toBeInTheDocument();
    expect(localStorage.getItem('token')).toBeNull();
  });
});
