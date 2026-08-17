export enum WorkspaceType {
  PERSONAL = "PERSONAL",
  TEAM = "TEAM",
}

export enum WorkspaceRole {
  ADMIN = "ADMIN",
  MEMBER = "MEMBER",
  VIEWER = "VIEWER",
}

export interface Workspace {
  id: string;
  name: string;
  type: WorkspaceType;
  owner_id: string;
}

export interface User {
  id: string;
  email: string;
  is_active: boolean;
  full_name?: string;
}

export interface WorkspaceMember {
  workspace_id: string;
  user_id: string;
  role: WorkspaceRole;
  user: {
    id: string;
    email: string;
    full_name?: string;
  };
}

// Teams
export enum TeamRole {
  ADMIN = "ADMIN",
  MEMBER = "MEMBER",
}

export interface Team {
  id: string;
  name: string;
  description?: string;
  workspace_id: string;
  created_at: string;
  projects?: Project[];
}

export interface TeamMember {
  team_id: string;
  user_id: string;
  role: TeamRole;
  email: string; // Flattened for display
}

// Projects
export enum ProjectRole {
  ADMIN = "ADMIN",
  MEMBER = "MEMBER",
  VIEWER = "VIEWER",
}

export interface Project {
  id: string;
  name: string;
  description?: string;
  workspace_id: string;
  created_at: string;
  is_archived: boolean;
}

export interface ProjectMember {
  project_id: string;
  user_id: string;
  role: ProjectRole;
  email: string; // Flattened for display
  full_name?: string;
  is_direct: boolean;
}

export interface ProjectTeam {
  project_id: string;
  team_id: string;
  team_name: string;
}

// Tasks
export enum TaskStatus {
  TODO = "TODO",
  IN_PROGRESS = "IN_PROGRESS",
  DONE = "DONE",
  BACKLOG = "BACKLOG",
}

export enum TaskPriority {
  P0 = "P0", // Critical
  P1 = "P1", // High
  P2 = "P2", // Medium
  P3 = "P3", // Low
}

export interface Task {
  id: string;
  number: number;
  title: string;
  description?: string;
  status: TaskStatus;
  priority: TaskPriority;
  project_id: string; // Updated from workspace_id
  assignee_id?: string;
  assignee?: {
    id: string;
    email: string;
    full_name?: string;
  };
  author?: {
    id: string;
    email: string;
    full_name?: string;
  };
  due_date?: string;
  completed_at?: string;
  created_at: string;
  updated_at: string;
  comment_count?: number;
  estimated_minutes?: number;
}

export interface ProjectInfo {
  id: string;
  name: string;
}

export interface TaskWithProject extends Task {
  project: ProjectInfo;
}

export interface Comment {
  id: string;
  content: string;
  task_id: string;
  author_id: string;
  author?: {
    id: string;
    email: string;
    full_name?: string;
  };
  created_at: string;
  updated_at: string;
}

export interface JSendResponse<T> {
  status: "success" | "fail" | "error";
  data?: T;
  message?: string;
}

export interface ApiError {
  response?: {
    data?: JSendResponse<unknown>;
  };
}

export enum NotificationEntityType {
  TASK = "TASK",
  COMMENT = "COMMENT",
  PROJECT = "PROJECT",
}

export enum NotificationType {
  ASSIGNED = "ASSIGNED",
  MENTIONED = "MENTIONED",
  STATUS_CHANGED = "STATUS_CHANGED",
  COMMENT_ADDED = "COMMENT_ADDED",
}

export interface Notification {
  id: string;
  user_id: string;
  actor_id: string;
  entity_id: string;
  entity_type: NotificationEntityType;
  type: NotificationType;
  title: string;
  body: string | null;
  is_read: boolean;
  is_cleared: boolean;
  created_at: string;
  updated_at: string;
}

// API Key Types
export interface ApiKey {
  id: string;
  user_id: string;
  name: string;
  token_prefix: string;
  last_used_at: string | null;
  expires_at: string | null;
  created_at: string;
  updated_at: string;
}

export interface CreateApiKeyResponse {
  raw_token: string;
  key: ApiKey;
}

// Time Tracking
export interface TimeEntry {
  id: string;
  task_id: string;
  user_id: string;
  user_full_name: string;
  description: string;
  start_time: string;
  end_time: string | null;
  duration_minutes: number | null;
  is_manual: boolean;
  created_at: string;
}

export interface ProjectTimeEntry extends TimeEntry {
  task_title: string;
  task_number: number;
}

export interface ActiveTimer {
  id: string;
  task_id: string;
  task_title: string;
  task_number: number;
  start_time: string;
}

export interface TaskTimeEntries {
  entries: TimeEntry[];
  total_logged_minutes: number;
  estimated_minutes: number | null;
}

export interface TaskTimeSummary {
  task_id: string;
  task_number: number;
  task_title: string;
  estimated_minutes: number | null;
  logged_minutes: number;
  is_over_budget: boolean;
}

export interface MemberTimeSummary {
  user_id: string;
  full_name: string;
  logged_minutes: number;
}

export interface DailyTimeEntry {
  date: string;
  logged_minutes: number;
}

export interface ProjectTimeAnalytics {
  total_estimated_minutes: number;
  total_logged_minutes: number;
  estimate_accuracy_percent: number;
  over_budget_task_count: number;
  by_task: TaskTimeSummary[];
  by_member: MemberTimeSummary[];
  daily_trend: DailyTimeEntry[];
}
