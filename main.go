package main

import (
	"go_final_project/pkg/api"
	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
	"log"
)

func main() {
	err := db.Init()
	if err != nil {
		log.Fatal(err)
	}

	api.Init()
	err = server.Run()
	if err != nil {
		return
	}
}
