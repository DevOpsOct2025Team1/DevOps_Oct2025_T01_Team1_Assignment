import { useEffect } from "react";
import { useNavigate } from "react-router";
import { clearAuth } from "../utils/auth";

export default function Logout() {
  const navigate = useNavigate();

  useEffect(() => {
    clearAuth();
    navigate("/login");
  }, [navigate]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="text-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto mb-4"></div>
        <p className="text-gray-600">Logging out...</p>
      </div>
    </div>
  );
}
