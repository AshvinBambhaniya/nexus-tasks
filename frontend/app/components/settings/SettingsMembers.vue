<script setup lang="ts">
import {
  Loader2,
  Trash2,
  UserPlus,
  Mail,
  Shield,
  User as UserIcon,
  Search,
} from "lucide-vue-next";
import type { WorkspaceMember } from "~/types";
import { useWorkspaces } from "~/composables/useWorkspaces";
import { useUsersStore } from "~/stores/user";

const { activeWorkspace } = useWorkspaces();
const userStore = useUsersStore();

const inviteEmail = ref("");
const isInviting = ref(false);
const searchQuery = ref("");

const {
  data: members,
  pending: isLoading,
  refresh,
} = useApi<WorkspaceMember[]>(
  () =>
    activeWorkspace.value
      ? `/api/v2/workspaces/${activeWorkspace.value.id}/members`
      : "",
  {
    key: "workspace-members",
    watch: [activeWorkspace],
    immediate: !!activeWorkspace.value,
  }
);

const filteredMembers = computed(() => {
  if (!members.value) return [];
  if (!searchQuery.value) return members.value;

  const query = searchQuery.value.toLowerCase();
  return members.value.filter(
    (m) =>
      m.user.email.toLowerCase().includes(query) ||
      m.user.full_name?.toLowerCase().includes(query)
  );
});

const handleInvite = async () => {
  if (!activeWorkspace.value || !inviteEmail.value) return;

  isInviting.value = true;
  try {
    await useMutation(
      `/api/v2/workspaces/${activeWorkspace.value.id}/members`,
      {
        method: "POST",
        body: { email: inviteEmail.value },
      }
    );
    inviteEmail.value = "";
    await refresh();
  } catch (err) {
    alert(getApiErrorMessage(err, "Failed to invite user"));
  } finally {
    isInviting.value = false;
  }
};

const handleRemove = async (userId: string) => {
  if (
    !activeWorkspace.value ||
    !confirm("Are you sure you want to remove this member?")
  )
    return;

  try {
    await useMutation(
      `/api/v2/workspaces/${activeWorkspace.value.id}/members/${userId}`,
      {
        method: "DELETE",
      }
    );
    await refresh();
  } catch (err) {
    alert(getApiErrorMessage(err, "Failed to remove member"));
  }
};
</script>

<template>
  <div>
    <div
      v-if="!activeWorkspace"
      class="border-border bg-muted/50 flex flex-col items-center justify-center rounded-xl border border-dashed p-16 text-center"
    >
      <div
        class="bg-card ring-border mb-4 flex h-16 w-16 items-center justify-center rounded-full shadow-sm ring-1"
      >
        <UserIcon class="text-muted-foreground/60 h-8 w-8" />
      </div>
      <h3 class="text-card-foreground text-lg font-bold">
        No Workspace Selected
      </h3>
      <p class="text-muted-foreground mt-1 max-w-xs text-sm text-pretty">
        Select a workspace from the dashboard to manage its members and
        permissions.
      </p>
    </div>

    <div
      v-else
      class="animate-in fade-in slide-in-from-right-4 space-y-8 duration-300"
    >
      <!-- Invite Section -->
      <section>
        <div class="mb-4">
          <h2 class="text-foreground text-lg font-semibold">
            Invite New Members
          </h2>
          <p class="text-muted-foreground text-sm">
            Expand your team by inviting collaborators to this workspace.
          </p>
        </div>

        <UiBaseCard
          class="border-border shadow-sm transition-shadow hover:shadow-md"
        >
          <div class="p-8">
            <form
              class="flex flex-col gap-4 sm:flex-row"
              @submit.prevent="handleInvite"
            >
              <div class="relative flex-1">
                <Mail
                  class="text-muted-foreground/60 absolute top-1/2 left-3.5 h-4 w-4 -translate-y-1/2"
                />
                <UiBaseInput
                  v-model="inviteEmail"
                  placeholder="colleague@company.com"
                  type="email"
                  required
                  class="focus:ring-primary/10 h-11 pl-10.5 transition-all focus:ring-2"
                />
              </div>
              <UiBaseButton
                type="submit"
                :disabled="isInviting"
                class="h-11 px-8 shadow-sm active:scale-95"
              >
                <Loader2 v-if="isInviting" class="mr-2 h-4 w-4 animate-spin" />
                <UserPlus v-else class="mr-2 h-4 w-4" />
                Send Invite
              </UiBaseButton>
            </form>
            <div
              class="text-muted-foreground/60 mt-4 flex items-center gap-2 text-[11px] font-medium"
            >
              <Shield class="h-3 w-3" />
              <span>Invited members will join as "Member" by default.</span>
            </div>
          </div>
        </UiBaseCard>
      </section>

      <!-- Members List Section -->
      <section>
        <div class="mb-4 flex items-end justify-between">
          <div>
            <h2 class="text-foreground text-lg font-semibold">
              Workspace Members
            </h2>
            <p class="text-muted-foreground text-sm">
              A list of all users with access to this environment.
            </p>
          </div>
          <div
            class="bg-primary/10 text-primary rounded-md px-2 py-1 text-xs font-bold"
          >
            {{ members?.length || 0 }} Total
          </div>
        </div>

        <UiBaseCard class="border-border overflow-hidden shadow-sm">
          <div class="border-border bg-muted/30 border-b p-4">
            <div class="relative max-w-md">
              <Search
                class="text-muted-foreground/60 absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2"
              />
              <UiBaseInput
                v-model="searchQuery"
                placeholder="Filter members..."
                class="bg-card focus:ring-primary/10 h-9 border-none pl-9 shadow-sm focus:ring-1"
              />
            </div>
          </div>

          <div
            v-if="isLoading"
            class="flex flex-col items-center justify-center p-20"
          >
            <Loader2 class="text-primary/20 h-10 w-10 animate-spin" />
            <p class="text-muted-foreground/60 mt-4 text-sm font-medium">
              Loading directory...
            </p>
          </div>

          <div v-else class="divide-border divide-y">
            <div
              v-for="member in filteredMembers"
              :key="member.user_id"
              class="group hover:bg-primary/5 flex items-center justify-between p-6 transition-all"
            >
              <div class="flex items-center gap-5">
                <div class="relative">
                  <UiBaseAvatar
                    :fallback="
                      (member.user.full_name ||
                        member.user.email)[0].toUpperCase()
                    "
                    class="border-border bg-card text-primary group-hover:ring-primary/20 h-12 w-12 border font-bold shadow-sm ring-2 ring-transparent transition-all"
                  />
                  <div
                    v-if="member.role === 'ADMIN'"
                    class="bg-primary text-primary-foreground ring-card absolute -right-1 -bottom-1 flex h-5 w-5 items-center justify-center rounded-full ring-2"
                  >
                    <Shield class="h-3 w-3" />
                  </div>
                </div>
                <div>
                  <div class="flex items-center gap-2">
                    <span class="text-card-foreground font-bold">{{
                      member.user.full_name || member.user.email.split("@")[0]
                    }}</span>
                    <UiBaseBadge
                      v-if="member.user.id === userStore.userData?.id"
                      class="bg-muted text-muted-foreground text-[9px] font-black uppercase"
                      >You</UiBaseBadge
                    >
                  </div>
                  <div
                    class="text-muted-foreground group-hover:text-primary/70 text-sm transition-colors"
                  >
                    {{ member.user.email }}
                  </div>
                </div>
              </div>

              <div class="flex items-center gap-6">
                <UiBaseBadge
                  :variant="member.role === 'ADMIN' ? 'default' : 'secondary'"
                  :class="[
                    'h-6 px-2.5 text-[10px] font-bold tracking-wider uppercase ring-1',
                    member.role === 'ADMIN'
                      ? 'bg-primary ring-primary-foreground/20'
                      : 'bg-card text-muted-foreground ring-border',
                  ]"
                >
                  {{ member.role.toLowerCase() }}
                </UiBaseBadge>

                <div class="flex w-10 justify-center">
                  <UiBaseButton
                    v-if="
                      member.role !== 'ADMIN' &&
                      member.user.id !== userStore.userData?.id
                    "
                    variant="ghost"
                    size="sm"
                    class="text-muted/60 hover:bg-destructive/10 hover:text-destructive h-9 w-9 rounded-full p-0 transition-all active:scale-90"
                    title="Remove from workspace"
                    @click="handleRemove(member.user_id)"
                  >
                    <Trash2 class="h-4.5 w-4.5" />
                  </UiBaseButton>
                </div>
              </div>
            </div>

            <div
              v-if="!filteredMembers || filteredMembers.length === 0"
              class="flex flex-col items-center justify-center p-20 text-center"
            >
              <div class="bg-muted mb-4 rounded-full p-4">
                <Search class="text-muted h-8 w-8" />
              </div>
              <h3 class="text-card-foreground font-medium">No members found</h3>
              <p class="text-muted-foreground mt-1 text-sm">
                Try adjusting your search or inviting new team members.
              </p>
            </div>
          </div>
        </UiBaseCard>
      </section>
    </div>
  </div>
</template>
