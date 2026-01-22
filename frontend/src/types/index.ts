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
}

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
  workspace_id: number;
  assignee_id?: number;
  due_date?: string;
  created_at: string;
  updated_at: string;
}
