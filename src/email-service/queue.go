package main

import (
	"os"
	"fmt"
	"os/signal"
	"syscall"
	"net/smtp"
	"log"
	"net/http"
	"strings"
	"io/ioutil"
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
	notifyUrl string
	msg []byte
}

var(
	emailQueue chan *EmailObj
	StopQueue = false
)



func sendEmail(email *EmailObj){
	auth:=smtp.PlainAuth("",email.user,email.password,email.smtpHost)
	err:=smtp.SendMail(email.smtpHost+":"+email.smtpPort,auth,email.user,email.to,email.msg)

	externalId:=email.id
	status:="SUCCESS"
	message:="OK"
	if err!=nil{
		status="FAILED"
		message=err.Error()
	}

	// todo logger
	log.Printf("[externalId:%s] Done! status:[%s], message:[%s]\n",externalId,status,message)

	// todo send result to notifyUrl
	if email.notifyUrl==""{
		return
	}
	client:=&http.Client{}
	request,err:=http.NewRequest("POST",email.notifyUrl,strings.NewReader("externalId="+externalId+"&status="+status+"&message="+message))

	if err!=nil{
		log.Printf("[externalId:%s] httpNewRequest Error! method:%s,notifyUrl:%s,error:%v\n",email.id,"POST",email.notifyUrl,err)
		return
	}

	// todo set headers
	request.Header.Set("Content-Type","application/x-www-form-urlencoded")
	// todo set more ...


	response,err:=client.Do(request)

	if err!=nil{
		log.Printf("[externalId:%s] sendRequest Error! notifyUrl:%s,error:%v\n",email.id,email.notifyUrl,err)
		retryQueue<-&NotifyRequest{email.notifyUrl,externalId,status,message}
		return
	}

	defer response.Body.Close()

	if response.StatusCode!=200{
		log.Printf("[externalId:%s] notifyResponseError statusCode:%v, will retry later.\n",email.id,response.StatusCode)
		retryQueue<-&NotifyRequest{email.notifyUrl,externalId,status,message}
		return
	}

	body,err:=ioutil.ReadAll(response.Body)

	if err!=nil{
		fmt.Printf("[externalId:%s] ResponseBody Error!  error:%v",email.id,err)
		retryQueue<-&NotifyRequest{email.notifyUrl,externalId,status,message}
		return
	}

	if string(body)!="success"{
		log.Printf("[externalId:%s] notifyResponseError body:%s,will retry later.\n",email.id,string(body))
		retryQueue<-&NotifyRequest{email.notifyUrl,externalId,status,message}
		return
	}

	log.Printf("[externalId:%s] notifyResponse Success!\n",email.id)

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