package model

import (
	"encoding/json"
	"fmt"

	"github.com/garyburd/redigo/redis"
)

// UserDao实例， 全局唯一
var CurrentUserDao *UserDao

type UserDao struct {
	pool *redis.Pool
}

// 初始化一个UserDao结构体实例
func InitUserDao(pool *redis.Pool) (currentUserDao *UserDao) {
	currentUserDao = &UserDao{pool: pool}
	return
}

// 用户id自增
func idIncr(conn redis.Conn) (id int, err error) {
	res, err := conn.Do("incr", "user_id")
	if err != nil {
		fmt.Printf("Id incr error : %v\n", err)
		return
	}
	// 类型断言
	if _, ok := res.(int64); ok {
		id = int(res.(int64))
	}

	return
}

// 根据 用户ID 来获取 用户信息
// 获取成功返回 user 信息，err nil
// 获取失败返回 err，user 为 nil
func (this *UserDao) GetUserById(id int) (user User, err error) {
	conn := this.pool.Get()
	defer conn.Close()

	res, err := redis.String(conn.Do("hget", "users", id))
	if err != nil {
		err = ERROR_USER_DOES_NOT_EXISTS
		return
	}

	err = json.Unmarshal([]byte(res), &user)
	if err != nil {
		fmt.Printf("Unmarshal user info error: %v\n", err)
		return
	}

	return
}

// 根据 用户UserName 来获取 用户信息
// 获取成功返回 user 信息，err nil
// 获取失败返回 err，user 为 nil
func (this *UserDao) GetUserByUserName(userName string) (user User, err error) {
	conn := this.pool.Get()
	defer conn.Close()

	res, err := redis.String(conn.Do("hget", "users", userName))
	if err != nil {
		err = ERROR_USER_DOES_NOT_EXISTS
		return
	}

	err = json.Unmarshal([]byte(res), &user)
	if err != nil {
		fmt.Printf("Unmarshal user info error: %v\n", err)
		return
	}

	return
}

// 注册
// 用户名不能重复
func (this *UserDao) Register(userName, password, passwordConfirm string) (user User, err error) {
	// 判断密码二次验证是否通过
	if passwordConfirm != password {
		err = ERROR_PASSWORD_DOES_NOT_MATCH
		return
	}

	// 用户名不能重复
	_, err = this.GetUserByUserName(userName)
	if err == nil {
		fmt.Printf("User already exists!\n")
		err = ERROR_USER_ALREADY_EXISTS
		return
	}

	conn := this.pool.Get()
	defer conn.Close()

	id, err := idIncr(conn)
	if err != nil {
		return
	}

	user = User{id, userName, password}
	info, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("Marshal user info error: %v\n", err)
		return
	}

	_, err = conn.Do("hset", "users", userName, info)
	if err != nil {
		fmt.Printf("set user info to redis error: %v\n", err)
		return
	}

	return
}

func (this *UserDao) Login(userName, password string) (user User, err error) {
	user, err = this.GetUserByUserName(userName)
	if err != nil {
		fmt.Printf("Get user by userName error: %v\n", err)
		return
	}

	if user.Password != password {
		err = ERROR_USER_PASSWORD
		return
	}

	return
}
