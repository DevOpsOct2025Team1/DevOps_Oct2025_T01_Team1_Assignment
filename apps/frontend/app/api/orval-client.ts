import { getStoredToken, clearAuth } from "../utils/auth";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:3001";

export type ApiError = {
  message: string;
  status?: number;
};

const isBrowser = typeof window !== "undefined";

export const getAuthHeaders = (): HeadersInit => {
  const token = getStoredToken();
  const headers: HeadersInit = {
    "Content-Type": "application/json",
  };

  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  return headers;
};

export const resolveUrl = (url: string): string => {
  if (url.startsWith("http://") || url.startsWith("https://")) {
    return url;
  }

  const trimmedBaseUrl = API_BASE_URL.endsWith("/")
    ? API_BASE_URL.slice(0, -1) // get rid of trailing /
    : API_BASE_URL;
  const normalizedPath = url.startsWith("/") ? url : `/${url}`;

  return `${trimmedBaseUrl}${normalizedPath}`;
};

export const customFetch = async <T>(
  url: string,
  options: RequestInit = {}
): Promise<T> => {
  // Start with default auth headers
  const headers = new Headers(getAuthHeaders());

  // Apply any headers provided in options, regardless of their HeadersInit form
  if (options.headers) {
    const optionHeaders = new Headers(options.headers);
    optionHeaders.forEach((value, key) => {
      headers.set(key, value);
    });
  }

  // Let the browser set the correct Content-Type boundary for FormData
  if (options.body instanceof FormData) {
    headers.delete("Content-Type");
  }

  const response = await fetch(resolveUrl(url), {
    ...options,
    headers,
  });

  let data: unknown = undefined;

  if (response.status !== 204) {
    try {
      data = await response.json();
    } catch {
      data = undefined;
    }
  }

  if (!response.ok) {
    const errorData = data as { error?: string; message?: string } | undefined;

    if (response.status === 401 && isBrowser) {
      const currentPath = window.location.pathname;
      if (currentPath !== "/login") {
        clearAuth();
        window.location.href = "/login";
      }
    }

    throw {
      message: errorData?.error || errorData?.message || "Request failed",
      status: response.status,
    } satisfies ApiError;
  }

  return {
    data,
    status: response.status,
    headers: response.headers,
  } as T;
};
