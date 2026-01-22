"use client";

import { WorkspaceSwitcher } from "@/components/workspace/workspace-switcher";
import { LayoutDashboard, CheckSquare, Inbox, Settings } from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";

const navigation = [
  { name: "My Focus", href: "/inbox", icon: Inbox },
  { name: "Boards", href: "/boards", icon: LayoutDashboard },
  { name: "All Tasks", href: "/tasks", icon: CheckSquare },
  { name: "Settings", href: "/settings", icon: Settings },
];

export function Sidebar() {
  const pathname = usePathname();

  return (
    <div className="flex h-full w-64 flex-col border-r border-gray-200 bg-gray-50/50">
      <div className="p-4">
        <WorkspaceSwitcher />
      </div>
      
      <div className="flex-1 overflow-y-auto px-3 py-2">
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
      </div>

      <div className="border-t border-gray-200 p-4">
        <div className="flex items-center gap-3">
            <div className="h-8 w-8 rounded-full bg-blue-600 flex items-center justify-center text-white font-semibold text-xs">
                US
            </div>
            <div className="text-sm">
                <p className="font-medium text-gray-700">User Name</p>
                <p className="text-xs text-gray-500">user@example.com</p>
            </div>
        </div>
      </div>
    </div>
  );
}
