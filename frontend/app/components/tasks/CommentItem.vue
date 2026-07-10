<script setup lang="ts">
import { formatDistanceToNow } from "date-fns";
import type { Comment } from "~/types";

interface Props {
  comment: Comment;
  currentUserId?: string;
}

const { comment, currentUserId = "" } = defineProps<Props>();

const emit = defineEmits(["delete"]);

const author = computed(
  () =>
    comment.author || {
      id: "",
      email: "unknown@user.com",
      full_name: "Unknown User",
    }
);

const authorName = computed(
  () => author.value.full_name || author.value.email.split("@")[0]
);
const initial = computed(
  () =>
    (author.value.full_name || author.value.email)?.[0]?.toUpperCase() || "?"
);
const isAuthor = computed(() => currentUserId === author.value.id);
const formattedContent = computed(() => {
  if (!currentUserId || !comment.content) return comment.content || "";
  return comment.content.replaceAll(
    `data-id="${currentUserId}"`,
    `data-id="${currentUserId}" data-is-me="true"`
  );
});
</script>

<template>
  <!-- eslint-disable vue/no-v-html -->
  <div class="group flex gap-4">
    <UiBaseAvatar
      :fallback="initial"
      class-name="mt-1 h-10 w-10 border border-border shadow-sm"
    />
    <div class="min-w-0 flex-1">
      <div class="mb-1 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <span
            class="text-foreground cursor-pointer text-sm font-semibold hover:underline"
          >
            {{ authorName }}
          </span>
          <span class="text-muted-foreground text-xs">
            commented
            {{ formatDistanceToNow(new Date(comment.created_at)) }} ago
          </span>
        </div>
        <button
          v-if="isAuthor"
          class="text-muted-foreground/60 hover:text-destructive text-xs opacity-0 transition-opacity group-hover:opacity-100"
          @click="emit('delete', comment.id)"
        >
          Delete
        </button>
      </div>
      <!-- eslint-disable vue/no-v-html -->
      <div
        class="prose prose-sm prose-pre:bg-muted prose-pre:border prose-pre:border-border border-border bg-card max-w-none rounded-lg border p-4 shadow-sm transition-colors"
        v-html="formattedContent"
      />
      <!-- eslint-enable vue/no-v-html -->
    </div>
  </div>
</template>
