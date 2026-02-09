"use client";

import { use, useState } from "react";
import { useRouter } from "next/navigation";
import { useProjectMembers } from "@/hooks/use-projects";
import { useTasks } from "@/hooks/use-tasks";
import { TaskStatus, TaskPriority } from "@/types";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { MarkdownEditor } from "@/components/ui/markdown-editor";
import { Card, CardContent } from "@/components/ui/card";
import { CheckSquare, ArrowLeft, Loader2 } from "lucide-react";
import Link from "next/link";
import { AssigneeSelector } from "@/components/tasks/selectors/assignee-selector";
import { StatusSelector } from "@/components/tasks/selectors/status-selector";
import { PrioritySelector } from "@/components/tasks/selectors/priority-selector";

export default function NewTaskPage({
  params,
}: {
  params: Promise<{ projectId: string }>;
}) {
  const resolvedParams = use(params);
  const projectId = parseInt(resolvedParams.projectId);
  const router = useRouter();

  const { createTask } = useTasks(projectId);
  const { members } = useProjectMembers(projectId);

  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [status, setStatus] = useState<TaskStatus>(TaskStatus.TODO);
  const [priority, setPriority] = useState<TaskPriority>(TaskPriority.P2);
  const [assigneeId, setAssigneeId] = useState<number | undefined>(undefined);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;

    setIsSubmitting(true);
    try {
      await createTask({
        title,
        description,
        status,
        priority,
        assignee_id: assigneeId,
      });
      router.push(`/projects/${projectId}`);
    } catch (error) {
      console.error(error);
      alert("Failed to create task");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="mx-auto max-w-5xl space-y-6">
      <div className="flex items-center gap-4">
        <Link
          href={`/projects/${projectId}`}
          className="rounded-full p-2 transition-colors hover:bg-gray-100"
        >
          <ArrowLeft className="h-5 w-5 text-gray-500" />
        </Link>
        <h1 className="flex items-center gap-2 text-2xl font-bold text-gray-900">
          <CheckSquare className="h-6 w-6 text-blue-600" />
          Create new task
        </h1>
      </div>

      <form
        onSubmit={handleSubmit}
        className="grid grid-cols-1 gap-8 lg:grid-cols-4"
      >
        {/* Main Content */}
        <div className="space-y-6 lg:col-span-3">
          <Card className="border-gray-200">
            <CardContent className="space-y-4 p-6">
              <div className="space-y-2">
                <Label htmlFor="title" className="text-sm font-semibold">
                  Title
                </Label>
                <Input
                  id="title"
                  placeholder="Task title"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  className="h-12 text-lg focus-visible:ring-blue-500"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label className="text-sm font-semibold">Description</Label>
                <MarkdownEditor
                  value={description}
                  onChange={setDescription}
                  placeholder="Add a description..."
                />
              </div>
            </CardContent>
          </Card>

          <div className="flex justify-end gap-3">
            <Button type="button" variant="ghost" onClick={() => router.back()}>
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={isSubmitting || !title.trim()}
              className="px-8"
            >
              {isSubmitting && (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              )}
              Submit New Task
            </Button>
          </div>
        </div>

        {/* Sidebar */}
        <div className="space-y-6 lg:col-span-1">
          <Card className="border-gray-200 shadow-none">
            <CardContent className="space-y-6 p-4">
              <div className="space-y-2 border-b border-gray-100 pb-4">
                <Label className="text-xs font-bold tracking-wider text-gray-500 uppercase">
                  Assignee
                </Label>
                <AssigneeSelector
                  members={members}
                  value={assigneeId}
                  onChange={setAssigneeId}
                />
              </div>

              <div className="space-y-2 border-b border-gray-100 pb-4">
                <Label className="text-xs font-bold tracking-wider text-gray-500 uppercase">
                  Status
                </Label>
                <StatusSelector value={status} onChange={setStatus} />
              </div>

              <div className="space-y-2">
                <Label className="text-xs font-bold tracking-wider text-gray-500 uppercase">
                  Priority
                </Label>
                <PrioritySelector value={priority} onChange={setPriority} />
              </div>
            </CardContent>
          </Card>
        </div>
      </form>
    </div>
  );
}
