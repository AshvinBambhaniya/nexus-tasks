"use client";

import Link from "next/link";
import useSWR from "swr";
import { format } from "date-fns";
import { 
  Folder, 
  Users, 
  CheckSquare, 
  Plus, 
  ArrowRight,
  Activity,
  Briefcase
} from "lucide-react";

import api from "@/lib/api";
import { cn } from "@/lib/utils";
import { useUser } from "@/hooks/use-user";
import { useWorkspaces } from "@/hooks/use-workspaces";
import { useProjects } from "@/hooks/use-projects";
import { useMyTasks } from "@/hooks/use-my-tasks";
import { WorkspaceMember, TaskStatus } from "@/types";

import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { ProjectDialog } from "@/components/project/project-dialog";
import { TaskDialog } from "@/components/tasks/task-dialog";
import { useState } from "react";

const fetcher = (url: string) => api.get(url).then((res) => res.data);

export default function DashboardPage() {
  const { user } = useUser();
  const { activeWorkspace } = useWorkspaces();
  const { projects, isLoading: projectsLoading } = useProjects();
  const { tasks: myTasks, isLoading: tasksLoading } = useMyTasks();
  
  // Fetch members for stats
  const { data: members } = useSWR<WorkspaceMember[]>(
    activeWorkspace ? `/api/v1/workspaces/${activeWorkspace.id}/members` : null,
    fetcher
  );

  const [isTaskDialogOpen, setIsTaskDialogOpen] = useState(false);
  const [isProjectDialogOpen, setIsProjectDialogOpen] = useState(false); // Controlled by Dialog Trigger usually, but for button we might need state

  // Filter tasks for this workspace (since /me returns all)
  // Assuming MyTasks API might return workspace_id in the future, currently we filter by project which has workspace_id
  // Actually, TaskWithProject has `project: { id, name }`. 
  // We need to know which workspace the project belongs to. 
  // The current API might not return workspace_id for the project in /tasks/me. 
  // Optimization: Just show "My Tasks" generally or filter if we can match project IDs with `projects` list.
  const workspaceProjectIds = new Set(projects.map(p => p.id));
  const workspaceTasks = myTasks.filter(t => workspaceProjectIds.has(t.project_id));
  
  const pendingTasksCount = workspaceTasks.filter(t => t.status !== TaskStatus.DONE).length;
  const recentTasks = workspaceTasks
    .filter(t => t.status !== TaskStatus.DONE)
    .sort((a, b) => {
       // Sort by priority (P0 > P1...) then due date
       if (a.priority !== b.priority) return a.priority.localeCompare(b.priority);
       return new Date(a.due_date || "2100-01-01").getTime() - new Date(b.due_date || "2100-01-01").getTime();
    })
    .slice(0, 5);

  const getTimeGreeting = () => {
    const hour = new Date().getHours();
    if (hour < 12) return "Good morning";
    if (hour < 18) return "Good afternoon";
    return "Good evening";
  };

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex flex-col gap-1 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">
            {getTimeGreeting()}, {user?.full_name?.split(" ")[0] || "there"}
          </h1>
          <p className="text-gray-500">
            Here&apos;s what&apos;s happening in <span className="font-medium text-gray-900">{activeWorkspace?.name}</span> today.
          </p>
        </div>
        <div className="flex gap-2 mt-4 sm:mt-0">
          <Button 
            onClick={() => setIsTaskDialogOpen(true)} 
            size="sm" 
            className="hidden sm:flex"
            disabled={projects.length === 0}
            title={projects.length === 0 ? "Create a project first" : "Create a new task"}
          >
            <Plus className="mr-2 h-4 w-4" /> New Task
          </Button>
          <TaskDialog 
            isOpen={isTaskDialogOpen} 
            onClose={() => setIsTaskDialogOpen(false)} 
            projectId={projects[0]?.id} // Default to first project if any
          />
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatsCard 
          title="Active Projects" 
          value={projects.filter(p => !p.is_archived).length} 
          icon={Folder}
          isLoading={projectsLoading}
          color="blue"
        />
        <StatsCard 
          title="Team Members" 
          value={members?.length || 0} 
          icon={Users}
          isLoading={!members}
          color="indigo"
        />
        <StatsCard 
          title="Pending Tasks" 
          value={pendingTasksCount} 
          icon={CheckSquare}
          isLoading={tasksLoading}
          color="orange"
        />
        <StatsCard 
          title="Workspace Activity" 
          value="Unknown" // Placeholder for now
          label="Last 7 days"
          icon={Activity}
          isLoading={false}
          color="green"
          simple
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Main Column: Projects */}
        <div className="lg:col-span-2 space-y-6">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold text-gray-900">Recent Projects</h2>
            <Link href="/projects/new" className="text-sm text-blue-600 hover:text-blue-700 font-medium">
               View all
            </Link>
          </div>
          
          {projectsLoading ? (
             <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
               {[1,2,3,4].map(i => <SkeletonCard key={i} />)}
             </div>
          ) : projects.length === 0 ? (
            <div className="rounded-lg border-2 border-dashed border-gray-200 p-8 text-center bg-gray-50">
               <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-gray-100">
                 <Briefcase className="h-6 w-6 text-gray-400" />
               </div>
               <h3 className="mt-2 text-sm font-semibold text-gray-900">No projects</h3>
               <p className="mt-1 text-sm text-gray-500">Get started by creating a new project.</p>
               <div className="mt-6">
                 {/* This would be the ProjectDialog trigger usually */}
                 <Button variant="outline" size="sm">Create Project</Button>
               </div>
            </div>
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              {projects.slice(0, 4).map((project) => (
                <Link key={project.id} href={`/projects/${project.id}`} className="block h-full">
                  <Card className="h-full p-5 hover:shadow-md transition-shadow cursor-pointer border-gray-200 hover:border-blue-200 group">
                    <div className="flex justify-between items-start mb-4">
                      <div className="p-2 bg-blue-50 rounded-lg group-hover:bg-blue-100 transition-colors">
                        <Folder className="h-5 w-5 text-blue-600" />
                      </div>
                      {project.is_archived && <Badge variant="secondary">Archived</Badge>}
                    </div>
                    <h3 className="font-semibold text-gray-900 mb-1 group-hover:text-blue-700 transition-colors">
                      {project.name}
                    </h3>
                    <p className="text-sm text-gray-500 line-clamp-2">
                      {project.description || "No description provided."}
                    </p>
                    <div className="mt-4 pt-4 border-t border-gray-100 flex items-center justify-between text-xs text-gray-500">
                      <span>Updated {format(new Date(project.created_at), "MMM d")}</span>
                      <ArrowRight className="h-3 w-3 opacity-0 group-hover:opacity-100 transition-opacity transform group-hover:translate-x-1" />
                    </div>
                  </Card>
                </Link>
              ))}
            </div>
          )}
        </div>

        {/* Side Column: My Tasks */}
        <div className="space-y-6">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold text-gray-900">My Priorities</h2>
            <Link href="/inbox" className="text-sm text-blue-600 hover:text-blue-700 font-medium">
               Go to Inbox
            </Link>
          </div>

          <div className="space-y-3">
            {tasksLoading ? (
               <div className="space-y-3">
                 {[1,2,3].map(i => <div key={i} className="h-16 bg-gray-100 rounded-lg animate-pulse" />)}
               </div>
            ) : recentTasks.length === 0 ? (
              <div className="rounded-lg border border-gray-200 p-6 text-center bg-gray-50">
                 <p className="text-sm text-gray-500">No pending tasks in this workspace.</p>
              </div>
            ) : (
              recentTasks.map(task => (
                <div key={task.id} className="group flex items-start gap-3 p-3 rounded-lg border border-gray-200 bg-white hover:border-blue-200 hover:shadow-sm transition-all">
                  <div className={cn("mt-0.5 h-2 w-2 rounded-full flex-shrink-0", 
                    task.priority === "P0" ? "bg-red-500" : 
                    task.priority === "P1" ? "bg-orange-500" : "bg-blue-500"
                  )} />
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-gray-900 truncate">{task.title}</p>
                    <div className="flex items-center gap-2 mt-1">
                       <span className="text-[10px] text-gray-500 uppercase tracking-wider">{task.project?.name}</span>
                       {task.due_date && (
                         <span className="text-[10px] text-gray-400">
                           {format(new Date(task.due_date), "MMM d")}
                         </span>
                       )}
                    </div>
                  </div>
                </div>
              ))
            )}
            
            {workspaceTasks.length > 5 && (
               <Link href="/inbox" className="block text-center text-xs text-gray-500 hover:text-gray-900 mt-2">
                 + {workspaceTasks.length - 5} more tasks
               </Link>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

function StatsCard({ title, value, icon: Icon, isLoading, color, label, simple }: any) {
  const colors = {
    blue: "bg-blue-50 text-blue-600",
    indigo: "bg-indigo-50 text-indigo-600",
    orange: "bg-orange-50 text-orange-600",
    green: "bg-green-50 text-green-600",
  };

  return (
    <Card className="p-5 border-gray-200 shadow-sm">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm font-medium text-gray-500">{title}</p>
          {isLoading ? (
             <div className="h-8 w-16 bg-gray-100 rounded animate-pulse mt-1" />
          ) : (
             <div className="flex items-baseline gap-2 mt-1">
                <span className="text-2xl font-bold text-gray-900">{value}</span>
                {label && <span className="text-xs text-gray-500 font-normal">{label}</span>}
             </div>
          )}
        </div>
        <div className={cn("p-2 rounded-lg", (colors as any)[color] || colors.blue)}>
          <Icon className="h-5 w-5" />
        </div>
      </div>
    </Card>
  );
}

function SkeletonCard() {
    return <div className="h-40 bg-gray-100 rounded-lg animate-pulse border border-gray-200" />
}