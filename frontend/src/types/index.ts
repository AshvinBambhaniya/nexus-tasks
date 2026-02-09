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
  id: number;
  name: string;
  type: WorkspaceType;
  owner_id: number;
}

export interface User {
  id: number;
  email: string;
  is_active: boolean;
  full_name?: string;
}

export interface WorkspaceMember {
  workspace_id: number;
  user_id: number;
  role: WorkspaceRole;
  user: {
    id: number;
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
  id: number;
  name: string;
  description?: string;
  workspace_id: number;
  created_at: string;
}

export interface TeamMember {
  team_id: number;
  user_id: number;
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
  id: number;
  name: string;
  description?: string;
  workspace_id: number;
  created_at: string;
  is_archived: boolean;
}

export interface ProjectMember {
  project_id: number;
  user_id: number;
  role: ProjectRole;
  email: string; // Flattened for display
  is_direct: boolean;
}

export interface ProjectTeam {
  project_id: number;
  team_id: number;
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
  id: number;
  title: string;
  description?: string;
  status: TaskStatus;
  priority: TaskPriority;
  project_id: number; // Updated from workspace_id
  assignee_id?: number;
  assignee?: {
    id: number;
    email: string;
    full_name?: string;
  };
  author?: {
    id: number;
    email: string;
    full_name?: string;
  };
  due_date?: string;
  created_at: string;
  updated_at: string;
  comment_count?: number;
}

export interface ProjectInfo {
  id: number;
  name: string;
}

export interface TaskWithProject extends Task {
  project: ProjectInfo;
}