"use client";

import { useTeamProjects } from "@/hooks/use-teams";
import { Card, CardContent } from "@/components/ui/card";
import { Folder, ArrowRight, Loader2 } from "lucide-react";
import Link from "next/link";
import { format } from "date-fns";

interface TeamProjectsProps {
  teamId: number;
}

export function TeamProjects({ teamId }: TeamProjectsProps) {
  const { projects, isLoading } = useTeamProjects(teamId);

  if (isLoading) {
    return (
      <div className="flex justify-center p-12">
        <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
      </div>
    );
  }

  if (projects.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center rounded-lg border-2 border-dashed border-gray-200 bg-gray-50 p-12 text-center text-gray-500">
        <Folder className="mb-3 h-10 w-10 text-gray-300" />
        <p>No projects assigned to this team yet.</p>
      </div>
    );
  }

  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      {projects.map((project) => (
        <Link
          key={project.id}
          href={`/projects/${project.id}`}
          className="group"
        >
          <Card className="cursor-pointer border-gray-200 transition-shadow group-hover:border-indigo-200 hover:shadow-md">
            <CardContent className="p-5">
              <div className="mb-3 flex items-center justify-between">
                <div className="rounded-lg bg-indigo-50 p-2 text-indigo-600 transition-colors group-hover:bg-indigo-100">
                  <Folder className="h-5 w-5" />
                </div>
                <ArrowRight className="h-4 w-4 transform text-gray-300 transition-colors group-hover:translate-x-1 group-hover:text-indigo-500" />
              </div>
              <h3 className="truncate font-semibold text-gray-900 transition-colors group-hover:text-indigo-700">
                {project.name}
              </h3>
              {project.description && (
                <p className="mt-1 line-clamp-2 text-sm leading-relaxed text-gray-500">
                  {project.description}
                </p>
              )}
              <div className="mt-4 border-t border-gray-50 pt-4 text-[10px] tracking-wider text-gray-400 uppercase">
                Created {format(new Date(project.created_at), "MMM d, yyyy")}
              </div>
            </CardContent>
          </Card>
        </Link>
      ))}
    </div>
  );
}
