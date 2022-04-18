package common

const (
	LoginMessageType             = "LoginMessage"
	RegisterMessageType          = "RegisterMessage"
	LoginResponseMessageType     = "LoginResponseMessageType"
	RegisterResponseMessageType  = "RegisterResponseMessageType"
	UserSendGroupMessageType     = "UserSendGroupMessageType"
	SendGroupMessageToClientType = "SendGroupMessageToClientType"
	ShowAllOnlineUsersType       = "ShowAllOnlineUsersType"
	PointToPointMessageType      = "PointToPointMessageType"

	ServerError = 500

	// status code for login
	LoginError   = 403
	NotExit      = 404
	LoginSucceed = 200

	// status code for register
	HasExited         = 403
	RegisterSuccessed = 200
	PasswordNotMatch  = 402
)

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type LoginMessage struct {
	UserName string
	Password string
}

type RegisterMessage struct {
	UserName        string
	Password        string
	PasswordConfirm string
}

type ResponseMessage struct {
	Type  string
	Code  int
	Error string // 错误消息
	Data  string
}

type UserSendGroupMessage struct {
	GroupID  int    // target group id => 0: all users
	UserName string // current user name
	Content  string
}

type PointToPointMessage struct {
	SourceUserID   int
	SourceUserName string
	TargetUserID   int
	TargetUserName string
	Content        string
}

// online user info
type UserInfo struct {
	ID       int
	UserName string
}
