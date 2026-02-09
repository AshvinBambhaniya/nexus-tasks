"use client";

import { useTasks } from "@/hooks/use-tasks";
import { BoardView } from "@/components/tasks/board/board-view";
import { Loader2 } from "lucide-react";
import { TaskStatus, Task } from "@/types";
import { useRouter } from "next/navigation";

interface ProjectBoardProps {
  projectId: number;
}

export function ProjectBoard({ projectId }: ProjectBoardProps) {
  const { tasks, isLoading, updateTask } = useTasks(projectId);
  const router = useRouter();

  const handleTaskMove = (taskId: number, newStatus: TaskStatus) => {
    updateTask(taskId, { status: newStatus });
  };

  const handleTaskClick = (task: Task) => {
    router.push(`/projects/${projectId}/tasks/${task.id}`);
  };

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
      </div>
    );
  }

  return (
    <div className="h-full">
      <BoardView
        tasks={tasks}
        onTaskMove={handleTaskMove}
        onTaskClick={handleTaskClick}
      />
    </div>
  );
}
