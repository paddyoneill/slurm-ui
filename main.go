package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/paddyoneill/slurm-ui/internal/db"
	"github.com/paddyoneill/slurm-ui/internal/server"
)

func main() {
	_ = godotenv.Load()
	database, err := db.Open("./local.db")
	if err != nil {
		panic(err)
	}

	s := server.New(database)
	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			panic(err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	<-shutdown

	if err := s.Shutdown(ctx); err != nil {
		panic(err)
	}
}
