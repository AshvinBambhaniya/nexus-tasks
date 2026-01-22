"use client";

import { useDraggable } from "@dnd-kit/core";
import { Task } from "@/types";
import { TaskCard } from "./task-card";

interface DraggableTaskCardProps {
  task: Task;
  onTaskClick?: (task: Task) => void;
}

export function DraggableTaskCard({ task, onTaskClick }: DraggableTaskCardProps) {
  const { attributes, listeners, setNodeRef, transform, isDragging } = useDraggable({
    id: task.id,
    data: { status: task.status },
  });

  const style = transform ? {
    transform: `translate3d(${transform.x}px, ${transform.y}px, 0)`,
  } : undefined;

  if (isDragging) {
      return (
        <div 
            ref={setNodeRef} 
            style={style} 
            className="opacity-50"
        >
             <TaskCard task={task} />
        </div>
      )
  }

  return (
    <div ref={setNodeRef} style={style} {...listeners} {...attributes} onClick={() => onTaskClick?.(task)}>
      <TaskCard task={task} />
    </div>
  );
}
