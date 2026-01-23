import { useEffect, useRef } from "react";
import Cookies from "js-cookie";
import { useSWRConfig } from "swr";
import { useWorkspaceStore } from "@/store/workspace-store";
import { API_URL } from "@/lib/api";
import { Task } from "@/types";

export function useSocket() {
  const { activeWorkspaceId } = useWorkspaceStore();
  const { mutate } = useSWRConfig();
  const socketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    if (!activeWorkspaceId) return;

    const token = Cookies.get("token");
    if (!token) return;

    // Construct WS URL
    const wsProtocol = API_URL.startsWith("https") ? "wss" : "ws";
    const wsBaseUrl = API_URL.replace(/^https?:\/\//, "");
    const wsUrl = `${wsProtocol}://${wsBaseUrl}/ws/${activeWorkspaceId}?token=${token}`;

    const socket = new WebSocket(wsUrl);
    socketRef.current = socket;

    socket.onopen = () => {
      console.log("WebSocket Connected");
    };

    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        handleMessage(data, activeWorkspaceId, mutate);
      } catch (err) {
        console.error("Failed to parse WS message", err);
      }
    };

    socket.onclose = () => {
      console.log("WebSocket Disconnected");
    };

    return () => {
      socket.close();
    };
  }, [activeWorkspaceId, mutate]);

  return socketRef.current;
}

function handleMessage(message: any, workspaceId: number, mutate: any) {
  const tasksKey = `/api/v1/workspaces/${workspaceId}/tasks`;
  const myTasksKey = "/api/v1/tasks/me";

  switch (message.type) {
    case "TASK_CREATED":
    case "TASK_UPDATED":
      const task: Task = message.task;
      
      // Update Workspace Tasks
      mutate(tasksKey, (currentTasks: Task[] = []) => {
        if (message.type === "TASK_CREATED") {
          return [...currentTasks, task];
        }
        return currentTasks.map((t) => (t.id === task.id ? task : t));
      }, false); // false = do not revalidate immediately

      // Update My Tasks (Inbox) - only if assigned to me (but we don't know current user id here easily without store)
      // For simplicity, we can revalidate the inbox or optimistic update if we knew the user ID.
      // Let's just revalidate inbox to be safe and simple.
      mutate(myTasksKey);
      break;

    case "TASK_DELETED":
      const taskId = message.task_id;
      
      mutate(tasksKey, (currentTasks: Task[] = []) => {
        return currentTasks.filter((t) => t.id !== taskId);
      }, false);

      mutate(myTasksKey);
      break;
  }
}
