import { useState, useEffect } from "react";
import { useNavigate } from "react-router";
import { useAuth } from "../contexts/AuthContext";
import { useLogin } from "../api/generated";
import { Button } from "../components/ui/button";
import { Input } from "../components/ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../components/ui/card";
import { authApi } from "../utils/api";

export default function Login() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [healthStatus, setHealthStatus] = useState("");
  const [healthChecking, setHealthChecking] = useState(false);
  const navigate = useNavigate();
  const { setAuth, isAuthenticated } = useAuth();
  const loginMutation = useLogin();
  const isLoading = loginMutation.isPending;

  useEffect(() => {
    if (isAuthenticated) {
      navigate("/dashboard");
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isAuthenticated]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setHealthStatus("");

    if (!username.trim() || !password.trim()) {
      setError("Username and password are required");
      return;
    }

    // TODO: ENDPOINTS - Login endpoint: POST /api/login
    try {
      const response = await loginMutation.mutateAsync({
        data: {
          username,
          password,
        },
      });
      const authData = response.data;

      if (
        !authData ||
        !("user" in authData) ||
        !authData.user ||
        !authData.user.id ||
        !authData.user.username ||
        !authData.user.role ||
        !authData.token
      ) {
        throw new Error("Login failed. Please try again.");
      }

      console.log("Login response:", authData);
      console.log("User role:", authData.user.role);

      const normalizedUser = {
        id: authData.user.id,
        username: authData.user.username,
        role: authData.user.role
      };

      setAuth(normalizedUser, authData.token);
      navigate("/dashboard");
    } catch (err: unknown) {
      const err2 = err as Error;
      if (err2.message) {
        setError(err2.message);
      } else {
        setError("Login failed. Please try again.");
      }
    }
  };

  // TODO: ENDPOINTS - Health check endpoint: GET /health
  const handleHealthCheck = async () => {
    setHealthChecking(true);
    setHealthStatus("");
    setError("");
    try {
      const response = await authApi.health();
      setHealthStatus(`Backend is healthy: ${response.status}`);
    } catch (err) {
      setHealthStatus("Backend is unavailable or unhealthy");
    } finally {
      setHealthChecking(false);
    }
  };

  return (
    <div className="flex-1 flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="text-3xl text-center">Login</CardTitle>
          <CardDescription className="text-center">
            Enter your credentials to access your account
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <label htmlFor="username" className="text-sm font-medium">
                Username
              </label>
              <Input
                id="username"
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                disabled={isLoading}
                placeholder="Enter your username"
              />
            </div>

            <div className="space-y-2">
              <label htmlFor="password" className="text-sm font-medium">
                Password
              </label>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                disabled={isLoading}
                placeholder="Enter your password"
              />
            </div>

            {error && (
              <div className="text-sm text-destructive bg-destructive/10 p-3 rounded-md">
                {error}
              </div>
            )}

            {healthStatus && (
              <div className="text-sm text-primary bg-primary/10 p-3 rounded-md">
                {healthStatus}
              </div>
            )}

            <Button
              type="submit"
              disabled={loginMutation.isPending}
              className="w-full"
            >
              {loginMutation.isPending ? "Logging in..." : "Login"}
            </Button>

            <Button
              type="button"
              variant="outline"
              onClick={handleHealthCheck}
              disabled={healthChecking}
              className="w-full"
            >
              {healthChecking ? "Checking..." : "Check Backend Health"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
