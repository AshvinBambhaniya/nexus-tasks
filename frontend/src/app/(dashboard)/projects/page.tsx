"use client";

import { useState } from "react";
import Link from "next/link";
import { format } from "date-fns";
import { 
  Folder, 
  Search, 
  Plus, 
  ArrowRight,
  Filter,
  Inbox
} from "lucide-react";

import { useProjects } from "@/hooks/use-projects";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { ProjectModal } from "@/components/project/project-dialog";

type FilterStatus = "active" | "archived" | "all";

export default function ProjectsPage() {
  const { projects, isLoading } = useProjects();
  const [searchQuery, setSearchQuery] = useState("");
  const [filterStatus, setFilterStatus] = useState<FilterStatus>("active");
  const [isDialogOpen, setIsDialogOpen] = useState(false);

  // Filter Logic
  const filteredProjects = projects.filter((project) => {
    // 1. Search Filter
    const matchesSearch = project.name.toLowerCase().includes(searchQuery.toLowerCase());
    
    // 2. Status Filter
    let matchesStatus = true;
    if (filterStatus === "active") {
        matchesStatus = !project.is_archived;
    } else if (filterStatus === "archived") {
        matchesStatus = project.is_archived;
    }
    
    return matchesSearch && matchesStatus;
  });

  return (
    <div className="space-y-8 h-full flex flex-col">
      {/* Header & Controls */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Projects</h1>
          <p className="text-gray-500 mt-1">Manage and organize your team&apos;s work.</p>
        </div>
        <Button onClick={() => setIsDialogOpen(true)}>
          <Plus className="mr-2 h-4 w-4" /> New Project
        </Button>
      </div>

      <div className="flex flex-col sm:flex-row gap-4 items-center justify-between border-b border-gray-200 pb-4">
          {/* Search */}
          <div className="relative w-full sm:w-72">
            <Search className="absolute left-3 top-2.5 h-4 w-4 text-gray-400" />
            <Input 
                placeholder="Search projects..." 
                className="pl-9"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
            />
          </div>

          {/* Filter Tabs */}
          <div className="flex p-1 bg-gray-100 rounded-lg self-start sm:self-auto">
             {(["active", "archived", "all"] as const).map((status) => (
                 <button
                    key={status}
                    onClick={() => setFilterStatus(status)}
                    className={cn(
                        "px-4 py-1.5 text-sm font-medium rounded-md capitalize transition-all",
                        filterStatus === status 
                            ? "bg-white text-gray-900 shadow-sm" 
                            : "text-gray-500 hover:text-gray-700"
                    )}
                 >
                    {status}
                 </button>
             ))}
          </div>
      </div>

      {/* Grid Content */}
      <div className="flex-1 overflow-y-auto">
         {isLoading ? (
             <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {[1,2,3,4,5,6].map(i => <SkeletonCard key={i} />)}
             </div>
         ) : filteredProjects.length === 0 ? (
             <div className="flex flex-col items-center justify-center h-64 text-center rounded-lg border-2 border-dashed border-gray-200 bg-gray-50">
                <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-gray-100 mb-4">
                    {searchQuery ? <Search className="h-6 w-6 text-gray-400" /> : <Folder className="h-6 w-6 text-gray-400" />}
                </div>
                <h3 className="text-lg font-medium text-gray-900">No projects found</h3>
                <p className="text-sm text-gray-500 mt-1 max-w-sm">
                    {searchQuery 
                        ? `No projects matching "${searchQuery}"` 
                        : filterStatus === "archived" 
                            ? "No archived projects." 
                            : "Get started by creating a new project."}
                </p>
                {!searchQuery && filterStatus !== "archived" && (
                     <Button variant="outline" className="mt-4" onClick={() => setIsDialogOpen(true)}>
                        Create Project
                     </Button>
                )}
             </div>
         ) : (
             <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {filteredProjects.map((project) => (
                    <Link key={project.id} href={`/projects/${project.id}`} className="block h-full">
                        <Card className="h-full p-6 hover:shadow-md transition-all cursor-pointer border-gray-200 hover:border-blue-200 group flex flex-col">
                            <div className="flex justify-between items-start mb-4">
                                <div className="p-2.5 bg-blue-50 rounded-xl group-hover:bg-blue-100 transition-colors">
                                    <Folder className="h-5 w-5 text-blue-600" />
                                </div>
                                {project.is_archived && <Badge variant="secondary">Archived</Badge>}
                            </div>
                            
                            <div className="flex-1">
                                <h3 className="font-semibold text-gray-900 mb-2 group-hover:text-blue-700 transition-colors text-lg">
                                    {project.name}
                                </h3>
                                <p className="text-sm text-gray-500 line-clamp-2 leading-relaxed">
                                    {project.description || "No description provided."}
                                </p>
                            </div>

                            <div className="mt-6 pt-4 border-t border-gray-100 flex items-center justify-between text-xs text-gray-500">
                                <div className="flex items-center gap-1">
                                    <span>Created {format(new Date(project.created_at), "MMM d, yyyy")}</span>
                                </div>
                                <div className="flex items-center gap-1 text-blue-600 font-medium opacity-0 group-hover:opacity-100 transition-opacity">
                                    Open <ArrowRight className="h-3 w-3" />
                                </div>
                            </div>
                        </Card>
                    </Link>
                ))}
             </div>
         )}
      </div>

      <ProjectModal 
        isOpen={isDialogOpen} 
        onClose={() => setIsDialogOpen(false)} 
      />
    </div>
  );
}

function SkeletonCard() {
    return (
        <div className="h-48 bg-gray-100 rounded-xl animate-pulse border border-gray-200" />
    )
}
