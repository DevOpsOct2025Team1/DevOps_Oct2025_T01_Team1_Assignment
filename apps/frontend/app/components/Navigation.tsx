import { useLocation } from "react-router";
import { useState, useEffect } from "react";
import { useAuth } from "../contexts/AuthContext";
import { Badge } from "./ui/badge";
import { Button } from "./ui/button";
import { Separator } from "./ui/separator";

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
    <nav className="bg-background shadow-sm border-b">
      <div className="max-w-7xl mx-auto px-4 py-3 flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <span className="font-semibold text-lg">DevOps App</span>
          <Separator orientation="vertical" className="h-6" />
          <a href="/dashboard" className="text-primary hover:text-primary/80 font-medium">
            {userRole === "admin" ? "Admin Dashboard" : "Dashboard"}
          </a>
        </div>
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <span className="text-sm">{user.username}</span>
            <Badge variant={userRole === "admin" ? "default" : "secondary"}>
              {userRole}
            </Badge>
          </div>
          <Button
            variant="ghost"
            size="sm"
            asChild
          >
            <a href="/logout">Logout</a>
          </Button>
        </div>
      </div>
    </nav>
  );
}