"use client";

import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { Task, TaskPriority } from "@/types";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Avatar } from "@/components/ui/avatar";
import { MessageSquare } from "lucide-react";
import { cn } from "@/lib/utils";

interface BoardCardProps {
  task: Task;
  onClick?: (task: Task) => void;
}

export function BoardCard({ task, onClick }: BoardCardProps) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: task.id,
    data: { task },
  });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  const priorityColors: Record<TaskPriority, string> = {
    [TaskPriority.P0]: "text-red-600 bg-red-50 border-red-100",
    [TaskPriority.P1]: "text-orange-600 bg-orange-50 border-orange-100",
    [TaskPriority.P2]: "text-blue-600 bg-blue-50 border-blue-100",
    [TaskPriority.P3]: "text-gray-500 bg-gray-50 border-gray-100",
  };

  if (isDragging) {
    return (
      <div
        ref={setNodeRef}
        style={style}
        className="h-[120px] rounded-lg border-2 border-dashed border-gray-300 bg-gray-100 opacity-50"
      />
    );
  }

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...listeners}
      {...attributes}
      className="touch-none"
    >
      <Card
        className="cursor-grab border-gray-200 transition-all hover:shadow-md active:cursor-grabbing"
        onClick={() => onClick?.(task)}
      >
        <CardContent className="space-y-3 p-3">
          <div className="flex items-start justify-between gap-2">
            <span className="line-clamp-2 text-sm leading-tight font-medium text-gray-900">
              {task.title}
            </span>
          </div>

          <div className="flex items-center justify-between">
            <Badge
              variant="outline"
              className={cn(
                "h-5 border px-1.5 py-0 text-[10px] font-bold",
                priorityColors[task.priority]
              )}
            >
              {task.priority}
            </Badge>

            <div className="flex items-center gap-2">
              {(task.comment_count ?? 0) > 0 && (
                <div className="flex items-center gap-1 text-xs text-gray-400">
                  <MessageSquare className="h-3 w-3" />
                  <span>{task.comment_count}</span>
                </div>
              )}
              {task.assignee && (
                <Avatar
                  fallback={task.assignee.email[0].toUpperCase()}
                  className="h-5 w-5 border border-gray-100 text-[9px]"
                />
              )}
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
