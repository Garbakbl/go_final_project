package server

import (
	"fmt"
	"go_final_project/pkg/api"
	"net/http"
	"os"
)

func Run() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	fmt.Printf("Listening on port %s...\n", port)
	err := http.ListenAndServe(":"+port, api.Router)
	if err != nil {
		fmt.Printf("Error starting server on port %s: %s\n", port, err)
		return err
	}
	return nil
}
