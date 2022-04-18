package process

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"gochat2/client/logger"
	"gochat2/client/model"
	"gochat2/client/utils"
	common "gochat2/common/message"
	"gochat2/config"
	"net"
	"os"
)

type UserProcess struct{}

func showAfterLoginMenu() {
	logger.Info("\n----------------login succeed!!----------------\n")
	logger.Info("\t\tselect the options:\n")
	logger.Info("\t\t\t1. Show all online users\n")
	logger.Info("\t\t\t2. Send group message\n")
	logger.Info("\t\t\t3. Point to point communicate\n")
	logger.Info("\t\t\t4. Exit \n")

	var (
		key         int
		content     string
		inputReader *bufio.Reader
		err         error
	)
	inputReader = bufio.NewReader(os.Stdin)

	fmt.Scanf("%d\n", &key)

	logger.Notice("key = %v\n", key)

	switch key {
	case 1:
		messageProcess := MessageProcess{}
		err = messageProcess.GetOnlineUserList()
		if err != nil {
			logger.Error("Some error occured when get online user list, error: %v\n", err)
		}
	case 2:
		logger.Notice("Say something:\n")
		content, err = inputReader.ReadString('\n')
		if err != nil {
			logger.Error("Some error occured when you input, error: %v\n", err)
		}

		messageProcess := MessageProcess{}
		err = messageProcess.SendGroupMessageToServer(0, model.CurrentUser.UserName, content)
		if err != nil {
			logger.Error("Some error occured when send group message, error: %v\n", err)
		} else {
			logger.Success("Send group message succeed!!")
		}

	case 3:
		var targetUserName string
		logger.Notice("Select one friend by user name:\n")
		fmt.Scanf("%s\n", &targetUserName)
		logger.Notice("input message: \n")
		content, err = inputReader.ReadString('\n')
		if err != nil {
			logger.Error("Some error occurred when you input, error: %v\n", err)
		}

		messageProcess := MessageProcess{}
		err = messageProcess.PointToPointCommunication(targetUserName, model.CurrentUser.UserName, content)
		if err != nil {
			logger.Error("Some error occurred when point to point comunication: %v\n", err)
			return
		}

		// errMsg := make(chan error)
		// go Response(conn, errMsg)
		// err = <-errMsg

		// if err != nil {
		// 	logger.Error("Send message error: %v\n", err)
		// }

	case 4:
		logger.Warn("Exit...\n")
		os.Exit(0)
	default:
		logger.Info("Selected invalid!!!\n")
	}
}

// 用户登陆
func (this *UserProcess) Login(userName string, password string) (err error) {
	// connect server
	serverInfo := config.Configuration.ServerInfo
	conn, err := net.Dial("tcp", serverInfo.Host)
	if err != nil {
		logger.Error("Connect server error: %v\n", err)
		return
	}

	var message common.Message
	message.Type = common.LoginMessageType

	var loginMessage common.LoginMessage
	loginMessage = common.LoginMessage{
		UserName: userName,
		Password: password,
	}

	data, err := json.Marshal(loginMessage)
	if err != nil {
		logger.Error("some error occured when you parse the login message, error: %v\n", err)
	}

	message.Data = string(data)

	data, _ = json.Marshal(message)

	dispatcher := utils.Dispatcher{Conn: conn}
	err = dispatcher.SendData(data)
	if err != nil {
		return
	}

	errMsg := make(chan error)
	go Response(conn, errMsg)
	err = <-errMsg
	if err != nil {
		return
	}

	for {
		showAfterLoginMenu()
	}
}

// 处理用户注册
func (this *UserProcess) Register(userName, password, passwordConfirm string) (err error) {
	if password != passwordConfirm {
		err = errors.New("confirm password not match")
		return
	}

	serverInfo := config.Configuration.ServerInfo
	conn, err := net.Dial("tcp", serverInfo.Host)
	if err != nil {
		logger.Error("Connect server error: %v", err)
		return
	}

	var message common.Message
	message.Type = common.RegisterMessageType

	var registerInfo common.RegisterMessage
	registerInfo = common.RegisterMessage{
		UserName:        userName,
		Password:        password,
		PasswordConfirm: passwordConfirm,
	}

	registermessage, err := json.Marshal(registerInfo)
	if err != nil {
		logger.Error("Some error occured when you parse the register message, error: %v\n", err)
	}

	message.Data = string(registermessage)

	data, err := json.Marshal(message)
	if err != nil {
		logger.Error("RegisterMessage json Marshal error: %v\n", err)
		return
	}

	dispatcher := utils.Dispatcher{Conn: conn}
	err = dispatcher.SendData(data)
	if err != nil {
		logger.Error("Send register data error: %v\n", err)
	}

	errMsg := make(chan error)
	go Response(conn, errMsg)
	err = <-errMsg

	return
}
