package main

import (
	"flag"
	"fmt"
	"log"
	"readyGo/generate"
	"readyGo/generate/configure"
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

	templateConfig := "configs/template_config.json"
	tc, err := configure.New(&templateConfig)
	fmt.Println(tc, err)

	tg, err := generate.New(configFile, tm, tc)
	if err != nil {
		log.Fatal("seems , things went wrong.. -->", err)
	}
	//tg.Gen = tm           //Assign Template Map to the Template Generator
	tg.Type = projectType // TODO define the flow based on project type
	tg.DBType = dbType    // TODO define the flow based on database type

	err = tg.GenerateAll("http_mongo")
	if err != nil {
		log.Fatal("seems , things went wrong.. -->", err)
	}

}
