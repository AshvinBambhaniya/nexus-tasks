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
  <div class="space-y-12 pb-10">
    <div
      v-if="!activeWorkspace"
      class="border-border bg-card flex flex-col items-center justify-center rounded-xl border p-12 text-center"
    >
      <div
        class="bg-muted mb-4 flex h-12 w-12 items-center justify-center rounded-full"
      >
        <UserIcon class="text-muted-foreground h-6 w-6" />
      </div>
      <h3 class="text-card-foreground text-sm font-semibold">
        No Workspace Selected
      </h3>
      <p class="text-muted-foreground mt-1 max-w-xs text-sm">
        Select a workspace from the dashboard to manage its members.
      </p>
    </div>

    <template v-else>
      <!-- Invite Section -->
      <section
        class="flex flex-col gap-6 lg:flex-row lg:items-start lg:justify-between"
      >
        <div class="lg:w-1/3">
          <h2 class="text-foreground text-lg font-semibold tracking-tight">
            Invite Members
          </h2>
          <p class="text-muted-foreground mt-1 text-sm">
            Expand your team by inviting collaborators to this workspace.
          </p>
        </div>

        <UiBaseCard class="border-border bg-card lg:w-2/3">
          <div class="p-6">
            <form
              class="flex flex-col gap-3 sm:flex-row"
              @submit.prevent="handleInvite"
            >
              <div class="relative flex-1">
                <Mail
                  class="text-muted-foreground/60 absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2"
                />
                <UiBaseInput
                  v-model="inviteEmail"
                  placeholder="colleague@company.com"
                  type="email"
                  required
                  class="pl-9"
                />
              </div>
              <UiBaseButton
                type="submit"
                :disabled="isInviting"
                class="min-w-[120px]"
              >
                <Loader2 v-if="isInviting" class="mr-2 h-4 w-4 animate-spin" />
                <UserPlus v-else class="mr-2 h-4 w-4" />
                Send Invite
              </UiBaseButton>
            </form>
            <div class="text-muted-foreground/60 mt-3 text-xs">
              Invited members will join as "Member" by default.
            </div>
          </div>
        </UiBaseCard>
      </section>

      <div class="border-border border-t" />

      <!-- Members List Section -->
      <section
        class="flex flex-col gap-6 lg:flex-row lg:items-start lg:justify-between"
      >
        <div class="lg:w-1/3">
          <h2 class="text-foreground text-lg font-semibold tracking-tight">
            Workspace Members
          </h2>
          <p class="text-muted-foreground mt-1 text-sm">
            A list of all users with access to this environment.
          </p>
        </div>

        <UiBaseCard class="border-border bg-card overflow-hidden lg:w-2/3">
          <div class="border-border bg-muted/30 border-b p-3">
            <div class="relative w-full sm:max-w-xs">
              <Search
                class="text-muted-foreground/60 absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2"
              />
              <input
                v-model="searchQuery"
                type="text"
                placeholder="Filter members..."
                class="bg-background border-border text-foreground placeholder:text-muted-foreground/60 focus:ring-primary/20 w-full rounded-md border py-1.5 pr-3 pl-9 text-sm transition-colors focus:ring-2 focus:outline-none"
              />
            </div>
          </div>

          <div
            v-if="isLoading"
            class="flex flex-col items-center justify-center p-12"
          >
            <Loader2 class="text-muted-foreground/40 h-8 w-8 animate-spin" />
          </div>

          <div v-else class="divide-border divide-y">
            <div
              v-for="member in filteredMembers"
              :key="member.user_id"
              class="hover:bg-muted/30 flex items-center justify-between p-4 transition-colors"
            >
              <div class="flex items-center gap-4">
                <div class="relative">
                  <UiBaseAvatar
                    :fallback="
                      (member.user.full_name ||
                        member.user.email)[0].toUpperCase()
                    "
                    class="h-10 w-10 text-sm font-medium"
                  />
                  <div
                    v-if="member.role === 'ADMIN'"
                    class="bg-primary text-primary-foreground ring-card absolute -right-1 -bottom-1 flex h-4 w-4 items-center justify-center rounded-full ring-2"
                  >
                    <Shield class="h-2 w-2" />
                  </div>
                </div>
                <div>
                  <div class="flex items-center gap-2">
                    <span class="text-card-foreground text-sm font-semibold">
                      {{
                        member.user.full_name || member.user.email.split("@")[0]
                      }}
                    </span>
                    <UiBaseBadge
                      v-if="member.user.id === userStore.userData?.id"
                      variant="secondary"
                      class="px-1.5 py-0.5 text-[9px] leading-none font-bold uppercase"
                    >
                      You
                    </UiBaseBadge>
                  </div>
                  <div class="text-muted-foreground text-xs">
                    {{ member.user.email }}
                  </div>
                </div>
              </div>

              <div class="flex items-center gap-4">
                <UiBaseBadge
                  :variant="member.role === 'ADMIN' ? 'default' : 'secondary'"
                  class="text-[10px] font-bold uppercase"
                >
                  {{ member.role.toLowerCase() }}
                </UiBaseBadge>

                <div class="flex w-8 justify-end">
                  <UiBaseButton
                    v-if="
                      member.role !== 'ADMIN' &&
                      member.user.id !== userStore.userData?.id
                    "
                    variant="ghost"
                    size="sm"
                    class="text-muted-foreground hover:bg-destructive/10 hover:text-destructive h-8 w-8 p-0"
                    title="Remove member"
                    @click="handleRemove(member.user_id)"
                  >
                    <Trash2 class="h-4 w-4" />
                  </UiBaseButton>
                </div>
              </div>
            </div>

            <div
              v-if="!filteredMembers || filteredMembers.length === 0"
              class="flex flex-col items-center justify-center p-12 text-center"
            >
              <div class="bg-muted mb-3 rounded-full p-3">
                <Search class="text-muted-foreground/50 h-6 w-6" />
              </div>
              <h3 class="text-card-foreground text-sm font-medium">
                No members found
              </h3>
              <p class="text-muted-foreground mt-1 text-xs">
                Try adjusting your search criteria.
              </p>
            </div>
          </div>
        </UiBaseCard>
      </section>
    </template>
  </div>
</template>
