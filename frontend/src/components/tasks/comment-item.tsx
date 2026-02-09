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

export function CommentItem({ comment, currentUserId, onDelete }: CommentItemProps) {
  const authorName = comment.author.full_name || comment.author.email.split("@")[0];
  const initial = (comment.author.full_name || comment.author.email)[0].toUpperCase();
  const isAuthor = currentUserId === comment.author.id;

  return (
    <div className="flex gap-4 group">
      <Avatar fallback={initial} className="h-10 w-10 mt-1 shadow-sm border border-gray-100" />
      <div className="flex-1 min-w-0">
        <div className="flex items-center justify-between mb-1">
            <div className="flex items-center gap-2">
                <span className="font-semibold text-gray-900 text-sm hover:underline cursor-pointer">
                    {authorName}
                </span>
                <span className="text-gray-500 text-xs">
                    commented {formatDistanceToNow(new Date(comment.created_at))} ago
                </span>
            </div>
            {isAuthor && onDelete && (
                <button 
                    onClick={() => onDelete(comment.id)}
                    className="text-gray-400 hover:text-red-600 text-xs opacity-0 group-hover:opacity-100 transition-opacity"
                >
                    Delete
                </button>
            )}
        </div>
        <div className="bg-white border border-gray-200 rounded-lg p-4 shadow-sm prose prose-sm max-w-none prose-pre:bg-gray-50 prose-pre:border prose-pre:border-gray-100">
          <ReactMarkdown remarkPlugins={[remarkGfm]}>{comment.content}</ReactMarkdown>
        </div>
      </div>
    </div>
  );
}
