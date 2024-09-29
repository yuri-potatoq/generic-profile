package main

import (
	"context"
	"fmt"
	pgx "github.com/jackc/pgx/v5"
	"os"
)

type application struct {
}

func main() {
	urlExample := "postgres://furry-profile:pass@172.30.42.187:5432/furry-profile"
	conn, err := pgx.Connect(context.Background(), urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var ID int
	err = conn.QueryRow(context.Background(), "select ID from enrollments").Scan(&ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(ID)

}
