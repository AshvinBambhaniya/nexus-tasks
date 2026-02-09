"use client";

import { useDroppable } from "@dnd-kit/core";
import {
  SortableContext,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { Task, TaskStatus } from "@/types";
import { BoardCard } from "./board-card";
import { cn } from "@/lib/utils";

interface BoardColumnProps {
  id: TaskStatus;
  title: string;
  tasks: Task[];
  onTaskClick?: (task: Task) => void;
  color?: string;
}

export function BoardColumn({
  id,
  title,
  tasks,
  onTaskClick,
  color,
}: BoardColumnProps) {
  const { setNodeRef, isOver } = useDroppable({
    id: id,
  });

  return (
    <div className="flex h-full w-80 min-w-[20rem] flex-col rounded-lg border border-gray-200 bg-gray-50/50">
      {/* Header */}
      <div className="flex items-center justify-between border-b border-gray-100 p-3">
        <div className="flex items-center gap-2">
          <div
            className={cn("h-2.5 w-2.5 rounded-full", color || "bg-gray-400")}
          />
          <h3 className="text-sm font-semibold text-gray-900">{title}</h3>
          <span className="rounded-full bg-gray-200 px-2 py-0.5 text-xs font-medium text-gray-500">
            {tasks.length}
          </span>
        </div>
      </div>

      {/* Content */}
      <div
        ref={setNodeRef}
        className={cn(
          "flex-1 space-y-2 overflow-y-auto p-2 transition-colors",
          isOver ? "bg-blue-50/50" : ""
        )}
      >
        <SortableContext
          items={tasks.map((t) => t.id)}
          strategy={verticalListSortingStrategy}
        >
          {tasks.map((task) => (
            <BoardCard key={task.id} task={task} onClick={onTaskClick} />
          ))}
        </SortableContext>

        {tasks.length === 0 && (
          <div className="flex h-24 items-center justify-center rounded-lg border-2 border-dashed border-gray-200 text-xs text-gray-400">
            Empty
          </div>
        )}
      </div>
    </div>
  );
}
