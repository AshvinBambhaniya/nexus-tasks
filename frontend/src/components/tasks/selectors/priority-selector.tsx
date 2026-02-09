"use client";

import { useState, useRef } from "react";
import {
  Check,
  ChevronDown,
  SignalHigh,
  SignalMedium,
  SignalLow,
  AlertOctagon,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { useClickOutside } from "@/hooks/use-click-outside";
import { TaskPriority } from "@/types";

interface PrioritySelectorProps {
  value: TaskPriority;
  onChange: (value: TaskPriority) => void;
}

const priorityConfig = {
  [TaskPriority.P0]: {
    label: "Critical",
    icon: AlertOctagon,
    color: "text-red-600",
    bg: "bg-red-50",
  },
  [TaskPriority.P1]: {
    label: "High",
    icon: SignalHigh,
    color: "text-orange-500",
    bg: "bg-orange-50",
  },
  [TaskPriority.P2]: {
    label: "Medium",
    icon: SignalMedium,
    color: "text-blue-500",
    bg: "bg-blue-50",
  },
  [TaskPriority.P3]: {
    label: "Low",
    icon: SignalLow,
    color: "text-gray-400",
    bg: "bg-gray-50",
  },
};

export function PrioritySelector({ value, onChange }: PrioritySelectorProps) {
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useClickOutside(ref, () => setOpen(false));

  const current = priorityConfig[value];
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
          {Object.entries(priorityConfig).map(([key, config]) => {
            const PriorityIcon = config.icon;
            const isSelected = value === key;
            return (
              <button
                key={key}
                type="button"
                onClick={() => {
                  onChange(key as TaskPriority);
                  setOpen(false);
                }}
                className={cn(
                  "flex w-full items-center justify-between rounded-sm px-2 py-1.5 text-sm transition-colors outline-none hover:bg-gray-100",
                  isSelected && "bg-gray-50"
                )}
              >
                <div className="flex items-center gap-2">
                  <PriorityIcon className={cn("h-4 w-4", config.color)} />
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
