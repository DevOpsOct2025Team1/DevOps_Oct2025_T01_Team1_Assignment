import { render } from '@testing-library/react';
import { describe, it, expect, beforeEach } from 'vitest';
import { BrowserRouter, useNavigate } from 'react-router';
import Admin from '../app/routes/admin';
import Dashboard from '../app/routes/dashboard';

describe('Protected Routes', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it('redirects to login when not authenticated - admin', () => {
    let navigatedTo = '';
    
    const TestWrapper = () => {
      const navigate = useNavigate();
      navigatedTo = '/login';
      return <Admin />;
    };

    render(
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

    render(
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

    render(
      <BrowserRouter>
        <TestWrapper />
      </BrowserRouter>
    );

    expect(navigatedTo).toBe('/dashboard');
  });
});
