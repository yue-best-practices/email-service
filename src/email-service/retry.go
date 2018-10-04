package main

import (
	"net/http"
	"log"
	"io/ioutil"
	"time"
	"strings"
	"errors"
	"fmt"
)

type NotifyRequest struct {
	notifyUrl string
	externalId string
	status string
	message string
}

var (
	retryQueue chan *NotifyRequest
)


func startRetryQueue(){
	retryQueue=make(chan *NotifyRequest,RetryQueueLen)
	for !StopQueue{
		request :=<-retryQueue

		if request==nil{
			continue
		}
		go retryRequest(request)
	}
}



func retryRequest(nr *NotifyRequest){
	times:=0
	for i:=0;i<len(RetryTimeDuration);i++{   //todo max retryTime to config
		time.Sleep(time.Duration(RetryTimeDuration[i])*time.Second)

		err:=SendRetryRequest(nr)
		times++
		if err!=nil{
			log.Printf("==Retry Request Error:%v==\n",err)
			continue
		}else{
			log.Printf("==Retry Request Success! ExternalId:%s==",nr.externalId)
			break
		}
	}

	log.Printf("==Retry Done! ExternalId:%s,Retry Times:%d==\n",nr.externalId,times)
}


func SendRetryRequest(nr *NotifyRequest) error{
	client:=&http.Client{}
	request,err:=http.NewRequest("POST",nr.notifyUrl,strings.NewReader("externalId="+nr.externalId+"&status="+nr.status+"&message="+nr.message))

	if err!=nil{
		return errors.New(fmt.Sprintf("httpNewRequest Error:%v",err))
	}

	// todo set headers
	request.Header.Set("Content-Type","application/x-www-form-urlencoded")

	response,err:=client.Do(request)
	if err!=nil{
		return errors.New(fmt.Sprintf("clientDo Error:%v",err))
	}



	defer response.Body.Close()

	if response.StatusCode!=200{
		return errors.New(fmt.Sprintf("Invalid StatusCode:%d",response.StatusCode))
	}

	body,e:=ioutil.ReadAll(response.Body)

	if e!=nil{
		return errors.New(fmt.Sprintf("ResponseBody Error:%v",e))
	}
	if string(body)!="success"{
		return errors.New(fmt.Sprintf("Invalid ResponseBody:%s",string(body)))
	}
	return nil
}