package main

import (
	"log"
	"readyGo/boxops"
	"readyGo/cmd"
	"readyGo/scaler"
)

func main() {
	// Muruga bless me.

	ops := boxops.New("../box")
	m, err := scaler.New(ops, "configs/scalers.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(len(m))
	log.Println(m.GetScaler("int").GoType)
	cmd.Execute()
}
