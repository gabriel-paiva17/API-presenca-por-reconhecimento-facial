package user

import (
	"encoding/json"
	"net/http"
)

type UserController struct {
	service *UserService
}

func NewUserController(service *UserService) *UserController {
	return &UserController{
		service: service,
	}
}

func (c *UserController) CreateUserHandler(res http.ResponseWriter, req *http.Request) {
	
	// TODO: implementar mensagens de erro documentadas no projeto
	// ou entao alterar documentacao
	
	var createUserRequest CreateUserRequest
	if err := json.NewDecoder(req.Body).Decode(&createUserRequest); err != nil {
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := req.Context()

	createUserResponse, err := c.service.CreateUser(ctx, &createUserRequest)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// Serializa a resposta como JSON e envia
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(map[string]interface{}{
		"message": "Usuário criado com sucesso.",
		"user":    createUserResponse,
	})
}