package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"platzi.com/go/rest-ws/handlers"
	"platzi.com/go/rest-ws/middleware"
	"platzi.com/go/rest-ws/server"
)

func main() {
	dir := flag.String("d", ".", "dir")
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*dir)))
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	s, err := server.NewServer(context.Background(), &server.Config{
		Port:        PORT,
		JWTSecret:   JWT_SECRET,
		DatabaseUrl: DATABASE_URL,
	})

	if err != nil {
		log.Fatal(err)
	}

	s.Start(BindRoutes)
}

func BindRoutes(s server.Server, r *mux.Router) {
	api := r.PathPrefix("/api/v1").Subrouter()

	api.Use(middleware.CheckAuthMiddleware(s))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// do NOT do this. (see below)
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	//r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) { http.FileServer(http.Dir("./views")) }).Methods(http.MethodGet)
	//r.HandleFunc("/", handlers.HomeHandler(s)).Methods(http.MethodGet)

	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/login", handlers.LoginHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/me", handlers.MeHandler(s)).Methods(http.MethodGet)
	api.HandleFunc("/posts", handlers.InsertPostHandler(s)).Methods(http.MethodPost)
	api.HandleFunc("/posts/{id}", handlers.GetPostByIdHandler(s)).Methods(http.MethodGet)
	api.HandleFunc("/posts/{id}", handlers.UpdatePostHandler(s)).Methods(http.MethodPut)
	api.HandleFunc("/posts/{id}", handlers.DeletePostHandler(s)).Methods(http.MethodDelete)
	r.HandleFunc("/posts", handlers.ListPostHandler(s)).Methods(http.MethodGet)

	r.HandleFunc("/ws", s.Hub().HandleWebSocket)
}
