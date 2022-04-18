package main

import (
	"fmt"
	"gochat2/client/logger"
	"gochat2/client/process"
)

func main() {
	var (
		key             int
		loop            bool = true
		userName        string
		password        string
		passwordConfirm string
	)

	for loop {
		logger.Info("\n----------------Welcome to the chat room----------------\n")
		logger.Info("\t\tSelect the options\n")
		logger.Info("\t\t\t1. Sign in\n")
		logger.Info("\t\t\t2. Sign up\n")
		logger.Info("\t\t\t3. Exit the system\n")

		fmt.Scanf("%v\n", &key)

		switch key {
		case 1:
			logger.Info("sign in please\n")
			logger.Notice("UserName: \n")
			fmt.Scanf("%v\n", &userName)
			logger.Notice("Password: \n")
			fmt.Scanf("%v\n", &password)

			p := &process.UserProcess{}
			err := p.Login(userName, password)
			if err != nil {
				logger.Error("Login failed, error: %v\n", err)
			} else {
				logger.Success("Login succeed!!!")
			}
		case 2:
			logger.Info("sign up please\n")
			logger.Notice("UserName: \n")
			fmt.Scanf("%s\n", &userName)
			logger.Notice("Password: \n")
			fmt.Scanf("%s\n", &password)
			logger.Notice("Password Confirm: \n")
			fmt.Scanf("%s\n", &passwordConfirm)

			p := &process.UserProcess{}
			err := p.Register(userName, password, passwordConfirm)
			if err != nil {
				logger.Error("Register failed, error: %v\n", err)
			}
		case 3:
			logger.Warn("Exit ...\n")
			loop = false // this is equal to 'os.Exit(0)'
		default:
			logger.Error("Select is invalid!!!\n")
		}
	}

}
