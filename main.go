package main

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/yuri-potatoq/generic-profile/enrollment"
	"os"
)

type application struct {
}

type Enroll struct {
	ID int `db:"id"`
}

func main() {
	urlExample := "postgres://furry-profile:pass@172.30.42.187:5432/furry-profile?sslmode=disable"
	ctx := context.Background()
	pool, err := sqlx.ConnectContext(ctx, "postgres", urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	r := enrollment.NewEnrollmentRepository(pool)
	s := enrollment.NewEnrollmentService(r)
	s.GetEnrollmentState(nil, 0)

	//var enroll = Enroll{}
	//row := pool.QueryRowxContext(context.Background(), "select ID from enrollments")
	//if err := row.Err(); err != nil {
	//	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	//	os.Exit(1)
	//}
	//
	//err = row.StructScan(&enroll)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(enroll)

}
