package process

import (
	"encoding/json"
	"fmt"
	common "gochat2/common/message"
	"gochat2/server/model"
	"gochat2/server/utils"
	"net"
)

type PointToPointMessageProcess struct{}

func (this *PointToPointMessageProcess) sendMessageToTargetUser(message string) (err error) {
	var (
		pointMessage common.PointToPointMessage
		code         int
	)
	err = json.Unmarshal([]byte(message), &pointMessage)
	if err != nil {
		fmt.Printf("json unmarshal point to point message error: %v\n", err)
		return
	}

	// find conn by targetUserName
	clientConn := model.ClientConn{}
	conn, err := clientConn.SearchByUserName(pointMessage.TargetUserName)

	var responseMessage common.ResponseMessage
	responseMessage.Type = common.PointToPointMessageType

	responseMessageData := common.PointToPointMessage{
		SourceUserName: pointMessage.SourceUserName,
		TargetUserName: pointMessage.TargetUserName,
		Content:        pointMessage.Content,
	}

	data, err := json.Marshal(responseMessageData)
	if err != nil {
		code = 200
	} else {
		code = common.ServerError
	}

	err = responseClient(conn, code, string(data), err.Error())
	if err != nil {
		fmt.Printf("point to point communicate, response client error: %v\n", err)
	}

	return
}

func responseClient(conn net.Conn, code int, data string, popErr string) (err error) {
	responseMessage := common.ResponseMessage{
		Type:  common.PointToPointMessageType,
		Code:  code,
		Data:  data,
		Error: popErr,
	}

	responseData, err := json.Marshal(responseMessage)
	if err != nil {
		fmt.Printf("json marshal response message error: %v\n", err)
	}

	dispatcher := utils.Dispatcher{Conn: conn}
	dispatcher.WriteData(responseData)

	return
}
