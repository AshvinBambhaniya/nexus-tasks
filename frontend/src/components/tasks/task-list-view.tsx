"use client";

import { Task, TaskPriority, TaskStatus } from "@/types";
import { Badge } from "@/components/ui/badge";

interface TaskListViewProps {
  tasks: Task[];
  onTaskClick?: (task: Task) => void;
}

export function TaskListView({ tasks, onTaskClick }: TaskListViewProps) {
  if (tasks.length === 0) {
    return (
      <div className="flex h-64 items-center justify-center rounded-lg border border-dashed border-gray-200 bg-gray-50 text-gray-500">
        No tasks found. Create one to get started!
      </div>
    );
  }

  return (
    <div className="rounded-md border border-gray-200">
      <table className="w-full text-sm text-left">
        <thead className="bg-gray-50 text-xs uppercase text-gray-700">
          <tr>
            <th className="px-4 py-3 font-medium">Title</th>
            <th className="px-4 py-3 font-medium">Status</th>
            <th className="px-4 py-3 font-medium">Priority</th>
            <th className="px-4 py-3 font-medium">Created</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-200">
          {tasks.map((task) => (
            <tr 
              key={task.id} 
              className="bg-white hover:bg-gray-50 cursor-pointer transition-colors"
              onClick={() => onTaskClick?.(task)}
            >
              <td className="px-4 py-3 font-medium text-gray-900">{task.title}</td>
              <td className="px-4 py-3">
                <StatusBadge status={task.status} />
              </td>
              <td className="px-4 py-3">
                <PriorityBadge priority={task.priority} />
              </td>
              <td className="px-4 py-3 text-gray-500">
                {new Date(task.created_at).toLocaleDateString()}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

function StatusBadge({ status }: { status: TaskStatus }) {
  const variants: Record<TaskStatus, "default" | "secondary" | "success" | "warning"> = {
    [TaskStatus.TODO]: "secondary",
    [TaskStatus.IN_PROGRESS]: "default",
    [TaskStatus.DONE]: "success",
    [TaskStatus.BACKLOG]: "secondary", // Fallback
  };
  
  return <Badge variant={variants[status] || "outline"}>{status.replace("_", " ")}</Badge>;
}

function PriorityBadge({ priority }: { priority: TaskPriority }) {
    const colors: Record<TaskPriority, string> = {
        [TaskPriority.P0]: "text-red-600 bg-red-50 border-red-200",
        [TaskPriority.P1]: "text-orange-600 bg-orange-50 border-orange-200",
        [TaskPriority.P2]: "text-blue-600 bg-blue-50 border-blue-200",
        [TaskPriority.P3]: "text-gray-600 bg-gray-50 border-gray-200",
    };

    return (
        <span className={`inline-flex items-center rounded-md border px-2 py-1 text-xs font-medium ${colors[priority]}`}>
            {priority}
        </span>
    );
}
