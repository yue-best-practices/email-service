package main

import (
	"strconv"
	"os"
	"fmt"
)

var (
	Host = ""   // email service http host
	Port = 8001 // email service http port
	MaxQueueLen=5
	DefaultUser = ""
	DefaultPassword = ""
	DefaultSmtpHost = ""
	DefaultSmtpPort = ""
	DefaultNickName = ""
)

func CheckEnvConf() error{

	p:=os.Getenv("PORT")
	if p!= ""{
		port,_ :=strconv.Atoi(p)
		if port>0 {
			Port=port
		}
	}

	queueLen:=os.Getenv("MAX_QUEUE_LEN")

	if queueLen!=""{
		l,_:=strconv.Atoi(queueLen)
		if l>0{
			MaxQueueLen=l
		}
	}


	defaultUser:=os.Getenv("USER")
	if defaultUser!=""{
		DefaultUser=defaultUser
	}

	defaultPassword:=os.Getenv("PASSWORD")
	if defaultPassword!=""{
		DefaultPassword=defaultPassword
	}

	defaultSmtpHost:=os.Getenv("SMTP_HOST")
	if defaultSmtpHost!=""{
		DefaultSmtpHost=defaultSmtpHost
	}

	defaultSmtpPort:=os.Getenv("SMTP_PORT")
	if defaultSmtpPort!=""{
		DefaultSmtpPort=defaultSmtpPort
	}

	defaultNickName:=os.Getenv("NICK_NAME")
	if defaultNickName!=""{
		DefaultNickName=defaultNickName
	}

	return nil
}


func DumpEnvConf(){
	fmt.Printf("listen host: %s \n",Host)
	fmt.Printf("listen port: %v \n",Port)
	fmt.Printf("MaxQueueLen: %v \n",MaxQueueLen)
	fmt.Printf("DefaultUser: %v \n",DefaultUser)
	fmt.Printf("DefaultPassword: %v \n",DefaultPassword)
	fmt.Printf("DefaultSmtpHost: %v \n",DefaultSmtpHost)
	fmt.Printf("DefaultSmtpPort: %v \n",DefaultSmtpPort)
	fmt.Printf("DefaultNickName: %v \n",DefaultNickName)
}

