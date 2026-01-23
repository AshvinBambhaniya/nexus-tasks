"use client";

import { useMyTasks } from "@/hooks/use-my-tasks";
import { TaskCard } from "@/components/tasks/task-card";
import { Loader2 } from "lucide-react";
import { TaskWithWorkspace } from "@/types";
import { format } from "date-fns";
import { Badge } from "@/components/ui/badge";

export default function InboxPage() {
  const { tasks, isLoading } = useMyTasks();

  if (isLoading) {
    return (
      <div className="flex h-full items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-blue-500" />
      </div>
    );
  }

  // Group tasks by status or just list them? 
  // The requirement says: "Display a unified list sorted by Priority/Due Date."
  // So a simple list is fine.
  
  return (
    <div className="h-full flex flex-col p-6 space-y-6 overflow-hidden">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">My Focus</h1>
        <p className="text-gray-500">All your tasks across all workspaces.</p>
      </div>

      <div className="flex-1 overflow-y-auto pr-2 space-y-3">
        {tasks.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-64 text-gray-400">
            <p>No tasks assigned to you yet.</p>
          </div>
        ) : (
          tasks.map((task) => (
            <InboxTaskCard key={task.id} task={task} />
          ))
        )}
      </div>
    </div>
  );
}

function InboxTaskCard({ task }: { task: TaskWithWorkspace }) {
  // We can reuse TaskCard but we might want to show workspace name.
  // Existing TaskCard might not support showing workspace name.
  // Let's create a simple card here or check TaskCard.
  
  // For now, I'll inline a simple card that shows workspace info.
  return (
      <div className="bg-white p-4 rounded-lg border border-gray-200 shadow-sm hover:shadow-md transition-shadow flex flex-col gap-2">
          <div className="flex items-start justify-between">
              <div className="space-y-1">
                  <div className="flex items-center gap-2">
                      <span className="text-xs font-medium text-gray-500 uppercase tracking-wider">
                          {task.workspace.name}
                      </span>
                      <Badge variant="outline" className={getPrioColor(task.priority)}>
                          {task.priority}
                      </Badge>
                  </div>
                  <h3 className="font-medium text-gray-900">{task.title}</h3>
              </div>
              <Badge className={getStatusColor(task.status)}>{task.status.replace("_", " ")}</Badge>
          </div>
          
          {task.description && (
              <p className="text-sm text-gray-500 line-clamp-2">{task.description}</p>
          )}

          {task.due_date && (
              <div className="text-xs text-gray-400 mt-2">
                  Due {format(new Date(task.due_date), "PPP")}
              </div>
          )}
      </div>
  )
}

function getPrioColor(prio: string) {
    switch (prio) {
        case "P0": return "text-red-600 border-red-200 bg-red-50";
        case "P1": return "text-orange-600 border-orange-200 bg-orange-50";
        case "P2": return "text-blue-600 border-blue-200 bg-blue-50";
        default: return "text-gray-600 border-gray-200 bg-gray-50";
    }
}

function getStatusColor(status: string) {
    switch (status) {
        case "TODO": return "bg-gray-100 text-gray-700 hover:bg-gray-200";
        case "IN_PROGRESS": return "bg-blue-100 text-blue-700 hover:bg-blue-200";
        case "DONE": return "bg-green-100 text-green-700 hover:bg-green-200";
        default: return "bg-gray-100 text-gray-700";
    }
}
