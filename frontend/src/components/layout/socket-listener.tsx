"use client";

import { useSocket } from "@/hooks/use-socket";

export function SocketListener() {
  useSocket();
  return null;
}
