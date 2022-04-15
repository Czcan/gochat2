package main

import (
	"fmt"
	"gochat2/config"
	"gochat2/server/model"
	"gochat2/server/process"
	"net"
	"time"
)

func init() {
	// 初始化Redis连接池，全局唯一
	redisInfo := config.Configuration.RedisInfo
	fmt.Println("redisInfo :", redisInfo)
	initRedisPool(redisInfo.MaxIdle, redisInfo.MaxActive, time.Second*(redisInfo.IdleTimeOut), redisInfo.Host)

	// 创建UserDao， 用于操作用户信息
	// UserDao实例，全局唯一
	model.CurrentUserDao = model.InitUserDao(pool)
}

func dialogue(conn net.Conn) {
	defer conn.Close()
	processor := process.Processor{Conn: conn}
	processor.MainProcess()
}

func main() {
	fmt.Println("Server is already!")

	serverInfo := config.Configuration.ServerInfo
	fmt.Println("serverInfo :", serverInfo)
	listener, err := net.Listen("tcp", serverInfo.Host)
	defer listener.Close()
	if err != nil {
		fmt.Printf("some error when run server, error: %v", err)
	}

	for {
		fmt.Println("Waiting for client!")

		conn, err := listener.Accept()
		defer conn.Close()
		if err != nil {
			fmt.Printf("some error when accepted server, error: %v", err)
		}

		// 一旦连接成功， 则启动一个协程 来负责服务端 与 该客户端的 通讯
		go dialogue(conn)
	}
}
