"use client";

import { use } from "react";
import { useState } from "react";
import { useTeam } from "@/hooks/use-teams";
import { TeamProjects } from "@/components/team/tabs/team-projects";
import { TeamMembers } from "@/components/team/tabs/team-members";
import { TeamSettings } from "@/components/team/tabs/team-settings";
import { cn } from "@/lib/utils";
import { Loader2 } from "lucide-react";

export default function TeamPage({ params }: { params: Promise<{ teamId: string }> }) {
  const resolvedParams = use(params);
  const teamId = parseInt(resolvedParams.teamId);
  
  const { team, isLoading } = useTeam(teamId); 
  
  const [activeTab, setActiveTab] = useState<"projects" | "members" | "settings">("projects");

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
    <div className="flex flex-col h-full space-y-6">
      {/* Header */}
      <div>
         <div className="flex items-center gap-2 mb-1">
             <h1 className="text-2xl font-bold tracking-tight text-gray-900">
               {team.name}
             </h1>
         </div>
         {team.description && (
           <p className="text-sm text-gray-500 max-w-2xl">{team.description}</p>
         )}
         
         {/* Tab Navigation */}
         <div className="flex gap-6 mt-6 border-b border-gray-200">
            {[
                { id: "projects", label: "Overview" },
                { id: "members", label: "Members" },
                { id: "settings", label: "Settings" }
            ].map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id as any)}
                className={cn(
                  "pb-3 text-sm font-medium transition-colors border-b-2",
                  activeTab === tab.id
                    ? "border-blue-600 text-blue-600"
                    : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
                )}
              >
                {tab.label}
              </button>
            ))}
         </div>
      </div>

      {/* Tab Content */}
      <div className="flex-1 min-h-0">
          {activeTab === "projects" && <TeamProjects teamId={teamId} />}
          {activeTab === "members" && <TeamMembers teamId={teamId} />}
          {activeTab === "settings" && <TeamSettings team={team} />}
      </div>
    </div>
  );
}