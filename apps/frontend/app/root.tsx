import {
  isRouteErrorResponse,
  Links,
  Meta,
  Outlet,
  Scripts,
  ScrollRestoration,
  useLocation,
} from "react-router";
import { useState, useEffect } from "react";
import { QueryClientProvider } from "@tanstack/react-query";

import type { Route } from "./+types/root";
import "./app.css";
import { createQueryClient } from "./api/query";
import { getStoredUser, isAuthenticated, type User } from "./utils/auth";

export const links: Route.LinksFunction = () => [
  { rel: "preconnect", href: "https://fonts.googleapis.com" },
  {
    rel: "preconnect",
    href: "https://fonts.gstatic.com",
    crossOrigin: "anonymous",
  },
  {
    rel: "stylesheet",
    href: "https://fonts.googleapis.com/css2?family=Inter:ital,opsz,wght@0,14..32,100..900;1,14..32,100..900&display=swap",
  },
];

export function Layout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" suppressHydrationWarning>
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <Meta />
        <Links />
      </head>
      <body>
        {children}
        <ScrollRestoration />
        <Scripts />
      </body>
    </html>
  );
}

function Navigation() {
  const location = useLocation();
  const [user, setUser] = useState<User | null>(null);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
    if (isAuthenticated()) {
      setUser(getStoredUser());
    } else {
      setUser(null);
    }
  }, [location.pathname]);

  if (!mounted || location.pathname === '/login') {
    return null;
  }

  if (!user) {
    return null;
  }

  return (
    <nav className="bg-white shadow-sm border-b">
      <div className="max-w-7xl mx-auto px-4 py-3 flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <span className="font-semibold text-gray-900">DevOps App</span>
          {user.role === 'admin' ? (
            <a href="/admin" className="text-indigo-600 hover:text-indigo-800 font-medium">
              Admin Dashboard
            </a>
          ) : (
            <a href="/dashboard" className="text-indigo-600 hover:text-indigo-800 font-medium">
              Dashboard
            </a>
          )}
        </div>
        <div className="flex items-center space-x-4">
          <span className="text-sm text-gray-600">
            {user.username} <span className="text-gray-400">({user.role})</span>
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

export default function App() {
  const [queryClient] = useState(() => createQueryClient());

  return (
    <QueryClientProvider client={queryClient}>
      <Navigation />
      <Outlet />
    </QueryClientProvider>
  );
}

export function ErrorBoundary({ error }: Route.ErrorBoundaryProps) {
  let message = "Oops!";
  let details = "An unexpected error occurred.";
  let stack: string | undefined;

  if (isRouteErrorResponse(error)) {
    message = error.status === 404 ? "404" : "Error";
    details =
      error.status === 404
        ? "The requested page could not be found."
        : error.statusText || details;
  } else if (import.meta.env.DEV && error && error instanceof Error) {
    details = error.message;
    stack = error.stack;
  }

  return (
    <main className="pt-16 p-4 container mx-auto">
      <h1>{message}</h1>
      <p>{details}</p>
      {stack && (
        <pre className="w-full p-4 overflow-x-auto">
          <code>{stack}</code>
        </pre>
      )}
    </main>
  );
}
