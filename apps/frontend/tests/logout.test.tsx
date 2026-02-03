import { render, screen } from '@testing-library/react';
import { describe, it, expect, beforeEach } from 'vitest';
import { BrowserRouter } from 'react-router';
import Logout from '../app/routes/logout';
import { clearAuthCache } from '../app/utils/auth';

describe('Logout', () => {
  beforeEach(() => {
    clearAuthCache();
    localStorage.setItem('token', 'fake-token');
    localStorage.setItem('user', JSON.stringify({
      id: '1',
      username: 'testuser',
      role: 'user'
    }));
  });

  it('clears localStorage on logout', () => {
    render(
      <BrowserRouter>
        <Logout />
      </BrowserRouter>
    );

    expect(localStorage.getItem('token')).toBeNull();
    expect(localStorage.getItem('user')).toBeNull();
  });

  it('shows logging out message', () => {
    render(
      <BrowserRouter>
        <Logout />
      </BrowserRouter>
    );

    expect(screen.getByText(/logging out/i)).toBeDefined();
  });
});
