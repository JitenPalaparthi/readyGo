package main

import (
	"fmt"
	"log"
	"readyGo/generate"
	"readyGo/generate/template"
)

func main() {

	fmt.Println("Hello Muruga")

	tm, err := template.New("templates")
	if err != nil {
		log.Fatal("error occured loading templates..---->", err)
	}

	tg, err := generate.New("config.json")
	if err != nil {
		log.Fatal("seems , things went wrong.. -->", err)

	}
	tg.Gen = tm

	err = tg.CreateMain(tm["main"])

	err = tg.GenerateAllModelFiles(tm["models"])
	fmt.Println(err)

}
