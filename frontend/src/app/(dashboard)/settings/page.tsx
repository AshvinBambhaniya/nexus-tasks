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
        <h1 className="text-3xl font-bold tracking-tight text-gray-900">Settings</h1>
        <p className="text-gray-500 mt-2">Manage your workspace and profile preferences.</p>
      </div>

      <div className="flex flex-col md:flex-row gap-8">
        {/* Sidebar Navigation for Settings */}
        <aside className="w-full md:w-64 flex-shrink-0">
          <nav className="flex flex-row md:flex-col space-x-2 md:space-x-0 md:space-y-1">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={cn(
                  "flex items-center w-full px-3 py-2 text-sm font-medium rounded-md transition-colors",
                  activeTab === tab.id
                    ? "bg-blue-50 text-blue-700"
                    : "text-gray-700 hover:bg-gray-100 hover:text-gray-900"
                )}
              >
                <tab.icon
                  className={cn(
                    "mr-3 h-4 w-4 flex-shrink-0",
                    activeTab === tab.id ? "text-blue-700" : "text-gray-400 group-hover:text-gray-500"
                  )}
                />
                {tab.label}
              </button>
            ))}
          </nav>
        </aside>

        {/* Content Area */}
        <div className="flex-1 min-w-0">
          {activeTab === "general" ? <GeneralSettings /> : <MemberSettings />}
        </div>
      </div>
    </div>
  );
}
