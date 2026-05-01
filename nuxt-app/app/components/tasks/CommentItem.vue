<script setup lang="ts">
import { formatDistanceToNow } from "date-fns";
import VueMarkdown from "vue-markdown-render";
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
const initial = computed(() =>
  (author.value.full_name || author.value.email)[0].toUpperCase()
);
const isAuthor = computed(() => currentUserId === author.value.id);
</script>

<template>
  <div class="group flex gap-4">
    <UiBaseAvatar
      :fallback="initial"
      class-name="mt-1 h-10 w-10 border border-gray-100 shadow-sm"
    />
    <div class="min-w-0 flex-1">
      <div class="mb-1 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <span
            class="cursor-pointer text-sm font-semibold text-gray-900 hover:underline"
          >
            {{ authorName }}
          </span>
          <span class="text-xs text-gray-500">
            commented
            {{ formatDistanceToNow(new Date(comment.created_at)) }} ago
          </span>
        </div>
        <button
          v-if="isAuthor"
          class="text-xs text-gray-400 opacity-0 transition-opacity group-hover:opacity-100 hover:text-red-600"
          @click="emit('delete', comment.id)"
        >
          Delete
        </button>
      </div>
      <VueMarkdown
        :source="comment.content"
        class="prose prose-sm prose-pre:bg-gray-50 prose-pre:border prose-pre:border-gray-100 max-w-none rounded-lg border border-gray-200 bg-white p-4 shadow-sm"
      />
    </div>
  </div>
</template>
