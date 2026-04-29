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
