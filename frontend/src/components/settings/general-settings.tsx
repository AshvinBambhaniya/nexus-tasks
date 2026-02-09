"use client";

import { useState, useEffect } from "react";
import { useUser } from "@/hooks/use-user";
import { useWorkspaces } from "@/hooks/use-workspaces";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { LogOut, Building, Loader2, Check } from "lucide-react";
import { ApiError } from "@/types";
import { useAuth } from "@/hooks/use-auth";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Avatar } from "@/components/ui/avatar";

export function GeneralSettings() {
  const { user, updateUser } = useUser();
  const { activeWorkspace } = useWorkspaces();
  const { logout } = useAuth();

  const [fullName, setFullName] = useState("");
  const [email, setEmail] = useState("");
  const [isUpdating, setIsUpdating] = useState(false);
  const [isSuccess, setIsSuccess] = useState(false);

  useEffect(() => {
    if (user) {
      setFullName(user.full_name || "");
      setEmail(user.email || "");
    }
  }, [user]);

  const handleUpdateProfile = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsUpdating(true);
    setIsSuccess(false);
    try {
      await updateUser({ full_name: fullName, email });
      setIsSuccess(true);
      setTimeout(() => setIsSuccess(false), 3000);
    } catch (err) {
      alert(
        (err as ApiError).response?.data?.detail || "Failed to update profile"
      );
    } finally {
      setIsUpdating(false);
    }
  };

  const userInitial = (user?.full_name || user?.email || "??")
    .substring(0, 2)
    .toUpperCase();

  return (
    <div className="space-y-6">
      {/* Profile Card */}
      <Card>
        <CardHeader>
          <CardTitle>Profile</CardTitle>
          <CardDescription>
            Manage your personal information and how others see you.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="flex flex-col items-center gap-6 border-b border-gray-100 pb-6 sm:flex-row">
            <Avatar className="h-24 w-24 text-2xl" fallback={userInitial} />
            <div className="space-y-1 text-center sm:text-left">
              <h3 className="text-xl font-semibold text-gray-900">
                {user?.full_name || "Nexus User"}
              </h3>
              <p className="text-gray-500">{user?.email}</p>
              <div className="pt-1 font-mono text-xs text-gray-400">
                User ID: {user?.id}
              </div>
            </div>
          </div>

          <form onSubmit={handleUpdateProfile} className="space-y-4 pt-2">
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="fullName">Full Name</Label>
                <Input
                  id="fullName"
                  value={fullName}
                  onChange={(e) => setFullName(e.target.value)}
                  placeholder="John Doe"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="email">Email Address</Label>
                <Input
                  id="email"
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="john@example.com"
                />
              </div>
            </div>

            <div className="flex items-center gap-4 pt-2">
              <Button type="submit" disabled={isUpdating}>
                {isUpdating ? (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                ) : isSuccess ? (
                  <Check className="mr-2 h-4 w-4 text-green-500" />
                ) : null}
                {isUpdating
                  ? "Saving..."
                  : isSuccess
                    ? "Saved!"
                    : "Save Changes"}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>

      {/* Workspace Card */}
      <Card>
        <CardHeader>
          <CardTitle>Workspace</CardTitle>
          <CardDescription>
            Details about your currently active workspace.
          </CardDescription>
        </CardHeader>
        <CardContent>
          {activeWorkspace ? (
            <div className="grid gap-6 sm:grid-cols-2">
              <div className="space-y-2">
                <Label className="text-gray-500">Workspace Name</Label>
                <div className="flex items-center gap-3 rounded-lg border border-gray-200 bg-gray-50 p-4 text-gray-900">
                  <div className="rounded-md border border-gray-200 bg-white p-2">
                    <Building className="h-4 w-4 text-blue-600" />
                  </div>
                  <span className="font-semibold">{activeWorkspace.name}</span>
                </div>
              </div>
              <div className="space-y-2">
                <Label className="text-gray-500">Type</Label>
                <div className="flex items-center gap-3 rounded-lg border border-gray-200 bg-gray-50 p-4 text-gray-900">
                  <span className="inline-flex items-center rounded-full bg-blue-100 px-2.5 py-0.5 text-xs font-semibold text-blue-800 capitalize">
                    {activeWorkspace.type.toLowerCase()}
                  </span>
                </div>
              </div>
            </div>
          ) : (
            <div className="rounded-lg border border-dashed border-gray-200 bg-gray-50 p-8 text-center text-sm text-gray-500">
              No active workspace selected.
            </div>
          )}
        </CardContent>
      </Card>

      {/* Danger Zone */}
      <Card className="border-red-100 bg-red-50/30">
        <CardHeader>
          <CardTitle className="text-red-900">Danger Zone</CardTitle>
          <CardDescription className="text-red-700">
            Actions that are irreversible or sign you out.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Button
            variant="outline"
            onClick={logout}
            className="border-red-200 bg-white text-red-600 hover:bg-red-100 hover:text-red-700"
          >
            <LogOut className="mr-2 h-4 w-4" /> Log out of account
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}
