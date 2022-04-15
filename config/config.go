package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"
)

type configuration struct {
	ServerInfo serverInfo
	RedisInfo  redisInfo
}

type serverInfo struct {
	Host string
}

type redisInfo struct {
	Host        string
	MaxIdle     int // 最大等待连接数量,0表示没限制
	MaxActive   int // 最大连接数，0表示无限制，一般设置为可能的最大并发量
	IdleTimeOut time.Duration
}

var Configuration = configuration{}

func init() {
	filePath := path.Join(os.Getenv("GOPATH"), "src/github.com/Czcan/gochat2/config/config.json")
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		fmt.Printf("Open file error : %v\n", err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Configuration)
	if err != nil {
		fmt.Printf("Init config error: %v\n", err)
	}
}
