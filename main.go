package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"readyGo/generate"
)

func main() {

	fmt.Println("Hello Muruga")

	tg, err := generate.New("config.json")
	if err != nil {
		fmt.Println("seems , things went wrong.. -->", err)
		os.Exit(1)
	}

	templates := make(map[string]string)

	files, err := ioutil.ReadDir("templates")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		content, err := ioutil.ReadFile("templates/" + file.Name())

		if err != nil {
			log.Fatal(err)
		}

		templates[file.Name()] = string(content)
	}

	fmt.Println(templates)

	err = tg.CreateMain(templates["main"])

	err = tg.GenerateAllModelFiles(templates["models"])
	fmt.Println(err)

}
