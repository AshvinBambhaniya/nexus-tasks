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
    isDragging 
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
            className="opacity-50 h-[120px] rounded-lg bg-gray-100 border-2 border-dashed border-gray-300"
          />
      )
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
        className="cursor-grab hover:shadow-md transition-all active:cursor-grabbing border-gray-200"
        onClick={() => onClick?.(task)}
      >
        <CardContent className="p-3 space-y-3">
          <div className="flex justify-between items-start gap-2">
             <span className="text-sm font-medium text-gray-900 line-clamp-2 leading-tight">
                {task.title}
             </span>
          </div>
          
          <div className="flex items-center justify-between">
             <Badge variant="outline" className={cn("text-[10px] px-1.5 py-0 h-5 font-bold border", priorityColors[task.priority])}>
                {task.priority}
             </Badge>
             
             <div className="flex items-center gap-2">
                {(task as any).comment_count > 0 && (
                    <div className="flex items-center gap-1 text-xs text-gray-400">
                        <MessageSquare className="h-3 w-3" />
                        <span>{(task as any).comment_count}</span>
                    </div>
                )}
                {task.assignee && (
                    <Avatar 
                        fallback={task.assignee.email[0].toUpperCase()} 
                        className="h-5 w-5 text-[9px] border border-gray-100"
                    />
                )}
             </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
