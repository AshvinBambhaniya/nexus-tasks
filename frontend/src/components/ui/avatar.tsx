import * as React from "react"
import { cn } from "@/lib/utils"

interface AvatarProps extends React.HTMLAttributes<HTMLDivElement> {
  src?: string
  alt?: string
  fallback: string
}

export function Avatar({ className, src, alt, fallback, ...props }: AvatarProps) {
  return (
    <div
      className={cn(
        "relative flex h-8 w-8 shrink-0 overflow-hidden rounded-full bg-gray-100",
        className
      )}
      {...props}
    >
      {src ? (
        <img
          className="aspect-square h-full w-full object-cover"
          src={src}
          alt={alt || ""}
        />
      ) : (
        <div className="flex h-full w-full items-center justify-center bg-blue-100 text-xs font-medium text-blue-600">
          {fallback}
        </div>
      )}
    </div>
  )
}
