package main

import (
	"context"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	ctx := context.Background()
	app, err := injectApi(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = app.initApi()
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
}
