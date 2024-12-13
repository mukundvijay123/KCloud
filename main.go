package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/mukundvijay123/KCloud/router"
)

func main() {
	conStr := "host=localhost port=5432 user=postgres password=123456789 dbname=kcloud sslmode=disable"
	db1, err := sql.Open("postgres", conStr)
	if err != nil {
		log.Fatal("error connecting to postgres: ", err)
	}
	defer db1.Close()

	err = db1.Ping()
	if err != nil {
		log.Fatal("No db connection")
	}
	fmt.Println("Connected to DB!!")
	server := router.NewAPIServer(":8080", db1)
	err = server.Run()
	if err != nil {
		log.Fatal("Not Working")
	}
}
