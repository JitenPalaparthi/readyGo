package main

import (
	"flag"
	"fmt"
	"log"
	"readyGo/generate"
	"readyGo/generate/template"
)

func main() {

	fmt.Println("Hello Muruga")

	configFile := flag.String("f", "config.json", "pass -f followed by the json file")

	flag.Parse()

	// Todo: implement json schema and its validation related things.

	log.Println("Loading all templates into in-memory")
	tm, err := template.New("templates")

	if err != nil {
		log.Fatal("error occured loading templates..---->", err)
	}

	log.Println("Generating files and dependencies based on config file.Here is the config file :", *configFile)

	tg, err := generate.New(configFile)
	if err != nil {
		log.Fatal("seems , things went wrong.. -->", err)

	}
	tg.Gen = tm // Assign Template Map to the Template Generator

	log.Println("Generating main.go file in the root directory")

	err = tg.CreateMain(tm["main"]) // Generate main.go
	if err != nil {
		log.Fatal("seems , things went wrong.. -->", err)

	}
	log.Println("Generating all model files in model directory")

	err = tg.GenerateAllModelFiles(tm["models"])
	if err != nil {
		log.Fatal("seems , things went wrong.. -->", err)

	}

}
