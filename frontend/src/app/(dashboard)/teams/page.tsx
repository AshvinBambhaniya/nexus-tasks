"use client";

import { useState } from "react";
import Link from "next/link";
import { format } from "date-fns";
import { 
  Users, 
  Search, 
  Plus, 
  ArrowRight,
} from "lucide-react";

import { useTeams } from "@/hooks/use-teams";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card } from "@/components/ui/card";
import { TeamModal } from "@/components/team/team-dialog";

export default function TeamsPage() {
  const { teams, isLoading } = useTeams();
  const [searchQuery, setSearchQuery] = useState("");
  const [isDialogOpen, setIsDialogOpen] = useState(false);

  // Filter Logic
  const filteredTeams = teams.filter((team) => {
    return team.name.toLowerCase().includes(searchQuery.toLowerCase());
  });

  return (
    <div className="space-y-8 h-full flex flex-col">
      {/* Header & Controls */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Teams</h1>
          <p className="text-gray-500 mt-1">Organize people into groups for easier management.</p>
        </div>
        <Button onClick={() => setIsDialogOpen(true)}>
          <Plus className="mr-2 h-4 w-4" /> New Team
        </Button>
      </div>

      <div className="flex flex-col sm:flex-row gap-4 items-center justify-between border-b border-gray-200 pb-4">
          {/* Search */}
          <div className="relative w-full sm:w-72">
            <Search className="absolute left-3 top-2.5 h-4 w-4 text-gray-400" />
            <Input 
                placeholder="Search teams..." 
                className="pl-9"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
            />
          </div>
      </div>

      {/* Grid Content */}
      <div className="flex-1 overflow-y-auto">
         {isLoading ? (
             <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {[1,2,3,4].map(i => <div key={i} className="h-40 bg-gray-100 rounded-xl animate-pulse border border-gray-200" />)}
             </div>
         ) : filteredTeams.length === 0 ? (
             <div className="flex flex-col items-center justify-center h-64 text-center rounded-lg border-2 border-dashed border-gray-200 bg-gray-50">
                <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-gray-100 mb-4">
                    <Users className="h-6 w-6 text-gray-400" />
                </div>
                <h3 className="text-lg font-medium text-gray-900">No teams found</h3>
                <p className="text-sm text-gray-500 mt-1 max-w-sm">
                    {searchQuery 
                        ? `No teams matching "${searchQuery}"` 
                        : "Create a team to group users together."}
                </p>
                {!searchQuery && (
                     <Button variant="outline" className="mt-4" onClick={() => setIsDialogOpen(true)}>
                        Create Team
                     </Button>
                )}
             </div>
         ) : (
             <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {filteredTeams.map((team) => (
                    <Link key={team.id} href={`/teams/${team.id}`} className="block h-full">
                        <Card className="h-full p-6 hover:shadow-md transition-all cursor-pointer border-gray-200 hover:border-blue-200 group flex flex-col">
                            <div className="flex justify-between items-start mb-4">
                                <div className="p-2.5 bg-indigo-50 rounded-xl group-hover:bg-indigo-100 transition-colors">
                                    <Users className="h-5 w-5 text-indigo-600" />
                                </div>
                            </div>
                            
                            <div className="flex-1">
                                <h3 className="font-semibold text-gray-900 mb-2 group-hover:text-indigo-700 transition-colors text-lg">
                                    {team.name}
                                </h3>
                                <p className="text-sm text-gray-500 line-clamp-2 leading-relaxed">
                                    {team.description || "No description provided."}
                                </p>
                            </div>

                            <div className="mt-6 pt-4 border-t border-gray-100 flex items-center justify-between text-xs text-gray-500">
                                <div className="flex items-center gap-1">
                                    <span>Created {format(new Date(team.created_at), "MMM d, yyyy")}</span>
                                </div>
                                <div className="flex items-center gap-1 text-indigo-600 font-medium opacity-0 group-hover:opacity-100 transition-opacity">
                                    Manage <ArrowRight className="h-3 w-3" />
                                </div>
                            </div>
                        </Card>
                    </Link>
                ))}
             </div>
         )}
      </div>

      <TeamModal 
        isOpen={isDialogOpen}
        onClose={() => setIsDialogOpen(false)}
      />
    </div>
  );
}
