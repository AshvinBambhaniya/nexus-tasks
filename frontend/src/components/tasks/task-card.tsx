import { Task } from "@/types";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

interface TaskCardProps {
  task: Task;
}

export function TaskCard({ task }: TaskCardProps) {
  return (
    <Card className="cursor-grab p-3 hover:shadow-md transition-shadow bg-white">
      <div className="flex flex-col gap-2">
        <span className="font-medium text-sm text-gray-900 line-clamp-2">{task.title}</span>
        <div className="flex items-center justify-between mt-1">
             <Badge variant="outline" className="text-[10px] px-1.5 py-0 h-5">
                {task.priority}
             </Badge>
             {/* Assignee Avatar could go here */}
        </div>
      </div>
    </Card>
  );
}
