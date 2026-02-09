"use client";

import { use, useState } from "react";
import { useProject } from "@/hooks/use-projects";
import { ProjectTasks } from "@/components/project/tabs/project-tasks";
import { ProjectBoard } from "@/components/project/tabs/project-board";
import { ProjectMembers } from "@/components/project/tabs/project-members";
import { ProjectSettings } from "@/components/project/tabs/project-settings";
import { cn } from "@/lib/utils";
import { Loader2 } from "lucide-react";

export default function ProjectPage({ params }: { params: Promise<{ projectId: string }> }) {
  const resolvedParams = use(params);
  const projectId = parseInt(resolvedParams.projectId);
  
  const { project, isLoading } = useProject(projectId); 
  
  const [activeTab, setActiveTab] = useState<"tasks" | "board" | "members" | "settings">("tasks");

  if (isLoading) {
    return (
      <div className="flex h-full items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
      </div>
    );
  }

  if (!project) {
    return (
      <div className="flex h-full items-center justify-center text-gray-500">
        Project not found.
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full space-y-6">
      {/* Header */}
      <div>
         <div className="flex items-center gap-2 mb-1">
             <h1 className="text-2xl font-bold tracking-tight text-gray-900">
               {project.name}
             </h1>
             {project.is_archived && (
                 <span className="bg-yellow-100 text-yellow-800 text-xs px-2 py-0.5 rounded-full border border-yellow-200">Archived</span>
             )}
         </div>
         {project.description && (
           <p className="text-sm text-gray-500 max-w-2xl">{project.description}</p>
         )}
         
         {/* Tab Navigation */}
         <div className="flex gap-6 mt-6 border-b border-gray-200">
            {(["tasks", "board", "members", "settings"] as const).map((tab) => (
              <button
                key={tab}
                onClick={() => setActiveTab(tab)}
                className={cn(
                  "pb-3 text-sm font-medium transition-colors border-b-2 capitalize",
                  activeTab === tab
                    ? "border-blue-600 text-blue-600"
                    : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
                )}
              >
                {tab}
              </button>
            ))}
         </div>
      </div>

      {/* Tab Content */}
      <div className="flex-1 min-h-0">
          {activeTab === "tasks" && <ProjectTasks projectId={projectId} />}
          {activeTab === "board" && <ProjectBoard projectId={projectId} />}
          {activeTab === "members" && <ProjectMembers projectId={projectId} />}
          {activeTab === "settings" && <ProjectSettings project={project} />}
      </div>
    </div>
  );
}
