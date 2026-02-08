"use client";

import { use, useState } from "react";
import { useTasks } from "@/hooks/use-tasks";
import { TaskListView } from "@/components/tasks/task-list-view";
import { TaskBoardView } from "@/components/tasks/task-board-view";
import { TaskDialog } from "@/components/tasks/task-dialog";
import { ProjectMemberDialog } from "@/components/project/project-member-dialog";
import { ProjectModal } from "@/components/project/project-dialog";
import { Button } from "@/components/ui/button";
import { Plus, LayoutList, LayoutDashboard, Users, MoreVertical, Edit, Archive } from "lucide-react";
import { Task, TaskStatus } from "@/types";
import { useProjects } from "@/hooks/use-projects";
import { cn } from "@/lib/utils";

export default function ProjectPage({ params }: { params: Promise<{ projectId: string }> }) {
  const resolvedParams = use(params);
  const projectId = parseInt(resolvedParams.projectId);
  
  const { tasks, isLoading: tasksLoading, updateTask } = useTasks(projectId);
  const { projects, updateProject } = useProjects(); 
  
  const project = projects.find(p => p.id === projectId);

  const [viewMode, setViewMode] = useState<"list" | "board">("list");
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [isMemberDialogOpen, setIsMemberDialogOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [isSettingsOpen, setIsSettingsOpen] = useState(false);
  const [selectedTask, setSelectedTask] = useState<Task | undefined>(undefined);

  const handleCreateClick = () => {
    setSelectedTask(undefined);
    setIsDialogOpen(true);
  };

  const handleTaskClick = (task: Task) => {
    setSelectedTask(task);
    setIsDialogOpen(true);
  };

  const handleCloseDialog = () => {
    setIsDialogOpen(false);
    setSelectedTask(undefined);
  };

  const handleTaskMove = (taskId: number, newStatus: TaskStatus) => {
    updateTask(taskId, { status: newStatus });
  };

  const handleArchive = async () => {
    if (!confirm("Are you sure you want to archive this project?")) return;
    try {
        await updateProject(projectId, { is_archived: true });
        // Optional: Redirect to dashboard or show toast
    } catch (e) {}
  }

  return (
    <div className="flex flex-col h-full space-y-6">
      <div className="flex items-center justify-between">
        <div>
           <div className="flex items-center gap-2">
               <h1 className="text-2xl font-bold tracking-tight text-gray-900">
                 {project ? project.name : "Project"}
               </h1>
               {project?.is_archived && (
                   <span className="bg-yellow-100 text-yellow-800 text-xs px-2 py-0.5 rounded-full border border-yellow-200">Archived</span>
               )}
           </div>
           {project?.description && (
             <p className="text-sm text-gray-500 mt-1">{project.description}</p>
           )}
        </div>
        <div className="flex items-center gap-2">
            <div className="relative">
                <Button variant="outline" size="icon" className="h-9 w-9" onClick={() => setIsSettingsOpen(!isSettingsOpen)}>
                    <MoreVertical className="h-4 w-4" />
                </Button>
                
                {isSettingsOpen && (
                    <>
                    <div className="fixed inset-0 z-40" onClick={() => setIsSettingsOpen(false)} />
                    <div className="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg border border-gray-200 z-50 py-1">
                        <button 
                            onClick={() => { setIsSettingsOpen(false); setIsEditModalOpen(true); }}
                            className="flex items-center w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 text-left"
                        >
                            <Edit className="h-4 w-4 mr-2" /> Edit Project
                        </button>
                        <button 
                            onClick={() => { setIsSettingsOpen(false); handleArchive(); }}
                            className="flex items-center w-full px-4 py-2 text-sm text-red-600 hover:bg-red-50 text-left"
                        >
                            <Archive className="h-4 w-4 mr-2" /> Archive Project
                        </button>
                    </div>
                    </>
                )}
            </div>

            <Button variant="outline" onClick={() => setIsMemberDialogOpen(true)}>
                <Users className="h-4 w-4 mr-2" />
                Members
            </Button>

            <div className="flex items-center rounded-lg border border-gray-200 bg-white p-1">
                <button
                    onClick={() => setViewMode("list")}
                    className={cn(
                        "flex items-center gap-2 rounded-md px-3 py-1.5 text-sm font-medium transition-colors",
                        viewMode === "list" 
                            ? "bg-blue-50 text-blue-600" 
                            : "text-gray-500 hover:bg-gray-50 hover:text-gray-900"
                    )}
                >
                    <LayoutList className="h-4 w-4" />
                    List
                </button>
                <button
                    onClick={() => setViewMode("board")}
                    className={cn(
                        "flex items-center gap-2 rounded-md px-3 py-1.5 text-sm font-medium transition-colors",
                        viewMode === "board" 
                            ? "bg-blue-50 text-blue-600" 
                            : "text-gray-500 hover:bg-gray-50 hover:text-gray-900"
                    )}
                >
                    <LayoutDashboard className="h-4 w-4" />
                    Board
                </button>
            </div>
            <Button onClick={handleCreateClick}>
                <Plus className="mr-2 h-4 w-4" /> Create Task
            </Button>
        </div>
      </div>

      {tasksLoading ? (
        <div className="flex h-64 items-center justify-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
        </div>
      ) : (
        <div className="flex-1 overflow-hidden">
            {viewMode === "list" ? (
                <TaskListView tasks={tasks} onTaskClick={handleTaskClick} />
            ) : (
                <TaskBoardView 
                    tasks={tasks} 
                    onTaskMove={handleTaskMove}
                    onTaskClick={handleTaskClick}
                />
            )}
        </div>
      )}

      <TaskDialog 
        isOpen={isDialogOpen} 
        onClose={handleCloseDialog} 
        task={selectedTask}
        projectId={projectId}
      />

      <ProjectMemberDialog 
        isOpen={isMemberDialogOpen}
        onClose={() => setIsMemberDialogOpen(false)}
        projectId={projectId}
      />

      <ProjectModal 
        isOpen={isEditModalOpen} 
        onClose={() => setIsEditModalOpen(false)} 
        project={project} 
      />
    </div>
  );
}
