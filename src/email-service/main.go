package main

import (
	"os"
	"fmt"
)

func main() {
	err :=CheckEnvConf()

	if err!=nil{
		fmt.Printf("Failed to check conf:%v\n",err)
		os.Exit(3)
		return
	}

	DumpEnvConf()

	StartServer()
}
