"use client";

import { useState } from "react";
import { Modal } from "@/components/ui/modal";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useTeams } from "@/hooks/use-teams";

export function TeamDialog() {
  const [isOpen, setIsOpen] = useState(false);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const { createTeam } = useTeams();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;

    try {
      await createTeam(name, description);
      setIsOpen(false);
      setName("");
      setDescription("");
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <>
      <Button 
        variant="ghost" 
        size="sm" 
        className="w-full justify-start text-xs text-gray-500 hover:text-gray-900 px-2"
        onClick={() => setIsOpen(true)}
      >
        + Add Team
      </Button>
      <Modal title="Create Team" isOpen={isOpen} onClose={() => setIsOpen(false)}>
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
            <Button type="button" variant="ghost" onClick={() => setIsOpen(false)}>
              Cancel
            </Button>
            <Button type="submit">Create Team</Button>
          </div>
        </form>
      </Modal>
    </>
  );
}
