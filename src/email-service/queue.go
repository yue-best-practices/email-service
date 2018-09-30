package main

import (
	"os"
	"fmt"
	"os/signal"
	"syscall"
	"net/smtp"
)

type EmailObj struct{
	id string
	user string
	password string
	smtpHost string
	smtpPort string
	to []string
	nickname string
	subject string
	body string
	contentType string
	msg []byte
}

var(
	emailQueue chan *EmailObj
	StopQueue = false
)



func sendEmail(email *EmailObj){
	auth:=smtp.PlainAuth("",email.user,email.password,email.smtpHost)
	err:=smtp.SendMail(email.smtpHost+":"+email.smtpPort,auth,email.user,email.to,email.msg)
	if err!=nil{
		fmt.Printf("[task:%s] Send Mail Error:%v,",email.id,err)
	}else{
		fmt.Printf("[task:%s] Send mail Success!",email.id)
	}
}


func setSignalHandler(){
	sign:=make(chan os.Signal,1)
	signal.Notify(sign,syscall.SIGTERM,syscall.SIGINT,syscall.SIGQUIT)
	go func() {
		s:=<-sign
		fmt.Printf("==signal:%v==\n",s)
		StopQueue=true
		emailQueue<-nil
	}()
}


func StartQueue()  {

	emailQueue=make(chan *EmailObj,MaxQueueLen)

	setSignalHandler()

	for !StopQueue{
		email :=<-emailQueue
		if email==nil{
			continue
		}
		sendEmail(email)
	}

	fmt.Println("====emailQueue stop!")
	os.Exit(0)
}