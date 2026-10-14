package main

import (
	"context"
	"docker/handlers"
	"docker/postgres"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	fmt.Println("PID", os.Getpid())
	ctx, cansel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cansel()
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
	server := &http.Server{Addr: ":8080", Handler: router}
	errch := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			return
		} else if err != nil {
			errch <- err
			return
		}
	}()
	select {
	case <-ctx.Done():
		fmt.Println("signal received , started stoped program ")
	case problem := <-errch:
		fmt.Println("The problem", problem)
		return
	}
	ctxtime, stop := context.WithTimeout(context.Background(), 5*time.Second)
	defer stop()
	if err := server.Shutdown(ctxtime); err != nil {
		fmt.Println("err", err)
		return
	}
	fmt.Println("Server is stopped")
}
