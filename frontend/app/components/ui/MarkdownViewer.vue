<script setup lang="ts">
import { computed } from "vue";
import { marked } from "marked";
import DOMPurify from "isomorphic-dompurify";

const props = defineProps<{
  content: string;
}>();

// Compute the safe HTML string whenever the content prop changes
const safeHtml = computed(() => {
  if (!props.content) return "";

  // Convert Markdown string to raw HTML
  const rawHtml = marked.parse(props.content, { async: false }) as string;

  // Sanitize the HTML to remove any potential XSS vectors
  return DOMPurify.sanitize(rawHtml, {
    ALLOWED_TAGS: [
      "b",
      "i",
      "em",
      "strong",
      "a",
      "p",
      "h1",
      "h2",
      "h3",
      "h4",
      "h5",
      "h6",
      "ul",
      "ol",
      "li",
      "code",
      "pre",
      "blockquote",
      "br",
      "hr",
      "del",
      "table",
      "thead",
      "tbody",
      "tr",
      "th",
      "td",
      "img",
      "span",
    ],
    ALLOWED_ATTR: [
      "href",
      "title",
      "target",
      "class",
      "src",
      "alt",
      "data-id",
      "data-is-me",
    ],
  });
});
</script>

<template>
  <!-- eslint-disable-next-line vue/no-v-html -->
  <div class="markdown-body" v-html="safeHtml" />
</template>

<style scoped>
@reference "../../assets/css/main.css";

/* Basic typography scoping so markdown doesn't break global CSS */
.markdown-body :deep(p) {
  margin-bottom: 0.75rem;
  line-height: 1.5;
}
.markdown-body :deep(h1),
.markdown-body :deep(h2),
.markdown-body :deep(h3),
.markdown-body :deep(h4),
.markdown-body :deep(h5),
.markdown-body :deep(h6) {
  margin-top: 1.5rem;
  margin-bottom: 0.75rem;
  font-weight: 600;
  line-height: 1.25;
}
.markdown-body :deep(h1) {
  font-size: 1.5rem;
}
.markdown-body :deep(h2) {
  font-size: 1.25rem;
}
.markdown-body :deep(h3) {
  font-size: 1.125rem;
}

.markdown-body :deep(ul),
.markdown-body :deep(ol) {
  margin-top: 0;
  margin-bottom: 0.75rem;
  padding-left: 2rem;
}

.markdown-body :deep(li) {
  margin-top: 0.25rem;
}

.markdown-body :deep(pre) {
  @apply bg-muted text-foreground mb-3 overflow-x-auto rounded-md p-4;
}

.markdown-body :deep(code) {
  @apply bg-muted text-foreground rounded px-1.5 py-0.5 font-mono text-[0.875em];
}

.markdown-body :deep(pre code) {
  @apply bg-transparent p-0 text-inherit;
}

.markdown-body :deep(blockquote) {
  @apply border-border text-muted-foreground my-3 border-l-4 pl-4;
}

.markdown-body :deep(a) {
  @apply text-primary hover:text-primary/80 underline transition-colors;
}

.markdown-body :deep(hr) {
  @apply border-border my-6 border-t;
}
</style>
