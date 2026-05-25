<script setup lang="ts">
import {
  Trash2,
  UserPlus,
  Users,
  Link as LinkIcon,
  Mail,
  Loader2,
} from "lucide-vue-next";

interface Props {
  projectId: string;
}

const { projectId } = defineProps<Props>();

const { members, isLoading, addMember, removeMember } =
  useProjectMembers(projectId);
const { teams: projectTeams, addTeam, removeTeam } = useProjectTeams(projectId);
const { teams: workspaceTeams } = useTeams();

const inviteEmail = ref("");
const isInviting = ref(false);
const selectedTeamId = ref("");
const isLinking = ref(false);

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
  if (!confirm("Remove this member from the project?")) return;
  try {
    await removeMember(userId);
  } catch (err) {
    alert(getApiErrorMessage(err, "Failed to remove member"));
  }
};

const handleLinkTeam = async () => {
  if (!selectedTeamId.value) return;
  isLinking.value = true;
  try {
    await addTeam(selectedTeamId.value);
    selectedTeamId.value = "";
  } catch (err) {
    alert(getApiErrorMessage(err, "Failed to link team"));
  } finally {
    isLinking.value = false;
  }
};

const handleUnlinkTeam = async (teamId: string) => {
  if (!confirm("Unlink this team from the project?")) return;
  try {
    await removeTeam(teamId);
  } catch (err) {
    alert(getApiErrorMessage(err, "Failed to unlink team"));
  }
};

const availableTeams = computed(() =>
  workspaceTeams.value.filter(
    (wt) => !projectTeams.value.some((pt) => pt.team_id === wt.id)
  )
);
</script>

<template>
  <div class="space-y-8">
    <UiBaseCard>
      <div class="border-border border-b p-6">
        <h3 class="text-card-foreground text-lg font-semibold">
          Project Members
        </h3>
        <p class="text-muted-foreground text-sm">
          All users with access to this project, including team members.
        </p>
      </div>
      <div class="space-y-6 p-6">
        <form class="flex gap-3" @submit.prevent="handleInvite">
          <div class="relative flex-1">
            <Mail
              class="text-muted-foreground/60 absolute top-2.5 left-3 h-4 w-4"
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
          <Loader2 class="text-muted-foreground/60 h-6 w-6 animate-spin" />
        </div>
        <div v-else class="divide-border divide-y">
          <div
            v-for="member in members"
            :key="member.user_id"
            class="hover:bg-muted flex items-center justify-between rounded-lg p-4 transition-colors"
          >
            <div class="flex items-center gap-3">
              <UiBaseAvatar
                :fallback="member.email[0].toUpperCase()"
                class-name="h-9 w-9 bg-primary/10 text-xs text-primary"
              />
              <div>
                <div class="text-card-foreground text-sm font-medium">
                  {{ member.email }}
                </div>
                <div
                  class="text-muted-foreground flex items-center gap-2 text-xs"
                >
                  {{ member.is_direct ? "Direct Member" : "via Team" }}
                </div>
              </div>
            </div>

            <div class="flex items-center gap-3">
              <UiBaseBadge
                :variant="member.role === 'ADMIN' ? 'default' : 'secondary'"
                class-name="h-5 text-[10px]"
              >
                {{ member.role }}
              </UiBaseBadge>
              <UiBaseButton
                v-if="member.is_direct"
                variant="ghost"
                size="sm"
                class-name="h-8 w-8 p-0 text-muted-foreground/60 hover:bg-destructive/10 hover:text-destructive"
                title="Remove direct access"
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

    <UiBaseCard>
      <div class="border-border border-b p-6">
        <h3 class="text-card-foreground text-lg font-semibold">Teams</h3>
        <p class="text-muted-foreground text-sm">
          Workspace teams with access to this project.
        </p>
      </div>
      <div class="space-y-6 p-6">
        <form class="flex gap-3" @submit.prevent="handleLinkTeam">
          <div class="relative flex-1">
            <Users
              class="text-muted-foreground/60 absolute top-2.5 left-3 h-4 w-4"
            />
            <select
              v-model="selectedTeamId"
              class="border-border bg-background ring-offset-background focus-visible:ring-primary flex h-10 w-full rounded-md border px-3 py-2 pl-9 text-sm focus-visible:ring-2 focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-50"
              required
            >
              <option value="" disabled>Select a team to link...</option>
              <option v-for="t in availableTeams" :key="t.id" :value="t.id">
                {{ t.name }}
              </option>
            </select>
          </div>
          <UiBaseButton
            type="submit"
            :disabled="isLinking || availableTeams.length === 0"
          >
            <template v-if="isLinking">
              <Loader2 class="mr-2 h-4 w-4 animate-spin" />
            </template>
            <template v-else>
              <LinkIcon class="mr-2 h-4 w-4" />
            </template>
            Link Team
          </UiBaseButton>
        </form>

        <div v-if="isLoading" class="flex justify-center p-8">
          <Loader2 class="text-muted-foreground/60 h-6 w-6 animate-spin" />
        </div>
        <div v-else class="divide-border divide-y">
          <div
            v-for="team in projectTeams"
            :key="team.team_id"
            class="hover:bg-muted flex items-center justify-between rounded-lg p-4 transition-colors"
          >
            <div class="flex items-center gap-3">
              <div
                class="flex h-9 w-9 items-center justify-center rounded-lg bg-indigo-500/10 text-indigo-500"
              >
                <Users class="h-4 w-4" />
              </div>
              <div class="text-card-foreground text-sm font-medium">
                {{ team.team_name }}
              </div>
            </div>
            <UiBaseButton
              variant="ghost"
              size="sm"
              class-name="h-8 w-8 p-0 text-muted-foreground/60 hover:bg-destructive/10 hover:text-destructive"
              @click="handleUnlinkTeam(team.team_id)"
            >
              <Trash2 class="h-4 w-4" />
            </UiBaseButton>
          </div>
          <div
            v-if="projectTeams.length === 0"
            class="text-muted-foreground py-6 text-center text-sm"
          >
            No teams linked.
          </div>
        </div>
      </div>
    </UiBaseCard>
  </div>
</template>
