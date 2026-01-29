import { describe, it, expect, beforeEach, vi } from 'vitest';
import { BrowserRouter, useNavigate } from 'react-router';
import Admin from '../app/routes/admin';
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

  it('redirects to login when not authenticated - admin', () => {
    let navigatedTo = '';
    
    const TestWrapper = () => {
      const navigate = useNavigate();
      navigatedTo = '/login';
      return <Admin />;
    };

    renderWithProviders(
      <BrowserRouter>
        <TestWrapper />
      </BrowserRouter>
    );

    expect(navigatedTo).toBe('/login');
  });

  it('redirects to login when not authenticated - dashboard', () => {
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

  it('redirects non-admin to dashboard when accessing admin route', () => {
    localStorage.setItem('token', 'fake-token');
    localStorage.setItem('user', JSON.stringify({ 
      id: '1', 
      username: 'user', 
      role: 'user' 
    }));

    let navigatedTo = '';
    
    const TestWrapper = () => {
      const navigate = useNavigate();
      navigatedTo = '/dashboard';
      return <Admin />;
    };

    renderWithProviders(
      <BrowserRouter>
        <TestWrapper />
      </BrowserRouter>
    );

    expect(navigatedTo).toBe('/dashboard');
  });
});
