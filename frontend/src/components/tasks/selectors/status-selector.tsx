"use client";

import { useState, useRef } from "react";
import {
  Check,
  ChevronDown,
  Circle,
  PlayCircle,
  CheckCircle2,
  Archive,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { useClickOutside } from "@/hooks/use-click-outside";
import { TaskStatus } from "@/types";

interface StatusSelectorProps {
  value: TaskStatus;
  onChange: (value: TaskStatus) => void;
}

const statusConfig = {
  [TaskStatus.TODO]: {
    label: "Todo",
    icon: Circle,
    color: "text-gray-500",
    bg: "bg-gray-100",
  },
  [TaskStatus.IN_PROGRESS]: {
    label: "In Progress",
    icon: PlayCircle,
    color: "text-blue-500",
    bg: "bg-blue-50",
  },
  [TaskStatus.DONE]: {
    label: "Done",
    icon: CheckCircle2,
    color: "text-green-500",
    bg: "bg-green-50",
  },
  [TaskStatus.BACKLOG]: {
    label: "Backlog",
    icon: Archive,
    color: "text-gray-400",
    bg: "bg-gray-50",
  },
};

export function StatusSelector({ value, onChange }: StatusSelectorProps) {
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useClickOutside(ref, () => setOpen(false));

  const current = statusConfig[value];
  const Icon = current.icon;

  return (
    <div className="relative" ref={ref}>
      <button
        type="button"
        onClick={() => setOpen(!open)}
        className={cn(
          "flex w-full items-center justify-between rounded-md border border-gray-200 bg-white px-3 py-2 text-sm ring-offset-white transition-colors hover:bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none"
        )}
      >
        <div className="flex items-center gap-2">
          <Icon className={cn("h-4 w-4", current.color)} />
          <span className="font-medium text-gray-700">{current.label}</span>
        </div>
        <ChevronDown className="h-4 w-4 text-gray-400" />
      </button>

      {open && (
        <div className="animate-in fade-in-0 zoom-in-95 absolute top-full z-50 mt-2 w-full rounded-md border border-gray-200 bg-white p-1 shadow-lg duration-100">
          {Object.entries(statusConfig).map(([key, config]) => {
            const StatusIcon = config.icon;
            const isSelected = value === key;
            return (
              <button
                key={key}
                type="button"
                onClick={() => {
                  onChange(key as TaskStatus);
                  setOpen(false);
                }}
                className={cn(
                  "flex w-full items-center justify-between rounded-sm px-2 py-1.5 text-sm transition-colors outline-none hover:bg-gray-100",
                  isSelected && "bg-gray-50"
                )}
              >
                <div className="flex items-center gap-2">
                  <StatusIcon className={cn("h-4 w-4", config.color)} />
                  <span className="text-gray-700">{config.label}</span>
                </div>
                {isSelected && <Check className="h-4 w-4 text-blue-600" />}
              </button>
            );
          })}
        </div>
      )}
    </div>
  );
}
