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
	)
	err = json.Unmarshal([]byte(message), &pointMessage)
	if err != nil {
		fmt.Printf("json unmarshal point to point message error: %v\n", err)
		return
	}

	// find conn by targetUserName
	clientConn := model.ClientConn{}
	conn, err := clientConn.SearchByUserName(pointMessage.TargetUserName)
	if err != nil {
		return
	}

	var responseMessage common.ResponseMessage
	responseMessage.Type = common.PointToPointMessageType

	responseMessageData := common.PointToPointMessage{
		SourceUserName: pointMessage.SourceUserName,
		TargetUserName: pointMessage.TargetUserName,
		Content:        pointMessage.Content,
	}

	data, err := json.Marshal(responseMessageData)
	if err != nil {
		return
	}

	responseMessage.Data = string(data)
	responseMessage.Code = 200

	responseData, _ := json.Marshal(responseMessage)
	if err != nil {
		return
	}

	dispatcher := utils.Dispatcher{Conn: conn}
	err = dispatcher.WriteData(responseData)

	return
}

func (this *PointToPointMessageProcess) responseClient(conn net.Conn, code int, data string, popErr string) (err error) {
	responseMessage := common.ResponseMessage{
		Code:  code,
		Type:  common.PointToPointMessageType,
		Error: popErr,
		Data:  data,
	}

	responseData, err := json.Marshal(responseMessage)
	if err != nil {
		fmt.Printf("some error when generate response message, error: %v", err)
	}

	dispatcher := utils.Dispatcher{Conn: conn}
	err = dispatcher.WriteData(responseData)

	return
}
