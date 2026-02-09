"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useTasks } from "@/hooks/use-tasks";
import { TaskIssueItem } from "@/components/tasks/task-issue-item";
import { Button } from "@/components/ui/button";
import {
  Plus,
  Search,
  Filter,
  CheckCircle2,
  Circle,
  Loader2,
} from "lucide-react";
import { TaskStatus } from "@/types";
import { cn } from "@/lib/utils";
import { Input } from "@/components/ui/input";

interface ProjectTasksProps {
  projectId: number;
}

export function ProjectTasks({ projectId }: ProjectTasksProps) {
  const router = useRouter();
  const { tasks, isLoading } = useTasks(projectId);
  const [statusFilter, setStatusFilter] = useState<"open" | "done">("open");
  const [searchQuery, setSearchQuery] = useState("");

  const openTasks = tasks.filter((t) => t.status !== TaskStatus.DONE);
  const doneTasks = tasks.filter((t) => t.status === TaskStatus.DONE);

  const displayedTasks = (
    statusFilter === "open" ? openTasks : doneTasks
  ).filter((t) => t.title.toLowerCase().includes(searchQuery.toLowerCase()));

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Toolbar */}
      <div className="flex flex-col items-center justify-between gap-4 sm:flex-row">
        <div className="relative w-full sm:w-96">
          <Search className="absolute top-2.5 left-3 h-4 w-4 text-gray-400" />
          <Input
            placeholder="Search all tasks"
            className="border-gray-200 bg-gray-50 pl-9"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>
        <Button onClick={() => router.push(`/projects/${projectId}/tasks/new`)}>
          <Plus className="mr-2 h-4 w-4" /> New Task
        </Button>
      </div>

      {/* GitHub Style List Container */}
      <div className="overflow-hidden rounded-lg border border-gray-200 bg-white shadow-sm">
        {/* List Header */}
        <div className="flex items-center justify-between border-b border-gray-200 bg-gray-50 px-4 py-3">
          <div className="flex items-center gap-4">
            <button
              onClick={() => setStatusFilter("open")}
              className={cn(
                "flex items-center gap-1.5 text-sm font-medium transition-colors",
                statusFilter === "open"
                  ? "text-gray-900"
                  : "text-gray-500 hover:text-gray-900"
              )}
            >
              <Circle className="h-4 w-4" />
              {openTasks.length} Open
            </button>
            <button
              onClick={() => setStatusFilter("done")}
              className={cn(
                "flex items-center gap-1.5 text-sm font-medium transition-colors",
                statusFilter === "done"
                  ? "text-gray-900"
                  : "text-gray-500 hover:text-gray-900"
              )}
            >
              <CheckCircle2 className="h-4 w-4" />
              {doneTasks.length} Done
            </button>
          </div>

          <div className="flex items-center gap-4 text-sm text-gray-500">
            {/* Future filters: Author, Label, Projects, Milestone, Assignee, Sort */}
            <button className="flex items-center gap-1 hover:text-gray-900">
              Sort <Filter className="h-3 w-3" />
            </button>
          </div>
        </div>

        {/* List Content */}
        <div className="divide-y divide-gray-100">
          {displayedTasks.map((task) => (
            <TaskIssueItem key={task.id} task={task} projectId={projectId} />
          ))}

          {displayedTasks.length === 0 && (
            <div className="p-12 text-center">
              <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center text-gray-200">
                <Plus className="h-10 w-10" />
              </div>
              <h3 className="text-lg font-medium text-gray-900">
                No tasks found
              </h3>
              <p className="mt-1 text-sm text-gray-500">
                Try adjusting your filters or search query.
              </p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
