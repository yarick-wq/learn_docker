package main

import (
	"context"
	"docker/handlers"
	"docker/postgres"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	str := os.Getenv("CONN_STRING")
	if str == "" {
		fmt.Println("не введена строка подключения к базе данных")
		return
	}
	pool, err := pgxpool.New(context.Background(), str)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	defer pool.Close()
	if err := pool.Ping(context.Background()); err != nil {
		fmt.Println("err", err)
		return
	}
	st := postgres.NewStorage(pool)
	//stor := memorystorage.NewStorage()
	hand := handlers.NewHandler(st)
	router := http.NewServeMux()
	router.HandleFunc("POST /employees", hand.CreateEmployee)
	router.HandleFunc("GET /employees", hand.CollectEmployee)
	router.HandleFunc("DELETE /employees", hand.DeleteEmployee)

	err = http.ListenAndServe(":8080", router)
	if errors.Is(err, http.ErrServerClosed) {
		return
	} else if err != nil {
		fmt.Println("err", err)
		return
	}

}
