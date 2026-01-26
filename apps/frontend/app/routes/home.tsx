import { useEffect } from 'react';
import { useNavigate } from 'react-router';
import type { Route } from "./+types/home";
import { isAuthenticated, isAdmin } from "../utils/auth";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "DevOps App" },
    { name: "description", content: "DevOps Assignment Application" },
  ];
}

export default function Home() {
  const navigate = useNavigate();

  useEffect(() => {
    if (isAuthenticated()) {
      if (isAdmin()) {
        navigate('/admin');
      } else {
        navigate('/dashboard');
      }
    } else {
      navigate('/login');
    }
  }, [navigate]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="text-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto mb-4"></div>
        <p className="text-gray-600">Loading...</p>
      </div>
    </div>
  );
}
