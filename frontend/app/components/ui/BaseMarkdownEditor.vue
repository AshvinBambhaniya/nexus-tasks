<script setup lang="ts">
import VueMarkdown from "vue-markdown-render";
import { cn } from "~/utils/cn";

interface Props {
  modelValue: string;
  placeholder?: string;
  className?: string;
  disabled?: boolean;
}

const {
  modelValue,
  placeholder = "Write something...",
  className = "",
  disabled = false,
} = defineProps<Props>();

const emit = defineEmits(["update:modelValue"]);

const tab = ref<"write" | "preview">("write");

const onInput = (event: Event) => {
  const target = event.target as HTMLTextAreaElement;
  emit("update:modelValue", target.value);
};
</script>

<template>
  <div
    :class="
      cn(
        'border-border bg-card focus-within:border-ring focus-within:ring-ring/20 flex flex-col overflow-hidden rounded-md border shadow-sm transition-all focus-within:ring-2',
        className
      )
    "
  >
    <!-- Header Tabs -->
    <div class="border-border bg-card flex items-center border-b px-1">
      <button
        type="button"
        :class="
          cn(
            '-mb-px border-b-2 px-4 py-2 text-sm font-medium transition-all',
            tab === 'write'
              ? 'border-primary text-primary'
              : 'text-muted-foreground hover:border-border hover:text-foreground border-transparent'
          )
        "
        @click="tab = 'write'"
      >
        Write
      </button>
      <button
        type="button"
        :class="
          cn(
            '-mb-px border-b-2 px-4 py-2 text-sm font-medium transition-all',
            tab === 'preview'
              ? 'border-primary text-primary'
              : 'text-muted-foreground hover:border-border hover:text-foreground border-transparent'
          )
        "
        @click="tab = 'preview'"
      >
        Preview
      </button>
    </div>

    <!-- Editor Content -->
    <div class="relative">
      <div v-if="tab === 'write'">
        <textarea
          :value="modelValue"
          :placeholder="placeholder"
          :disabled="disabled"
          class="text-foreground min-h-[200px] w-full resize-y rounded-none border-0 bg-transparent px-4 py-3 font-sans text-sm focus:ring-0 focus:outline-none"
          @input="onInput"
        />
      </div>
      <VueMarkdown
        v-else-if="modelValue"
        :source="modelValue"
        class="prose prose-sm bg-muted/30 min-h-[200px] max-w-none p-4"
      />
      <div v-else class="bg-muted/30 min-h-[200px] p-4">
        <span class="text-muted-foreground text-sm italic">
          Nothing to preview
        </span>
      </div>
    </div>
  </div>
</template>
