<script setup lang="ts">
import {
  Listbox,
  ListboxButton,
  ListboxOptions,
  ListboxOption,
  TransitionRoot,
} from "@headlessui/vue";
import { Check, ChevronsUpDown, User as UserIcon } from "lucide-vue-next";
import type { ProjectMember } from "~/types";

interface Props {
  members?: ProjectMember[];
  modelValue?: string | null;
  disabled?: boolean;
}

const {
  members = [],
  modelValue = null,
  disabled = false,
} = defineProps<Props>();
const emit = defineEmits(["update:modelValue", "change"]);

const selected = computed(
  () => members.find((m) => m.user_id === modelValue) || null
);

const handleChange = (value: string | null) => {
  emit("update:modelValue", value);
  emit("change", value);
};
</script>

<template>
  <Listbox
    :model-value="modelValue"
    :disabled="disabled"
    @update:model-value="handleChange"
  >
    <div class="relative mt-1">
      <ListboxButton
        class="border-border bg-background focus-visible:ring-ring relative w-full cursor-default rounded-md border py-2 pr-10 pl-3 text-left text-sm focus:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 sm:text-sm"
      >
        <span class="flex items-center gap-2 truncate">
          <UiBaseAvatar
            v-if="selected"
            :fallback="selected.email?.[0]?.toUpperCase() || '?'"
            class-name="h-5 w-5"
          />
          <UserIcon v-else class="text-muted-foreground/70 h-4 w-4" />
          <span class="block truncate">
            {{ selected ? selected.email : "Unassigned" }}
          </span>
        </span>
        <span
          class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2"
        >
          <ChevronsUpDown
            class="text-muted-foreground/70 h-4 w-4"
            aria-hidden="true"
          />
        </span>
      </ListboxButton>

      <TransitionRoot
        leave="transition ease-in duration-100"
        leave-from="opacity-100"
        leave-to="opacity-0"
      >
        <ListboxOptions
          class="bg-popover ring-border absolute z-50 mt-1 max-h-60 w-full overflow-auto rounded-md py-1 text-base shadow-lg ring-1 focus:outline-none sm:text-sm"
        >
          <ListboxOption
            v-slot="{ active, selected: isSelected }"
            :value="null"
            as="template"
          >
            <li
              :class="[
                active ? 'bg-accent text-accent-foreground' : 'text-foreground',
                'relative cursor-default py-2 pr-10 pl-10 select-none',
              ]"
            >
              <span class="flex items-center gap-2">
                <UserIcon class="text-muted-foreground/70 h-4 w-4" />
                <span
                  :class="[
                    !isSelected ? 'font-medium' : 'font-normal',
                    'block truncate',
                  ]"
                >
                  Unassigned
                </span>
              </span>
              <span
                v-if="!modelValue"
                class="text-primary absolute inset-y-0 left-0 flex items-center pl-3"
              >
                <Check class="h-4 w-4" aria-hidden="true" />
              </span>
            </li>
          </ListboxOption>

          <ListboxOption
            v-for="member in members"
            :key="member.user_id"
            v-slot="{ active, selected: isSelected }"
            :value="member.user_id"
            as="template"
          >
            <li
              :class="[
                active ? 'bg-accent text-accent-foreground' : 'text-foreground',
                'relative cursor-default py-2 pr-10 pl-10 select-none',
              ]"
            >
              <span class="flex items-center gap-2">
                <UiBaseAvatar
                  :fallback="member.email?.[0]?.toUpperCase() || '?'"
                  class-name="h-5 w-5"
                />
                <span
                  :class="[
                    isSelected ? 'font-medium' : 'font-normal',
                    'block truncate',
                  ]"
                >
                  {{ member.email }}
                </span>
              </span>
              <span
                v-if="isSelected"
                class="text-primary absolute inset-y-0 left-0 flex items-center pl-3"
              >
                <Check class="h-4 w-4" aria-hidden="true" />
              </span>
            </li>
          </ListboxOption>
        </ListboxOptions>
      </TransitionRoot>
    </div>
  </Listbox>
</template>
