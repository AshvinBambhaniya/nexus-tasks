"use client";

import { useState } from "react";
import { Modal } from "@/components/ui/modal";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useTeams } from "@/hooks/use-teams";

interface TeamModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export function TeamModal({ isOpen, onClose }: TeamModalProps) {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const { createTeam } = useTeams();
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;

    setIsLoading(true);
    try {
      await createTeam(name, description);
      onClose();
      setName("");
      setDescription("");
    } catch (error) {
      console.error(error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Modal title="Create Team" isOpen={isOpen} onClose={onClose}>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <Label htmlFor="team-name">Team Name</Label>
          <Input
            id="team-name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="e.g. Engineering"
            required
          />
        </div>
        <div>
          <Label htmlFor="team-desc">Description</Label>
          <Textarea
            id="team-desc"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Brief description of the team..."
          />
        </div>
        <div className="flex justify-end gap-2">
          <Button
            type="button"
            variant="ghost"
            onClick={onClose}
            disabled={isLoading}
          >
            Cancel
          </Button>
          <Button type="submit" disabled={isLoading}>
            {isLoading ? "Creating..." : "Create Team"}
          </Button>
        </div>
      </form>
    </Modal>
  );
}

export function TeamDialog() {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <>
      <Button
        variant="ghost"
        size="sm"
        className="w-full justify-start px-2 text-xs text-gray-500 hover:text-gray-900"
        onClick={() => setIsOpen(true)}
      >
        + Add Team
      </Button>
      <TeamModal isOpen={isOpen} onClose={() => setIsOpen(false)} />
    </>
  );
}
