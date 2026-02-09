"use client";

import { useState, useRef } from "react";
import { Command } from "cmdk";
import { Check, ChevronsUpDown, User as UserIcon, Search } from "lucide-react";
import { cn } from "@/lib/utils";
import { useClickOutside } from "@/hooks/use-click-outside";
import { Avatar } from "@/components/ui/avatar";
import { ProjectMemberResponse } from "@/types";

interface AssigneeSelectorProps {
  members: ProjectMemberResponse[];
  value?: number;
  onChange: (value: number | undefined) => void;
}

export function AssigneeSelector({ members, value, onChange }: AssigneeSelectorProps) {
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useClickOutside(ref, () => setOpen(false));

  const selectedMember = members.find((m) => m.user_id === value);

  return (
    <div className="relative w-full" ref={ref}>
      <button
        type="button"
        onClick={() => setOpen(!open)}
        className={cn(
          "flex w-full items-center justify-between rounded-md border border-gray-200 bg-white px-3 py-2 text-sm ring-offset-white placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:cursor-not-allowed disabled:opacity-50 hover:bg-gray-50 transition-colors",
          !value && "text-gray-500"
        )}
      >
        <div className="flex items-center gap-2 flex-1 min-w-0">
          {selectedMember ? (
            <>
               <Avatar 
                 fallback={selectedMember.email[0].toUpperCase()} 
                 className="h-5 w-5 bg-blue-100 text-[10px] shrink-0"
               />
               <span className="truncate">{selectedMember.email}</span>
            </>
          ) : (
            <>
              <div className="flex h-5 w-5 items-center justify-center rounded-full border border-dashed border-gray-400">
                <UserIcon className="h-3 w-3 text-gray-400" />
              </div>
              <span className="flex-1 text-left">Unassigned</span>
            </>
          )}
        </div>
        <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
      </button>

      {open && (
        <div className="absolute top-full mt-2 z-50 w-full min-w-[200px] rounded-md border border-gray-200 bg-white shadow-lg animate-in fade-in-0 zoom-in-95 duration-100">
          <Command className="w-full overflow-hidden rounded-md">
            <div className="flex items-center border-b px-3">
              <Search className="mr-2 h-4 w-4 shrink-0 opacity-50" />
              <Command.Input
                placeholder="Search member..."
                className="flex h-9 w-full rounded-md bg-transparent py-3 text-sm outline-none placeholder:text-gray-500 disabled:cursor-not-allowed disabled:opacity-50"
              />
            </div>
            <Command.List className="max-h-[300px] overflow-y-auto p-1">
              <Command.Empty className="py-2 text-center text-xs text-gray-500">
                No member found.
              </Command.Empty>
              
              <Command.Item
                onSelect={() => {
                  onChange(undefined);
                  setOpen(false);
                }}
                className={cn(
                  "relative flex cursor-default select-none items-center rounded-sm px-2 py-2 text-sm outline-none data-[selected=true]:bg-gray-100 data-[selected=true]:text-gray-900",
                  !value && "bg-gray-50 text-gray-900"
                )}
              >
                <div className="flex h-5 w-5 items-center justify-center rounded-full border border-dashed border-gray-400 mr-2 shrink-0">
                  <UserIcon className="h-3 w-3 text-gray-400" />
                </div>
                <span className="flex-1 truncate">Unassigned</span>
                {!value && <Check className="ml-2 h-4 w-4 shrink-0 text-blue-600" />}
              </Command.Item>

              {members.map((member) => (
                <Command.Item
                  key={member.user_id}
                  onSelect={() => {
                    onChange(member.user_id);
                    setOpen(false);
                  }}
                  className={cn(
                    "relative flex cursor-default select-none items-center rounded-sm px-2 py-2 text-sm outline-none data-[selected=true]:bg-gray-100 data-[selected=true]:text-gray-900",
                    value === member.user_id && "bg-gray-50 text-gray-900"
                  )}
                >
                  <Avatar 
                    fallback={member.email[0].toUpperCase()} 
                    className="mr-2 h-5 w-5 shrink-0 bg-blue-100 text-[10px]"
                  />
                  <span className="flex-1 truncate">{member.email}</span>
                  {value === member.user_id && (
                    <Check className="ml-2 h-4 w-4 shrink-0 text-blue-600" />
                  )}
                </Command.Item>
              ))}
            </Command.List>
          </Command>
        </div>
      )}
    </div>
  );
}
