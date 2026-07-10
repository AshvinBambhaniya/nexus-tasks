<script setup lang="ts">
import {
  Loader2,
  Inbox,
  Search,
  Check,
  CheckSquare,
  MessageCircle,
  MoreHorizontal,
} from "lucide-vue-next";
import {
  type Notification,
  NotificationType,
  NotificationEntityType,
} from "~/types";
import { useInbox } from "~/composables/useInbox";
import { formatDistanceToNow } from "date-fns";

definePageMeta({ layout: "dashboard" });

const {
  notifications,
  isLoading,
  fetchInbox,
  markAsRead,
  clearNotification,
  clearAll,
} = useInbox();

onMounted(() => {
  fetchInbox();
});

// Search and Filter State
const searchQuery = ref("");
const activeNotificationId = ref<string | null>(null);

const activeNotification = computed(() => {
  if (!activeNotificationId.value) return null;
  return (
    notifications.value?.find((n) => n.id === activeNotificationId.value) ||
    null
  );
});

// Watch notifications to auto-select the first one if none is selected
watch(notifications, (newNotifications) => {
  if (
    newNotifications &&
    newNotifications.length > 0 &&
    !activeNotificationId.value
  ) {
    activeNotificationId.value = newNotifications[0]?.id || null;
    if (!newNotifications[0].is_read) {
      markAsRead(newNotifications[0].id);
    }
  }
});

const filteredNotifications = computed(() => {
  if (!notifications.value) return [];
  let result = [...notifications.value];
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase();
    result = result.filter(
      (n) =>
        n.title.toLowerCase().includes(q) || n.body?.toLowerCase().includes(q)
    );
  }
  return result;
});

const selectNotification = (n: Notification) => {
  activeNotificationId.value = n.id;
  if (!n.is_read) {
    markAsRead(n.id);
  }
};

const handleClear = (e: Event, id: string) => {
  e.stopPropagation();
  clearNotification(id);
  if (activeNotificationId.value === id) {
    activeNotificationId.value = null;
  }
};

const getNotificationIcon = (type: NotificationType) => {
  if (
    type === NotificationType.COMMENT_ADDED ||
    type === NotificationType.MENTIONED
  )
    return MessageCircle;
  return CheckSquare;
};
</script>

<template>
  <div
    class="border-border bg-card animate-in fade-in mx-auto flex h-[calc(100vh-8rem)] max-h-[1000px] w-full max-w-[1600px] overflow-hidden rounded-xl border shadow-sm duration-500"
  >
    <!-- Left Pane: List View -->
    <div
      class="border-border bg-card/50 flex w-[350px] shrink-0 flex-col border-r xl:w-[400px]"
    >
      <!-- Toolbar -->
      <div class="border-border flex flex-col gap-3 border-b p-4">
        <div class="flex items-center justify-between">
          <h2 class="text-foreground text-sm font-semibold tracking-tight">
            Inbox
          </h2>
          <div class="flex items-center gap-2">
            <button
              class="text-muted-foreground hover:text-foreground hover:bg-muted rounded-md p-1.5 transition-colors"
              title="Clear all read"
              @click="clearAll"
            >
              <Check class="h-4 w-4" />
            </button>
            <button
              class="text-muted-foreground hover:text-foreground hover:bg-muted rounded-md p-1.5 transition-colors"
            >
              <MoreHorizontal class="h-4 w-4" />
            </button>
          </div>
        </div>

        <div class="relative">
          <Search
            class="text-muted-foreground absolute top-2 left-2.5 h-4 w-4"
          />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search notifications..."
            class="border-border bg-muted/50 text-foreground focus:border-primary focus:ring-primary h-8 w-full rounded-md border pr-3 pl-9 text-sm transition-all focus:ring-1 focus:outline-none"
          />
        </div>
      </div>

      <!-- Feed -->
      <div class="custom-scrollbar flex-1 overflow-y-auto">
        <div
          v-if="isLoading && notifications.length === 0"
          class="flex items-center justify-center py-10"
        >
          <Loader2 class="text-muted-foreground h-6 w-6 animate-spin" />
        </div>

        <div
          v-else-if="filteredNotifications.length === 0"
          class="flex flex-col items-center justify-center px-6 py-20 text-center"
        >
          <Inbox class="text-muted-foreground/30 mb-4 h-12 w-12" />
          <h3 class="text-foreground text-base font-semibold">Inbox Zero</h3>
          <p class="text-muted-foreground mt-1 text-sm">
            You're all caught up! Take a breather.
          </p>
        </div>

        <div v-else class="flex flex-col">
          <button
            v-for="notification in filteredNotifications"
            :key="notification.id"
            :class="[
              'border-border/50 group relative flex flex-col items-start gap-1.5 border-b p-4 text-left transition-colors',
              activeNotificationId === notification.id
                ? 'bg-muted'
                : 'hover:bg-muted/50',
            ]"
            @click="selectNotification(notification)"
          >
            <!-- Active Indicator -->
            <div
              v-if="activeNotificationId === notification.id"
              class="bg-primary absolute top-0 left-0 h-full w-[3px]"
            />

            <!-- Unread Indicator -->
            <div
              v-if="!notification.is_read"
              class="absolute top-4 right-4 h-2.5 w-2.5 rounded-full bg-blue-500 shadow-[0_0_8px_rgba(59,130,246,0.6)]"
            />

            <!-- Clear Action (Hover) -->
            <button
              v-if="notification.is_read"
              class="bg-background/80 hover:bg-muted text-muted-foreground hover:text-foreground absolute top-3 right-3 rounded-md p-1.5 opacity-0 shadow-sm transition-all group-hover:opacity-100"
              title="Clear notification"
              @click="handleClear($event, notification.id)"
            >
              <Check class="h-3.5 w-3.5" />
            </button>

            <div class="flex w-full items-center gap-2 pr-6">
              <component
                :is="getNotificationIcon(notification.type)"
                class="text-muted-foreground h-3.5 w-3.5"
              />
              <span
                class="text-muted-foreground/70 flex-1 truncate text-[10px] font-bold tracking-widest uppercase"
              >
                {{ notification.type.replace("_", " ") }}
              </span>
              <span class="text-muted-foreground/50 text-[10px]"
                >{{
                  formatDistanceToNow(new Date(notification.created_at))
                }}
                ago</span
              >
            </div>

            <h4
              :class="[
                'w-full truncate pr-4 text-sm tracking-tight',
                !notification.is_read
                  ? 'text-foreground font-bold'
                  : 'text-foreground/80 font-medium',
              ]"
            >
              {{ notification.title }}
            </h4>

            <p
              v-if="notification.body"
              class="text-muted-foreground line-clamp-2 w-full text-xs"
            >
              {{ notification.body }}
            </p>
          </button>
        </div>
      </div>
    </div>

    <!-- Right Pane: Detail View -->
    <div class="bg-background/30 flex flex-1 flex-col">
      <div
        v-if="isLoading && !activeNotification"
        class="flex h-full items-center justify-center"
      >
        <Loader2 class="text-muted-foreground h-8 w-8 animate-spin" />
      </div>

      <div
        v-else-if="!activeNotification"
        class="flex h-full flex-col items-center justify-center text-center"
      >
        <Inbox class="text-muted-foreground/30 mb-4 h-12 w-12" />
        <h3 class="text-muted-foreground text-lg font-medium">
          Select a notification to view details
        </h3>
      </div>

      <!-- Conditionally render context based on entity type -->
      <template
        v-else-if="
          activeNotification.entity_type === NotificationEntityType.TASK ||
          activeNotification.entity_type === NotificationEntityType.COMMENT
        "
      >
        <TaskDetail
          :key="activeNotification.entity_id"
          :task-id="activeNotification.entity_id"
          :project-id="''"
        />
      </template>
      <div
        v-else
        class="flex h-full flex-col items-center justify-center text-center"
      >
        <p class="text-muted-foreground text-sm">
          Context preview not available for this entity type.
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: var(--muted);
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: var(--border);
}
</style>
