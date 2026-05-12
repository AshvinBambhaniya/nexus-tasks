export const useSocket = () => {
  const workspaceStore = useWorkspaceStore();
  const config = useRuntimeConfig();
  const socket = ref<WebSocket | null>(null);

  const connect = () => {
    if (!workspaceStore.activeWorkspaceId) return;

    const wsProtocol = config.public.apiUrl.startsWith("https") ? "wss" : "ws";
    const wsBaseUrl = config.public.apiUrl.replace(/^https?:\/\//, "");
    const wsUrl = `${wsProtocol}://${wsBaseUrl}/ws/${workspaceStore.activeWorkspaceId}`;

    socket.value = new WebSocket(wsUrl);

    socket.value.onopen = () => {
      console.log("WebSocket Connected");
    };

    socket.value.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        handleSocketMessage(data);
      } catch (err) {
        console.error("Failed to parse WS message", err);
      }
    };

    socket.value.onclose = () => {
      console.log("WebSocket Disconnected");
    };
  };

  const disconnect = () => {
    if (socket.value) {
      socket.value.close();
      socket.value = null;
    }
  };

  const handleSocketMessage = (message: { type: string }) => {
    // For now, we'll use clearNuxtData or specific refresh
    // depending on the message type. In a full implementation,
    // we would trigger refreshes for specific SWR-like keys.

    // Example: if TASK_UPDATED, refresh task-related useFetch calls
    if (
      message.type === "TASK_CREATED" ||
      message.type === "TASK_UPDATED" ||
      message.type === "TASK_DELETED"
    ) {
      // Global refresh for any active useFetch with these keys
      refreshNuxtData(
        (key) => key.includes("/api/v2/projects/") || key === "my-tasks"
      );
    }
  };

  // Watch for workspace changes
  watch(
    () => workspaceStore.activeWorkspaceId,
    () => {
      disconnect();
      connect();
    }
  );

  onMounted(() => {
    connect();
  });

  onUnmounted(() => {
    disconnect();
  });

  return {
    socket,
  };
};
