package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"go-http-test/db"
	"go-http-test/handler"
	"go-http-test/mdw"

	"github.com/jackc/pgx/v5"
)

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	db := db.New(conn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connexion à la base de données établie avec succès")

	server := handler.New(true, db)
	mdw := mdw.New(true)


	http.Handle("/", mdw.Logger(http.HandlerFunc(server.Home)))
	http.Handle("/auth", mdw.Logger(mdw.BasicAuth(http.HandlerFunc(server.Login))))
	http.Handle("/register", mdw.Logger(http.HandlerFunc(server.Register)))
	http.Handle("/users", mdw.Logger(http.HandlerFunc(server.GetAllUsers)))
	http.Handle("/users/search", mdw.Logger(http.HandlerFunc(server.SearchUsers)))
	http.Handle("/users/login", mdw.Logger(http.HandlerFunc(server.Login)))

	//group routes for admin dashboard
	http.Handle("/admin", mdw.Logger(mdw.BasicAuth(http.HandlerFunc(server.AdminDashboard))))

	fmt.Println("Server starting on port 8080...")	

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}