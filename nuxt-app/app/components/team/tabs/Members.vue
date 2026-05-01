<script setup lang="ts">
import { Loader2, Trash2, UserPlus, Mail } from "lucide-vue-next";

interface Props {
  teamId: string;
}

const { teamId } = defineProps<Props>();

const { members, isLoading, addMember, removeMember } = useTeamMembers(teamId);
const inviteEmail = ref("");
const isInviting = ref(false);

const handleInvite = async () => {
  if (!inviteEmail.value) return;

  isInviting.value = true;
  try {
    await addMember(inviteEmail.value);
    inviteEmail.value = "";
  } catch (err) {
    alert(getApiErrorMessage(err, "Failed to add member"));
  } finally {
    isInviting.value = false;
  }
};

const handleRemoveMember = async (userId: string) => {
  if (!confirm("Remove this member from the team?")) return;
  try {
    await removeMember(userId);
  } catch (err) {
    alert(getApiErrorMessage(err, "Failed to remove member"));
  }
};
</script>

<template>
  <UiBaseCard>
    <div class="border-b border-gray-100 p-6">
      <h3 class="text-lg font-semibold text-gray-900">Team Members</h3>
      <p class="text-sm text-gray-500">Manage the roster for this team.</p>
    </div>
    <div class="space-y-6 p-6">
      <form class="flex gap-3" @submit.prevent="handleInvite">
        <div class="relative flex-1">
          <Mail class="absolute top-2.5 left-3 h-4 w-4 text-gray-400" />
          <UiBaseInput
            v-model="inviteEmail"
            type="email"
            placeholder="name@example.com"
            required
            class-name="pl-9"
            :disabled="isInviting"
          />
        </div>
        <UiBaseButton type="submit" :disabled="isInviting">
          <template v-if="isInviting">
            <Loader2 class="mr-2 h-4 w-4 animate-spin" />
          </template>
          <template v-else>
            <UserPlus class="mr-2 h-4 w-4" />
          </template>
          Add
        </UiBaseButton>
      </form>

      <div v-if="isLoading" class="flex justify-center p-8">
        <Loader2 class="h-6 w-6 animate-spin text-gray-400" />
      </div>
      <div v-else class="divide-y divide-gray-100">
        <div
          v-for="member in members"
          :key="member.user_id"
          class="flex items-center justify-between rounded-lg p-4 transition-colors hover:bg-gray-50"
        >
          <div class="flex items-center gap-3">
            <UiBaseAvatar
              :fallback="(member.email || '?')[0].toUpperCase()"
              class-name="h-9 w-9 bg-indigo-100 text-xs text-indigo-700"
            />
            <div>
              <div class="text-sm font-medium text-gray-900">
                {{ member.email }}
              </div>
              <div class="text-xs text-gray-500 capitalize">
                {{ (member.role || "").toLowerCase() }}
              </div>
            </div>
          </div>

          <div class="flex items-center gap-3">
            <UiBaseButton
              v-if="member.role !== 'ADMIN'"
              variant="ghost"
              size="sm"
              class-name="h-8 w-8 p-0 text-gray-400 hover:bg-red-50 hover:text-red-600"
              @click="handleRemoveMember(member.user_id)"
            >
              <Trash2 class="h-4 w-4" />
            </UiBaseButton>
          </div>
        </div>
        <div
          v-if="members.length === 0"
          class="py-6 text-center text-sm text-gray-500"
        >
          No members found.
        </div>
      </div>
    </div>
  </UiBaseCard>
</template>
