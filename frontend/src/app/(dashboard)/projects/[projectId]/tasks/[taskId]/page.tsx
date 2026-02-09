"use client";

import { use, useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useTask, useTasks } from "@/hooks/use-tasks";
import { useProjectMembers } from "@/hooks/use-projects";
import { TaskStatus, TaskPriority } from "@/types";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { MarkdownEditor } from "@/components/ui/markdown-editor";
import { Card, CardContent } from "@/components/ui/card";
import { Avatar } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { CommentItem } from "@/components/tasks/comment-item";
import { 
  ArrowLeft, 
  Loader2, 
  Trash2, 
  Settings2, 
  MessageSquare,
  Clock,
  User as UserIcon,
  CheckCircle2,
  Circle,
  AlertCircle
} from "lucide-react";
import Link from "next/link";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { cn } from "@/lib/utils";
import { formatDistanceToNow } from "date-fns";

import { useUser } from "@/hooks/use-user";

export default function TaskDetailPage({ params }: { params: Promise<{ projectId: string, taskId: string }> }) {
  const resolvedParams = use(params);
  const projectId = parseInt(resolvedParams.projectId);
  const taskId = parseInt(resolvedParams.taskId);
  const router = useRouter();
  
  const { user } = useUser();
  const { task, comments, isLoading, mutateTask, createComment, deleteComment } = useTask(taskId);
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

  const handleUpdateAssignee = async (assigneeId: number | null) => {
    await updateTask(taskId, { assignee_id: assigneeId as any });
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
    } catch (error) {
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
          alert("Failed to delete comment");
      }
  }

  if (isLoading) return <div className="flex h-full items-center justify-center"><Loader2 className="h-8 w-8 animate-spin text-blue-500" /></div>;
  if (!task) return <div className="p-8 text-center text-gray-500">Task not found</div>;

  const statusIcons: Record<TaskStatus, any> = {
    [TaskStatus.BACKLOG]: Clock,
    [TaskStatus.TODO]: Circle,
    [TaskStatus.IN_PROGRESS]: AlertCircle,
    [TaskStatus.DONE]: CheckCircle2,
  };
  const StatusIcon = statusIcons[task.status] || Circle;

  return (
    <div className="max-w-6xl mx-auto space-y-6 pb-20">
      {/* Header */}
      <div className="space-y-4">
        <Link href={`/projects/${projectId}`} className="inline-flex items-center text-sm text-gray-500 hover:text-gray-900 gap-1 mb-2">
          <ArrowLeft className="h-4 w-4" /> Back to project
        </Link>

        <div className="flex flex-col gap-2">
          <div className="flex items-start justify-between gap-4">
            {isEditingTitle ? (
              <div className="flex-1 flex gap-2">
                <Input 
                  value={titleValue} 
                  onChange={(e) => setTitleValue(e.target.value)} 
                  className="text-2xl font-bold h-12"
                  autoFocus
                />
                <Button onClick={handleUpdateTitle}>Save</Button>
                <Button variant="ghost" onClick={() => setIsEditTitle(false)}>Cancel</Button>
              </div>
            ) : (
              <h1 className="text-3xl font-bold text-gray-900 group flex items-center gap-2">
                {task.title}
                <span className="text-gray-400 font-normal">#{task.id}</span>
                <button 
                  onClick={() => setIsEditTitle(true)}
                  className="opacity-0 group-hover:opacity-100 p-1 hover:bg-gray-100 rounded transition-all"
                >
                  <Settings2 className="h-4 w-4 text-gray-500" />
                </button>
              </h1>
            )}
            <div className="flex gap-2">
               <Button variant="outline" size="sm" onClick={() => handleUpdateStatus(task.status === TaskStatus.DONE ? TaskStatus.TODO : TaskStatus.DONE)}>
                  {task.status === TaskStatus.DONE ? "Reopen" : "Complete"}
               </Button>
            </div>
          </div>

          <div className="flex items-center gap-3 flex-wrap border-b border-gray-200 pb-6">
            <Badge className={cn("flex items-center gap-1.5 px-3 py-1 text-sm rounded-full", 
                task.status === TaskStatus.DONE ? "bg-green-600 hover:bg-green-700" : "bg-blue-600 hover:bg-blue-700"
            )}>
              <StatusIcon className="h-4 w-4" />
              {task.status.replace("_", " ")}
            </Badge>
            <span className="text-gray-500 text-sm flex items-center gap-1.5">
               <UserIcon className="h-4 w-4" />
               <span className="font-semibold text-gray-700">{task.author?.full_name || task.author?.email || "Unknown"}</span> opened this task {formatDistanceToNow(new Date(task.created_at))} ago • {comments.length} comments
            </span>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
        {/* Main Stream */}
        <div className="lg:col-span-3 space-y-8">
          {/* Description */}
          <div className="flex gap-4">
             <Avatar fallback="U" className="h-10 w-10 mt-1 border border-gray-100" />
             <div className="flex-1">
                <div className="bg-white border border-gray-200 rounded-lg shadow-sm">
                   <div className="bg-gray-50 px-4 py-2 border-b border-gray-200 flex justify-between items-center rounded-t-lg">
                      <span className="text-sm font-medium text-gray-700">Description</span>
                   </div>
                   <div className="p-4 prose prose-sm max-w-none prose-pre:bg-gray-50 prose-pre:border prose-pre:border-gray-100">
                      <ReactMarkdown remarkPlugins={[remarkGfm]}>{task.description || "_No description provided._"}</ReactMarkdown>
                   </div>
                </div>
             </div>
          </div>

          {/* Comments List */}
          <div className="relative space-y-8 before:absolute before:left-[1.25rem] before:top-0 before:bottom-0 before:w-0.5 before:bg-gray-100">
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
             <Avatar fallback="U" className="h-10 w-10 mt-1 border border-gray-100 shadow-sm" />
             <div className="flex-1">
                <form onSubmit={handleSubmitComment} className="space-y-4">
                   <MarkdownEditor 
                      value={commentContent} 
                      onChange={setCommentContent}
                      placeholder="Add a comment..."
                      className="shadow-sm"
                   />
                   <div className="flex justify-end">
                      <Button type="submit" disabled={isSubmittingComment || !commentContent.trim()}>
                         {isSubmittingComment && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                         Comment
                      </Button>
                   </div>
                </form>
             </div>
          </div>
        </div>

        {/* Sidebar Controls */}
        <div className="lg:col-span-1 space-y-8">
           <div className="space-y-6">
              <div className="space-y-2 pb-4 border-b border-gray-100">
                 <Label className="text-xs font-bold text-gray-500 uppercase tracking-wider">Assignees</Label>
                 <select 
                    className="w-full text-sm border-gray-200 rounded-md focus:ring-blue-500 focus:border-blue-500"
                    value={task.assignee_id || ""}
                    onChange={(e) => handleUpdateAssignee(e.target.value ? parseInt(e.target.value) : null)}
                 >
                    <option value="">Unassigned</option>
                    {members.map(m => <option key={m.user_id} value={m.user_id}>{m.email}</option>)}
                 </select>
              </div>

              <div className="space-y-2 pb-4 border-b border-gray-100">
                 <Label className="text-xs font-bold text-gray-500 uppercase tracking-wider">Status</Label>
                 <select 
                    className="w-full text-sm border-gray-200 rounded-md focus:ring-blue-500 focus:border-blue-500"
                    value={task.status}
                    onChange={(e) => handleUpdateStatus(e.target.value as TaskStatus)}
                 >
                    {Object.values(TaskStatus).map(s => (
                        <option key={s} value={s}>{s.replace("_", " ")}</option>
                    ))}
                 </select>
              </div>

              <div className="space-y-2 pb-4 border-b border-gray-100">
                 <Label className="text-xs font-bold text-gray-500 uppercase tracking-wider">Priority</Label>
                 <select 
                    className="w-full text-sm border-gray-200 rounded-md focus:ring-blue-500 focus:border-blue-500"
                    value={task.priority}
                    onChange={(e) => handleUpdatePriority(e.target.value as TaskPriority)}
                 >
                    <option value={TaskPriority.P0}>Critical (P0)</option>
                    <option value={TaskPriority.P1}>High (P1)</option>
                    <option value={TaskPriority.P2}>Medium (P2)</option>
                    <option value={TaskPriority.P3}>Low (P3)</option>
                 </select>
              </div>

              <div className="pt-4">
                 <Button variant="outline" className="w-full text-red-600 hover:bg-red-50 hover:text-red-700 border-red-100" onClick={handleDelete}>
                    <Trash2 className="mr-2 h-4 w-4" />
                    Delete Task
                 </Button>
              </div>
           </div>
        </div>
      </div>
    </div>
  );
}
