"use client";

import { useState } from "react";
import { useProjectMembers, useProjectTeams } from "@/hooks/use-projects";
import { useWorkspaces } from "@/hooks/use-workspaces";
import { useTeams } from "@/hooks/use-teams";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Avatar } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Loader2, Trash2, UserPlus, Users, Link as LinkIcon, Mail } from "lucide-react";

interface ProjectMembersProps {
  projectId: number;
}

export function ProjectMembers({ projectId }: ProjectMembersProps) {
  return (
    <div className="space-y-8">
      <DirectMembersSection projectId={projectId} />
      <TeamMembersSection projectId={projectId} />
    </div>
  );
}

function DirectMembersSection({ projectId }: { projectId: number }) {
  const { members, isLoading, addMember, removeMember } = useProjectMembers(projectId);
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
    if (!confirm("Remove this member from the project?")) return;
    try {
      await removeMember(userId);
    } catch (err: any) {
        alert(err.response?.data?.detail || "Failed to remove member");
    }
  };

  // Filter only direct members for this list
  // const directMembers = members.filter(m => m.is_direct);
  // Requested change: Show all members
  const allMembers = members;

  return (
    <Card>
      <CardHeader>
        <CardTitle>Project Members</CardTitle>
        <CardDescription>All users with access to this project, including team members.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Invite Form */}
        <form onSubmit={handleInvite} className="flex gap-3">
          <div className="relative flex-1">
            <Mail className="absolute left-3 top-2.5 h-4 w-4 text-gray-400" />
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
            {isInviting ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <UserPlus className="mr-2 h-4 w-4" />}
            Add
          </Button>
        </form>

        {/* List */}
        {isLoading ? (
          <div className="flex justify-center p-8"><Loader2 className="h-6 w-6 animate-spin text-gray-400" /></div>
        ) : (
          <div className="divide-y divide-gray-100">
             {allMembers.map((member) => (
                <div key={member.user_id} className="flex items-center justify-between p-4 hover:bg-gray-50 rounded-lg transition-colors">
                  <div className="flex items-center gap-3">
                    <Avatar
                       fallback={(member.email)[0].toUpperCase()}
                       className="bg-blue-100 text-blue-700 h-9 w-9 text-xs"
                    />
                    <div>
                      <div className="font-medium text-sm text-gray-900">{member.email}</div>
                      <div className="text-xs text-gray-500 flex items-center gap-2">
                          {member.is_direct ? "Direct Member" : "via Team"}
                      </div>
                    </div>
                  </div>
                  
                  <div className="flex items-center gap-3">
                    <Badge variant={member.role === "ADMIN" ? "default" : "secondary"} className="text-[10px] h-5">
                      {member.role}
                    </Badge>
                    {member.is_direct && (
                        <Button
                            variant="ghost"
                            size="icon"
                            className="h-8 w-8 text-gray-400 hover:text-red-600 hover:bg-red-50"
                            onClick={() => handleRemove(member.user_id)}
                            title="Remove direct access"
                        >
                            <Trash2 className="h-4 w-4" />
                        </Button>
                    )}
                  </div>
                </div>
             ))}
             {allMembers.length === 0 && (
                 <div className="text-center py-6 text-sm text-gray-500">No members found.</div>
             )}
          </div>
        )}
      </CardContent>
    </Card>
  );
}

function TeamMembersSection({ projectId }: { projectId: number }) {
  const { teams: projectTeams, isLoading, addTeam, removeTeam } = useProjectTeams(projectId);
  const { teams: workspaceTeams } = useTeams(); // Fetch all available teams to link
  
  const [selectedTeamId, setSelectedTeamId] = useState<string>("");
  const [isLinking, setIsLinking] = useState(false);

  const handleLinkTeam = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedTeamId) return;

    setIsLinking(true);
    try {
        await addTeam(parseInt(selectedTeamId));
        setSelectedTeamId("");
    } catch (err: any) {
        alert(err.response?.data?.detail || "Failed to link team");
    } finally {
        setIsLinking(false);
    }
  };

  const handleUnlink = async (teamId: number) => {
      if (!confirm("Unlink this team from the project?")) return;
      try {
          await removeTeam(teamId);
      } catch (err: any) {
          alert(err.response?.data?.detail || "Failed to unlink team");
      }
  }

  // Filter out teams already linked
  const availableTeams = workspaceTeams.filter(wt => !projectTeams.some(pt => pt.team_id === wt.id));

  return (
    <Card>
      <CardHeader>
        <CardTitle>Teams</CardTitle>
        <CardDescription>Workspace teams with access to this project.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
         {/* Link Form */}
         <form onSubmit={handleLinkTeam} className="flex gap-3">
            <div className="relative flex-1">
                <Users className="absolute left-3 top-2.5 h-4 w-4 text-gray-400" />
                <select
                    className="flex h-10 w-full rounded-md border border-gray-200 bg-white px-3 py-2 pl-9 text-sm ring-offset-white file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-gray-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                    value={selectedTeamId}
                    onChange={(e) => setSelectedTeamId(e.target.value)}
                    required
                >
                    <option value="" disabled>Select a team to link...</option>
                    {availableTeams.map(t => (
                        <option key={t.id} value={t.id}>{t.name}</option>
                    ))}
                </select>
            </div>
            <Button type="submit" disabled={isLinking || availableTeams.length === 0}>
                {isLinking ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <LinkIcon className="mr-2 h-4 w-4" />}
                Link Team
            </Button>
         </form>

         {/* List */}
         {isLoading ? (
            <div className="flex justify-center p-8"><Loader2 className="h-6 w-6 animate-spin text-gray-400" /></div>
         ) : (
            <div className="divide-y divide-gray-100">
                {projectTeams.map(team => (
                    <div key={team.team_id} className="flex items-center justify-between p-4 hover:bg-gray-50 rounded-lg transition-colors">
                        <div className="flex items-center gap-3">
                            <div className="h-9 w-9 rounded-lg bg-indigo-100 flex items-center justify-center text-indigo-700">
                                <Users className="h-4 w-4" />
                            </div>
                            <div className="font-medium text-sm text-gray-900">{team.team_name}</div>
                        </div>
                        <Button
                            variant="ghost"
                            size="icon"
                            className="h-8 w-8 text-gray-400 hover:text-red-600 hover:bg-red-50"
                            onClick={() => handleUnlink(team.team_id)}
                        >
                            <Trash2 className="h-4 w-4" />
                        </Button>
                    </div>
                ))}
                {projectTeams.length === 0 && (
                    <div className="text-center py-6 text-sm text-gray-500">No teams linked.</div>
                )}
            </div>
         )}
      </CardContent>
    </Card>
  );
}
