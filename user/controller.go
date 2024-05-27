package user

import (
	"encoding/json"
	"net/http"
	"myproject/utils"
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
	
	var createUserRequest CreateUserRequest
	if err := json.NewDecoder(req.Body).Decode(&createUserRequest); err != nil {
		utils.WriteErrorResponse(res, http.StatusBadRequest, "Request Body invalido.")
		return
	}

	if !utils.IsValidEmail(createUserRequest.Email) {
		utils.WriteErrorResponse(res, http.StatusBadRequest, "Email invalido.")
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