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
  projectId: number;
}

const { projectId } = defineProps<Props>();

const { members, isLoading, addMember, removeMember } =
  useProjectMembers(projectId);
const {
  teams: projectTeams,
  addTeam,
  removeTeam,
} = useProjectTeams(projectId);
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

const handleRemoveMember = async (userId: number) => {
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
    await addTeam(parseInt(selectedTeamId.value));
    selectedTeamId.value = "";
  } catch (err) {
    alert(getApiErrorMessage(err, "Failed to link team"));
  } finally {
    isLinking.value = false;
  }
};

const handleUnlinkTeam = async (teamId: number) => {
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
      <div class="border-b border-gray-100 p-6">
        <h3 class="text-lg font-semibold text-gray-900">Project Members</h3>
        <p class="text-sm text-gray-500">
          All users with access to this project, including team members.
        </p>
      </div>
      <div class="p-6 space-y-6">
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
                :fallback="member.email[0].toUpperCase()"
                class-name="h-9 w-9 bg-blue-100 text-xs text-blue-700"
              />
              <div>
                <div class="text-sm font-medium text-gray-900">
                  {{ member.email }}
                </div>
                <div class="flex items-center gap-2 text-xs text-gray-500">
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
                class-name="h-8 w-8 p-0 text-gray-400 hover:bg-red-50 hover:text-red-600"
                title="Remove direct access"
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

    <UiBaseCard>
      <div class="border-b border-gray-100 p-6">
        <h3 class="text-lg font-semibold text-gray-900">Teams</h3>
        <p class="text-sm text-gray-500">
          Workspace teams with access to this project.
        </p>
      </div>
      <div class="p-6 space-y-6">
        <form class="flex gap-3" @submit.prevent="handleLinkTeam">
          <div class="relative flex-1">
            <Users class="absolute top-2.5 left-3 h-4 w-4 text-gray-400" />
            <select
              v-model="selectedTeamId"
              class="flex h-10 w-full rounded-md border border-gray-200 bg-white px-3 py-2 pl-9 text-sm ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 disabled:cursor-not-allowed disabled:opacity-50"
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
          <Loader2 class="h-6 w-6 animate-spin text-gray-400" />
        </div>
        <div v-else class="divide-y divide-gray-100">
          <div
            v-for="team in projectTeams"
            :key="team.team_id"
            class="flex items-center justify-between rounded-lg p-4 transition-colors hover:bg-gray-50"
          >
            <div class="flex items-center gap-3">
              <div
                class="flex h-9 w-9 items-center justify-center rounded-lg bg-indigo-100 text-indigo-700"
              >
                <Users class="h-4 w-4" />
              </div>
              <div class="text-sm font-medium text-gray-900">
                {{ team.team_name }}
              </div>
            </div>
            <UiBaseButton
              variant="ghost"
              size="sm"
              class-name="h-8 w-8 p-0 text-gray-400 hover:bg-red-50 hover:text-red-600"
              @click="handleUnlinkTeam(team.team_id)"
            >
              <Trash2 class="h-4 w-4" />
            </UiBaseButton>
          </div>
          <div
            v-if="projectTeams.length === 0"
            class="py-6 text-center text-sm text-gray-500"
          >
            No teams linked.
          </div>
        </div>
      </div>
    </UiBaseCard>
  </div>
</template>
