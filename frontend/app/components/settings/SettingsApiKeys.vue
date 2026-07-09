<script setup lang="ts">
import {
  Key,
  Plus,
  Loader2,
  Trash2,
  Copy,
  Check,
  AlertTriangle,
  Clock,
  ShieldCheck,
} from "lucide-vue-next";
import { useApiKeys } from "~/composables/useApiKeys";
import type { CreateApiKeyResponse } from "~/types";

const { keys, isLoading, createKey, revokeKey } = useApiKeys();

// Generate key state
const newKeyName = ref("");
const isGenerating = ref(false);
const generatedResponse = ref<CreateApiKeyResponse | null>(null);
const hasCopied = ref(false);

// Revoke state
const revokingKeyId = ref<string | null>(null);

const handleGenerate = async () => {
  if (!newKeyName.value.trim()) return;

  isGenerating.value = true;
  try {
    const result = await createKey(newKeyName.value.trim());
    generatedResponse.value = result;
  } catch (err) {
    alert(getApiErrorMessage(err, "Failed to generate API key"));
  } finally {
    isGenerating.value = false;
  }
};

const handleCopy = async () => {
  if (!generatedResponse.value) return;

  try {
    await navigator.clipboard.writeText(generatedResponse.value.raw_token);
    hasCopied.value = true;
    setTimeout(() => (hasCopied.value = false), 3000);
  } catch {
    // Fallback for non-HTTPS environments
    const el = document.createElement("textarea");
    el.value = generatedResponse.value.raw_token;
    document.body.appendChild(el);
    el.select();
    document.execCommand("copy");
    document.body.removeChild(el);
    hasCopied.value = true;
    setTimeout(() => (hasCopied.value = false), 3000);
  }
};

const handleCloseDialog = () => {
  newKeyName.value = "";
  generatedResponse.value = null;
  hasCopied.value = false;
};

const handleRevoke = async (keyId: string) => {
  if (
    !confirm(
      "Are you sure you want to revoke this API key? Any integrations using it will stop working immediately."
    )
  )
    return;

  revokingKeyId.value = keyId;
  try {
    await revokeKey(keyId);
  } catch (err) {
    alert(getApiErrorMessage(err, "Failed to revoke API key"));
  } finally {
    revokingKeyId.value = null;
  }
};

const formatDate = (dateStr: string | null) => {
  if (!dateStr) return "Never";
  const date = new Date(dateStr);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMins = Math.floor(diffMs / 60000);
  const diffHours = Math.floor(diffMs / 3600000);
  const diffDays = Math.floor(diffMs / 86400000);

  if (diffMins < 1) return "Just now";
  if (diffMins < 60) return `${diffMins}m ago`;
  if (diffHours < 24) return `${diffHours}h ago`;
  if (diffDays < 30) return `${diffDays}d ago`;
  return date.toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
};
</script>

<template>
  <div class="space-y-12 pb-10">
    <!-- Create API Key Section -->
    <section
      class="flex flex-col gap-6 lg:flex-row lg:items-start lg:justify-between"
    >
      <div class="lg:w-1/3">
        <h2 class="text-foreground text-lg font-semibold tracking-tight">
          API Keys
        </h2>
        <p class="text-muted-foreground mt-1 text-sm">
          Create personal access tokens for MCP servers, CI/CD pipelines, and
          external integrations.
        </p>
      </div>

      <UiBaseCard class="border-border bg-card lg:w-2/3">
        <div class="p-6">
          <div class="flex flex-col gap-4 sm:flex-row sm:items-end">
            <div class="flex-1 space-y-2">
              <UiBaseLabel
                for="keyName"
                class="text-muted-foreground text-xs font-semibold tracking-wide uppercase"
                >Token Name</UiBaseLabel
              >
              <UiBaseInput
                id="keyName"
                v-model="newKeyName"
                placeholder="e.g., Cursor MCP, CI Pipeline"
                @keydown.enter="handleGenerate"
              />
            </div>
            <UiBaseButton
              :disabled="isGenerating || !newKeyName.trim()"
              class="min-w-[160px]"
              @click="handleGenerate"
            >
              <Loader2 v-if="isGenerating" class="mr-2 h-4 w-4 animate-spin" />
              <Plus v-else class="mr-2 h-4 w-4" />
              Generate Key
            </UiBaseButton>
          </div>
          <p class="text-muted-foreground/60 mt-3 text-xs">
            Keys authenticate as your user account. Treat them like passwords.
          </p>
        </div>
      </UiBaseCard>
    </section>

    <!-- Generated Token Display (inline, appears after generation) -->
    <Transition
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-2"
    >
      <section
        v-if="generatedResponse"
        class="flex flex-col gap-6 lg:flex-row lg:items-start lg:justify-between"
      >
        <div class="lg:w-1/3" />

        <div
          class="overflow-hidden rounded-xl border border-emerald-500/30 bg-emerald-500/5 lg:w-2/3"
        >
          <div class="p-6">
            <div class="mb-4 flex items-center gap-2">
              <ShieldCheck class="h-5 w-5 text-emerald-500" />
              <h3 class="font-semibold text-emerald-600 dark:text-emerald-400">
                Your API Key
              </h3>
            </div>

            <div
              class="bg-background/60 border-border flex items-center gap-2 rounded-lg border p-3"
            >
              <code
                class="text-foreground flex-1 font-mono text-sm break-all select-all"
              >
                {{ generatedResponse.raw_token }}
              </code>
              <button
                class="text-muted-foreground hover:text-foreground shrink-0 rounded-md p-2 transition-colors"
                title="Copy token"
                @click="handleCopy"
              >
                <Check v-if="hasCopied" class="h-4 w-4 text-emerald-500" />
                <Copy v-else class="h-4 w-4" />
              </button>
            </div>

            <div class="mt-3 flex items-start gap-2">
              <AlertTriangle class="mt-0.5 h-4 w-4 shrink-0 text-amber-500" />
              <p class="text-muted-foreground text-xs leading-relaxed">
                Make sure to copy your token now.
                <span class="text-foreground font-medium"
                  >You won't be able to see it again.</span
                >
              </p>
            </div>

            <div class="mt-4 flex justify-end">
              <UiBaseButton
                variant="outline"
                size="sm"
                @click="handleCloseDialog"
              >
                Done
              </UiBaseButton>
            </div>
          </div>
        </div>
      </section>
    </Transition>

    <div class="border-border border-t" />

    <!-- Active Keys Section -->
    <section
      class="flex flex-col gap-6 lg:flex-row lg:items-start lg:justify-between"
    >
      <div class="lg:w-1/3">
        <h2 class="text-foreground text-lg font-semibold tracking-tight">
          Active Keys
        </h2>
        <p class="text-muted-foreground mt-1 text-sm">
          Manage your existing tokens. Revoking a key is immediate and
          irreversible.
        </p>
      </div>

      <UiBaseCard class="border-border bg-card overflow-hidden lg:w-2/3">
        <!-- Loading State -->
        <div
          v-if="isLoading"
          class="flex flex-col items-center justify-center p-12"
        >
          <Loader2 class="text-muted-foreground/40 h-8 w-8 animate-spin" />
        </div>

        <!-- Keys List -->
        <div v-else-if="keys.length > 0" class="divide-border divide-y">
          <div
            v-for="key in keys"
            :key="key.id"
            class="hover:bg-muted/30 flex items-center justify-between p-4 transition-colors"
          >
            <div class="flex items-center gap-4">
              <div
                class="bg-muted flex h-10 w-10 shrink-0 items-center justify-center rounded-lg"
              >
                <Key class="text-muted-foreground h-4 w-4" />
              </div>
              <div>
                <div class="flex items-center gap-2">
                  <span class="text-card-foreground text-sm font-semibold">
                    {{ key.name }}
                  </span>
                </div>
                <div
                  class="text-muted-foreground mt-0.5 flex items-center gap-3 text-xs"
                >
                  <code
                    class="bg-muted rounded px-1.5 py-0.5 font-mono text-[11px]"
                  >
                    {{ key.token_prefix }}••••••••
                  </code>
                  <span class="flex items-center gap-1">
                    <Clock class="h-3 w-3" />
                    {{ formatDate(key.last_used_at) }}
                  </span>
                </div>
              </div>
            </div>

            <div class="flex items-center gap-3">
              <span class="text-muted-foreground/60 hidden text-xs sm:inline">
                Created {{ formatDate(key.created_at) }}
              </span>
              <UiBaseButton
                variant="ghost"
                size="sm"
                class="text-muted-foreground hover:bg-destructive/10 hover:text-destructive h-8 w-8 p-0"
                title="Revoke key"
                :disabled="revokingKeyId === key.id"
                @click="handleRevoke(key.id)"
              >
                <Loader2
                  v-if="revokingKeyId === key.id"
                  class="h-4 w-4 animate-spin"
                />
                <Trash2 v-else class="h-4 w-4" />
              </UiBaseButton>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <div
          v-else
          class="flex flex-col items-center justify-center p-12 text-center"
        >
          <div class="bg-muted mb-3 rounded-full p-3">
            <Key class="text-muted-foreground/50 h-6 w-6" />
          </div>
          <h3 class="text-card-foreground text-sm font-medium">
            No API keys yet
          </h3>
          <p class="text-muted-foreground mt-1 max-w-xs text-sm">
            Generate your first key to connect AI tools like Cursor or Claude
            Desktop to your Nexus Tasks.
          </p>
        </div>
      </UiBaseCard>
    </section>
  </div>
</template>
