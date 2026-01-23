"use client";

import { useEffect, useState } from "react";
import { Modal } from "@/components/ui/modal";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Task, TaskPriority, TaskStatus } from "@/types";
import { useTasks } from "@/hooks/use-tasks";

// Assuming react-hook-form is NOT available based on package.json, 
// I will use simple state.

interface TaskDialogProps {
  isOpen: boolean;
  onClose: () => void;
  task?: Task; // If provided, we are editing
}

export function TaskDialog({ isOpen, onClose, task }: TaskDialogProps) {
  const { createTask, updateTask, deleteTask } = useTasks();
  const [isLoading, setIsLoading] = useState(false);
  
  const [formData, setFormData] = useState({
    title: "",
    description: "",
    status: TaskStatus.TODO,
    priority: TaskPriority.P2,
  });

  useEffect(() => {
    if (task) {
      setFormData({
        title: task.title,
        description: task.description || "",
        status: task.status,
        priority: task.priority,
      });
    } else {
      setFormData({
        title: "",
        description: "",
        status: TaskStatus.TODO,
        priority: TaskPriority.P2, // Default Medium
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
      description={task ? "Update task details." : "Add a new task to your workspace."}
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="title">Title</Label>
          <Input
            id="title"
            required
            placeholder="e.g. Implement authentication"
            value={formData.title}
            onChange={(e) => setFormData({ ...formData, title: e.target.value })}
          />
        </div>

        <div className="space-y-2">
          <Label htmlFor="description">Description</Label>
          <Textarea
            id="description"
            placeholder="Add more details..."
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label htmlFor="status">Status</Label>
            <select
              id="status"
              className="flex h-10 w-full rounded-md border border-gray-200 bg-white px-3 py-2 text-sm text-gray-900 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
              value={formData.status}
              onChange={(e) => setFormData({ ...formData, status: e.target.value as TaskStatus })}
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
              className="flex h-10 w-full rounded-md border border-gray-200 bg-white px-3 py-2 text-sm text-gray-900 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
              value={formData.priority}
              onChange={(e) => setFormData({ ...formData, priority: e.target.value as TaskPriority })}
            >
              <option value={TaskPriority.P0}>P0 - Critical</option>
              <option value={TaskPriority.P1}>P1 - High</option>
              <option value={TaskPriority.P2}>P2 - Medium</option>
              <option value={TaskPriority.P3}>P3 - Low</option>
            </select>
          </div>
        </div>

        <div className="flex justify-between pt-4">
          {task ? (
             <Button type="button" variant="destructive" onClick={handleDelete} disabled={isLoading}>
                Delete
             </Button>
          ) : <div></div>}
          
          <div className="flex gap-2">
             <Button type="button" variant="ghost" onClick={onClose} disabled={isLoading}>
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
