"use client";

import { useDroppable } from "@dnd-kit/core";
import { Task } from "@/types";
import { DraggableTaskCard } from "./draggable-task-card";

interface BoardColumnProps {
  id: string;
  title: string;
  tasks: Task[];
  onTaskClick?: (task: Task) => void;
}

export function BoardColumn({ id, title, tasks, onTaskClick }: BoardColumnProps) {
  const { setNodeRef } = useDroppable({
    id: id,
  });

  return (
    <div ref={setNodeRef} className="flex h-full w-80 flex-col rounded-lg bg-gray-100/50 border border-gray-200">
      <div className="p-4 font-medium text-gray-700 flex items-center justify-between">
        {title}
        <span className="text-xs text-gray-400 bg-gray-200 px-2 py-0.5 rounded-full">{tasks.length}</span>
      </div>
      <div className="flex-1 space-y-3 overflow-y-auto p-3">
        {tasks.map((task) => (
          <DraggableTaskCard key={task.id} task={task} onTaskClick={onTaskClick} />
        ))}
      </div>
    </div>
  );
}
