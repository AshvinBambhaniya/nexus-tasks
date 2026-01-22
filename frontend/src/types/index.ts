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
