import { useLocation } from "react-router";
import { useState, useEffect } from "react";
import { useAuth } from "../contexts/AuthContext";

export default function Navigation() {
  const location = useLocation();
  const { user } = useAuth();
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted || location.pathname === "/login") {
    return null;
  }

  if (!user) {
    return null;
  }

  const userRole = typeof user.role === "number"
    ? (user.role === 2 ? "admin" : "user")
    : user.role;

  return (
    <nav className="bg-white shadow-sm border-b">
      <div className="max-w-7xl mx-auto px-4 py-3 flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <span className="font-semibold text-gray-900">DevOps App</span>
          <a href="/dashboard" className="text-indigo-600 hover:text-indigo-800 font-medium">
            {userRole === "admin" ? "Admin Dashboard" : "Dashboard"}
          </a>
        </div>
        <div className="flex items-center space-x-4">
          <span className="text-sm text-gray-600">
            {user.username} <span className="text-gray-400">({userRole})</span>
          </span>
          <a
            href="/logout"
            className="text-sm text-red-600 hover:text-red-800 font-medium"
          >
            Logout
          </a>
        </div>
      </div>
    </nav>
  );
}