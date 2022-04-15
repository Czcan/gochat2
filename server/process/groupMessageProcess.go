package process

import (
	"encoding/json"
	"fmt"
	common "gochat2/common/message"
	"gochat2/server/model"
	"gochat2/server/utils"
)

type SendGroupMessageProcess struct{}

func (this *SendGroupMessageProcess) sendGroupMessage(message string) (err error) {
	var groupMessage common.UserSendGroupMessage

	err = json.Unmarshal([]byte(message), groupMessage)
	if err != nil {
		fmt.Printf("json unmarshal group message error: %v\n", err)
	}

	// group message sender
	sourceUserName := groupMessage.UserName

	var toClientMessage common.ResponseMessage
	toClientMessage.Type = common.SendGroupMessageToClientType
	toClientMessage.Data = message

	data, err := json.Marshal(toClientMessage)
	if err != nil {
		fmt.Printf("json marshal group message to client error : %v\n", err)
	}

	for id, connInfo := range model.ClientConnMap {
		// continue the group message sender
		if connInfo.UserName == sourceUserName {
			continue
		}

		fmt.Printf("client id: %v\n", id)

		dispatcher := utils.Dispatcher{Conn: connInfo.Conn}

		err = dispatcher.WriteData(data)
		if err != nil {
			fmt.Printf("send group message error: %v\n", err)
		} else {
			fmt.Printf("send succeed!!\n")
		}
	}

	return
}
