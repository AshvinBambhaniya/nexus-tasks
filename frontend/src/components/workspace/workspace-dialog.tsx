"use client";

import { useState } from "react";
import { Modal } from "@/components/ui/modal";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { useWorkspaces } from "@/hooks/use-workspaces";

interface WorkspaceDialogProps {
  isOpen: boolean;
  onClose: () => void;
}

export function WorkspaceDialog({ isOpen, onClose }: WorkspaceDialogProps) {
  const { createWorkspace } = useWorkspaces();
  const [name, setName] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;

    setIsLoading(true);
    try {
      await createWorkspace(name);
      setName("");
      onClose();
    } catch (error) {
      console.error("Failed to create workspace:", error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title="Create Workspace"
      description="Add a new team workspace to collaborate with others."
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="name" className="text-gray-700">Workspace Name</Label>
          <Input
            id="name"
            required
            placeholder="e.g. Marketing Team"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="text-gray-900"
          />
        </div>

        <div className="flex justify-end gap-2 pt-4 text-gray-700">
          <Button type="button" variant="ghost" onClick={onClose} disabled={isLoading}>
            Cancel
          </Button>
          <Button type="submit" disabled={isLoading || !name.trim()}>
            {isLoading ? "Creating..." : "Create Workspace"}
          </Button>
        </div>
      </form>
    </Modal>
  );
}
