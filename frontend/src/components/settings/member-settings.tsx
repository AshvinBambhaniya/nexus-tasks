"use client";

import { useState } from "react";
import useSWR from "swr";
import api from "@/lib/api";
import { useWorkspaces } from "@/hooks/use-workspaces";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Avatar } from "@/components/ui/avatar";
import {
  Loader2,
  Trash2,
  UserPlus,
  Mail,
  Shield,
  User as UserIcon,
} from "lucide-react";
import { WorkspaceMember, ApiError } from "@/types";

export function MemberSettings() {
  const { activeWorkspace } = useWorkspaces();
  const [inviteEmail, setInviteEmail] = useState("");
  const [isInviting, setIsInviting] = useState(false);

  const {
    data: members,
    error,
    mutate,
  } = useSWR<WorkspaceMember[]>(
    activeWorkspace ? `/api/v1/workspaces/${activeWorkspace.id}/members` : null,
    (url: string) => api.get(url).then((res) => res.data)
  );

  const isLoading = !members && !error;

  const handleInvite = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!activeWorkspace || !inviteEmail) return;

    setIsInviting(true);
    try {
      await api.post(`/api/v1/workspaces/${activeWorkspace.id}/members`, {
        email: inviteEmail,
      });
      setInviteEmail("");
      mutate(); // Refresh list
    } catch (err) {
      alert(
        (err as ApiError).response?.data?.detail || "Failed to invite user"
      );
    } finally {
      setIsInviting(false);
    }
  };

  const handleRemove = async (userId: number) => {
    if (
      !activeWorkspace ||
      !confirm("Are you sure you want to remove this member?")
    )
      return;

    try {
      await api.delete(
        `/api/v1/workspaces/${activeWorkspace.id}/members/${userId}`
      );
      mutate();
    } catch (err) {
      alert(
        (err as ApiError).response?.data?.detail || "Failed to remove member"
      );
    }
  };

  if (!activeWorkspace) {
    return (
      <Card>
        <CardContent className="p-8 text-center text-gray-500">
          Please select a workspace to manage members.
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      {/* Invite Section */}
      <Card>
        <CardHeader>
          <CardTitle>Invite Members</CardTitle>
          <CardDescription>
            Add new members to your team by email.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleInvite} className="flex gap-3">
            <div className="relative flex-1">
              <Mail className="absolute top-2.5 left-3 h-4 w-4 text-gray-400" />
              <Input
                placeholder="name@example.com"
                type="email"
                value={inviteEmail}
                onChange={(e) => setInviteEmail(e.target.value)}
                required
                className="pl-9"
              />
            </div>
            <Button type="submit" disabled={isInviting}>
              {isInviting ? (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              ) : (
                <UserPlus className="mr-2 h-4 w-4" />
              )}
              Invite
            </Button>
          </form>
        </CardContent>
      </Card>

      {/* Members List */}
      <Card>
        <CardHeader>
          <CardTitle>Team Members</CardTitle>
          <CardDescription>
            Manage existing members and their roles.
          </CardDescription>
        </CardHeader>
        <CardContent className="p-0">
          {isLoading ? (
            <div className="flex justify-center p-12">
              <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
            </div>
          ) : (
            <div className="divide-y divide-gray-100">
              {members?.map((member) => (
                <div
                  key={member.user_id}
                  className="flex items-center justify-between p-6 transition-colors hover:bg-gray-50/50"
                >
                  <div className="flex items-center gap-4">
                    <Avatar
                      fallback={(member.user.full_name ||
                        member.user.email)[0].toUpperCase()}
                      className="bg-blue-100 font-medium text-blue-700"
                    />
                    <div>
                      <div className="font-medium text-gray-900">
                        {member.user.full_name ||
                          member.user.email.split("@")[0]}
                      </div>
                      <div className="flex items-center gap-2 text-sm text-gray-500">
                        {member.user.email}
                      </div>
                    </div>
                  </div>

                  <div className="flex items-center gap-4">
                    <Badge
                      variant={
                        member.role === "ADMIN" ? "default" : "secondary"
                      }
                      className="capitalize"
                    >
                      {member.role === "ADMIN" && (
                        <Shield className="mr-1 h-3 w-3" />
                      )}
                      {member.role.toLowerCase()}
                    </Badge>

                    {member.role !== "ADMIN" && (
                      <Button
                        variant="ghost"
                        size="icon"
                        className="text-gray-400 hover:bg-red-50 hover:text-red-600"
                        onClick={() => handleRemove(member.user_id)}
                        title="Remove member"
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    )}
                  </div>
                </div>
              ))}

              {members?.length === 0 && (
                <div className="p-12 text-center text-gray-500">
                  <UserIcon className="mx-auto mb-3 h-12 w-12 text-gray-200" />
                  <p>No members found in this workspace.</p>
                </div>
              )}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
