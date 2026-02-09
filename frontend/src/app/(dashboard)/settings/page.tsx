"use client";

import { useState } from "react";
import { GeneralSettings } from "@/components/settings/general-settings";
import { MemberSettings } from "@/components/settings/member-settings";
import { cn } from "@/lib/utils";
import { User, Users } from "lucide-react";

export default function SettingsPage() {
  const [activeTab, setActiveTab] = useState<"general" | "members">("general");

  const tabs = [
    { id: "general", label: "General", icon: User },
    { id: "members", label: "Members", icon: Users },
  ] as const;

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-bold tracking-tight text-gray-900">
          Settings
        </h1>
        <p className="mt-2 text-gray-500">
          Manage your workspace and profile preferences.
        </p>
      </div>

      <div className="flex flex-col gap-8 md:flex-row">
        {/* Sidebar Navigation for Settings */}
        <aside className="w-full flex-shrink-0 md:w-64">
          <nav className="flex flex-row space-x-2 md:flex-col md:space-y-1 md:space-x-0">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={cn(
                  "flex w-full items-center rounded-md px-3 py-2 text-sm font-medium transition-colors",
                  activeTab === tab.id
                    ? "bg-blue-50 text-blue-700"
                    : "text-gray-700 hover:bg-gray-100 hover:text-gray-900"
                )}
              >
                <tab.icon
                  className={cn(
                    "mr-3 h-4 w-4 flex-shrink-0",
                    activeTab === tab.id
                      ? "text-blue-700"
                      : "text-gray-400 group-hover:text-gray-500"
                  )}
                />
                {tab.label}
              </button>
            ))}
          </nav>
        </aside>

        {/* Content Area */}
        <div className="min-w-0 flex-1">
          {activeTab === "general" ? <GeneralSettings /> : <MemberSettings />}
        </div>
      </div>
    </div>
  );
}
