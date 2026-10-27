import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { afterEach, describe, expect, it } from 'vitest';

import App from './App';

describe('App routing', () => {
  afterEach(() => {
    localStorage.clear();
  });

  it('renders the login page on /login', () => {
    render(
      <MemoryRouter initialEntries={['/login']}>
        <App />
      </MemoryRouter>,
    );

    expect(screen.getByRole('heading', { name: 'Вход' })).toBeInTheDocument();
  });

  it('redirects unauthenticated users from / to the login page', () => {
    render(
      <MemoryRouter initialEntries={['/']}>
        <App />
      </MemoryRouter>,
    );

    expect(screen.getByRole('heading', { name: 'Вход' })).toBeInTheDocument();
  });
});
