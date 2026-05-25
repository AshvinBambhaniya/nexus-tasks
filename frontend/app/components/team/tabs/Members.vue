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
    <div class="border-border border-b p-6">
      <h3 class="text-foreground text-lg font-semibold">Team Members</h3>
      <p class="text-muted-foreground text-sm">
        Manage the roster for this team.
      </p>
    </div>
    <div class="space-y-6 p-6">
      <form class="flex gap-3" @submit.prevent="handleInvite">
        <div class="relative flex-1">
          <Mail
            class="text-muted-foreground/70 absolute top-2.5 left-3 h-4 w-4"
          />
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
        <Loader2 class="text-muted-foreground/70 h-6 w-6 animate-spin" />
      </div>
      <div v-else class="divide-border divide-y">
        <div
          v-for="member in members"
          :key="member.user_id"
          class="hover:bg-muted flex items-center justify-between rounded-lg p-4 transition-colors"
        >
          <div class="flex items-center gap-3">
            <UiBaseAvatar
              :fallback="(member.email || '?')[0].toUpperCase()"
              class-name="h-9 w-9 bg-primary/10 text-xs text-primary"
            />
            <div>
              <div class="text-foreground text-sm font-medium">
                {{ member.email }}
              </div>
              <div class="text-muted-foreground text-xs capitalize">
                {{ (member.role || "").toLowerCase() }}
              </div>
            </div>
          </div>

          <div class="flex items-center gap-3">
            <UiBaseButton
              v-if="member.role !== 'ADMIN'"
              variant="ghost"
              size="sm"
              class-name="h-8 w-8 p-0 text-muted-foreground/70 hover:bg-destructive/10 hover:text-destructive"
              @click="handleRemoveMember(member.user_id)"
            >
              <Trash2 class="h-4 w-4" />
            </UiBaseButton>
          </div>
        </div>
        <div
          v-if="members.length === 0"
          class="text-muted-foreground py-6 text-center text-sm"
        >
          No members found.
        </div>
      </div>
    </div>
  </UiBaseCard>
</template>
