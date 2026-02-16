"use client";

import { use, useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useTask, useTasks } from "@/hooks/use-tasks";
import { useProjectMembers } from "@/hooks/use-projects";
import { TaskStatus, TaskPriority, ApiError } from "@/types";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { MarkdownEditor } from "@/components/ui/markdown-editor";
import { Avatar } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { CommentItem } from "@/components/tasks/comment-item";
import {
  ArrowLeft,
  Loader2,
  Trash2,
  Settings2,
  Clock,
  User as UserIcon,
  CheckCircle2,
  Circle,
  AlertCircle,
} from "lucide-react";
import Link from "next/link";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { cn } from "@/lib/utils";
import { formatDistanceToNow } from "date-fns";

import { useUser } from "@/hooks/use-user";

import { AssigneeSelector } from "@/components/tasks/selectors/assignee-selector";
import { StatusSelector } from "@/components/tasks/selectors/status-selector";
import { PrioritySelector } from "@/components/tasks/selectors/priority-selector";

export default function TaskDetailPage({
  params,
}: {
  params: Promise<{ projectId: string; taskId: string }>;
}) {
  const resolvedParams = use(params);
  const projectId = parseInt(resolvedParams.projectId);
  const taskId = parseInt(resolvedParams.taskId);
  const router = useRouter();

  const { user } = useUser();
  const {
    task,
    comments,
    isLoading,
    mutateTask,
    createComment,
    deleteComment,
  } = useTask(taskId);
  const { updateTask, deleteTask } = useTasks(projectId);
  const { members } = useProjectMembers(projectId);

  const [commentContent, setCommentContent] = useState("");
  const [isSubmittingComment, setIsSubmittingComment] = useState(false);
  const [isEditingTitle, setIsEditTitle] = useState(false);
  const [titleValue, setTitleValue] = useState("");

  useEffect(() => {
    if (task) setTitleValue(task.title);
  }, [task]);

  const handleUpdateStatus = async (status: TaskStatus) => {
    await updateTask(taskId, { status });
    mutateTask();
  };

  const handleUpdatePriority = async (priority: TaskPriority) => {
    await updateTask(taskId, { priority });
    mutateTask();
  };

  const handleUpdateAssignee = async (assigneeId: number | undefined) => {
    await updateTask(taskId, { assignee_id: assigneeId });
    mutateTask();
  };

  const handleUpdateDueDate = async (dateStr: string) => {
    await updateTask(taskId, { due_date: (dateStr || null) as string });
    mutateTask();
  };

  const handleUpdateTitle = async () => {
    if (!titleValue.trim() || titleValue === task?.title) {
      setIsEditTitle(false);
      return;
    }
    await updateTask(taskId, { title: titleValue });
    setIsEditTitle(false);
    mutateTask();
  };

  const handleDelete = async () => {
    if (!confirm("Are you sure you want to delete this task?")) return;
    await deleteTask(taskId);
    router.push(`/projects/${projectId}`);
  };

  const handleSubmitComment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!commentContent.trim()) return;

    setIsSubmittingComment(true);
    try {
      await createComment(commentContent);
      setCommentContent("");
    } catch {
      alert("Failed to add comment");
    } finally {
      setIsSubmittingComment(false);
    }
  };

  const handleDeleteComment = async (commentId: number) => {
    if (!confirm("Delete this comment?")) return;
    try {
      await deleteComment(commentId);
    } catch (err) {
      alert(
        (err as ApiError).response?.data?.detail || "Failed to delete comment"
      );
    }
  };

  if (isLoading)
    return (
      <div className="flex h-full items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-blue-500" />
      </div>
    );
  if (!task)
    return <div className="p-8 text-center text-gray-500">Task not found</div>;

  const statusIcons: Record<TaskStatus, React.ElementType> = {
    [TaskStatus.BACKLOG]: Clock,
    [TaskStatus.TODO]: Circle,
    [TaskStatus.IN_PROGRESS]: AlertCircle,
    [TaskStatus.DONE]: CheckCircle2,
  };
  const StatusIcon = statusIcons[task.status] || Circle;

  return (
    <div className="mx-auto max-w-6xl space-y-6 pb-20">
      {/* Header */}
      <div className="space-y-4">
        <Link
          href={`/projects/${projectId}`}
          className="mb-2 inline-flex items-center gap-1 text-sm text-gray-500 hover:text-gray-900"
        >
          <ArrowLeft className="h-4 w-4" /> Back to project
        </Link>

        <div className="flex flex-col gap-2">
          <div className="flex items-start justify-between gap-4">
            {isEditingTitle ? (
              <div className="flex flex-1 gap-2">
                <Input
                  value={titleValue}
                  onChange={(e) => setTitleValue(e.target.value)}
                  className="h-12 text-2xl font-bold"
                  autoFocus
                />
                <Button onClick={handleUpdateTitle}>Save</Button>
                <Button variant="ghost" onClick={() => setIsEditTitle(false)}>
                  Cancel
                </Button>
              </div>
            ) : (
              <h1 className="group flex items-center gap-2 text-3xl font-bold text-gray-900">
                {task.title}
                <span className="font-normal text-gray-400">#{task.id}</span>
                <button
                  onClick={() => setIsEditTitle(true)}
                  className="rounded p-1 opacity-0 transition-all group-hover:opacity-100 hover:bg-gray-100"
                >
                  <Settings2 className="h-4 w-4 text-gray-500" />
                </button>
              </h1>
            )}
            <div className="flex gap-2">
              {task.status === TaskStatus.DONE ? (
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handleUpdateStatus(TaskStatus.TODO)}
                  className="border-gray-200 text-gray-600 hover:text-gray-900"
                >
                  <ArrowLeft className="mr-2 h-4 w-4" />
                  Reopen Task
                </Button>
              ) : (
                <Button
                  variant="default"
                  size="sm"
                  onClick={() => handleUpdateStatus(TaskStatus.DONE)}
                  className="border-transparent bg-green-600 text-white hover:bg-green-700"
                >
                  <CheckCircle2 className="mr-2 h-4 w-4" />
                  Complete Task
                </Button>
              )}
              <Button
                variant="destructive"
                className="w-small"
                onClick={handleDelete}
              >
                <Trash2 className="mr-2 h-4 w-4" />
                Delete Task
              </Button>
            </div>
          </div>

          <div className="flex flex-wrap items-center gap-3 border-b border-gray-200 pb-6">
            <Badge
              className={cn(
                "flex items-center gap-1.5 rounded-full px-3 py-1 text-sm font-medium transition-colors",
                task.status === TaskStatus.DONE
                  ? "bg-green-100 text-green-700 hover:bg-green-200"
                  : "bg-blue-100 text-blue-700 hover:bg-blue-200"
              )}
            >
              <StatusIcon className="h-4 w-4" />
              {task.status.replace("_", " ")}
            </Badge>
            <span className="flex items-center gap-1.5 text-sm font-medium text-gray-500">
              <UserIcon className="h-4 w-4" />
              <span className="font-semibold text-gray-900">
                {task.author?.full_name || task.author?.email || "Unknown"}
              </span>{" "}
              opened this task {formatDistanceToNow(new Date(task.created_at))}{" "}
              ago • {comments.length} comments
            </span>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-8 lg:grid-cols-4">
        {/* Main Stream */}
        <div className="space-y-8 lg:col-span-3">
          {/* Description */}
          <div className="flex gap-4">
            <Avatar
              fallback={
                task.author?.full_name?.[0] || task.author?.email?.[0] || "U"
              }
              className="mt-1 h-10 w-10 border border-gray-100 shadow-sm"
            />
            <div className="flex-1">
              <div className="overflow-hidden rounded-lg border border-gray-200 bg-white shadow-sm">
                <div className="flex items-center justify-between border-b border-gray-200 bg-gray-50/50 px-4 py-2">
                  <span className="text-sm font-semibold text-gray-700">
                    Description
                  </span>
                </div>
                <div className="prose prose-sm prose-pre:bg-gray-50 prose-pre:border prose-pre:border-gray-100 max-w-none p-4">
                  <ReactMarkdown remarkPlugins={[remarkGfm]}>
                    {task.description || "_No description provided._"}
                  </ReactMarkdown>
                </div>
              </div>
            </div>
          </div>

          {/* Comments List */}
          <div className="relative space-y-8 before:absolute before:top-0 before:bottom-0 before:left-[1.25rem] before:w-0.5 before:bg-gray-100">
            {comments.map((comment) => (
              <CommentItem
                key={comment.id}
                comment={comment}
                currentUserId={user?.id}
                onDelete={handleDeleteComment}
              />
            ))}
          </div>

          {/* New Comment Box */}
          <div className="flex gap-4 border-t border-gray-200 pt-8">
            <Avatar
              fallback={user?.full_name?.[0] || user?.email?.[0] || "U"}
              className="mt-1 h-10 w-10 border border-gray-100 shadow-sm"
            />
            <div className="flex-1">
              <form onSubmit={handleSubmitComment} className="space-y-4">
                <MarkdownEditor
                  value={commentContent}
                  onChange={setCommentContent}
                  placeholder="Add a comment..."
                  className="border-gray-200 shadow-sm"
                />
                <div className="flex justify-end">
                  <Button
                    type="submit"
                    disabled={isSubmittingComment || !commentContent.trim()}
                    className="bg-blue-600 px-6 text-white hover:bg-blue-700"
                  >
                    {isSubmittingComment && (
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    )}
                    Comment
                  </Button>
                </div>
              </form>
            </div>
          </div>
        </div>

        {/* Sidebar Controls */}
        <div className="space-y-8 lg:col-span-1">
          <div className="space-y-6">
            <div className="space-y-2 border-b border-gray-100 pb-4">
              <Label className="text-xs font-bold tracking-wider text-gray-500 uppercase">
                Due Date
              </Label>
              <Input
                type="date"
                value={task.due_date ? task.due_date.split("T")[0] : ""}
                onChange={(e) => handleUpdateDueDate(e.target.value)}
              />
            </div>

            <div className="space-y-2 border-b border-gray-100 pb-4">
              <Label className="text-xs font-bold tracking-wider text-gray-500 uppercase">
                Assignees
              </Label>
              <AssigneeSelector
                members={members}
                value={task.assignee_id || undefined}
                onChange={handleUpdateAssignee}
              />
            </div>

            <div className="space-y-2 border-b border-gray-100 pb-4">
              <Label className="text-xs font-bold tracking-wider text-gray-500 uppercase">
                Status
              </Label>
              <StatusSelector
                value={task.status}
                onChange={handleUpdateStatus}
              />
            </div>

            {task.completed_at && (
              <div className="space-y-2 border-b border-gray-100 pb-4">
                <Label className="text-xs font-bold tracking-wider text-gray-500 uppercase">
                  Completed On
                </Label>
                <div className="flex items-center gap-2 text-sm text-gray-700">
                  <CheckCircle2 className="h-4 w-4 text-green-600" />
                  {new Date(task.completed_at).toLocaleDateString()}
                </div>
              </div>
            )}

            <div className="space-y-2 border-b border-gray-100 pb-4">
              <Label className="text-xs font-bold tracking-wider text-gray-500 uppercase">
                Priority
              </Label>
              <PrioritySelector
                value={task.priority}
                onChange={handleUpdatePriority}
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
