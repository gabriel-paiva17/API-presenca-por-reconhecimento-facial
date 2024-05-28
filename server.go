package main

import (
	"context"
	"log"
	"net/http"
	"time"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"myproject/user"
)

func CreateServer() {

	client, ctx, cancel, err := connectDB("mongodb://localhost:27017")
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal("Failed to disconnect from MongoDB:", err)
		}
		cancel()
	}()

	r := mux.NewRouter()

	userRepo := user.NewMongoUserRepository(client, "web_presenca", "users")
	userService := user.NewUserService(userRepo)
	userController := user.NewUserController(userService)

	r.HandleFunc("/auth/register", userController.CreateUserHandler).Methods("POST")

	cors := configureCORS()

    s := &http.Server{
        Addr:         "localhost:8080",
        Handler:      cors(r),
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }
	
	log.Fatal(s.ListenAndServe())

}


func connectDB(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		cancel()
		return nil, nil, nil, err
	}
	
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		cancel()
		return nil, nil, nil, err
	}

	fmt.Println("Conectado ao MongoDB!")
	return client, ctx, cancel, nil
}

func configureCORS() func(http.Handler) http.Handler {
    return handlers.CORS(
        handlers.AllowedOrigins([]string{"http://localhost:5173"}),
        handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
        handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
    )
}