"use client";

import { useState } from "react";
import useSWR from "swr";
import api from "@/lib/api";
import { useWorkspaces } from "@/hooks/use-workspaces";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Loader2, Trash2 } from "lucide-react";
import { WorkspaceMember } from "@/types";

export function MemberSettings() {
  const { activeWorkspace } = useWorkspaces();
  const [inviteEmail, setInviteEmail] = useState("");
  const [isInviting, setIsInviting] = useState(false);
  
  const { data: members, error, mutate } = useSWR<WorkspaceMember[]>(
    activeWorkspace ? `/api/v1/workspaces/${activeWorkspace.id}/members` : null,
    (url: string) => api.get(url).then(res => res.data)
  );
  
  const isLoading = !members && !error;

  const handleInvite = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!activeWorkspace || !inviteEmail) return;

    setIsInviting(true);
    try {
      await api.post(`/api/v1/workspaces/${activeWorkspace.id}/members`, { email: inviteEmail });
      setInviteEmail("");
      mutate(); // Refresh list
    } catch (err: any) {
      alert(err.response?.data?.detail || "Failed to invite user");
    } finally {
      setIsInviting(false);
    }
  };

  const handleRemove = async (userId: number) => {
    if (!activeWorkspace || !confirm("Are you sure you want to remove this member?")) return;

    try {
      await api.delete(`/api/v1/workspaces/${activeWorkspace.id}/members/${userId}`);
      mutate();
    } catch (err: any) {
      alert(err.response?.data?.detail || "Failed to remove member");
    }
  };

  if (!activeWorkspace) return <div>Please select a workspace.</div>;

  return (
    <div className="space-y-6">
      <div className="rounded-lg border border-gray-200 p-4">
        <h2 className="text-lg font-semibold mb-4 text-gray-900">Invite Member</h2>
        <form onSubmit={handleInvite} className="flex gap-2">
          <Input 
            placeholder="Enter email address" 
            type="email"
            value={inviteEmail}
            onChange={(e) => setInviteEmail(e.target.value)}
            required
            className="text-gray-900"
          />
          <Button type="submit" disabled={isInviting}>
            {isInviting ? <Loader2 className="h-4 w-4 animate-spin" /> : "Invite"}
          </Button>
        </form>
      </div>

      <div className="rounded-lg border border-gray-200 overflow-hidden">
        <div className="bg-gray-50 px-4 py-3 border-b border-gray-200">
          <h3 className="font-medium text-gray-900">Team Members</h3>
        </div>
        {isLoading ? (
          <div className="p-8 flex justify-center">
             <Loader2 className="h-6 w-6 animate-spin text-gray-400" />
          </div>
        ) : (
          <table className="w-full text-sm text-left">
            <thead className="bg-white text-gray-500 font-medium border-b border-gray-100">
              <tr>
                <th className="px-4 py-3 text-gray-700">User</th>
                <th className="px-4 py-3 text-gray-700">Role</th>
                <th className="px-4 py-3 text-right text-gray-700">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {members?.map((member) => (
                <tr key={member.user_id} className="hover:bg-gray-50">
                  <td className="px-4 py-3">
                    <div className="font-medium text-gray-900">{member.user.email}</div>
                  </td>
                  <td className="px-4 py-3">
                    <Badge variant={member.role === "ADMIN" ? "default" : "secondary"}>
                      {member.role}
                    </Badge>
                  </td>
                  <td className="px-4 py-3 text-right">
                    {member.role !== "ADMIN" && ( // Simple logic: can't remove admins for now to be safe
                       <Button 
                         variant="ghost" 
                         size="sm" 
                         className="h-8 w-8 p-0 text-red-500 hover:text-red-700 hover:bg-red-50"
                         onClick={() => handleRemove(member.user_id)}
                       >
                         <Trash2 className="h-4 w-4" />
                       </Button>
                    )}
                  </td>
                </tr>
              ))}
              {members?.length === 0 && (
                <tr>
                   <td colSpan={3} className="px-4 py-8 text-center text-gray-500">No members found.</td>
                </tr>
              )}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
