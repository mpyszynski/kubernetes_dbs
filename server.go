package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	name := os.Getenv("NAME")
	// Mongo
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongo-service:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	// MySQL
	dsn := "root:password@tcp(mysql-service:3306)/test-db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := client.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatal(err)
		}
		pingMySQL := db.Select("*")
		w.Write([]byte(fmt.Sprintf("Hello %v from World, rows affected %v", name, pingMySQL.RowsAffected)))
	})

	server := http.Server{
		Addr:    ":8000",
		Handler: mux,
	}

	fmt.Println("Server is listening on port :8000")
	server.ListenAndServe()
}
