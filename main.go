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

	configFile := flag.String("f", "configs/config.json", "pass -f followed by the json file.By default it takes pre defined config file")

	projectType := flag.String("project-type", "http", "pass project-type followed by the type of the project http | grpc | cloudEvents | cli.The default type is http.")

	dbType := flag.String("db-type", "mongo", "pass db-type followed by the type of the database mongo | sql | none.By default it takes mongo is the database type ")

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
	tg.Gen = tm           //Assign Template Map to the Template Generator
	tg.Type = projectType // TODO define the flow based on project type
	tg.DBType = dbType    // TODO define the flow based on database type

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
