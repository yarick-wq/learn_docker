package main

import (
	"docker/handlers"
	memorystorage "docker/memory_storage"
	"errors"
	"fmt"
	"net/http"
)

func main() {
	stor := memorystorage.NewStorage()
	hand := handlers.NewHandler(stor)
	router := http.NewServeMux()
	router.HandleFunc("POST /employees", hand.CreateEmployee)
	router.HandleFunc("GET /employees", hand.CollectEmployee)
	router.HandleFunc("DELETE /employees", hand.DeleteEmployee)
	err := http.ListenAndServe(":8080", router)
	if errors.Is(err, http.ErrServerClosed) {
		return
	} else if err != nil {
		fmt.Println("err", err)
		return
	}
}
