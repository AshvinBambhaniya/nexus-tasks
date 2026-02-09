"use client";

import { useEffect, useState } from "react";
import { Modal } from "@/components/ui/modal";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Task, TaskPriority, TaskStatus, ApiError } from "@/types";
import { useTasks } from "@/hooks/use-tasks";
import { useProjectMembers } from "@/hooks/use-projects";

// Assuming react-hook-form is NOT available based on package.json,
// I will use simple state.

interface TaskDialogProps {
  isOpen: boolean;
  onClose: () => void;
  task?: Task; // If provided, we are editing
  projectId?: number; // Required for creation
}

export function TaskDialog({
  isOpen,
  onClose,
  task,
  projectId,
}: TaskDialogProps) {
  const activeProjectId = projectId || task?.project_id;

  // Pass projectId to hook. If editing (task exists), we might not need projectId strictly if updateTask uses task.id,
  // but createTask needs it.
  const { createTask, updateTask, deleteTask } = useTasks(activeProjectId);
  const { members } = useProjectMembers(activeProjectId);

  const [isLoading, setIsLoading] = useState(false);

  const [formData, setFormData] = useState({
    title: "",
    description: "",
    status: TaskStatus.TODO,
    priority: TaskPriority.P2,
    assignee_id: undefined as number | undefined,
  });

  useEffect(() => {
    if (task) {
      setFormData({
        title: task.title,
        description: task.description || "",
        status: task.status,
        priority: task.priority,
        assignee_id: task.assignee_id,
      });
    } else {
      setFormData({
        title: "",
        description: "",
        status: TaskStatus.TODO,
        priority: TaskPriority.P2, // Default Medium
        assignee_id: undefined,
      });
    }
  }, [task, isOpen]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      if (task) {
        await updateTask(task.id, formData);
      } else {
        await createTask(formData);
      }
      onClose();
    } catch (error) {
      console.error("Failed to save task:", error);
      alert(
        (error as ApiError).response?.data?.detail || "Failed to save task"
      );
    } finally {
      setIsLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!task) return;
    if (!confirm("Are you sure you want to delete this task?")) return;

    setIsLoading(true);
    try {
      await deleteTask(task.id);
      onClose();
    } catch (error) {
      console.error("Failed to delete task:", error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={task ? "Edit Task" : "Create Task"}
      description={
        task ? "Update task details." : "Add a new task to your workspace."
      }
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="title">Title</Label>
          <Input
            id="title"
            required
            placeholder="e.g. Implement authentication"
            value={formData.title}
            onChange={(e) =>
              setFormData({ ...formData, title: e.target.value })
            }
          />
        </div>

        <div className="space-y-2">
          <Label htmlFor="description">Description</Label>
          <Textarea
            id="description"
            placeholder="Add more details..."
            value={formData.description}
            onChange={(e) =>
              setFormData({ ...formData, description: e.target.value })
            }
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label htmlFor="status">Status</Label>
            <select
              id="status"
              className="flex h-10 w-full rounded-md border border-gray-200 bg-white px-3 py-2 text-sm text-gray-900 focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:outline-none"
              value={formData.status}
              onChange={(e) =>
                setFormData({
                  ...formData,
                  status: e.target.value as TaskStatus,
                })
              }
            >
              {Object.values(TaskStatus).map((status) => (
                <option key={status} value={status}>
                  {status.replace("_", " ")}
                </option>
              ))}
            </select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="priority">Priority</Label>
            <select
              id="priority"
              className="flex h-10 w-full rounded-md border border-gray-200 bg-white px-3 py-2 text-sm text-gray-900 focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:outline-none"
              value={formData.priority}
              onChange={(e) =>
                setFormData({
                  ...formData,
                  priority: e.target.value as TaskPriority,
                })
              }
            >
              <option value={TaskPriority.P0}>P0 - Critical</option>
              <option value={TaskPriority.P1}>P1 - High</option>
              <option value={TaskPriority.P2}>P2 - Medium</option>
              <option value={TaskPriority.P3}>P3 - Low</option>
            </select>
          </div>
        </div>

        <div className="space-y-2">
          <Label htmlFor="assignee">Assignee</Label>
          <select
            id="assignee"
            className="flex h-10 w-full rounded-md border border-gray-200 bg-white px-3 py-2 text-sm text-gray-900 focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:outline-none"
            value={formData.assignee_id || ""}
            onChange={(e) =>
              setFormData({
                ...formData,
                assignee_id: e.target.value
                  ? parseInt(e.target.value)
                  : undefined,
              })
            }
          >
            <option value="">Unassigned</option>
            {members.map((member) => (
              <option key={member.user_id} value={member.user_id}>
                {member.email}
              </option>
            ))}
          </select>
        </div>

        <div className="flex justify-between pt-4">
          {task ? (
            <Button
              type="button"
              variant="destructive"
              onClick={handleDelete}
              disabled={isLoading}
            >
              Delete
            </Button>
          ) : (
            <div></div>
          )}

          <div className="flex gap-2">
            <Button
              type="button"
              variant="ghost"
              onClick={onClose}
              disabled={isLoading}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={isLoading}>
              {isLoading ? "Saving..." : task ? "Save Changes" : "Create Task"}
            </Button>
          </div>
        </div>
      </form>
    </Modal>
  );
}
