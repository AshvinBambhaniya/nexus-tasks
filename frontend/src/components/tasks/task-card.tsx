import { Task, TaskWithProject } from "@/types";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Avatar } from "@/components/ui/avatar";
import { Clock, CheckCircle2 } from "lucide-react";
import { format, isPast, isToday } from "date-fns";
import { cn } from "@/lib/utils";

interface TaskCardProps {
  task: Task | TaskWithProject;
}

export function TaskCard({ task }: TaskCardProps) {
  const project = (task as TaskWithProject).project;

  return (
    <Card className="cursor-grab bg-white p-3 transition-shadow hover:shadow-md">
      <div className="flex flex-col gap-2">
        <div className="space-y-1">
          {project && (
            <span className="block text-[10px] font-semibold tracking-wider text-gray-500 uppercase">
              {project.name}
            </span>
          )}
          <span className="line-clamp-2 text-sm font-medium text-gray-900">
            {task.title}
          </span>
        </div>
        <div className="mt-1 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Badge variant="outline" className="h-5 px-1.5 py-0 text-[10px]">
              {task.priority}
            </Badge>
            {task.status === "DONE" && task.completed_at ? (
              <div className="flex items-center gap-1 text-[10px] text-green-600 font-medium">
                <CheckCircle2 className="h-3 w-3" />
                <span>{format(new Date(task.completed_at), "MMM d")}</span>
              </div>
            ) : (
              task.due_date && (
                <div
                  className={cn(
                    "flex items-center gap-1 text-[10px]",
                    isPast(new Date(task.due_date)) &&
                      !isToday(new Date(task.due_date)) &&
                      task.status !== "DONE"
                      ? "font-medium text-red-600"
                      : "text-gray-500"
                  )}
                >
                  <Clock className="h-3 w-3" />
                  <span>{format(new Date(task.due_date), "MMM d")}</span>
                </div>
              )
            )}
          </div>
          {task.assignee && (
            <Avatar
              className="h-5 w-5"
              fallback={(task.assignee.full_name ||
                task.assignee.email)[0].toUpperCase()}
              title={task.assignee.full_name || task.assignee.email}
            />
          )}
        </div>
      </div>
    </Card>
  );
}
