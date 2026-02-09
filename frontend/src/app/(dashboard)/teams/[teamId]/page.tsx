"use client";

import { use } from "react";
import { useState } from "react";
import { useTeam } from "@/hooks/use-teams";
import { TeamProjects } from "@/components/team/tabs/team-projects";
import { TeamMembers } from "@/components/team/tabs/team-members";
import { TeamSettings } from "@/components/team/tabs/team-settings";
import { cn } from "@/lib/utils";
import { Loader2 } from "lucide-react";

export default function TeamPage({
  params,
}: {
  params: Promise<{ teamId: string }>;
}) {
  const resolvedParams = use(params);
  const teamId = parseInt(resolvedParams.teamId);

  const { team, isLoading } = useTeam(teamId);

  const [activeTab, setActiveTab] = useState<
    "projects" | "members" | "settings"
  >("projects");

  if (isLoading) {
    return (
      <div className="flex h-full items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
      </div>
    );
  }

  if (!team) {
    return (
      <div className="flex h-full items-center justify-center text-gray-500">
        Team not found.
      </div>
    );
  }

  return (
    <div className="flex h-full flex-col space-y-6">
      {/* Header */}
      <div>
        <div className="mb-1 flex items-center gap-2">
          <h1 className="text-2xl font-bold tracking-tight text-gray-900">
            {team.name}
          </h1>
        </div>
        {team.description && (
          <p className="max-w-2xl text-sm text-gray-500">{team.description}</p>
        )}

        {/* Tab Navigation */}
        <div className="mt-6 flex gap-6 border-b border-gray-200">
          {[
            { id: "projects", label: "Overview" },
            { id: "members", label: "Members" },
            { id: "settings", label: "Settings" },
          ].map((tab) => (
            <button
              key={tab.id}
              onClick={() =>
                setActiveTab(tab.id as "projects" | "members" | "settings")
              }
              className={cn(
                "border-b-2 pb-3 text-sm font-medium transition-colors",
                activeTab === tab.id
                  ? "border-blue-600 text-blue-600"
                  : "border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700"
              )}
            >
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      {/* Tab Content */}
      <div className="min-h-0 flex-1">
        {activeTab === "projects" && <TeamProjects teamId={teamId} />}
        {activeTab === "members" && <TeamMembers teamId={teamId} />}
        {activeTab === "settings" && <TeamSettings team={team} />}
      </div>
    </div>
  );
}
