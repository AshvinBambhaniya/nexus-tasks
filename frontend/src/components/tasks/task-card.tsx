import { Task, TaskWithProject } from "@/types";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Avatar } from "@/components/ui/avatar";

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
          <Badge variant="outline" className="h-5 px-1.5 py-0 text-[10px]">
            {task.priority}
          </Badge>
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
