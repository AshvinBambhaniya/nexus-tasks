package structs

type ReqRegisterUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
}

type ReqLoginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ReqCreateWorkspace struct {
	Name string `json:"name" validate:"required"`
}

type ReqCreateTeam struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type ReqInviteWorkspaceMember struct {
	Email string `json:"email" validate:"required,email"`
}

type ReqUpdateTeam struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ReqAddTeamMember struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required"` // Should be validated against ENUM
}
