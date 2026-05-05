<script setup lang="ts">
import {
  Listbox,
  ListboxButton,
  ListboxOptions,
  ListboxOption,
  TransitionRoot,
} from "@headlessui/vue";
import {
  Check,
  ChevronDown,
  SignalHigh,
  SignalMedium,
  SignalLow,
  AlertOctagon,
} from "lucide-vue-next";
import { TaskPriority } from "~/types";

interface Props {
  modelValue: TaskPriority;
}

const { modelValue } = defineProps<Props>();
const emit = defineEmits(["update:modelValue", "change"]);

const priorityConfig = {
  [TaskPriority.P0]: {
    label: "Critical",
    icon: AlertOctagon,
    color: "text-red-600",
    bg: "bg-red-50",
  },
  [TaskPriority.P1]: {
    label: "High",
    icon: SignalHigh,
    color: "text-orange-500",
    bg: "bg-orange-50",
  },
  [TaskPriority.P2]: {
    label: "Medium",
    icon: SignalMedium,
    color: "text-blue-500",
    bg: "bg-blue-50",
  },
  [TaskPriority.P3]: {
    label: "Low",
    icon: SignalLow,
    color: "text-gray-400",
    bg: "bg-gray-50",
  },
};

const current = computed(() => priorityConfig[modelValue]);

const handleChange = (value: TaskPriority) => {
  emit("update:modelValue", value);
  emit("change", value);
};
</script>

<template>
  <Listbox :model-value="modelValue" @update:model-value="handleChange">
    <div class="relative mt-1">
      <ListboxButton
        class="flex w-full items-center justify-between rounded-md border border-gray-200 bg-white px-3 py-2 text-sm ring-offset-white transition-colors hover:bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none"
      >
        <div class="flex items-center gap-2">
          <component :is="current.icon" :class="cn('h-4 w-4', current.color)" />
          <span class="font-medium text-gray-700">{{ current.label }}</span>
        </div>
        <ChevronDown class="h-4 w-4 text-gray-400" />
      </ListboxButton>

      <TransitionRoot
        leave="transition ease-in duration-100"
        leave-from="opacity-100"
        leave-to="opacity-0"
      >
        <ListboxOptions
          class="absolute z-50 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white p-1 text-base shadow-lg ring-1 ring-black/5 focus:outline-none sm:text-sm"
        >
          <ListboxOption
            v-for="(config, key) in priorityConfig"
            :key="key"
            v-slot="{ active, selected: isSelected }"
            :value="key"
            as="template"
          >
            <li
              :class="[
                active ? 'bg-gray-100' : '',
                isSelected ? 'bg-gray-50' : '',
                'relative flex cursor-default items-center justify-between rounded-sm px-2 py-1.5 transition-colors select-none',
              ]"
            >
              <div class="flex items-center gap-2">
                <component
                  :is="config.icon"
                  :class="cn('h-4 w-4', config.color)"
                />
                <span
                  :class="[
                    isSelected ? 'font-medium' : 'font-normal',
                    'text-gray-700',
                  ]"
                >
                  {{ config.label }}
                </span>
              </div>
              <span v-if="isSelected" class="text-blue-600">
                <Check class="h-4 w-4" aria-hidden="true" />
              </span>
            </li>
          </ListboxOption>
        </ListboxOptions>
      </TransitionRoot>
    </div>
  </Listbox>
</template>
