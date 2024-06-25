package main

import (
	"log"
	"myproject/api"
	"myproject/utils"
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

	userRepo := api.NewUserRepository(client, dbName, "users")
	groupRepo := api.NewGroupRepository(client, dbName, "groups")
	sessionRepo := api.NewSessionRepository(client, dbName, "sessions")

	userService := api.NewUserService(userRepo, groupRepo, sessionRepo, secretKey)
	groupService := api.NewGroupService(groupRepo, userRepo, sessionRepo)
	sessionService := api.NewSessionService(sessionRepo, groupRepo)

	userController := api.NewUserController(userService)
	groupController := api.NewGroupController(groupService)
	sessionController := api.NewSessionController(sessionService)

	r.HandleFunc("/auth/register", userController.CreateUser).Methods("POST")
	r.HandleFunc("/auth/login", userController.LoginUser).Methods("POST")
	r.HandleFunc("/auth/logout", utils.Authenticate(userController.LogoutUser)).Methods("POST")
	r.HandleFunc("/auth/delete", utils.Authenticate(userController.DeleteUser)).Methods("DELETE")  // frontend misising
	
	r.HandleFunc("/grupos", utils.Authenticate(groupController.GetAllGroupsByUserID)).Methods("GET")
	r.HandleFunc("/grupos/criar", utils.Authenticate(groupController.CreateGroup)).Methods("POST")
	r.HandleFunc("/grupos/{nome-do-grupo}/detalhes", utils.Authenticate(groupController.GetGroupDetails)).Methods("GET")
	r.HandleFunc("/grupos/{nome-do-grupo}/detalhes/adicionar", utils.Authenticate(groupController.AddMemberToGroup)).Methods("POST")
	r.HandleFunc("/grupos/{nome-do-grupo}/deletar", utils.Authenticate(groupController.DeleteOneGroup)).Methods("DELETE")
	r.HandleFunc("/grupos/deletar", utils.Authenticate(groupController.DeleteAllGroupsFromUser)).Methods("DELETE")
	r.HandleFunc("/grupos/{nome-do-grupo}/detalhes/{nome-do-membro}/deletar", utils.Authenticate(groupController.RemoveOneMemberFromGroup)).Methods("DELETE")
	r.HandleFunc("/grupos/{nome-do-grupo}/detalhes/deletar-membros", utils.Authenticate(groupController.RemoveAllMembersFromGroup)).Methods("DELETE")

	r.HandleFunc("/grupos/{nome-do-grupo}/sessoes/iniciar", utils.Authenticate(sessionController.StartNewSession)).Methods("POST")
	r.HandleFunc("/grupos/{nome-do-grupo}/sessoes/{nome-da-sessao}/validar-face", utils.Authenticate(sessionController.ValidateFace)).Methods("PATCH")
	r.HandleFunc("/grupos/{nome-do-grupo}/sessoes/{nome-da-sessao}/encerrar", utils.Authenticate(sessionController.EndSession)).Methods("POST")
	r.HandleFunc("/grupos/{nome-do-grupo}/sessoes/em-andamento", utils.Authenticate(sessionController.GetActiveSessions)).Methods("GET")
	r.HandleFunc("/grupos/{nome-do-grupo}/sessoes/encerradas", utils.Authenticate(sessionController.GetEndedSessions)).Methods("GET")
	r.HandleFunc("/grupos/{nome-do-grupo}/sessoes/{nome-da-sessao}/detalhes", utils.Authenticate(sessionController.GetSessionDetails)).Methods("GET")
	r.HandleFunc("/grupos/{nome-do-grupo}/sessoes/{nome-da-sessao}/deletar", utils.Authenticate(sessionController.DeleteOneSession)).Methods("DELETE")	
	
	// frontend missing
	r.HandleFunc("/grupos/{nome-do-grupo}/sessoes/em-andamento/deletar-todas", utils.Authenticate(sessionController.DeleteAllActiveSessionsOfAGroup)).Methods("DELETE")
	r.HandleFunc("/grupos/{nome-do-grupo}/sessoes/encerradas/deletar-todas", utils.Authenticate(sessionController.DeleteAllEndedSessionsOfAGroup)).Methods("DELETE")
    r.HandleFunc("/grupos/{nome-do-grupo}/sessoes/{nome-da-sessao}/detalhes/{nome-do-membro}/editar-presenca", utils.Authenticate(sessionController.UpdateMemberAttendance)).Methods("PATCH")

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
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)
}