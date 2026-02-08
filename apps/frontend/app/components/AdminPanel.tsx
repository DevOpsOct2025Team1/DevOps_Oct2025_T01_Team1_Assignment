import { useState } from "react";
import { getHours } from "date-fns";
import { useGetApiAdminListUsers } from "../api/generated";
import { Input } from "./ui/input";
import { Search, Loader2 } from "lucide-react";
import { CreateUserDialog } from "./CreateUserDialog";
import { UserRowActions } from "./UserRowActions";

interface User {
  id: string;
  username: string;
  role: string;
}

export default function AdminPanel() {
  const [search, setSearch] = useState("");
  
  const getGreeting = () => {
    const hours = getHours(new Date());
    if (hours < 12) return "Good morning";
    if (hours < 18) return "Good afternoon";
    return "Good evening";
  };
  
  const { data: users = [], isLoading } = useGetApiAdminListUsers(
    search ? { username: search } : {},
    {
      query: {
        select: (data) => {
          if (data.status === 200 && data.data) {
             return (data.data['users'] || []) as unknown as User[];
          }
          return [];
        }
      }
    }
  );

  const formatRole = (role: string) => {
    if (role === "ROLE_ADMIN" || role === "2") return "Admin";
    return "User";
  };

  return (
    <div className="space-y-8">
      <div className="space-y-2">
        <h1 className="text-4xl font-bold tracking-tight">{getGreeting()} Admin</h1>
        <p className="text-muted-foreground text-lg">
          What would you like to manage today?
        </p>
      </div>

      <div className="flex items-center justify-between gap-4">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            placeholder="Search Users"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="pl-9 bg-white rounded-full"
          />
        </div>

        <CreateUserDialog />
      </div>

      <div className="space-y-1">
        {isLoading ? (
          <div className="flex justify-center py-8">
            <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
          </div>
        ) : users.length === 0 ? (
          <p className="text-center text-muted-foreground py-8">No users found</p>
        ) : (
          users.map((user) => (
            <div 
              key={user.id} 
              className="flex items-center justify-between p-4 hover:bg-gray-50 rounded-lg transition-colors group"
            >
              <div className="flex items-center gap-4">
                <span className="font-medium text-lg">{user.username}</span>
              </div>
              
              <div className="flex items-center gap-6">
                <span className="text-sm">
                  {formatRole(user.role)}
                </span>
                
                <UserRowActions user={user} />
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
