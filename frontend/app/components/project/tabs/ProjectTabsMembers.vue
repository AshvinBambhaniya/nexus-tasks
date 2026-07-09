<script setup lang="ts">
import {
  Trash2,
  UserPlus,
  Users,
  Link as LinkIcon,
  Mail,
  Loader2,
  Shield,
  User,
} from "lucide-vue-next";
import { cn } from "~/utils/cn";

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
  <div class="mx-auto max-w-5xl space-y-8">
    <!-- Project Members Section -->
    <div class="space-y-4">
      <div class="flex flex-col justify-between gap-4 sm:flex-row sm:items-end">
        <div>
          <h3 class="text-foreground text-lg font-semibold tracking-tight">
            Project Members
          </h3>
          <p class="text-muted-foreground text-sm">
            Manage who has access to this project and their permission levels.
          </p>
        </div>
        <form class="flex items-center gap-2" @submit.prevent="handleInvite">
          <div class="relative w-full sm:w-64">
            <Mail
              class="text-muted-foreground absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2"
            />
            <input
              v-model="inviteEmail"
              type="email"
              placeholder="Email address"
              required
              class="bg-background border-border focus:ring-primary focus:border-primary h-9 w-full rounded-md border pr-3 pl-9 text-sm transition-all focus:ring-1 focus:outline-none disabled:opacity-50"
              :disabled="isInviting"
            />
          </div>
          <button
            type="submit"
            :disabled="isInviting"
            class="bg-primary text-primary-foreground hover:bg-primary/90 flex h-9 min-w-[90px] items-center justify-center gap-2 rounded-md px-4 text-sm font-medium transition-colors disabled:opacity-50"
          >
            <Loader2 v-if="isInviting" class="h-4 w-4 animate-spin" />
            <template v-else>
              <UserPlus class="h-4 w-4" />
              Add
            </template>
          </button>
        </form>
      </div>

      <div
        class="border-border bg-card overflow-hidden rounded-lg border shadow-sm"
      >
        <div v-if="isLoading" class="flex justify-center p-8">
          <Loader2 class="text-muted-foreground/50 h-6 w-6 animate-spin" />
        </div>
        <div
          v-else-if="members.length === 0"
          class="flex flex-col items-center justify-center px-4 py-12 text-center"
        >
          <div
            class="bg-muted/50 border-border mb-3 flex h-12 w-12 items-center justify-center rounded-full border"
          >
            <Users class="text-muted-foreground h-5 w-5" />
          </div>
          <p class="text-foreground text-sm font-medium">No members yet</p>
          <p class="text-muted-foreground mt-1 text-xs">
            Add members using their email address above.
          </p>
        </div>
        <div v-else class="divide-border divide-y">
          <div
            v-for="member in members"
            :key="member.user_id"
            class="hover:bg-muted/30 group flex items-center justify-between p-4 transition-colors"
          >
            <div class="flex items-center gap-4">
              <UiBaseAvatar
                :fallback="member.email?.[0]?.toUpperCase() || '?'"
                class-name="h-10 w-10 bg-primary/10 text-primary text-sm font-medium border border-primary/20"
              />
              <div class="flex flex-col">
                <span class="text-foreground text-sm font-medium">{{
                  member.email
                }}</span>
                <span
                  class="text-muted-foreground mt-0.5 flex items-center gap-1.5 text-xs"
                >
                  <component
                    :is="member.is_direct ? User : Users"
                    class="h-3 w-3"
                  />
                  {{ member.is_direct ? "Direct Member" : "via Team" }}
                </span>
              </div>
            </div>

            <div class="flex items-center gap-4">
              <span
                :class="
                  cn(
                    'rounded-full border px-2 py-0.5 text-[11px] font-medium',
                    member.role === 'ADMIN'
                      ? 'bg-primary/10 text-primary border-primary/20'
                      : 'bg-muted text-muted-foreground border-border'
                  )
                "
              >
                <Shield
                  v-if="member.role === 'ADMIN'"
                  class="mr-1 mb-0.5 inline h-3 w-3"
                />
                {{ member.role }}
              </span>

              <div class="flex w-8 justify-end">
                <button
                  v-if="member.is_direct"
                  class="text-muted-foreground hover:text-destructive hover:bg-destructive/10 rounded-md p-1.5 opacity-0 transition-all group-hover:opacity-100 focus:opacity-100"
                  title="Remove direct access"
                  @click="handleRemoveMember(member.user_id)"
                >
                  <Trash2 class="h-4 w-4" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Linked Teams Section -->
    <div class="border-border space-y-4 border-t pt-4">
      <div class="flex flex-col justify-between gap-4 sm:flex-row sm:items-end">
        <div>
          <h3 class="text-foreground text-lg font-semibold tracking-tight">
            Linked Teams
          </h3>
          <p class="text-muted-foreground text-sm">
            Grant project access to entire workspace teams.
          </p>
        </div>
        <form class="flex items-center gap-2" @submit.prevent="handleLinkTeam">
          <div class="relative w-full sm:w-64">
            <Users
              class="text-muted-foreground pointer-events-none absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2"
            />
            <select
              v-model="selectedTeamId"
              class="bg-background border-border focus:ring-primary focus:border-primary h-9 w-full appearance-none rounded-md border pr-8 pl-9 text-sm transition-all focus:ring-1 focus:outline-none disabled:opacity-50"
              required
              :disabled="availableTeams.length === 0"
            >
              <option value="" disabled>Select a team to link...</option>
              <option v-for="t in availableTeams" :key="t.id" :value="t.id">
                {{ t.name }}
              </option>
            </select>
            <!-- Custom chevron for select -->
            <svg
              class="text-muted-foreground pointer-events-none absolute top-1/2 right-3 h-4 w-4 -translate-y-1/2"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M19 9l-7 7-7-7"
              />
            </svg>
          </div>
          <button
            type="submit"
            :disabled="isLinking || availableTeams.length === 0"
            class="bg-secondary text-secondary-foreground hover:bg-secondary/80 border-border flex h-9 min-w-[90px] items-center justify-center gap-2 rounded-md border px-4 text-sm font-medium transition-colors disabled:opacity-50"
          >
            <Loader2 v-if="isLinking" class="h-4 w-4 animate-spin" />
            <template v-else>
              <LinkIcon class="h-4 w-4" />
              Link
            </template>
          </button>
        </form>
      </div>

      <div
        class="border-border bg-card overflow-hidden rounded-lg border shadow-sm"
      >
        <div v-if="isLoading" class="flex justify-center p-8">
          <Loader2 class="text-muted-foreground/50 h-6 w-6 animate-spin" />
        </div>
        <div
          v-else-if="projectTeams.length === 0"
          class="flex flex-col items-center justify-center px-4 py-12 text-center"
        >
          <div
            class="bg-muted/50 border-border mb-3 flex h-12 w-12 items-center justify-center rounded-full border"
          >
            <LinkIcon class="text-muted-foreground h-5 w-5" />
          </div>
          <p class="text-foreground text-sm font-medium">No linked teams</p>
          <p class="text-muted-foreground mt-1 text-xs">
            Link a workspace team to give all its members access.
          </p>
        </div>
        <div v-else class="divide-border divide-y">
          <div
            v-for="team in projectTeams"
            :key="team.team_id"
            class="hover:bg-muted/30 group flex items-center justify-between p-4 transition-colors"
          >
            <div class="flex items-center gap-4">
              <div
                class="flex h-10 w-10 items-center justify-center rounded-md border border-indigo-500/20 bg-indigo-500/10 text-indigo-500"
              >
                <Users class="h-5 w-5" />
              </div>
              <span class="text-foreground text-sm font-medium">{{
                team.team_name
              }}</span>
            </div>

            <div class="flex w-8 justify-end">
              <button
                class="text-muted-foreground hover:text-destructive hover:bg-destructive/10 rounded-md p-1.5 opacity-0 transition-all group-hover:opacity-100 focus:opacity-100"
                title="Unlink team"
                @click="handleUnlinkTeam(team.team_id)"
              >
                <Trash2 class="h-4 w-4" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
