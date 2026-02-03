import { useState, useEffect } from "react";
import { useNavigate } from "react-router";
import { getStoredUser, isAuthenticated, isAdmin, type User } from "../utils/auth";
import AdminPanel from "../components/AdminPanel";
import UserPanel from "../components/UserPanel";

export default function Dashboard() {
  const navigate = useNavigate();
  const [user, setUser] = useState<User | null>(null);
  const userIsAdmin = isAdmin();

  useEffect(() => {
    if (!isAuthenticated()) {
      navigate("/login");
      return;
    }

    const storedUser = getStoredUser();
    setUser(storedUser);
  }, [navigate]);

  if (!user) {
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-5xl mx-auto px-4">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">
            {userIsAdmin ? "Admin Dashboard" : "Dashboard"}
          </h1>
          <p className="text-gray-600 mt-1">Welcome, {user.username}</p>
        </div>

        {userIsAdmin ? <AdminPanel /> : <UserPanel />}
      </div>
    </div>
  );
}
