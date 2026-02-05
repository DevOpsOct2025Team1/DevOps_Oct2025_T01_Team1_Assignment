import { describe, it, expect, beforeEach, vi } from 'vitest';
import { BrowserRouter } from 'react-router';
import Dashboard from '../app/routes/dashboard';
import { renderWithProviders } from './test-utils';
import { clearAuthCache } from '../app/utils/auth';

describe('Protected Routes', () => {
  beforeEach(() => {
    clearAuthCache();
    localStorage.clear();
  });

  it('redirects to login when not authenticated', () => {
    renderWithProviders(
      <BrowserRouter>
        <Dashboard />
      </BrowserRouter>
    );

    expect(localStorage.getItem('token')).toBeNull();
  });

  it('shows user dashboard for regular users', async () => {
    localStorage.setItem('token', 'fake-token');
    localStorage.setItem('user', JSON.stringify({
      id: '1',
      username: 'user',
      role: 'user'
    }));

    const { findByText, getByPlaceholderText } = renderWithProviders(
      <BrowserRouter>
        <Dashboard />
      </BrowserRouter>
    );

    expect(await findByText('Hello, user!')).toBeTruthy();
    expect(getByPlaceholderText('Search uploaded files')).toBeTruthy();
    expect(await findByText('Presentation.pptx')).toBeTruthy();
  });
});
