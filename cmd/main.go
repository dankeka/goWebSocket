package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	s "github.com/dankeka/goWebSocket"
	hand "github.com/dankeka/goWebSocket/pkg/handler"
)

func main() {
	srv := new(s.Server)
	h := new(hand.Handler)

	go func() {
		if err := srv.Run("8080", h.InitRouters()); err != nil {
			log.Fatal("error server: ", err.Error())
		}
	}()

	fmt.Println("Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<- quit

	fmt.Println("Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		fmt.Printf("error %s", err.Error())
	}
}