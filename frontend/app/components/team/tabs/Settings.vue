<script setup lang="ts">
import {
  Loader2,
  Check,
  Trash2,
  Settings2,
  AlertTriangle,
} from "lucide-vue-next";
import type { Team } from "~/types";

interface Props {
  team: Team;
}

const { team } = defineProps<Props>();

const router = useRouter();
const { updateTeam, deleteTeam } = useTeams();

const name = ref(team.name);
const description = ref(team.description || "");
const isSaving = ref(false);
const isSuccess = ref(false);

watchEffect(() => {
  if (team) {
    name.value = team.name;
    description.value = team.description || "";
  }
});

const handleSave = async () => {
  if (!name.value.trim()) return;
  isSaving.value = true;
  isSuccess.value = false;

  try {
    await updateTeam(team.id, {
      name: name.value,
      description: description.value,
    });
    isSuccess.value = true;
    setTimeout(() => (isSuccess.value = false), 3000);
  } catch {
    alert("Failed to update team");
  } finally {
    isSaving.value = false;
  }
};

const handleDelete = async () => {
  if (
    !confirm(
      "Are you sure you want to delete this team? This action cannot be undone."
    )
  )
    return;
  try {
    await deleteTeam(team.id);
    router.push("/teams");
  } catch {
    alert("Failed to delete team");
  }
};
</script>

<template>
  <div class="max-w-4xl space-y-8 pb-12">
    <!-- General Settings Card -->
    <div
      class="border-border bg-card overflow-hidden rounded-xl border shadow-sm"
    >
      <div
        class="border-border bg-muted/20 flex items-center gap-3 border-b p-6"
      >
        <div class="bg-primary/10 text-primary rounded-lg p-2">
          <Settings2 class="h-5 w-5" />
        </div>
        <div>
          <h3 class="text-foreground text-lg font-semibold tracking-tight">
            General Settings
          </h3>
          <p class="text-muted-foreground text-sm">
            Update your team's profile and description.
          </p>
        </div>
      </div>

      <div class="p-6 sm:p-8">
        <form class="space-y-6" @submit.prevent="handleSave">
          <div class="space-y-2">
            <label for="name" class="text-foreground text-sm font-semibold"
              >Team Name</label
            >
            <input
              id="name"
              v-model="name"
              required
              class="border-border bg-background focus:border-primary focus:ring-primary/20 text-foreground w-full rounded-lg border px-4 py-2.5 text-sm shadow-sm transition-all focus:ring-4 focus:outline-none"
            />
          </div>

          <div class="space-y-2">
            <label
              for="description"
              class="text-foreground text-sm font-semibold"
              >Description</label
            >
            <textarea
              id="description"
              v-model="description"
              placeholder="What is this team responsible for?"
              :rows="4"
              class="border-border bg-background focus:border-primary focus:ring-primary/20 text-foreground w-full resize-y rounded-lg border px-4 py-3 text-sm shadow-sm transition-all focus:ring-4 focus:outline-none"
            />
          </div>

          <div
            class="border-border flex items-center justify-between border-t pt-4"
          >
            <p class="text-muted-foreground text-sm">
              Please use 32 characters at maximum.
            </p>
            <button
              type="submit"
              :disabled="isSaving"
              class="bg-foreground text-background hover:bg-foreground/90 flex items-center justify-center gap-2 rounded-lg px-6 py-2.5 text-sm font-semibold transition-colors disabled:opacity-50"
            >
              <template v-if="isSaving">
                <Loader2 class="h-4 w-4 animate-spin" /> Saving...
              </template>
              <template v-else-if="isSuccess">
                <Check class="h-4 w-4 text-green-500" /> Saved!
              </template>
              <template v-else> Save Changes </template>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Danger Zone Card -->
    <div
      class="border-destructive/30 bg-destructive/5 overflow-hidden rounded-xl border shadow-sm"
    >
      <div
        class="border-destructive/20 bg-destructive/10 flex items-center gap-3 border-b p-6"
      >
        <div class="bg-destructive/20 text-destructive rounded-lg p-2">
          <AlertTriangle class="h-5 w-5" />
        </div>
        <div>
          <h3 class="text-destructive text-lg font-semibold tracking-tight">
            Danger Zone
          </h3>
          <p class="text-destructive/80 text-sm">
            Irreversible actions for this team.
          </p>
        </div>
      </div>

      <div class="p-6 sm:p-8">
        <div
          class="flex flex-col justify-between gap-4 sm:flex-row sm:items-center"
        >
          <div>
            <h4 class="text-foreground font-semibold">Delete Team</h4>
            <p class="text-muted-foreground mt-1 max-w-md text-sm">
              Permanently remove this team and its memberships. This action
              cannot be undone and will immediately revoke access for all
              members.
            </p>
          </div>
          <button
            class="bg-destructive text-destructive-foreground hover:bg-destructive/90 flex shrink-0 items-center justify-center gap-2 rounded-lg px-6 py-2.5 text-sm font-semibold shadow-sm transition-colors"
            @click="handleDelete"
          >
            <Trash2 class="h-4 w-4" /> Delete Team
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
