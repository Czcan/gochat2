package process

import (
	"encoding/json"
	"errors"
	"gochat2/client/logger"
	"gochat2/client/model"
	"gochat2/client/utils"
	common "gochat2/common/message"
	"io"
	"net"
)

func dealLoginResponse(responseMsg common.ResponseMessage) (err error) {
	switch responseMsg.Code {
	case 200:
		// 解析当前用户信息
		var userInfo common.UserInfo
		err = json.Unmarshal([]byte(responseMsg.Data), &userInfo)
		if err != nil {
			return
		}

		user := model.User{}
		err = user.InitCurrentUser(userInfo.ID, userInfo.UserName)
		logger.Success("Login succeed!!!")
		logger.Notice("Current user ID: %v,  user Name: %v \n", model.CurrentUser.UserID, model.CurrentUser.UserName)
		if err != nil {
			return
		}
	case 403:
		err = errors.New("Password invalid!!")
	case 404:
		err = errors.New("User does not exist!")
	case 500:
		err = errors.New("Server Error!")
	default:
		err = errors.New("Some error, unknown response status")
	}

	return
}

func dealRegisterResponse(responseMessage common.ResponseMessage) (err error) {
	switch responseMessage.Code {
	case 200:
		logger.Success("Register succeed!!!\n")
	case 500:
		err = errors.New("Server Error!")
	case 403:
		err = errors.New("User already exists!")
	case 402:
		err = errors.New("Password invalid!")
	default:
		err = errors.New("Some error, unknown response status")
	}

	return
}

func showOnlineUserList(responseMsg common.ResponseMessage) (err error) {
	if responseMsg.Code != 200 {
		err = errors.New("Server Error!")
	}

	var userList []common.UserInfo
	err = json.Unmarshal([]byte(responseMsg.Data), &userList)

	logger.Success("Online user list(%v users)\n", len(userList))
	logger.Notice("\t\tID\t\tUserName\n")

	for _, userInfo := range userList {
		logger.Success("\t\t%v\t\t%v\n", userInfo.ID, userInfo.UserName)
	}

	return
}

func dealGroupMessage(responseMsg common.ResponseMessage) (err error) {
	var groupMessage common.UserSendGroupMessage
	err = json.Unmarshal([]byte(responseMsg.Data), &groupMessage)
	if err != nil {
		return
	}

	logger.Success("%v send to you: \n", groupMessage.UserName)
	logger.Notice("%v\n", groupMessage.Content)

	return
}

func dealPointToPointCommunicate(responseMsg common.ResponseMessage) (err error) {
	if responseMsg.Code != 200 {
		err = errors.New(responseMsg.Error)
	}

	var pointToPointMessage common.PointToPointMessage
	err = json.Unmarshal([]byte(responseMsg.Data), &pointToPointMessage)
	if err != nil {
		return
	}

	logger.Success("%v say\n", pointToPointMessage.SourceUserName)
	logger.Notice("\t%v\n", pointToPointMessage.Content)

	return
}

// 处理服务端的返回
var responseMsg common.ResponseMessage

func Response(conn net.Conn, errMsg chan error) (err error) {

	dispatcher := utils.Dispatcher{Conn: conn}

	for {
		responseMsg, err = dispatcher.ReadData()
		if err != nil && err != io.EOF {
			logger.Error("Waiting response error: %v\n", err)
		}

		logger.Error("responseMsg Type:  %v\n", responseMsg.Type)

		switch responseMsg.Type {
		case common.LoginResponseMessageType:
			err = dealLoginResponse(responseMsg)
			errMsg <- err
		case common.RegisterResponseMessageType:
			err = dealRegisterResponse(responseMsg)
			errMsg <- err
		case common.ShowAllOnlineUsersType:
			err = showOnlineUserList(responseMsg)
			errMsg <- err
		case common.SendGroupMessageToClientType:
			err = dealGroupMessage(responseMsg)
			errMsg <- err
		case common.PointToPointMessageType:
			err = dealPointToPointCommunicate(responseMsg)
			errMsg <- err
		default:
			logger.Error("Unknown message type!")
		}

		if err != nil {
			return
		}
	}
}
