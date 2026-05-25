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
import { cn } from "~/utils/cn";

interface Props {
  modelValue: TaskPriority;
}

const { modelValue } = defineProps<Props>();
const emit = defineEmits(["update:modelValue", "change"]);

const priorityConfig = {
  [TaskPriority.P0]: {
    label: "Critical",
    icon: AlertOctagon,
    color: "text-destructive",
    bg: "bg-destructive/10",
  },
  [TaskPriority.P1]: {
    label: "High",
    icon: SignalHigh,
    color: "text-orange-500",
    bg: "bg-orange-500/10",
  },
  [TaskPriority.P2]: {
    label: "Medium",
    icon: SignalMedium,
    color: "text-primary",
    bg: "bg-primary/10",
  },
  [TaskPriority.P3]: {
    label: "Low",
    icon: SignalLow,
    color: "text-muted-foreground/70",
    bg: "bg-muted",
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
        class="border-border bg-background ring-offset-background hover:bg-muted focus:ring-ring flex w-full items-center justify-between rounded-md border px-3 py-2 text-sm transition-colors focus:ring-2 focus:outline-none"
      >
        <div class="flex items-center gap-2">
          <component :is="current.icon" :class="cn('h-4 w-4', current.color)" />
          <span class="text-foreground/80 font-medium">{{
            current.label
          }}</span>
        </div>
        <ChevronDown class="text-muted-foreground/70 h-4 w-4" />
      </ListboxButton>

      <TransitionRoot
        leave="transition ease-in duration-100"
        leave-from="opacity-100"
        leave-to="opacity-0"
      >
        <ListboxOptions
          class="bg-popover ring-border absolute z-50 mt-1 max-h-60 w-full overflow-auto rounded-md p-1 text-base shadow-lg ring-1 focus:outline-none sm:text-sm"
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
                active ? 'bg-accent' : '',
                isSelected ? 'bg-muted' : '',
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
                    'text-foreground/80',
                  ]"
                >
                  {{ config.label }}
                </span>
              </div>
              <span v-if="isSelected" class="text-primary">
                <Check class="h-4 w-4" aria-hidden="true" />
              </span>
            </li>
          </ListboxOption>
        </ListboxOptions>
      </TransitionRoot>
    </div>
  </Listbox>
</template>
