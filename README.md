# 项目简介 
Golang网络编程 实现的 多人在线聊天系统
使用 gorutinue 来达到高并发的效果
使用 redis 来保存用户的 注册信息

## 项目概况
1. 用户注册，登陆
2. 显示在线用户列表
3. 发送群消息（目前是发送给所有在线用户）
4. 私发消息给某个用户
5.按照消息的类型(info, notice, warn, error, success) 使用不同的颜色打印消息（Unix 和 window 均支持）

## 项目结构
D:.
│  client.exe
│  go.mod
│  go.sum
│  gochat2.txt
│  main.exe
│  README.md
│  
├─.vscode
│      launch.json
│      
├─client       // 客户端代码
│  │  main.go   // 主Menu
│  │  
│  ├─logger
│  │      logger.go  // 自定义日志打印
│  │      
│  ├─model
│  │      user.go   // currentUser
│  │      
│  ├─process
│  │      messageProcess.go  //  客户端 对登陆成功后三种操作的处理(显示在线用户列表，群发消息，私发消息)
│  │      serverProcess.go  //  客户端 对 服务端返回数据 的处理
│  │      userProcess.go    //  用户菜单 及 用户登陆，注册的处理
│  │      
│  └─utils
│          utils.go        // 调度器, 封装 客户端 向Conn 读取， 写入数据(特点：数据传输增加了消息长度，以便判断消息是否完整)
│          
├─common
│  └─message
│          message.go      // 定义了一些全局类型，状态定义及对应的HTTP Code， data model
│          
├─config
│      config.go           // 定义了ServerInfo， RedisInfo等model，实现对 config.json配置文件 的 读取
│      config.json         // 配置 服务器，redis等
│      
└─server    // 服务端代码
    ├─main
    │      main.go          
    │      redis.go        // redis连接池，全局唯一， 初始化initRedisPool(), redis存取用户注册信息
    │      
    ├─model
    │      clientConn.go   // map[userID]ConnInfo 存储连接信息
    │      error.go        // 自定义一些逻辑错误
    │      user.go  
    │      userDao.go      // 全局唯一UserDao{pool *redis Pool}, 操作用户信息, initUserDao(),
    │      
    ├─process   // 处理与客户端的连接，收发消息
    │      groupMessageProcess.go           // 服务端对 群发消息请求的处理   
    │      onlineInfoProcess.go             // 服务端对 显示在线用户列表的处理 
    │      pointToPointMessageProcess.go    //服务端对 私发消息请求的处理 
    │      processor.go     // 消息处理器入口
    │      userProcess.go                   //服务端对 用户登陆， 注册请求的处理 
    │      
    └─utils
            utils.go    // 调度器, 封装 服务端 向Conn 读取， 写入数据(特点：数据传输增加了消息长度，以便判断消息是否完整)


# 本地运行本项目（Unix系统下)
下载到本地的GOPATH目录下， 这是Golang项目，需要你本地有配置Golang环境

cd ${GOPATH}/src
git clone git@github.com:Czcan/gochat2.git

## 导包
本项目用go mod的方式
根目录(即 ../gochat2/)下执行以下命令
go mod init gochat2  // gochat2是module Name，最好跟project Name一样，以便后续如果转回普通vendor模式
go mod tidy

## 编译和运行

编译并运行服务端
go build -o server gochat2/server/main
./server  
// 注意： 生成的.exe文件在gochat2根目录下，这样运行，否则需要./exe文件路径

编译并运行客户端
go build -o client gochat2/client
./client  

这样就可以了，你就可以在本地运行项目了
