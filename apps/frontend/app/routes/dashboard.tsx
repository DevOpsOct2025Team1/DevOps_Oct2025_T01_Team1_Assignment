import { lazy, Suspense, useEffect } from "react";
import { useNavigate } from "react-router";
import { useAuth } from "~/contexts/AuthContext";

const AdminPanel = lazy(() => import("../components/AdminPanel"));
const UserPanel = lazy(() => import("../components/UserPanel"));

export default function Dashboard() {
  const navigate = useNavigate();
  const { user, isAuthenticated } = useAuth();

  useEffect(() => {
    if (!isAuthenticated) {
      navigate("/login");
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isAuthenticated]);

  if (!user) {
    return null;
  }

  const userIsAdmin = typeof user.role === "number"
    ? user.role === 2
    : user.role.toLowerCase().includes("admin");

  return (
    <div className="flex-1 bg-gray-50 py-8">
      <div className="max-w-5xl mx-auto px-4">
        <Suspense fallback={<div className="text-center py-8 text-gray-600">Loading...</div>}>
          {userIsAdmin ? <AdminPanel /> : <UserPanel />}
        </Suspense>
      </div>
    </div>
  );
}
