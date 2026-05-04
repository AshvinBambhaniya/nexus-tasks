<script setup lang="ts">
import { Loader2, Check, Trash2 } from "lucide-vue-next";
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
  <div class="max-w-3xl space-y-8">
    <UiBaseCard>
      <div class="border-b border-gray-100 p-6">
        <h3 class="text-lg font-semibold text-gray-900">General Settings</h3>
        <p class="text-sm text-gray-500">Update your team details.</p>
      </div>
      <div class="p-6">
        <form class="space-y-4" @submit.prevent="handleSave">
          <div class="space-y-2">
            <UiBaseLabel for="name">Team Name</UiBaseLabel>
            <UiBaseInput id="name" v-model="name" required />
          </div>

          <div class="space-y-2">
            <UiBaseLabel for="description">Description</UiBaseLabel>
            <UiBaseTextArea
              id="description"
              v-model="description"
              placeholder="What is this team responsible for?"
              :rows="4"
            />
          </div>

          <div class="pt-2">
            <UiBaseButton type="submit" :disabled="isSaving">
              <template v-if="isSaving">
                <Loader2 class="mr-2 h-4 w-4 animate-spin" />
              </template>
              <template v-else-if="isSuccess">
                <Check class="mr-2 h-4 w-4" />
              </template>
              {{
                isSaving ? "Saving..." : isSuccess ? "Saved!" : "Save Changes"
              }}
            </UiBaseButton>
          </div>
        </form>
      </div>
    </UiBaseCard>

    <div class="rounded-lg border border-red-100 bg-red-50/30 p-6">
      <h3 class="text-lg font-medium text-red-900">Danger Zone</h3>
      <p class="text-sm text-red-700">Irreversible actions for this team.</p>
      <div class="mt-4 flex items-center justify-between">
        <div>
          <h4 class="font-medium text-gray-900">Delete Team</h4>
          <p class="text-sm text-gray-500">
            Permanently remove this team and its memberships.
          </p>
        </div>
        <UiBaseButton variant="destructive" size="sm" @click="handleDelete">
          <Trash2 class="mr-2 h-4 w-4" /> Delete Team
        </UiBaseButton>
      </div>
    </div>
  </div>
</template>
