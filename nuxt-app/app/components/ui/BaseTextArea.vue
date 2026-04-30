<script setup lang="ts">
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
        'flex min-h-[80px] w-full rounded-md border border-gray-200 bg-white px-3 py-2 text-sm text-gray-900 ring-offset-white placeholder:text-gray-500 focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2 focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-50',
        className
      )
    "
    @input="onInput"
  />
</template>
