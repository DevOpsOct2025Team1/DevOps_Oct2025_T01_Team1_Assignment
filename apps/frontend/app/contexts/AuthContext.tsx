import { createContext, useContext, useState, type ReactNode } from "react";
import {
  getStoredUser,
  setAuth as setAuthStorage,
  clearAuth as clearAuthStorage,
  type User,
} from "../utils/auth";

type AuthContextType = {
  user: User | null;
  setAuth: (user: User, token: string) => void;
  clearAuth: () => void;
  isAuthenticated: boolean;
};

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(() => getStoredUser());

  const handleSetAuth = (user: User, token: string) => {
    setAuthStorage(user, token);
    setUser(user);
  };

  const handleClearAuth = () => {
    clearAuthStorage();
    setUser(null);
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        setAuth: handleSetAuth,
        clearAuth: handleClearAuth,
        isAuthenticated: !!user,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}