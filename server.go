package main

import (
	"log"
	"myproject/group"
	"myproject/user"
	"myproject/utils"
	"myproject/session"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

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
	dbName := os.Getenv("MONGO_DBNAME")

	// conectando ao db

	client, ctx, cancel, err := utils.ConnectDB(mongoURI)
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

	userRepo := user.NewUserRepository(client, dbName, "users")
	groupRepo := group.NewGroupRepository(client, dbName, "groups")
	sessionRepo := session.NewSessionRepository(client, dbName, "sessions")

	userService := user.NewUserService(userRepo, secretKey)
	groupService := group.NewGroupService(groupRepo, userRepo)
	sessionService := session.NewSessionService(sessionRepo, groupRepo)

	userController := user.NewUserController(userService)
	groupController := group.NewGroupController(groupService)
	sessionController := session.NewSessionController(sessionService)

	r.HandleFunc("/auth/register", userController.CreateUserHandler).Methods("POST")
	r.HandleFunc("/auth/login", userController.LoginUserHandler).Methods("POST")
	r.HandleFunc("/auth/logout", utils.Authenticate(userController.LogoutUserHandler)).Methods("POST")

	r.HandleFunc("/grupos", utils.Authenticate(groupController.GetAllGroupsByUserID)).Methods("GET")
	r.HandleFunc("/grupos/criar", utils.Authenticate(groupController.CreateGroupHandler)).Methods("POST")
	r.HandleFunc("/grupos/{nome-do-grupo}/detalhes", utils.Authenticate(groupController.GetGroupDetails)).Methods("GET")
	r.HandleFunc("/grupos/{nome-do-grupo}/detalhes/adicionar", utils.Authenticate(groupController.AddMemberToGroup)).Methods("POST")

	r.HandleFunc("/grupos/{nome-do-grupo}/sessoes/iniciar", utils.Authenticate(sessionController.StartNewSession)).Methods("POST")
	r.HandleFunc("/grupos/{nome-do-grupo}/sessoes/{nome-da-sessao}/validar-face", utils.Authenticate(sessionController.ValidateFace)).Methods("PUT")
	r.HandleFunc("/grupos/{nome-do-grupo}/sessoes/{nome-da-sessao}/encerrar", utils.Authenticate(sessionController.EndSession)).Methods("POST")
	r.HandleFunc("/grupos/{nome-do-grupo}/sessoes/em-andamento", utils.Authenticate(sessionController.GetActiveSessions)).Methods("GET")
	r.HandleFunc("/grupos/{nome-do-grupo}/sessoes/encerradas", utils.Authenticate(sessionController.GetEndedSessions)).Methods("GET")
	r.HandleFunc("/grupos/{nome-do-grupo}/sessoes/{nome-da-sessao}/detalhes", utils.Authenticate(sessionController.GetSessionDetails)).Methods("GET")
	r.HandleFunc("/grupos/{nome-do-grupo}/sessoes/{nome-da-sessao}/deletar", utils.Authenticate(sessionController.DeleteOneSession)).Methods("DELETE")
	r.HandleFunc("/grupos/{nome-do-grupo}/sessoes/em-andamento/deletar", utils.Authenticate(sessionController.DeleteAllActiveSessionsOfAGroup)).Methods("DELETE")
	r.HandleFunc("/grupos/{nome-do-grupo}/sessoes/encerradas/deletar", utils.Authenticate(sessionController.DeleteAllEndedSessionsOfAGroup)).Methods("DELETE")

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

func configureCORS() func(http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)
}