package model

import "net"

type ClientConn struct{}

type ConnInfo struct {
	Conn     net.Conn
	UserName string
}

var ClientConnMap map[int]ConnInfo

func init() {
	ClientConnMap = make(map[int]ConnInfo)
}

func (cc ClientConn) Save(userID int, userName string, userConn net.Conn) {
	ClientConnMap[userID] = ConnInfo{Conn: userConn, UserName: userName}
}

func (cc ClientConn) Del(userConn net.Conn) {
	for id, connInfo := range ClientConnMap {
		if userConn == connInfo.Conn {
			delete(ClientConnMap, id)
		}
	}
}

func (cc ClientConn) SearchByUserName(userName string) (conn net.Conn, err error) {
	// for _, connInfo := range ClientConnMap {
	// 	if userName == connInfo.UserName {
	// 		return connInfo.Conn, nil
	// 	}
	// }

	// return

	user, err := CurrentUserDao.GetUserByName(userName)
	if err != nil {
		return
	}

	conn = ClientConnMap[user.ID].Conn

	return
}
