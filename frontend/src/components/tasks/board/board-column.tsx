"use client";

import { useDroppable } from "@dnd-kit/core";
import { SortableContext, verticalListSortingStrategy } from "@dnd-kit/sortable";
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

export function BoardColumn({ id, title, tasks, onTaskClick, color }: BoardColumnProps) {
  const { setNodeRef, isOver } = useDroppable({
    id: id,
  });

  return (
    <div className="flex h-full w-80 min-w-[20rem] flex-col rounded-lg bg-gray-50/50 border border-gray-200">
      {/* Header */}
      <div className="p-3 flex items-center justify-between border-b border-gray-100">
        <div className="flex items-center gap-2">
            <div className={cn("h-2.5 w-2.5 rounded-full", color || "bg-gray-400")} />
            <h3 className="font-semibold text-sm text-gray-900">{title}</h3>
            <span className="text-xs font-medium text-gray-500 bg-gray-200 px-2 py-0.5 rounded-full">
                {tasks.length}
            </span>
        </div>
      </div>

      {/* Content */}
      <div 
        ref={setNodeRef} 
        className={cn(
            "flex-1 overflow-y-auto p-2 space-y-2 transition-colors",
            isOver ? "bg-blue-50/50" : ""
        )}
      >
        <SortableContext items={tasks.map(t => t.id)} strategy={verticalListSortingStrategy}>
            {tasks.map((task) => (
              <BoardCard key={task.id} task={task} onClick={onTaskClick} />
            ))}
        </SortableContext>
        
        {tasks.length === 0 && (
            <div className="h-24 border-2 border-dashed border-gray-200 rounded-lg flex items-center justify-center text-xs text-gray-400">
                Empty
            </div>
        )}
      </div>
    </div>
  );
}
