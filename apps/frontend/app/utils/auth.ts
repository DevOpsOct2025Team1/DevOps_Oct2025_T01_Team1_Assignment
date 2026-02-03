export type User = {
  id: string;
  username: string;
  role: string | number;
};

let cachedUser: User | null | undefined = undefined;
let cachedToken: string | null | undefined = undefined;

function isBrowser(): boolean {
  return typeof window !== "undefined";
}

export function getStoredUser(): User | null {
  if (!isBrowser()) return null;
  if (cachedUser !== undefined) return cachedUser as User | null;

  const userStr = localStorage.getItem("user");
  if (!userStr) {
    cachedUser = null;
    return null;
  }
  try {
    cachedUser = JSON.parse(userStr) as User;
    return cachedUser;
  } catch {
    cachedUser = null;
    return null;
  }
}

export function getStoredToken(): string | null {
  if (!isBrowser()) return null;
  if (cachedToken !== undefined) return cachedToken;

  cachedToken = localStorage.getItem("token");
  return cachedToken;
}

export function setAuth(user: User, token: string): void {
  if (!isBrowser()) return;
  localStorage.setItem("user", JSON.stringify(user));
  localStorage.setItem("token", token);
  cachedUser = user;
  cachedToken = token;
}

export function clearAuth(): void {
  if (!isBrowser()) return;
  localStorage.removeItem("user");
  localStorage.removeItem("token");
  cachedUser = null;
  cachedToken = null;
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
