"use client";

import { useState } from "react";
import { Modal } from "@/components/ui/modal";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useProjectMembers, useProjectTeams } from "@/hooks/use-projects";
import { useTeams } from "@/hooks/use-teams";
import { Loader2, Trash2, User, Users } from "lucide-react";
import { cn } from "@/lib/utils";

interface ProjectMemberDialogProps {
  isOpen: boolean;
  onClose: () => void;
  projectId: number;
}

export function ProjectMemberDialog({
  isOpen,
  onClose,
  projectId,
}: ProjectMemberDialogProps) {
  const {
    members,
    addMember,
    removeMember,
    isLoading: loadingMembers,
  } = useProjectMembers(projectId);
  const {
    teams: projectTeams,
    addTeam,
    removeTeam,
    isLoading: loadingTeams,
  } = useProjectTeams(projectId);
  const { teams: allTeams } = useTeams();

  const [activeTab, setActiveTab] = useState<"members" | "teams">("members");

  // Member State
  const [email, setEmail] = useState("");
  const [isSubmittingMember, setIsSubmittingMember] = useState(false);

  // Team State
  const [selectedTeamId, setSelectedTeamId] = useState<string>("");
  const [isSubmittingTeam, setIsSubmittingTeam] = useState(false);

  const handleAddMember = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email) return;

    setIsSubmittingMember(true);
    try {
      await addMember(email);
      setEmail("");
    } catch (err) {
      alert("Failed to add member.");
    } finally {
      setIsSubmittingMember(false);
    }
  };

  const handleAddTeam = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedTeamId) return;

    setIsSubmittingTeam(true);
    try {
      await addTeam(parseInt(selectedTeamId));
      setSelectedTeamId("");
    } catch (err) {
      alert("Failed to add team.");
    } finally {
      setIsSubmittingTeam(false);
    }
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title="Project Access"
      description="Manage members and teams with access."
    >
      <div className="flex gap-2 mb-4 border-b pb-2">
        <button
          onClick={() => setActiveTab("members")}
          className={cn(
            "flex items-center gap-2 px-3 py-1.5 text-sm font-medium rounded-md transition-colors",
            activeTab === "members"
              ? "bg-blue-50 text-blue-600"
              : "text-gray-500 hover:text-gray-900"
          )}
        >
          <User className="h-4 w-4" />
          Members
        </button>
        <button
          onClick={() => setActiveTab("teams")}
          className={cn(
            "flex items-center gap-2 px-3 py-1.5 text-sm font-medium rounded-md transition-colors",
            activeTab === "teams"
              ? "bg-blue-50 text-blue-600"
              : "text-gray-500 hover:text-gray-900"
          )}
        >
          <Users className="h-4 w-4" />
          Teams
        </button>
      </div>

      {activeTab === "members" ? (
        <div className="space-y-6">
          <form onSubmit={handleAddMember} className="flex gap-2">
            <Input
              placeholder="Add member by email..."
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="text-gray-900"
            />
            <Button type="submit" disabled={isSubmittingMember}>
              {isSubmittingMember ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                "Add"
              )}
            </Button>
          </form>

          <div className="border rounded-md divide-y max-h-[300px] overflow-y-auto">
            {loadingMembers ? (
              <div className="p-4 text-center text-sm text-gray-500">
                Loading...
              </div>
            ) : members.length === 0 ? (
              <div className="p-4 text-center text-sm text-gray-500">
                No individual members found.
              </div>
            ) : (
              members.map((member) => (
                <div
                  key={member.user_id}
                  className="p-3 flex justify-between items-center hover:bg-gray-50"
                >
                  <span className="text-sm font-medium text-gray-900">
                    {member.email}
                  </span>
                  <div className="flex items-center gap-2">
                    <span className="text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded">
                      {member.role}
                    </span>
                    {!member.is_direct && (
                      <span className="text-[10px] text-blue-600 bg-blue-50 px-1.5 py-0.5 rounded border border-blue-100">
                        Team
                      </span>
                    )}
                    {member.role !== "ADMIN" && member.is_direct && (
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-8 w-8 p-0 text-red-500"
                        onClick={() => removeMember(member.user_id)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    )}
                  </div>
                </div>
              ))
            )}
          </div>
        </div>
      ) : (
        <div className="space-y-6">
          <form onSubmit={handleAddTeam} className="flex gap-2">
            <select
              className="flex h-10 w-full rounded-md border border-gray-200 bg-white px-3 py-2 text-sm text-gray-900 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
              value={selectedTeamId}
              onChange={(e) => setSelectedTeamId(e.target.value)}
            >
              <option value="">Select a team...</option>
              {allTeams.map((team) => (
                <option key={team.id} value={team.id}>
                  {team.name}
                </option>
              ))}
            </select>
            <Button type="submit" disabled={isSubmittingTeam}>
              {isSubmittingTeam ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                "Add"
              )}
            </Button>
          </form>

          <div className="border rounded-md divide-y max-h-[300px] overflow-y-auto">
            {loadingTeams ? (
              <div className="p-4 text-center text-sm text-gray-500">
                Loading...
              </div>
            ) : projectTeams.length === 0 ? (
              <div className="p-4 text-center text-sm text-gray-500">
                No teams assigned.
              </div>
            ) : (
              projectTeams.map((team) => (
                <div
                  key={team.team_id}
                  className="p-3 flex justify-between items-center hover:bg-gray-50"
                >
                  <div className="flex items-center gap-2">
                    <Users className="h-4 w-4 text-gray-400" />
                    <span className="text-sm font-medium text-gray-900">
                      {team.team_name}
                    </span>
                  </div>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="h-8 w-8 p-0 text-red-500"
                    onClick={() => removeTeam(team.team_id)}
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              ))
            )}
          </div>
        </div>
      )}
    </Modal>
  );
}
