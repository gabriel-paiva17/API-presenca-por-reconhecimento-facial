package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"myproject/group"
	"myproject/user"
	"myproject/utils"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func CreateServer() {

	// definindo variaves de ambiente

	env := os.Getenv("APP_ENV")
    if env == "" {
        env = "development"
    }
    utils.LoadEnv(env)

	mongoURI := os.Getenv("MONGO_URI")
    secretKey := os.Getenv("SECRET_KEY")
    serverAddr := os.Getenv("SERVER_ADDR")

	// conectando ao db

	client, ctx, cancel, err := connectDB(mongoURI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal("Failed to disconnect from MongoDB:", err)
		}
		cancel()
	}()

	// configurando router

	r := mux.NewRouter()

	userRepo := user.NewUserRepository(client, "web_presenca", "users")
	groupRepo := group.NewGroupRepository(client, "web_presenca", "groups")

	userService := user.NewUserService(userRepo, secretKey)
	groupService := group.NewGroupService(groupRepo)

	userController := user.NewUserController(userService)
	groupController := group.NewGroupController(groupService)

	r.HandleFunc("/auth/register", userController.CreateUserHandler).Methods("POST")
	r.HandleFunc("/auth/login", userController.LoginUserHandler).Methods("POST")
	r.HandleFunc("/auth/logout", utils.Authenticate(userController.LogoutUserHandler)).Methods("POST")
	
	r.HandleFunc("/grupos", utils.Authenticate(groupController.GetAllGroupsByUserID)).Methods("GET")
	r.HandleFunc("/grupos", utils.Authenticate(groupController.CreateGroupHandler)).Methods("POST")

	// configurando server e CORS

	cors := configureCORS()

    s := &http.Server{
        Addr:         serverAddr,
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