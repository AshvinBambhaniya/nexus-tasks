"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useTeams } from "@/hooks/use-teams";
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
import { Loader2, Check, Trash2 } from "lucide-react";
import { Team } from "@/types";

interface TeamSettingsProps {
  team: Team;
}

export function TeamSettings({ team }: TeamSettingsProps) {
  const router = useRouter();
  const { updateTeam, deleteTeam } = useTeams();

  const [name, setName] = useState(team.name);
  const [description, setDescription] = useState(team.description || "");
  const [isSaving, setIsSaving] = useState(false);
  const [isSuccess, setIsSuccess] = useState(false);

  useEffect(() => {
    setName(team.name);
    setDescription(team.description || "");
  }, [team]);

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSaving(true);
    setIsSuccess(false);

    try {
      await updateTeam(team.id, { name, description });
      setIsSuccess(true);
      setTimeout(() => setIsSuccess(false), 3000);
    } catch {
      alert("Failed to update team");
    } finally {
      setIsSaving(false);
    }
  };

  const handleDelete = async () => {
    if (
      !confirm(
        "Are you sure you want to delete this team? This action cannot be undone."
      )
    )
      return;
    try {
      await deleteTeam(team.id);
      router.push("/teams");
    } catch {
      alert("Failed to delete team");
    }
  };

  return (
    <div className="max-w-3xl space-y-8">
      <Card>
        <CardHeader>
          <CardTitle>General Settings</CardTitle>
          <CardDescription>Update your team details.</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSave} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">Team Name</Label>
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
                placeholder="What is this team responsible for?"
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
            Irreversible actions for this team.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center justify-between">
            <div>
              <h4 className="font-medium text-gray-900">Delete Team</h4>
              <p className="text-sm text-gray-500">
                Permanently remove this team and its memberships.
              </p>
            </div>
            <Button
              variant="outline"
              onClick={handleDelete}
              className="border-red-200 bg-white text-red-600 hover:bg-red-50 hover:text-red-700"
            >
              <Trash2 className="mr-2 h-4 w-4" /> Delete Team
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
