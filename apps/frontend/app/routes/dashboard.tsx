import { lazy, Suspense, useEffect } from "react";
import { useNavigate } from "react-router";
import { useAuth } from "../contexts/AuthContext";

const UserPanel = lazy(() => import("../components/UserPanel"));

export default function Dashboard() {
  const navigate = useNavigate();
  const { user, isAuthenticated } = useAuth();

  useEffect(() => {
    if (!isAuthenticated) {
      navigate("/login");
      return;
    }

    // Redirect admin users to the admin dashboard
    if (user) {
      const userIsAdmin = typeof user.role === "number"
        ? user.role === 2
        : user.role.toLowerCase().includes("admin");
      
      if (userIsAdmin) {
        navigate("/admin");
        return;
      }
    }
  }, [isAuthenticated, user, navigate]);

  if (!user) {
    return null;
  }

  return (
    <div className="flex-1 bg-gray-50 py-8">
      <div className="max-w-5xl mx-auto px-4">
        {/* Greeting Section */}
        <div className="mb-6">
          <h1 className="text-3xl font-bold text-gray-900">Hello, {user.username}!</h1>
          <p className="text-gray-600 mt-2">
            Welcome to your dashboard. Access your files and manage your content.
          </p>
        </div>

        <Suspense fallback={<div className="text-center py-8 text-gray-600">Loading...</div>}>
          <UserPanel />
        </Suspense>
      </div>
    </div>
  );
}
