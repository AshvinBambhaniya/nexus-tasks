<script setup lang="ts">
import {
  Loader2,
  Trash2,
  UserPlus,
  Mail,
  ShieldAlert,
  Shield,
  Users,
} from "lucide-vue-next";

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
  <div class="space-y-6">
    <!-- Inline Command Bar / Invite Form -->
    <div
      class="border-border bg-card focus-within:ring-primary/20 rounded-xl border p-2 shadow-sm transition-all focus-within:ring-2"
    >
      <form class="flex items-center gap-3" @submit.prevent="handleInvite">
        <div class="relative flex flex-1 items-center">
          <Mail class="text-muted-foreground/50 ml-3 h-5 w-5" />
          <input
            v-model="inviteEmail"
            type="email"
            placeholder="Invite by email address..."
            required
            class="text-foreground placeholder:text-muted-foreground w-full bg-transparent py-2 pr-4 pl-3 text-sm focus:outline-none"
            :disabled="isInviting"
          />
        </div>
        <button
          type="submit"
          :disabled="isInviting || !inviteEmail"
          class="bg-primary text-primary-foreground hover:bg-primary/90 flex items-center justify-center gap-2 rounded-lg px-4 py-2 text-sm font-medium transition-colors disabled:opacity-50"
        >
          <Loader2 v-if="isInviting" class="h-4 w-4 animate-spin" />
          <UserPlus v-else class="h-4 w-4" />
          Invite
        </button>
      </form>
    </div>

    <!-- Roster List -->
    <div
      class="border-border bg-card overflow-hidden rounded-xl border shadow-sm"
    >
      <div v-if="isLoading" class="flex justify-center p-12">
        <Loader2 class="text-muted-foreground/70 h-8 w-8 animate-spin" />
      </div>

      <div v-else class="divide-border divide-y">
        <div
          v-for="member in members"
          :key="member.user_id"
          class="group hover:bg-muted/50 flex items-center justify-between p-4 transition-colors"
        >
          <div class="flex items-center gap-4">
            <div
              class="bg-primary/10 text-primary flex h-10 w-10 shrink-0 items-center justify-center rounded-full font-bold"
            >
              {{ (member.email || "?")[0].toUpperCase() }}
            </div>
            <div>
              <div class="text-foreground text-sm font-semibold tracking-tight">
                {{ member.email }}
              </div>
              <div
                class="text-muted-foreground mt-0.5 flex items-center gap-1.5 text-xs capitalize"
              >
                <ShieldAlert
                  v-if="member.role === 'ADMIN'"
                  class="h-3 w-3 text-red-500"
                />
                <Shield v-else class="h-3 w-3" />
                {{ (member.role || "Member").toLowerCase() }}
              </div>
            </div>
          </div>

          <div class="flex items-center gap-3 pr-2">
            <button
              v-if="member.role !== 'ADMIN'"
              class="text-muted-foreground/40 hover:bg-destructive/10 hover:text-destructive flex h-8 w-8 items-center justify-center rounded-md opacity-0 transition-all group-hover:opacity-100 focus:opacity-100"
              title="Remove member"
              @click="handleRemoveMember(member.user_id)"
            >
              <Trash2 class="h-4 w-4" />
            </button>
            <span v-else class="text-muted-foreground/30 pr-2 text-xs italic"
              >Owner</span
            >
          </div>
        </div>

        <div
          v-if="members.length === 0"
          class="text-muted-foreground flex flex-col items-center gap-3 py-16 text-center"
        >
          <Users class="h-8 w-8 opacity-20" />
          <p class="text-sm">No members found.</p>
        </div>
      </div>
    </div>
  </div>
</template>
