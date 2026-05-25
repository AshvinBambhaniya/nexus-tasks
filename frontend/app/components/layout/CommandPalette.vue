<script setup lang="ts">
import {
  TransitionRoot,
  TransitionChild,
  Dialog,
  DialogPanel,
  Combobox,
  ComboboxInput,
  ComboboxOptions,
  ComboboxOption,
} from "@headlessui/vue";
import {
  Search,
  Inbox,
  CheckSquare,
  LayoutDashboard,
  Building2,
} from "lucide-vue-next";
import type { Component } from "vue";
import type { Workspace } from "~/types";

const isOpen = ref(false);
const query = ref("");
const router = useRouter();
const workspaceStore = useWorkspaceStore();
const { workspaces } = useWorkspaces();

// Toggle the menu when ⌘K is pressed
const onKeyDown = (e: KeyboardEvent) => {
  if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
    e.preventDefault();
    isOpen.value = !isOpen.value;
  }
};

onMounted(() => {
  window.addEventListener("keydown", onKeyDown);
});

onUnmounted(() => {
  window.removeEventListener("keydown", onKeyDown);
});

interface NavItem {
  name: string;
  href: string;
  icon: Component;
}

const navigation: NavItem[] = [
  { name: "Go to Inbox", href: "/inbox", icon: Inbox },
  { name: "All Tasks", href: "/tasks", icon: CheckSquare },
  { name: "Dashboard", href: "/dashboard", icon: LayoutDashboard },
];

const filteredNavigation = computed(() =>
  query.value === ""
    ? navigation
    : navigation.filter((item) =>
        item.name.toLowerCase().includes(query.value.toLowerCase())
      )
);

const filteredWorkspaces = computed(() =>
  query.value === ""
    ? workspaces.value
    : workspaces.value.filter((w) =>
        w.name.toLowerCase().includes(query.value.toLowerCase())
      )
);

const runCommand = (action: () => void) => {
  isOpen.value = false;
  action();
};

const handleSelect = (item: NavItem | Workspace) => {
  if ("href" in item) {
    runCommand(() => router.push(item.href));
  } else if ("id" in item) {
    // Workspace selection
    runCommand(() => {
      workspaceStore.setActiveWorkspaceId(item.id);
      router.push("/dashboard");
    });
  }
};
</script>

<template>
  <TransitionRoot :show="isOpen" as="template" appear>
    <Dialog as="div" class="relative z-50" @close="isOpen = false">
      <TransitionChild
        as="template"
        enter="duration-300 ease-out"
        enter-from="opacity-0"
        enter-to="opacity-100"
        leave="duration-200 ease-in"
        leave-from="opacity-100"
        leave-to="opacity-0"
      >
        <div
          class="fixed inset-0 bg-black/50 backdrop-blur-sm transition-opacity"
        />
      </TransitionChild>

      <div class="fixed inset-0 z-10 overflow-y-auto p-4 sm:p-6 md:p-20">
        <TransitionChild
          as="template"
          enter="duration-300 ease-out"
          enter-from="opacity-0 scale-95"
          enter-to="opacity-100 scale-100"
          leave="duration-200 ease-in"
          leave-from="opacity-100 scale-100"
          leave-to="opacity-0 scale-95"
        >
          <DialogPanel
            class="divide-border bg-card ring-border mx-auto max-w-2xl transform divide-y overflow-hidden rounded-xl shadow-2xl ring-1 transition-all"
          >
            <Combobox @update:model-value="handleSelect">
              <div class="relative">
                <Search
                  class="text-muted-foreground pointer-events-none absolute top-3.5 left-4 h-5 w-5"
                  aria-hidden="true"
                />
                <ComboboxInput
                  class="text-card-foreground placeholder:text-muted-foreground h-12 w-full border-0 bg-transparent pr-4 pl-11 outline-none focus:ring-0 sm:text-sm"
                  placeholder="Search..."
                  @change="query = $event.target.value"
                />
              </div>

              <ComboboxOptions
                v-if="
                  filteredNavigation.length > 0 || filteredWorkspaces.length > 0
                "
                static
                class="divide-border max-h-80 scroll-py-2 divide-y overflow-y-auto"
              >
                <li class="p-2">
                  <h2
                    v-if="filteredNavigation.length > 0"
                    class="text-muted-foreground mt-4 mb-2 px-3 text-xs font-semibold tracking-wider uppercase"
                  >
                    Navigation
                  </h2>
                  <ComboboxOption
                    v-for="item in filteredNavigation"
                    :key="item.href"
                    v-slot="{ active }"
                    :value="item"
                    as="template"
                  >
                    <div
                      :class="[
                        'flex cursor-default items-center rounded-md px-3 py-2 transition-colors select-none',
                        active
                          ? 'bg-primary text-primary-foreground'
                          : 'text-card-foreground',
                      ]"
                    >
                      <component
                        :is="item.icon"
                        :class="[
                          'h-5 w-5 flex-none',
                          active
                            ? 'text-primary-foreground'
                            : 'text-muted-foreground',
                        ]"
                        aria-hidden="true"
                      />
                      <span class="ml-3 flex-auto truncate">{{
                        item.name
                      }}</span>
                    </div>
                  </ComboboxOption>
                </li>

                <li v-if="filteredWorkspaces.length > 0" class="p-2">
                  <h2
                    class="text-muted-foreground mt-4 mb-2 px-3 text-xs font-semibold tracking-wider uppercase"
                  >
                    Workspaces
                  </h2>
                  <ComboboxOption
                    v-for="workspace in filteredWorkspaces"
                    :key="workspace.id"
                    v-slot="{ active }"
                    :value="workspace"
                    as="template"
                  >
                    <div
                      :class="[
                        'flex cursor-default items-center rounded-md px-3 py-2 transition-colors select-none',
                        active
                          ? 'bg-primary text-primary-foreground'
                          : 'text-card-foreground',
                      ]"
                    >
                      <Building2
                        :class="[
                          'h-5 w-5 flex-none',
                          active
                            ? 'text-primary-foreground'
                            : 'text-muted-foreground',
                        ]"
                        aria-hidden="true"
                      />
                      <span class="ml-3 flex-auto truncate">{{
                        workspace.name
                      }}</span>
                    </div>
                  </ComboboxOption>
                </li>
              </ComboboxOptions>

              <div
                v-if="
                  query !== '' &&
                  filteredNavigation.length === 0 &&
                  filteredWorkspaces.length === 0
                "
                class="px-6 py-14 text-center sm:px-14"
              >
                <Search
                  class="text-muted-foreground mx-auto h-6 w-6"
                  aria-hidden="true"
                />
                <p class="text-card-foreground mt-4 text-sm">
                  No results found for "{{ query }}".
                </p>
              </div>
            </Combobox>
          </DialogPanel>
        </TransitionChild>
      </div>
    </Dialog>
  </TransitionRoot>
</template>
