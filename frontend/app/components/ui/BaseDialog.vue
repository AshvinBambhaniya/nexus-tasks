<script setup lang="ts">
import {
  TransitionRoot,
  TransitionChild,
  Dialog,
  DialogPanel,
  DialogTitle,
} from "@headlessui/vue";
import { X } from "lucide-vue-next";
import { cn } from "~/utils/cn";

interface Props {
  isOpen?: boolean;
  title: string;
  description?: string;
  className?: string;
}

const {
  isOpen = false,
  title,
  description = "",
  className = "",
} = defineProps<Props>();

const emit = defineEmits(["close"]);

const onClose = () => {
  emit("close");
};
</script>

<template>
  <TransitionRoot :show="isOpen" as="template" appear>
    <Dialog as="div" class="relative z-50" @close="onClose">
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

      <div class="fixed inset-0 z-10 overflow-y-auto">
        <div
          class="flex min-h-full items-center justify-center p-4 text-center sm:p-0"
        >
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
              :class="
                cn(
                  'bg-card relative w-full max-w-lg transform overflow-hidden rounded-lg text-left shadow-xl transition-all',
                  className
                )
              "
            >
              <div
                class="border-border flex items-center justify-between border-b px-6 py-4"
              >
                <div class="space-y-1">
                  <DialogTitle
                    as="h2"
                    class="text-card-foreground text-lg font-semibold"
                  >
                    {{ title }}
                  </DialogTitle>
                  <p v-if="description" class="text-muted-foreground text-sm">
                    {{ description }}
                  </p>
                </div>
                <UiBaseButton
                  variant="ghost"
                  size="sm"
                  class="h-8 w-8 p-0"
                  @click="onClose"
                >
                  <X class="h-4 w-4" />
                  <span class="sr-only">Close</span>
                </UiBaseButton>
              </div>
              <div class="p-6">
                <slot />
              </div>
            </DialogPanel>
          </TransitionChild>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
</template>
