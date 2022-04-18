package process

import (
	"fmt"
	common "gochat2/common/message"
	"gochat2/server/model"
	"gochat2/server/utils"
	"io"
	"net"
)

type Processor struct {
	Conn net.Conn
}

// 处理消息
// 根据消息类型， 选择相对于的 处理方式
func (this *Processor) messageProcess(message common.Message) (err error) {
	fmt.Printf("message from client type: %v\n", message.Type)

	switch message.Type {
	case common.LoginMessageType:
		p := UserProcess{Conn: this.Conn}
		err = p.UserLogin(message.Data)
		if err != nil {
			fmt.Printf("some error when user login : %v\n", err)
		}
	case common.RegisterMessageType:
		p := UserProcess{Conn: this.Conn}
		err = p.UserRegister(message.Data)
		if err != nil {
			fmt.Printf("some error when user register : %v\n", err)
		}
	case common.PointToPointMessageType:
		fmt.Println("point to point communicate!!")
		p := PointToPointMessageProcess{}
		err = p.sendMessageToTargetUser(message.Data)
		if err != nil {
			fmt.Printf("some error when point to point communicate: %v\n", err)
		}
		// var code int
		// if err != nil {
		// 	code = 400
		// } else {
		// 	code = 100
		// }

		// // responseClient(conn net.Conn, code int, data string, err error) {
		// err := p.responseClient(this.Conn, code, "", err.Error())
		// if err != nil {
		// 	fmt.Printf("some err when popmessage: %v", err)
		// }
	case common.UserSendGroupMessageType:
		fmt.Println("user send group message!")
		p := SendGroupMessageProcess{}
		err = p.sendGroupMessage(message.Data)
		if err != nil {
			fmt.Printf("some error when send group message: %v\n", err)
		}
	case common.ShowAllOnlineUsersType:
		fmt.Println("yes in this!!!!!!!!!!!!!!!!!!!!!!!!!!")
		p := OnlineInfoProcess{this.Conn}
		err = p.showAllOnlineUserList()
		if err != nil {
			fmt.Printf("some error when get online user list: %v\n", err)
		}
	default:
		fmt.Println("Other type!")
	}

	return
}

// 处理和用户之间的 通讯
func (this *Processor) MainProcess() {
	// 循环读取 来自 客户端的 消息
	for {
		dispatcher := utils.Dispatcher{Conn: this.Conn}

		message, err := dispatcher.ReadData()
		if err != nil {
			if err == io.EOF {
				cc := model.ClientConn{}
				cc.Del(this.Conn)
				fmt.Println("Client closed!!")
				break
			}

			fmt.Printf("get client message error: %v\n", err)
		}

		// 处理来自客户端的消息
		// 根据消息的类型， 选择相对应的 处理方式
		err = this.messageProcess(message)
		if err != nil {
			fmt.Printf("fix client message error: %v\n", err)
			break
		}
	}
}
