"use client";

import { Avatar } from "@/components/ui/avatar";
import { formatDistanceToNow } from "date-fns";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

interface CommentItemProps {
  comment: {
    id: number;
    content: string;
    created_at: string;
    author: {
      id: number;
      email: string;
      full_name?: string;
    };
  };
  currentUserId?: number;
  onDelete?: (id: number) => void;
}

export function CommentItem({
  comment,
  currentUserId,
  onDelete,
}: CommentItemProps) {
  const authorName =
    comment.author.full_name || comment.author.email.split("@")[0];
  const initial = (comment.author.full_name ||
    comment.author.email)[0].toUpperCase();
  const isAuthor = currentUserId === comment.author.id;

  return (
    <div className="group flex gap-4">
      <Avatar
        fallback={initial}
        className="mt-1 h-10 w-10 border border-gray-100 shadow-sm"
      />
      <div className="min-w-0 flex-1">
        <div className="mb-1 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <span className="cursor-pointer text-sm font-semibold text-gray-900 hover:underline">
              {authorName}
            </span>
            <span className="text-xs text-gray-500">
              commented {formatDistanceToNow(new Date(comment.created_at))} ago
            </span>
          </div>
          {isAuthor && onDelete && (
            <button
              onClick={() => onDelete(comment.id)}
              className="text-xs text-gray-400 opacity-0 transition-opacity group-hover:opacity-100 hover:text-red-600"
            >
              Delete
            </button>
          )}
        </div>
        <div className="prose prose-sm prose-pre:bg-gray-50 prose-pre:border prose-pre:border-gray-100 max-w-none rounded-lg border border-gray-200 bg-white p-4 shadow-sm">
          <ReactMarkdown remarkPlugins={[remarkGfm]}>
            {comment.content}
          </ReactMarkdown>
        </div>
      </div>
    </div>
  );
}
