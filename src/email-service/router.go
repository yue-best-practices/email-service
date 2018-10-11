package main

import (
	"github.com/urfave/negroni"
	"github.com/gernest/alien"
	"net/http"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

func StartServer(){
	go StartQueue()
	go startRetryQueue()
	setSignalHandler()

	api:=negroni.New()
	api.Use(negroni.NewLogger())

	router:=alien.New()
	router.Put("/putTask/:id",putTask)
	api.UseHandler(router)
	listenParam:=fmt.Sprintf("%s:%d",Host,Port)

	http.ListenAndServe(listenParam,api)
}

func PathParam(r *http.Request, n string) string {
	p := alien.GetParams(r)
	return p.Get(n)
}

func WriteResponse(w http.ResponseWriter,code int,message string)  {
	v:=map[string]interface{}{"code":code, "message":message}
	b,_:=json.Marshal(v)
	w.Header().Set("Content-Type","application/json")
	w.Write(b)
}


func putTask(w http.ResponseWriter,r *http.Request){
	id:=PathParam(r,"id")

	if id==""{
		WriteResponse(w,101,"Invalid Request")
		return
	}


	if r.Body==nil{
		WriteResponse(w,101,"Invalid Request")
		return
	}

	defer r.Body.Close()

	Body,err:=ioutil.ReadAll(r.Body)

	if err!=nil{
		WriteResponse(w,102,"Read Body Failed")
		return
	}

	//解析post 参数
	var f interface{}

	json.Unmarshal(Body,&f)
	m:=f.(map[string]interface{})

	user :=DefaultUser
	_user:=m["user"]

	password:=DefaultPassword
	_password:=m["password"]

	smtpHost:=DefaultSmtpHost
	_smtpHost:=m["smtpHost"]

	smtpPort:=DefaultSmtpPort
	_smtpPort:=m["smtpPort"]

	if _user != nil{
		user=_user.(string)
	}

	if _password!=nil{
		password=_password.(string)
	}

	if _smtpHost!=nil{
		smtpHost=_smtpHost.(string)
	}

	if _smtpPort!=nil{
		smtpPort=_smtpPort.(string)
	}

	var to []string

	_to:=m["to"]

	if _to !=nil{
		if reflect.TypeOf(_to).String() =="string" {
			to=[]string{_to.(string)}
		}else{
			for _,t:=range _to.([] interface{}){
				to=append(to,t.(string))
			}
		}
	}


	nickname:=DefaultNickName
	_nickname:=m["nickname"]

	if _nickname!=nil{
		nickname=_nickname.(string)
	}

	var subject string
	_subject:=m["subject"]

	if _subject!=nil{
		subject=_subject.(string)
	}

	var notifyUrl string

	_notifyUrl:=m["notifyUrl"]
	if _notifyUrl!=nil{
		notifyUrl=_notifyUrl.(string)
	}


	var body string
	_body:=m["body"]

	if _body!=nil{
		body=_body.(string)
	}

	contentType:="Content-Type: text/plain; charset=UTF-8"

	_contentType:=m["contentType"]

	if _contentType!=nil{
		contentType=_contentType.(string)
	}

	isSSL:=false
	_isSSL:=m["isSSL"]
	if _isSSL!=nil{
		isSSL=_isSSL.(bool)
	}



	if user=="" || password=="" || smtpHost=="" || smtpPort=="" || len(to)==0 || nickname =="" || subject == "" || body==""{
		WriteResponse(w,104,"Invalid Parameter")
		return
	}

	msg:=[]byte("To: "+strings.Join(to,",")+"\r\nFrom: "+nickname+
		"<"+user+">\r\nSubject: "+subject+"\r\n"+contentType+"\r\n\r\n"+body)

	emailQueue<-&EmailObj{id,user,password,smtpHost,smtpPort,to,nickname,subject,body,contentType,notifyUrl,msg,isSSL}

	WriteResponse(w,200,"OK")
}