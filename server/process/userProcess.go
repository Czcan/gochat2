package process

import (
	"encoding/json"
	"fmt"
	common "gochat2/common/message"
	"gochat2/server/model"
	"gochat2/server/utils"
	"net"
)

type UserProcess struct {
	Conn net.Conn
}

func register(userName, password, passwordConfirm string) (user model.User, err error) {
	user, err = model.CurrentUserDao.Register(userName, password, passwordConfirm)
	return
}

func login(userName, password string) (user model.User, err error) {
	user, err = model.CurrentUserDao.Login(userName, password)
	return
}

// 响应客户端
func (this *UserProcess) responseClient(responseMessageType string, code int, data string, err error) {
	var responseMessage common.ResponseMessage
	responseMessage.Code = code
	responseMessage.Data = data
	responseMessage.Type = responseMessageType

	responseData, err := json.Marshal(responseMessage.Data)
	if err != nil {
		fmt.Printf("json marshal response message error: %v\n", err)
		return
	}

	dispatcher := utils.Dispatcher{Conn: this.Conn}

	err = dispatcher.WriteData(responseData)
}

func (this *UserProcess) UserRegister(message string) (err error) {
	var (
		registerMessage common.RegisterMessage
		code            int
		data            string
	)

	err = json.Unmarshal([]byte(message), &registerMessage)
	if err != nil {
		code = common.ServerError
	}

	_, err = register(registerMessage.UserName, registerMessage.Password, registerMessage.PasswordConfirm)
	switch err {
	case nil:
		code = common.RegisterSuccessed
	case model.ERROR_PASSWORD_DOES_NOT_MATCH:
		code = 402
	case model.ERROR_USER_ALREADY_EXISTS:
		code = 403
	default:
		code = 500
	}

	this.responseClient(common.RegisterResponseMessageType, code, data, err)
	return
}

func (this *UserProcess) UserLogin(message string) (err error) {
	var (
		loginMessage common.LoginMessage
		code         int
		data         string
	)

	err = json.Unmarshal([]byte(message), &loginMessage)
	if err != nil {
		code = common.ServerError
	}

	user, err := login(loginMessage.UserName, loginMessage.Password)
	switch err {
	case nil:
		code = common.LoginSuccess

		// save user conn status
		clientConn := model.ClientConn{}
		clientConn.Save(user.ID, user.Name, this.Conn)

		userInfo := &common.UserInfo{ID: user.ID, UserName: user.Name}
		info, _ := json.Marshal(userInfo)
		data = string(info)

	case model.ERROR_USER_DOES_NOT_EXISTS:
		code = 404
	case model.ERROR_USER_PASSWORD:
		code = 403
	default:
		code = 500
	}

	this.responseClient(common.LoginResponseMessageType, code, data, err)
	return
}
