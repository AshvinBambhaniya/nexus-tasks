"use client";

import { useState } from "react";
import { useMyTasks } from "@/hooks/use-my-tasks"; // Changed hook
import { TaskBoardView } from "@/components/tasks/task-board-view";
import { TaskDialog } from "@/components/tasks/task-dialog";
import { Button } from "@/components/ui/button";
import { Plus } from "lucide-react";
import { Task, TaskStatus } from "@/types";

export default function BoardsPage() {
  const { tasks, isLoading, updateTask } = useMyTasks(); // Use my tasks hook
  
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [selectedTask, setSelectedTask] = useState<Task | undefined>(undefined);

  const handleTaskMove = (taskId: number, newStatus: TaskStatus) => {
    updateTask(taskId, { status: newStatus });
  };

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

  return (
    <div className="flex flex-col h-full space-y-6">
      <div className="flex items-center justify-between">
        <div>
            <h1 className="text-2xl font-bold tracking-tight text-gray-900">My Board</h1>
            <p className="text-sm text-gray-500 mt-1">Board view of all tasks assigned to you across projects.</p>
        </div>
        {/* Creating a task from "My Board" is tricky because we need a project ID. 
            The TaskDialog might support selecting a project, or we might disable creation here.
            For now, let's keep the button but it might fail if TaskDialog requires projectId 
            and we don't pass it.
        */}
        {/* <Button onClick={handleCreateClick}>
            <Plus className="mr-2 h-4 w-4" /> Create Task
        </Button> */}
      </div>

      {isLoading ? (
        <div className="flex h-64 items-center justify-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
        </div>
      ) : (
        <div className="flex-1 overflow-hidden">
             <TaskBoardView 
                tasks={tasks} 
                onTaskMove={handleTaskMove} 
                onTaskClick={handleTaskClick} 
             />
        </div>
      )}

      {/* 
        TaskDialog requires a projectId if we are creating a new task.
        If we are editing (selectedTask is set), it has a project_id in it.
        So editing should work. Creating new might be disabled or need enhancement.
        For now, I'll comment out the create button above and only allow editing.
      */}
      <TaskDialog 
        isOpen={isDialogOpen} 
        onClose={handleCloseDialog} 
        task={selectedTask}
        projectId={selectedTask?.project_id || 0} // Fallback, creation disabled anyway
      />
    </div>
  );
}
