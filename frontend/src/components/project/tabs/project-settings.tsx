"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useProjects } from "@/hooks/use-projects";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Loader2, Check, Archive, Trash2 } from "lucide-react";
import { Project } from "@/types";

interface ProjectSettingsProps {
  project: Project;
}

export function ProjectSettings({ project }: ProjectSettingsProps) {
  const router = useRouter();
  const { updateProject } = useProjects();

  const [name, setName] = useState(project.name);
  const [description, setDescription] = useState(project.description || "");
  const [isSaving, setIsSaving] = useState(false);
  const [isSuccess, setIsSuccess] = useState(false);

  useEffect(() => {
    setName(project.name);
    setDescription(project.description || "");
  }, [project]);

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSaving(true);
    setIsSuccess(false);

    try {
      await updateProject(project.id, { name, description });
      setIsSuccess(true);
      setTimeout(() => setIsSuccess(false), 3000);
    } catch (err) {
      alert("Failed to update project");
    } finally {
      setIsSaving(false);
    }
  };

  const handleArchive = async () => {
    if (
      !confirm(
        "Are you sure you want to archive this project? It will be hidden from the active list."
      )
    )
      return;
    try {
      await updateProject(project.id, { is_archived: true });
      router.push("/dashboard");
    } catch (err) {
      alert("Failed to archive project");
    }
  };

  return (
    <div className="max-w-3xl space-y-8">
      <Card>
        <CardHeader>
          <CardTitle>General Settings</CardTitle>
          <CardDescription>Update your project details.</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSave} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">Project Name</Label>
              <Input
                id="name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="description">Description</Label>
              <Textarea
                id="description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="Describe the project goals..."
                rows={4}
              />
            </div>

            <div className="pt-2">
              <Button type="submit" disabled={isSaving}>
                {isSaving ? (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                ) : isSuccess ? (
                  <Check className="mr-2 h-4 w-4" />
                ) : null}
                {isSaving ? "Saving..." : isSuccess ? "Saved!" : "Save Changes"}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>

      <Card className="border-red-100 bg-red-50/30">
        <CardHeader>
          <CardTitle className="text-red-900">Danger Zone</CardTitle>
          <CardDescription className="text-red-700">
            Irreversible actions for this project.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center justify-between">
            <div>
              <h4 className="font-medium text-gray-900">Archive Project</h4>
              <p className="text-sm text-gray-500">
                Hide this project from your workspace.
              </p>
            </div>
            <Button
              variant="outline"
              onClick={handleArchive}
              className="border-red-200 bg-white text-red-600 hover:bg-red-50 hover:text-red-700"
            >
              <Archive className="mr-2 h-4 w-4" /> Archive
            </Button>
          </div>

          {/* Delete functionality not yet implemented in backend properly (hard delete), so we keep it simple or minimal */}
        </CardContent>
      </Card>
    </div>
  );
}
