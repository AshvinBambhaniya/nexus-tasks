"use client";

import { useState } from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { cn } from "@/lib/utils";
import { Textarea } from "./textarea";

interface MarkdownEditorProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  className?: string;
}

export function MarkdownEditor({
  value,
  onChange,
  placeholder,
  className,
}: MarkdownEditorProps) {
  const [tab, setTab] = useState<"write" | "preview">("write");

  return (
    <div
      className={cn(
        "flex flex-col overflow-hidden rounded-md border border-gray-200 bg-white shadow-sm transition-all focus-within:border-blue-500 focus-within:ring-2 focus-within:ring-blue-500/20",
        className
      )}
    >
      {/* Header Tabs */}
      <div className="flex items-center border-b border-gray-200 bg-white px-1">
        <button
          type="button"
          onClick={() => setTab("write")}
          className={cn(
            "-mb-px border-b-2 px-4 py-2 text-sm font-medium transition-all",
            tab === "write"
              ? "border-blue-600 text-blue-600"
              : "border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700"
          )}
        >
          Write
        </button>
        <button
          type="button"
          onClick={() => setTab("preview")}
          className={cn(
            "-mb-px border-b-2 px-4 py-2 text-sm font-medium transition-all",
            tab === "preview"
              ? "border-blue-600 text-blue-600"
              : "border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700"
          )}
        >
          Preview
        </button>
      </div>

      {/* Editor Content */}
      <div className="relative">
        {tab === "write" ? (
          <Textarea
            value={value}
            onChange={(e) => onChange(e.target.value)}
            placeholder={placeholder}
            className="min-h-[200px] w-full resize-y rounded-none border-0 px-4 py-3 font-sans text-sm focus-visible:ring-0 focus-visible:ring-offset-0"
          />
        ) : (
          <div className="prose prose-sm min-h-[200px] max-w-none bg-gray-50/30 p-4">
            {value ? (
              <ReactMarkdown remarkPlugins={[remarkGfm]}>{value}</ReactMarkdown>
            ) : (
              <span className="text-sm text-gray-400 italic">
                Nothing to preview
              </span>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
