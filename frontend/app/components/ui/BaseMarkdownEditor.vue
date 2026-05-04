<script setup lang="ts">
import VueMarkdown from "vue-markdown-render";

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
        'flex flex-col overflow-hidden rounded-md border border-gray-200 bg-white shadow-sm transition-all focus-within:border-blue-500 focus-within:ring-2 focus-within:ring-blue-500/20',
        className
      )
    "
  >
    <!-- Header Tabs -->
    <div class="flex items-center border-b border-gray-200 bg-white px-1">
      <button
        type="button"
        :class="
          cn(
            '-mb-px border-b-2 px-4 py-2 text-sm font-medium transition-all',
            tab === 'write'
              ? 'border-blue-600 text-blue-600'
              : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'
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
              ? 'border-blue-600 text-blue-600'
              : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'
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
          class="min-h-[200px] w-full resize-y rounded-none border-0 px-4 py-3 font-sans text-sm focus:ring-0 focus:outline-none"
          @input="onInput"
        />
      </div>
      <VueMarkdown
        v-else-if="modelValue"
        :source="modelValue"
        class="prose prose-sm min-h-[200px] max-w-none bg-gray-50/30 p-4"
      />
      <div v-else class="min-h-[200px] bg-gray-50/30 p-4">
        <span class="text-sm text-gray-400 italic">
          Nothing to preview
        </span>
      </div>
    </div>
  </div>
</template>
