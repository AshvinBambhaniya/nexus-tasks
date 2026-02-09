"use client";

import { useState } from "react";
import { useTeamMembers } from "@/hooks/use-teams";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Avatar } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Loader2, Trash2, UserPlus, Mail } from "lucide-react";

interface TeamMembersProps {
  teamId: number;
}

export function TeamMembers({ teamId }: TeamMembersProps) {
  const { members, isLoading, addMember, removeMember } =
    useTeamMembers(teamId);
  const [inviteEmail, setInviteEmail] = useState("");
  const [isInviting, setIsInviting] = useState(false);

  const handleInvite = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!inviteEmail) return;

    setIsInviting(true);
    try {
      await addMember(inviteEmail);
      setInviteEmail("");
    } catch (err: any) {
      alert(err.response?.data?.detail || "Failed to add member");
    } finally {
      setIsInviting(false);
    }
  };

  const handleRemove = async (userId: number) => {
    if (!confirm("Remove this member from the team?")) return;
    try {
      await removeMember(userId);
    } catch (err: any) {
      alert(err.response?.data?.detail || "Failed to remove member");
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Team Members</CardTitle>
        <CardDescription>Manage the roster for this team.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Invite Form */}
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
            Add
          </Button>
        </form>

        {/* List */}
        {isLoading ? (
          <div className="flex justify-center p-8">
            <Loader2 className="h-6 w-6 animate-spin text-gray-400" />
          </div>
        ) : (
          <div className="divide-y divide-gray-100">
            {members.map((member) => (
              <div
                key={member.user_id}
                className="flex items-center justify-between rounded-lg p-4 transition-colors hover:bg-gray-50"
              >
                <div className="flex items-center gap-3">
                  <Avatar
                    fallback={member.email[0].toUpperCase()}
                    className="h-9 w-9 bg-indigo-100 text-xs text-indigo-700"
                  />
                  <div>
                    <div className="text-sm font-medium text-gray-900">
                      {member.email}
                    </div>
                    <div className="text-xs text-gray-500 capitalize">
                      {member.role.toLowerCase()}
                    </div>
                  </div>
                </div>

                <div className="flex items-center gap-3">
                  {member.role !== "ADMIN" && (
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-8 w-8 text-gray-400 hover:bg-red-50 hover:text-red-600"
                      onClick={() => handleRemove(member.user_id)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  )}
                </div>
              </div>
            ))}
            {members.length === 0 && (
              <div className="py-6 text-center text-sm text-gray-500">
                No members found.
              </div>
            )}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
