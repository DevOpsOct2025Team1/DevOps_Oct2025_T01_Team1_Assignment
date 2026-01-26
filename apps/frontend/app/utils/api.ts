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
      throw error;
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

export const authApi = {
  login: (username: string, password: string) =>
    api.post<LoginResponse>('/api/login', { username, password }),
  
  createUser: (data: CreateUserRequest) =>
    api.post<CreateUserResponse>('/api/admin/create_user', data),
  
  deleteUser: (id: string) =>
    api.delete<DeleteUserResponse>('/api/admin/delete_user', { id }),
  
  health: () => api.get<{ status: string }>('/health'),
};
