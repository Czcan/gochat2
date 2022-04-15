package process

import (
	"encoding/json"
	"fmt"
	common "gochat2/common/message"
	"gochat2/server/model"
	"gochat2/server/utils"
	"net"
)

type OnlineInfoProcess struct {
	Conn net.Conn
}

func (this *OnlineInfoProcess) showAllOnlineUserList(message string) (err error) {
	var (
		onlineUserListInfo []common.UserInfo
		code               int
	)

	for _, connInfo := range model.ClientConnMap {
		user, err := model.CurrentUserDao.GetUserByName(connInfo.UserName)
		if err != nil {
			continue
		}

		userInfo := common.UserInfo{ID: user.ID, UserName: user.Name}
		onlineUserListInfo = append(onlineUserListInfo, userInfo)
	}

	data, err := json.Marshal(onlineUserListInfo)
	if err != nil {
		code = common.ServerError
	} else {
		code = 200
	}

	err = this.responseClient(this.Conn, code, string(data), err.Error())
	if err != nil {
		fmt.Printf("response online user list message to client error: %v\n", err)
	}

	return
}

func (this *OnlineInfoProcess) responseClient(conn net.Conn, code int, data string, popErr string) (err error) {
	responseMessage := common.ResponseMessage{
		Type:  common.ShowAllOnlineUsersType,
		Code:  code,
		Data:  data,
		Error: popErr,
	}

	responseData, err := json.Marshal(responseMessage)
	if err != nil {
		fmt.Printf("json marshal online user list response message error: %v\n", err)
	}

	dispatcher := utils.Dispatcher{Conn: this.Conn}
	err = dispatcher.WriteData(responseData)

	return
}
