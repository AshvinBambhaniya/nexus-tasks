<template>
  <div class="mention-list bg-popover text-popover-foreground flex flex-col overflow-hidden rounded-md border shadow-md">
    <template v-if="items.length">
      <button
        v-for="(item, index) in items"
        :key="item.user_id"
        class="mention-item flex items-center gap-2 px-3 py-2 text-sm text-left transition-colors"
        :class="{ 'bg-muted': index === selectedIndex, 'hover:bg-muted/50': index !== selectedIndex }"
        @click="selectItem(index)"
      >
        <div class="bg-primary/10 text-primary flex h-6 w-6 shrink-0 items-center justify-center rounded-full text-[10px] font-bold">
          {{ getInitial(item) }}
        </div>
        <span class="truncate">{{ item.full_name || item.email || "Unknown User" }}</span>
      </button>
    </template>
    <div v-else class="px-3 py-2 text-sm text-muted-foreground text-center">
      No users found
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  items: any[]
  command: (item: any) => void
}>()

const selectedIndex = ref(0)

watch(() => props.items, () => {
  selectedIndex.value = 0
})

const upHandler = () => {
  selectedIndex.value = ((selectedIndex.value + props.items.length) - 1) % props.items.length
}

const downHandler = () => {
  selectedIndex.value = (selectedIndex.value + 1) % props.items.length
}

const enterHandler = () => {
  selectItem(selectedIndex.value)
}

const selectItem = (index: number) => {
  const item = props.items[index]
  if (item) {
    props.command({ id: item.user_id, label: item.full_name || item.email || "Unknown" })
  }
}

const getInitial = (item: any) => {
  const name = item.full_name || item.email || "U"
  return name.substring(0, 2).toUpperCase()
}

// Expose these methods to the VueRenderer so it can call them on key events
defineExpose({
  onKeyDown: (props: any) => {
    if (props.event.key === 'ArrowUp') {
      upHandler()
      return true
    }
    if (props.event.key === 'ArrowDown') {
      downHandler()
      return true
    }
    if (props.event.key === 'Enter') {
      enterHandler()
      return true
    }
    return false
  }
})
</script>

<style scoped>
.mention-list {
  min-width: 200px;
  max-height: 250px;
  overflow-y: auto;
}
</style>
