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
      <div className="flex flex-col items-center justify-center p-12 text-center border-2 border-dashed border-gray-200 rounded-lg bg-gray-50 text-gray-500">
        <Folder className="h-10 w-10 text-gray-300 mb-3" />
        <p>No projects assigned to this team yet.</p>
      </div>
    );
  }

  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      {projects.map((project) => (
        <Link key={project.id} href={`/projects/${project.id}`} className="group">
          <Card className="hover:shadow-md transition-shadow cursor-pointer border-gray-200 group-hover:border-indigo-200">
            <CardContent className="p-5">
              <div className="flex items-center justify-between mb-3">
                <div className="p-2 bg-indigo-50 rounded-lg group-hover:bg-indigo-100 transition-colors text-indigo-600">
                  <Folder className="h-5 w-5" />
                </div>
                <ArrowRight className="h-4 w-4 text-gray-300 group-hover:text-indigo-500 transition-colors transform group-hover:translate-x-1" />
              </div>
              <h3 className="font-semibold text-gray-900 group-hover:text-indigo-700 transition-colors truncate">
                {project.name}
              </h3>
              {project.description && (
                <p className="text-sm text-gray-500 line-clamp-2 mt-1 leading-relaxed">
                  {project.description}
                </p>
              )}
              <div className="mt-4 pt-4 border-t border-gray-50 text-[10px] text-gray-400 uppercase tracking-wider">
                Created {format(new Date(project.created_at), "MMM d, yyyy")}
              </div>
            </CardContent>
          </Card>
        </Link>
      ))}
    </div>
  );
}