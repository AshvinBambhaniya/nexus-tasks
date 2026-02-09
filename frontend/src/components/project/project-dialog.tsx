"use client";

import { useEffect, useState } from "react";
import { Modal } from "@/components/ui/modal";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useProjects } from "@/hooks/use-projects";
import { Project } from "@/types";

interface ProjectModalProps {
  isOpen: boolean;
  onClose: () => void;
  project?: Project; // If provided, edit mode
}

export function ProjectModal({ isOpen, onClose, project }: ProjectModalProps) {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const { createProject, updateProject } = useProjects();
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (project) {
      setName(project.name);
      setDescription(project.description || "");
    } else {
      setName("");
      setDescription("");
    }
  }, [project, isOpen]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;

    setIsLoading(true);
    try {
      if (project) {
        await updateProject(project.id, { name, description });
      } else {
        await createProject(name, description);
      }
      onClose();
      if (!project) {
        // Clear form only on create success
        setName("");
        setDescription("");
      }
    } catch (error) {
      console.error(error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Modal
      title={project ? "Edit Project" : "Create Project"}
      isOpen={isOpen}
      onClose={onClose}
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <Label htmlFor="project-name">Project Name</Label>
          <Input
            id="project-name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="e.g. Website Redesign"
            required
          />
        </div>
        <div>
          <Label htmlFor="project-desc">Description</Label>
          <Textarea
            id="project-desc"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Brief description of the project..."
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
            {isLoading
              ? "Saving..."
              : project
                ? "Save Changes"
                : "Create Project"}
          </Button>
        </div>
      </form>
    </Modal>
  );
}

export function ProjectDialog() {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <>
      <Button
        variant="ghost"
        size="sm"
        className="w-full justify-start px-2 text-xs text-gray-500 hover:text-gray-900"
        onClick={() => setIsOpen(true)}
      >
        + Add Project
      </Button>
      <ProjectModal isOpen={isOpen} onClose={() => setIsOpen(false)} />
    </>
  );
}
