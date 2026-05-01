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
