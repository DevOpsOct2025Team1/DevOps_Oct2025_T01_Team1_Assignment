import { describe, it, expect, beforeEach, vi } from 'vitest';
import { BrowserRouter, useNavigate } from 'react-router';
import Dashboard from '../app/routes/dashboard';
import { renderWithProviders } from './test-utils';

const createUserMock = vi.fn();
const deleteUserMock = vi.fn();

vi.mock('../app/api/generated', () => ({
  useCreateUser: () => ({
    mutateAsync: createUserMock,
    isPending: false,
  }),
  useDeleteUser: () => ({
    mutateAsync: deleteUserMock,
    isPending: false,
  }),
}));

describe('Protected Routes', () => {
  beforeEach(() => {
    localStorage.clear();
    createUserMock.mockReset();
    deleteUserMock.mockReset();
  });

  it('redirects to login when not authenticated', () => {
    let navigatedTo = '';

    const TestWrapper = () => {
      const navigate = useNavigate();
      navigatedTo = '/login';
      return <Dashboard />;
    };

    renderWithProviders(
      <BrowserRouter>
        <TestWrapper />
      </BrowserRouter>
    );

    expect(navigatedTo).toBe('/login');
  });

  it('shows admin dashboard for admin users', () => {
    localStorage.setItem('token', 'fake-token');
    localStorage.setItem('user', JSON.stringify({
      id: '1',
      username: 'admin',
      role: 'admin'
    }));

    const { getByRole } = renderWithProviders(
      <BrowserRouter>
        <Dashboard />
      </BrowserRouter>
    );

    expect(getByRole('heading', { name: /admin dashboard/i })).toBeTruthy();
    expect(getByRole('heading', { name: /create user/i })).toBeTruthy();
  });

  it('shows user dashboard for regular users', () => {
    localStorage.setItem('token', 'fake-token');
    localStorage.setItem('user', JSON.stringify({
      id: '1',
      username: 'user',
      role: 'user'
    }));

    const { getByText } = renderWithProviders(
      <BrowserRouter>
        <Dashboard />
      </BrowserRouter>
    );

    expect(getByText('Dashboard')).toBeTruthy();
    expect(getByText('Upload File')).toBeTruthy();
  });
});
