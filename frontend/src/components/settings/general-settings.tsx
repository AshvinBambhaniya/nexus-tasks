"use client";

import { useUser } from "@/hooks/use-user";
import { useWorkspaces } from "@/hooks/use-workspaces";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { LogOut } from "lucide-react";
import { useAuth } from "@/hooks/use-auth";

export function GeneralSettings() {
  const { user } = useUser();
  const { activeWorkspace } = useWorkspaces();
  const { logout } = useAuth();

  return (
    <div className="space-y-6 max-w-xl">
      <div className="space-y-4 rounded-lg border border-gray-200 p-4">
        <h2 className="text-lg font-semibold text-gray-900">Profile</h2>
        <div className="flex items-center gap-4">
          <div className="h-12 w-12 rounded-full bg-blue-100 flex items-center justify-center text-blue-700 text-lg font-bold">
            {user?.email.substring(0, 2).toUpperCase()}
          </div>
          <div>
            <p className="font-medium text-gray-900">{user?.email}</p>
            <p className="text-sm text-gray-500">User ID: {user?.id}</p>
          </div>
        </div>
      </div>

      <div className="space-y-4 rounded-lg border border-gray-200 p-4">
        <h2 className="text-lg font-semibold text-gray-900">Active Workspace</h2>
        {activeWorkspace ? (
          <div className="space-y-2">
            <div>
              <Label className="text-xs text-gray-500">Name</Label>
              <p className="font-medium text-gray-900">{activeWorkspace.name}</p>
            </div>
            <div>
              <Label className="text-xs text-gray-500">Type</Label>
              <p className="font-medium capitalize text-gray-900">{activeWorkspace.type.toLowerCase()}</p>
            </div>
          </div>
        ) : (
          <p className="text-sm text-gray-500">No active workspace selected.</p>
        )}
      </div>

      <div className="pt-4">
         <Button variant="destructive" onClick={logout} className="w-full sm:w-auto">
            <LogOut className="mr-2 h-4 w-4" /> Log out
         </Button>
      </div>
    </div>
  );
}
