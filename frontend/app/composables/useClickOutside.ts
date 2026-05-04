import { onUnmounted, onMounted, type Ref } from "vue";

export const useClickOutside = (
  elRef: Ref<HTMLElement | null>,
  callback: () => void
) => {
  if (!elRef) return;

  const listener = (e: MouseEvent) => {
    if (elRef.value && !elRef.value.contains(e.target as Node)) {
      callback();
    }
  };

  onMounted(() => {
    document.addEventListener("mousedown", listener);
  });

  onUnmounted(() => {
    document.removeEventListener("mousedown", listener);
  });
};
