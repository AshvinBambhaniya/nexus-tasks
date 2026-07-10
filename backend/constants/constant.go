// Package constants defines constants used throughout the application.
package constants

// Cookie constants
const (
	CookieUser            = "access_token"
	WorkspaceTypePersonal = "Personal"
)

// Fiber context keys
const (
	ContextUID             = "userId"
	ContextWorkspaceMember = "workspaceMember"
)

// URL parameter keys
const (
	ParamUID         = "userId"
	ParamWorkspaceID = "workspaceId"
	ParamTeamID      = "teamId"
	ParamProjectID   = "projectId"
	ParamTaskID      = "taskId"
	ParamCommentID   = "commentId"
)

// Response property keys
const (
	// PropMessage is a key for messages
	PropMessage = "message"
	// PropError is a key for errors
	PropError = "error"
	// ResponseKeyContent is a key for content in responses
	ResponseKeyContent = "content"
)

// Fail messages
// ...
const (
	Unauthenticated = "unauthenticated to access resource"
	Unauthorized    = "unauthorized to access resource"
	// #nosec G101
	InvalidCredentials = "invalid credentials"
	UserNotExist       = "user does not exists"
)

// Invalid ID Messages
const (
	ErrInvalidUserID       = "invalid user id"
	ErrInvalidWorkspaceID  = "invalid workspace id"
	ErrInvalidTeamID       = "invalid team id"
	ErrInvalidProjectID    = "invalid project id"
	ErrInvalidTaskID       = "invalid task id"
	ErrInvalidCommentID    = "invalid comment id"
	ErrInvalidTargetUserID = "invalid target user id"
)

// Shared Success Messages
const (
	MsgMemberRemoved  = "Member removed"
	MsgMemberAdded    = "Member added"
	MsgMemberInvited  = "Member invited"
	MsgUpdated        = "Updated"
	MsgTaskDeleted    = "Task deleted"
	MsgCommentDeleted = "Comment deleted"
	MsgTeamRemoved    = "Team removed"
	MsgTeamDeleted    = "Team deleted"
)

// Error messages
const (
	ErrGetUser             = "error while get user"
	ErrLoginUser           = "error while login user"
	ErrInsertUser          = "error while creating user, please try after sometime"
	ErrHealthCheckDb       = "error while checking health of database"
	ErrUnauthenticated     = "error verifing user identity"
	ErrKratosAuth          = "error while fetching user from kratos"
	ErrKratosDataInsertion = "error while inserting user data came from kratos"
	ErrKratosIDEmpty       = "error no session_id found in kratos cookie"
	ErrKratosCookieTime    = "error while parsing the expiration time of the cookie"
)

// Events
const (
	EventUserRegistered = "event:userRegistered"
)

// Topics
const (
	TopicWorkspaceInvites = "workspace_invites"
	TopicWelcomeMail      = "welcome_mail"
)

// API Key constants
const (
	ParamKeyID      = "keyId"
	ErrInvalidKeyID = "invalid api key id"
	MsgKeyCreated   = "API key created"
	MsgKeyRevoked   = "API key revoked"
	APIKeyPrefix    = "ntx_"
)

// Query parameters
const (
	QueryStatus     = "status"
	QueryAssigneeID = "assignee_id"
)
