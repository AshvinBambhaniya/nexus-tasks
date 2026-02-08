"use client";

import { WorkspaceSwitcher } from "@/components/workspace/workspace-switcher";
import { LayoutDashboard, CheckSquare, Inbox, Settings, LogOut, Folder, Users } from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";
import { useAuth } from "@/hooks/use-auth";
import { useUser } from "@/hooks/use-user";
import { useProjects } from "@/hooks/use-projects";
import { useTeams } from "@/hooks/use-teams";
import { ProjectDialog } from "@/components/project/project-dialog";
import { TeamDialog } from "@/components/team/team-dialog";

const navigation = [
  { name: "My Focus", href: "/inbox", icon: Inbox },
  { name: "Boards", href: "/boards", icon: LayoutDashboard },
  { name: "All Tasks", href: "/tasks", icon: CheckSquare },
  { name: "Settings", href: "/settings", icon: Settings },
];

export function Sidebar() {
  const pathname = usePathname();
  const { logout } = useAuth();
  const { user, isLoading } = useUser();
  
  const { projects } = useProjects();
  const { teams } = useTeams();

  const userInitial = user?.full_name
    ? user.full_name.split(" ").map((n) => n[0]).join("").substring(0, 2).toUpperCase()
    : user?.email?.substring(0, 2).toUpperCase() || "??";

  return (
    <div className="flex h-full w-64 flex-col border-r border-gray-200 bg-gray-50/50">
      <div className="p-4">
        <WorkspaceSwitcher />
      </div>
      
      <div className="flex-1 overflow-y-auto px-3 py-2 space-y-6">
        {/* Main Nav */}
        <nav className="space-y-1">
          {navigation.map((item) => {
            const isActive = pathname === item.href;
            return (
              <Link
                key={item.name}
                href={item.href}
                className={cn(
                  "group flex items-center rounded-md px-3 py-2 text-sm font-medium transition-colors",
                  isActive
                    ? "bg-blue-50 text-blue-600"
                    : "text-gray-700 hover:bg-gray-100 hover:text-gray-900"
                )}
              >
                <item.icon
                  className={cn(
                    "mr-3 h-5 w-5 flex-shrink-0 transition-colors",
                    isActive ? "text-blue-600" : "text-gray-400 group-hover:text-gray-500"
                  )}
                />
                {item.name}
              </Link>
            );
          })}
        </nav>

        {/* Projects Section */}
        <div>
          <div className="px-3 mb-2 flex items-center justify-between">
            <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider">Projects</h3>
          </div>
          <div className="space-y-1">
            {projects.map((project) => (
              <Link
                key={project.id}
                href={`/projects/${project.id}`}
                className={cn(
                  "group flex items-center rounded-md px-3 py-2 text-sm font-medium transition-colors",
                  pathname === `/projects/${project.id}`
                    ? "bg-blue-50 text-blue-600"
                    : "text-gray-700 hover:bg-gray-100 hover:text-gray-900"
                )}
              >
                <Folder className="mr-3 h-4 w-4 text-gray-400" />
                <span className="truncate">{project.name}</span>
              </Link>
            ))}
            <div className="px-1 pt-1">
              <ProjectDialog />
            </div>
          </div>
        </div>

        {/* Teams Section */}
        <div>
          <div className="px-3 mb-2 flex items-center justify-between">
            <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider">Teams</h3>
          </div>
          <div className="space-y-1">
            {teams.map((team) => (
              <Link
                key={team.id}
                href={`/teams/${team.id}`}
                className={cn(
                  "group flex items-center rounded-md px-3 py-2 text-sm font-medium transition-colors",
                  pathname === `/teams/${team.id}`
                    ? "bg-blue-50 text-blue-600"
                    : "text-gray-700 hover:bg-gray-100 hover:text-gray-900"
                )}
              >
                <Users className="mr-3 h-4 w-4 text-gray-400" />
                <span className="truncate">{team.name}</span>
              </Link>
            ))}
            <div className="px-1 pt-1">
              <TeamDialog />
            </div>
          </div>
        </div>
      </div>

      <div className="border-t border-gray-200 p-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
              <div className="h-8 w-8 rounded-full bg-blue-600 flex items-center justify-center text-white font-semibold text-xs">
                  {isLoading ? "..." : userInitial}
              </div>
              <div className="text-sm">
                  <p className="font-medium text-gray-700">{isLoading ? "Loading..." : (user?.full_name || user?.email?.split("@")[0])}</p>
                  <p className="text-xs text-gray-500">{isLoading ? "Please wait" : user?.email}</p>
              </div>
          </div>
          <button 
            onClick={logout}
            className="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-md transition-colors"
            title="Logout"
          >
            <LogOut className="h-4 w-4" />
          </button>
        </div>
      </div>
    </div>
  );
}
