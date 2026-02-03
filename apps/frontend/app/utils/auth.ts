export type User = {
  id: string;
  username: string;
  role: string | number;
};

function isBrowser(): boolean {
  return typeof window !== "undefined";
}

export function getStoredUser(): User | null {
  if (!isBrowser()) return null;
  const userStr = localStorage.getItem("user");
  if (!userStr) return null;
  try {
    return JSON.parse(userStr);
  } catch {
    return null;
  }
}

export function getStoredToken(): string | null {
  if (!isBrowser()) return null;
  return localStorage.getItem("token");
}

export function setAuth(user: User, token: string): void {
  if (!isBrowser()) return;
  localStorage.setItem("user", JSON.stringify(user));
  localStorage.setItem("token", token);
}

export function clearAuth(): void {
  if (!isBrowser()) return;
  localStorage.removeItem("user");
  localStorage.removeItem("token");
}

export function isAuthenticated(): boolean {
  return !!getStoredToken();
}

export function isAdmin(): boolean {
  const user = getStoredUser();
  if (!user) {
    return false;
  }

  if (typeof user.role === "number") {
    return user.role === 2;
  }

  return user.role.toLowerCase().includes("admin");
}
