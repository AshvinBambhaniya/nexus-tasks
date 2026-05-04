package structs

// WorkspaceInvitePayload for email worker
type WorkspaceInvitePayload struct {
	Email         string `json:"email"`
	WorkspaceName string `json:"workspace_name"`
}
