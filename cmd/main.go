package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"file-storage/internal/container"
)

func main() {
	var err error

	defer func() {
		if r := recover(); r != nil {
			log.Fatalln("Recovered", r)
		}
		if err != nil {
			log.Fatalf("finished with err: %v\n", err)
		}
	}()

	c, err := container.New()
	if err != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go c.WeightsLoader.Watch(ctx)

	err = http.ListenAndServe(fmt.Sprintf(":%d", c.Config.Server.Port), c.HttpHandler)
}
