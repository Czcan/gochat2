package process

import (
	"encoding/json"
	"fmt"
	"gochat2/client/utils"
	common "gochat2/common/message"
	"gochat2/config"
	"net"
)

type MessageProcess struct{}

func (mp *MessageProcess) GetOnlineUserList() (err error) {
	serverInfo := config.Configuration.ServerInfo
	conn, err := net.Dial("tcp", serverInfo.Host)
	if err != nil {
		fmt.Printf("dial the getOnlineUserList conn error: %v\n", err)
	}

	var message common.Message
	message.Type = common.ShowAllOnlineUsersType

	requestBody, err := json.Marshal("")
	if err != nil {
		return
	}

	message.Data = string(requestBody)

	data, err := json.Marshal(message)
	if err != nil {
		return
	}

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

	return
}

func (mp *MessageProcess) SendGroupMessageToServer(groupID int, currentUserName string, content string) (err error) {
	serverInfo := config.Configuration.ServerInfo
	conn, err := net.Dial("tcp", serverInfo.Host)
	if err != nil {
		fmt.Printf("dial the sendgroupMessage conn error: %v\n", err)
	}

	var message common.Message
	message.Type = common.UserSendGroupMessageType

	userGroupMessage := common.UserSendGroupMessage{
		GroupID:  groupID,
		UserName: currentUserName,
		Content:  content,
	}

	data, err := json.Marshal(userGroupMessage)
	if err != nil {
		return
	}

	message.Data = string(data)
	groupData, _ := json.Marshal(message)

	dispatcher := utils.Dispatcher{Conn: conn}
	err = dispatcher.SendData(groupData)

	return
}

func (mp *MessageProcess) PointToPointCommunication(targetUserName string, sourceUserName string, content string) (err error) {
	serverInfo := config.Configuration.ServerInfo
	conn, err := net.Dial("tcp", serverInfo.Host)
	if err != nil {
		return
	}

	var message common.Message
	message.Type = common.PointToPointMessageType

	var pointToPointMessage common.PointToPointMessage

	pointToPointMessage = common.PointToPointMessage{
		SourceUserName: sourceUserName,
		TargetUserName: targetUserName,
		Content:        content,
	}

	data, err := json.Marshal(pointToPointMessage)
	if err != nil {
		return
	}

	message.Data = string(data)

	pointMessage, _ := json.Marshal(message)

	dispatcher := utils.Dispatcher{Conn: conn}
	err = dispatcher.SendData(pointMessage)
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

	return
}
