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
    <Card className="cursor-grab p-3 hover:shadow-md transition-shadow bg-white">
      <div className="flex flex-col gap-2">
        <div className="space-y-1">
            {project && (
                <span className="text-[10px] font-semibold text-gray-500 uppercase tracking-wider block">
                    {project.name}
                </span>
            )}
            <span className="font-medium text-sm text-gray-900 line-clamp-2">{task.title}</span>
        </div>
        <div className="flex items-center justify-between mt-1">
             <Badge variant="outline" className="text-[10px] px-1.5 py-0 h-5">
                {task.priority}
             </Badge>
             {task.assignee && (
                <Avatar 
                  className="h-5 w-5" 
                  fallback={(task.assignee.full_name || task.assignee.email)[0].toUpperCase()} 
                  title={task.assignee.full_name || task.assignee.email} 
                />
             )}
        </div>
      </div>
    </Card>
  );
}
