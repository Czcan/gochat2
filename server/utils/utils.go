package utils

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	common "gochat2/common/message"
	"net"
)

type Dispatcher struct {
	Conn net.Conn
	Buf  [10240]byte
}

func (dispatcher Dispatcher) ReadData() (message common.Message, err error) {
	buf := dispatcher.Buf

	// 读取消息长度(期望长度)
	_, err = dispatcher.Conn.Read(buf[:4])
	if err != nil {
		return
	}

	var dataLen uint32
	// []byte  -> uint32
	dataLen = binary.BigEndian.Uint32(buf[0:4])

	// 读取消息本身
	n, err := dispatcher.Conn.Read(buf[:dataLen])
	if err != nil {
		fmt.Printf("server read login data error: %v\n", err)
		return
	}

	// 对比 消息本身长度 与 期望长度 是否一致
	if n != int(dataLen) {
		err = errors.New("login message length error")
		return
	}

	// 从 conn中 解析数据， 并且放到 message中， 这里注意要传 message 地址
	err = json.Unmarshal(buf[:dataLen], &message)
	if err != nil {
		fmt.Printf("Unmarshal data from conn error: %v\n", err)
		return
	}

	return
}

func (dispatcher Dispatcher) WriteData(data []byte) (err error) {
	var dataLen uint32
	dataLen = uint32(len(data))

	var bytes [4]byte
	// uint32  -> []byte
	binary.BigEndian.PutUint32(bytes[:], dataLen)

	// 将 消息长度 发送给 客户端
	_, err = dispatcher.Conn.Write(bytes[:])
	if err != nil {
		fmt.Printf("send data len to server error: %v\n", err)
		return
	}

	// 将 消息 发送给 客户端
	_, err = dispatcher.Conn.Write(data)
	if err != nil {
		fmt.Printf("send data to server error: %v\n", err)
		return
	}

	return
}
