<template>
  <div
    class="mention-list bg-popover text-popover-foreground flex flex-col overflow-hidden rounded-md border shadow-md"
  >
    <template v-if="items.length">
      <button
        v-for="(item, index) in items"
        :key="item.user_id"
        class="mention-item flex items-center gap-2 px-3 py-2 text-left text-sm transition-colors"
        :class="{
          'bg-muted': index === selectedIndex,
          'hover:bg-muted/50': index !== selectedIndex,
        }"
        @click="selectItem(index)"
      >
        <div
          class="bg-primary/10 text-primary flex h-6 w-6 shrink-0 items-center justify-center rounded-full text-[10px] font-bold"
        >
          {{ getInitial(item) }}
        </div>
        <span class="truncate">{{
          item.full_name || item.email || "Unknown User"
        }}</span>
      </button>
    </template>
    <div v-else class="text-muted-foreground px-3 py-2 text-center text-sm">
      No users found
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";

export interface MentionItem {
  user_id: string;
  full_name?: string;
  email?: string;
}

const props = defineProps<{
  items: MentionItem[];
  command: (item: { id: string; label: string }) => void;
}>();

const selectedIndex = ref(0);

watch(
  () => props.items,
  () => {
    selectedIndex.value = 0;
  }
);

const upHandler = () => {
  selectedIndex.value =
    (selectedIndex.value + props.items.length - 1) % props.items.length;
};

const downHandler = () => {
  selectedIndex.value = (selectedIndex.value + 1) % props.items.length;
};

const enterHandler = () => {
  selectItem(selectedIndex.value);
};

const selectItem = (index: number) => {
  const item = props.items[index];
  if (item) {
    props.command({
      id: item.user_id,
      label: item.full_name || item.email || "Unknown",
    });
  }
};

const getInitial = (item: MentionItem) => {
  const name = item.full_name || item.email || "U";
  return name.substring(0, 2).toUpperCase();
};

// Expose these methods to the VueRenderer so it can call them on key events
defineExpose({
  onKeyDown: (props: { event: KeyboardEvent }) => {
    if (props.event.key === "ArrowUp") {
      upHandler();
      return true;
    }
    if (props.event.key === "ArrowDown") {
      downHandler();
      return true;
    }
    if (props.event.key === "Enter") {
      enterHandler();
      return true;
    }
    return false;
  },
});
</script>

<style scoped>
.mention-list {
  min-width: 200px;
  max-height: 250px;
  overflow-y: auto;
}
</style>
