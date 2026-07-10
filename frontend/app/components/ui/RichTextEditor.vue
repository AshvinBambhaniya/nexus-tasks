<template>
  <div
    class="rich-text-editor border-border bg-card focus-within:ring-primary relative rounded-xl border p-1 shadow-sm transition-all focus-within:ring-1"
  >
    <editor-content
      :editor="editor"
      class="custom-scrollbar max-h-[250px] min-h-[80px] overflow-y-auto px-3 py-2 text-sm"
    />
    <div
      class="border-border mt-2 flex items-center justify-between border-t px-2 pt-2 pb-1"
    >
      <div class="flex items-center gap-1">
        <button
          type="button"
          class="text-muted-foreground hover:bg-muted hover:text-foreground rounded-md p-1.5 transition-colors"
        >
          <Paperclip class="h-4 w-4" />
        </button>
        <button
          type="button"
          class="text-muted-foreground hover:bg-muted hover:text-foreground rounded-md p-1.5 px-2 text-xs font-bold transition-colors"
        >
          @
        </button>
      </div>
      <button
        type="button"
        class="bg-primary text-primary-foreground hover:bg-primary/90 flex h-8 items-center justify-center rounded-md px-4 text-xs font-medium transition-colors disabled:opacity-50"
        :disabled="isSubmitting || !hasContent"
        @click="handleSubmit"
      >
        <Loader2 v-if="isSubmitting" class="mr-1 h-3 w-3 animate-spin" />
        {{ isSubmitting ? "Posting..." : "Comment" }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onBeforeUnmount } from "vue";
import { useEditor, EditorContent, VueRenderer } from "@tiptap/vue-3";
import StarterKit from "@tiptap/starter-kit";
import Placeholder from "@tiptap/extension-placeholder";
import Mention from "@tiptap/extension-mention";
import tippy from "tippy.js";
import MentionList from "./MentionList.vue";
import { Paperclip, Loader2 } from "lucide-vue-next";
import { useProjectMembers } from "~/composables/useProjects";

const props = defineProps<{
  projectId: string;
  isSubmitting?: boolean;
}>();

const emit = defineEmits<{
  (e: "submit", content: string, mentionedUserIds: string[]): void;
}>();

const { members } = useProjectMembers(props.projectId);

// State
const contentHTML = ref("");
const hasContent = computed(
  () => contentHTML.value.trim().length > 0 && contentHTML.value !== "<p></p>"
);

// Mention Configuration
const suggestion = {
  items: ({ query }: { query: string }) => {
    if (!members.value) return [];

    return members.value
      .filter((member: { full_name?: string; email?: string }) => {
        const name = member.full_name?.toLowerCase() || "";
        const email = member.email?.toLowerCase() || "";
        const q = query.toLowerCase();
        return name.includes(q) || email.includes(q);
      })
      .slice(0, 5);
  },

  render: () => {
    let component: VueRenderer;
    let popup: Array<{
      hide: () => void;
      destroy: () => void;
      setProps: (p: unknown) => void;
    }>;

    return {
      onStart: (props: Record<string, unknown>) => {
        component = new VueRenderer(MentionList, {
          props,
          editor: props.editor,
        });

        if (!props.clientRect) {
          return;
        }

        popup = tippy("body", {
          getReferenceClientRect: props.clientRect,
          appendTo: () => document.body,
          content: component.element,
          showOnCreate: true,
          interactive: true,
          trigger: "manual",
          placement: "bottom-start",
        });
      },
      onUpdate(props: Record<string, unknown>) {
        component.updateProps(props);

        if (!props.clientRect) {
          return;
        }

        popup[0].setProps({
          getReferenceClientRect: props.clientRect,
        });
      },
      onKeyDown(props: { event: KeyboardEvent }) {
        if (props.event.key === "Escape") {
          popup[0].hide();
          return true;
        }
        return component.ref?.onKeyDown(props);
      },
      onExit() {
        popup[0].destroy();
        component.destroy();
      },
    };
  },
};

// Editor Initialization
const editor = useEditor({
  extensions: [
    StarterKit,
    Placeholder.configure({
      placeholder: "Ask a question or post an update... Type @ to mention.",
      emptyEditorClass: "is-editor-empty",
    }),
    Mention.configure({
      HTMLAttributes: {
        class:
          "mention-pill text-primary bg-primary/10 px-1 rounded-sm font-semibold inline-block cursor-pointer",
      },
      suggestion,
    }),
  ],
  onUpdate: ({ editor }) => {
    contentHTML.value = editor.getHTML();
  },
});

onBeforeUnmount(() => {
  if (editor.value) {
    editor.value.destroy();
  }
});

const handleSubmit = () => {
  if (!editor.value || !hasContent.value) return;

  // Extract mentioned users from JSON
  const json = editor.value.getJSON();
  const mentionedUserIds = new Set<string>();

  const traverse = (node: Record<string, unknown>) => {
    if (node.type === "mention" && node.attrs?.id) {
      mentionedUserIds.add(node.attrs.id);
    }
    if (node.content) {
      node.content.forEach(traverse);
    }
  };

  traverse(json);

  emit("submit", editor.value.getHTML(), Array.from(mentionedUserIds));

  // Clear editor on submit attempt (assuming success, parent handles failure)
  editor.value.commands.clearContent();
};
</script>

<style>
/* Tiptap standard styles */
.ProseMirror {
  outline: none;
}
.ProseMirror p {
  margin: 0;
}
.ProseMirror p.is-editor-empty:first-child::before {
  color: var(--muted-foreground);
  content: attr(data-placeholder);
  pointer-events: none;
  opacity: 0.7;
}
</style>
