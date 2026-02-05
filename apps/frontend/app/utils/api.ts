// TODO: DEPLOYMENT - Configure VITE_API_BASE_URL environment variable for production
// This should point to your production API gateway URL
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3001';

interface ApiError {
  message: string;
  status?: number;
}

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private getAuthHeaders(): HeadersInit {
    const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
    return headers;
  }

  async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    const headers = this.getAuthHeaders();

    try {
      const response = await fetch(url, {
        ...options,
        headers: {
          ...headers,
          ...options.headers,
        },
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        
        // TODO: SECURITY - 401 Unauthorized handler
        // Automatically clears session and redirects to login on authentication failure
        if (response.status === 401 && typeof window !== 'undefined') {
          const currentPath = window.location.pathname;
          if (currentPath !== '/login') {
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            window.location.href = '/login';
          }
        }
        
        throw {
          message: errorData.error || errorData.message || 'Request failed',
          status: response.status,
        } as ApiError;
      }

      return response.json();
    } catch (error) {
      // If error is already an ApiError, re-throw it
      if (error && typeof error === 'object' && 'status' in error) {
        throw error;
      }
      
      // Transform network errors into ApiError format
      console.error('Network error:', error);
      throw {
        message: error instanceof Error ? error.message : 'Network error occurred',
        status: 0,
      } as ApiError;
    }
  }

  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'GET' });
  }

  async post<T>(endpoint: string, body?: unknown): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: JSON.stringify(body),
    });
  }

  async delete<T>(endpoint: string, body?: unknown): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'DELETE',
      body: body ? JSON.stringify(body) : undefined,
    });
  }
}

export const api = new ApiClient(API_BASE_URL);

export type User = {
  id: string;
  username: string;
  role: string;
}

export type LoginResponse = {
  user: User;
  token: string;
}

export type CreateUserRequest = {
  username: string;
  password: string;
}

export type CreateUserResponse = {
  user: User;
  token: string;
}

export type DeleteUserRequest = {
  id: string;
}

export type DeleteUserResponse = {
  success: boolean;
}

export type ListUsersResponse = {
  users: User[];
}

export type UpdateUserRoleRequest = {
  id: string;
  role: string;
}

export type UpdateUserRoleResponse = {
  user: User;
}

// TODO: ENDPOINTS - These are the backend API endpoints
// Verify these match actual API gateway routes before deployment
export const authApi = {
  // POST /api/login - User authentication
  login: (username: string, password: string) =>
    api.post<LoginResponse>('/api/login', { username, password }),
  
  // POST /api/admin/create_user - Create new user (admin only)
  createUser: (data: CreateUserRequest) =>
    api.post<CreateUserResponse>('/api/admin/create_user', data),
  
  // DELETE /api/admin/delete_user/:id - Delete user by ID (admin only)
  deleteUser: (id: string) =>
    api.delete<DeleteUserResponse>(`/api/admin/delete_user/${encodeURIComponent(id)}`),
  
  // GET /api/admin - List all users (admin only, currently not implemented in backend)
  listUsers: () =>
    api.get<ListUsersResponse>('/api/admin'),
  
  // POST /api/admin/update_user_role - Update user role (admin only, currently not implemented in backend)
  updateUserRole: (data: UpdateUserRoleRequest) =>
    api.post<UpdateUserRoleResponse>('/api/admin/update_user_role', data),
  
  // GET /health - Backend health check
  health: () => api.get<{ status: string }>('/health'),
};
