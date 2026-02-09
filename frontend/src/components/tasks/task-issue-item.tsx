"use client";

import { Task, TaskStatus, TaskPriority } from "@/types";
import { Badge } from "@/components/ui/badge";
import {
  MessageSquare,
  Circle,
  CheckCircle2,
  Clock,
  AlertCircle,
  User as UserIcon,
} from "lucide-react";
import Link from "next/link";
import { formatDistanceToNow } from "date-fns";
import { cn } from "@/lib/utils";
import { Avatar } from "@/components/ui/avatar";

interface TaskIssueItemProps {
  task: Task;
  projectId: number;
}

export function TaskIssueItem({ task, projectId }: TaskIssueItemProps) {
  const statusIcons: Record<TaskStatus, any> = {
    [TaskStatus.BACKLOG]: Clock,
    [TaskStatus.TODO]: Circle,
    [TaskStatus.IN_PROGRESS]: AlertCircle,
    [TaskStatus.DONE]: CheckCircle2,
  };

  const statusColors: Record<TaskStatus, string> = {
    [TaskStatus.BACKLOG]: "text-gray-400",
    [TaskStatus.TODO]: "text-blue-500",
    [TaskStatus.IN_PROGRESS]: "text-purple-500",
    [TaskStatus.DONE]: "text-green-500",
  };

  const priorityColors: Record<TaskPriority, string> = {
    [TaskPriority.P0]: "text-red-600 bg-red-50 border-red-100",
    [TaskPriority.P1]: "text-orange-600 bg-orange-50 border-orange-100",
    [TaskPriority.P2]: "text-blue-600 bg-blue-50 border-blue-100",
    [TaskPriority.P3]: "text-gray-500 bg-gray-50 border-gray-100",
  };

  const Icon = statusIcons[task.status] || Circle;

  return (
    <div className="group flex items-start gap-3 border-b border-gray-100 p-4 transition-colors last:border-0 hover:bg-gray-50">
      <div className={cn("mt-1 shrink-0", statusColors[task.status])}>
        <Icon className="h-5 w-5" />
      </div>

      <div className="min-w-0 flex-1 space-y-1">
        <div className="flex flex-wrap items-center gap-2">
          <Link
            href={`/projects/${projectId}/tasks/${task.id}`}
            className="font-bold text-gray-900 transition-colors hover:text-blue-600"
          >
            {task.title}
          </Link>
          <Badge
            variant="outline"
            className={cn(
              "h-5 px-1.5 text-[10px] font-bold uppercase",
              priorityColors[task.priority]
            )}
          >
            {task.priority}
          </Badge>
        </div>

        <div className="flex flex-wrap items-center gap-2 text-xs text-gray-500">
          <span>#{task.id}</span>
          <span>
            opened {formatDistanceToNow(new Date(task.created_at))} ago
          </span>
          {task.assignee && (
            <div className="flex items-center gap-1">
              <Avatar
                fallback={task.assignee.email[0].toUpperCase()}
                className="h-4 w-4 text-[8px]"
              />
              <span className="cursor-pointer hover:text-blue-600">
                {task.assignee.email.split("@")[0]}
              </span>
            </div>
          )}
        </div>
      </div>

      <div className="flex shrink-0 items-center gap-4 text-gray-400">
        {(task as any).comment_count > 0 && (
          <div className="flex items-center gap-1 text-xs">
            <MessageSquare className="h-3.5 w-3.5" />
            <span>{(task as any).comment_count}</span>
          </div>
        )}
        <div className="hidden items-center gap-1 text-xs group-hover:flex">
          <MessageSquare className="h-3.5 w-3.5" />
          <span>Details</span>
        </div>
      </div>
    </div>
  );
}
