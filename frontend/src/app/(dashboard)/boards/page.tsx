"use client";

import { useState } from "react";
import { useMyTasks } from "@/hooks/use-my-tasks"; // Changed hook
import { BoardView } from "@/components/tasks/board/board-view";
import { TaskDialog } from "@/components/tasks/task-dialog";
import { Button } from "@/components/ui/button";
import { Plus } from "lucide-react";
import { Task, TaskStatus } from "@/types";

export default function BoardsPage() {
  const { tasks, isLoading, updateTask } = useMyTasks();
  const [isDialogOpen, setIsDialogOpen] = useState(false);

  const handleTaskMove = (taskId: number, newStatus: TaskStatus) => {
    updateTask(taskId, { status: newStatus });
  };

  if (isLoading) {
    return <div className="p-8 text-center">Loading...</div>;
  }

  return (
    <div className="flex h-full flex-col space-y-4">
      <div className="flex items-center justify-between px-2">
        <h1 className="text-2xl font-bold">My Boards</h1>
        <Button onClick={() => setIsDialogOpen(true)}>
          <Plus className="mr-2 h-4 w-4" /> New Task
        </Button>
      </div>

      <div className="flex-1 overflow-hidden">
        <BoardView tasks={tasks} onTaskMove={handleTaskMove} />
      </div>

      <TaskDialog
        isOpen={isDialogOpen}
        onClose={() => setIsDialogOpen(false)}
      />
    </div>
  );
}
