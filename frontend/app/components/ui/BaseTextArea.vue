<script setup lang="ts">
import { cn } from "~/utils/cn";

interface Props {
  modelValue?: string;
  placeholder?: string;
  disabled?: boolean;
  className?: string;
  rows?: number;
}

const {
  modelValue = "",
  placeholder = "",
  disabled = false,
  className = "",
  rows = 3,
} = defineProps<Props>();

const emit = defineEmits(["update:modelValue"]);

const onInput = (event: Event) => {
  const target = event.target as HTMLTextAreaElement;
  emit("update:modelValue", target.value);
};
</script>

<template>
  <textarea
    :value="modelValue"
    :placeholder="placeholder"
    :disabled="disabled"
    :rows="rows"
    :class="
      cn(
        'border-input bg-background text-foreground ring-offset-background placeholder:text-muted-foreground focus-visible:ring-ring flex min-h-[80px] w-full rounded-md border px-3 py-2 text-sm focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-50',
        className
      )
    "
    @input="onInput"
  />
</template>
