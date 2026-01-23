"use client";

import { useState } from "react";
import { GeneralSettings } from "@/components/settings/general-settings";
import { MemberSettings } from "@/components/settings/member-settings";

export default function SettingsPage() {
  const [activeTab, setActiveTab] = useState<"general" | "members">("general");

  return (
    <div className="space-y-6">
      <div className="border-b border-gray-200">
        <h1 className="text-2xl font-bold tracking-tight text-gray-900 pb-4">Settings</h1>
        <div className="flex gap-6 -mb-px">
          <button
            onClick={() => setActiveTab("general")}
            className={`pb-3 text-sm font-medium transition-colors border-b-2 ${
              activeTab === "general"
                ? "border-blue-600 text-blue-600"
                : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
            }`}
          >
            General
          </button>
          <button
            onClick={() => setActiveTab("members")}
            className={`pb-3 text-sm font-medium transition-colors border-b-2 ${
              activeTab === "members"
                ? "border-blue-600 text-blue-600"
                : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
            }`}
          >
            Members
          </button>
        </div>
      </div>

      <div className="py-2">
        {activeTab === "general" ? <GeneralSettings /> : <MemberSettings />}
      </div>
    </div>
  );
}