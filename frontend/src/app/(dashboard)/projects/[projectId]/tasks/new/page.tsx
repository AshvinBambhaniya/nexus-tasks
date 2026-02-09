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

export default function NewTaskPage({ params }: { params: Promise<{ projectId: string }> }) {
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
      } as any);
      router.push(`/projects/${projectId}`);
    } catch (error) {
      console.error(error);
      alert("Failed to create task");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="max-w-5xl mx-auto space-y-6">
      <div className="flex items-center gap-4">
        <Link href={`/projects/${projectId}`} className="p-2 hover:bg-gray-100 rounded-full transition-colors">
          <ArrowLeft className="h-5 w-5 text-gray-500" />
        </Link>
        <h1 className="text-2xl font-bold text-gray-900 flex items-center gap-2">
          <CheckSquare className="h-6 w-6 text-blue-600" />
          Create new task
        </h1>
      </div>

      <form onSubmit={handleSubmit} className="grid grid-cols-1 lg:grid-cols-4 gap-8">
        {/* Main Content */}
        <div className="lg:col-span-3 space-y-6">
          <Card className="border-gray-200">
            <CardContent className="p-6 space-y-4">
              <div className="space-y-2">
                <Label htmlFor="title" className="text-sm font-semibold">Title</Label>
                <Input
                  id="title"
                  placeholder="Task title"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  className="text-lg h-12 focus-visible:ring-blue-500"
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
             <Button type="button" variant="ghost" onClick={() => router.back()}>Cancel</Button>
             <Button type="submit" disabled={isSubmitting || !title.trim()} className="px-8">
                {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Submit New Task
             </Button>
          </div>
        </div>

        {/* Sidebar */}
        <div className="lg:col-span-1 space-y-6">
          <Card className="border-gray-200 shadow-none">
            <CardContent className="p-4 space-y-6">
              <div className="space-y-2 border-b border-gray-100 pb-4">
                <Label className="text-xs font-bold text-gray-500 uppercase tracking-wider">Assignees</Label>
                <select
                  className="w-full text-sm border-gray-200 rounded-md focus:ring-blue-500 focus:border-blue-500"
                  value={assigneeId || ""}
                  onChange={(e) => setAssigneeId(e.target.value ? parseInt(e.target.value) : undefined)}
                >
                  <option value="">Unassigned</option>
                  {members.map((m) => (
                    <option key={m.user_id} value={m.user_id}>{m.email}</option>
                  ))}
                </select>
              </div>

              <div className="space-y-2 border-b border-gray-100 pb-4">
                <Label className="text-xs font-bold text-gray-500 uppercase tracking-wider">Status</Label>
                <select
                  className="w-full text-sm border-gray-200 rounded-md focus:ring-blue-500 focus:border-blue-500"
                  value={status}
                  onChange={(e) => setStatus(e.target.value as TaskStatus)}
                >
                  {Object.values(TaskStatus).map((s) => (
                    <option key={s} value={s}>{s.replace("_", " ")}</option>
                  ))}
                </select>
              </div>

              <div className="space-y-2">
                <Label className="text-xs font-bold text-gray-500 uppercase tracking-wider">Priority</Label>
                <select
                  className="w-full text-sm border-gray-200 rounded-md focus:ring-blue-500 focus:border-blue-500"
                  value={priority}
                  onChange={(e) => setPriority(e.target.value as TaskPriority)}
                >
                  <option value={TaskPriority.P0}>P0 - Critical</option>
                  <option value={TaskPriority.P1}>P1 - High</option>
                  <option value={TaskPriority.P2}>P2 - Medium</option>
                  <option value={TaskPriority.P3}>P3 - Low</option>
                </select>
              </div>
            </CardContent>
          </Card>
        </div>
      </form>
    </div>
  );
}
