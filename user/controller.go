package user

import (
	"encoding/json"
	"errors"
	"myproject/utils"
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
	
	var createUserRequest CreateUserRequest
	if err := json.NewDecoder(req.Body).Decode(&createUserRequest); err != nil {
		utils.WriteErrorResponse(res, http.StatusBadRequest, "Request Body invalido.")
		return
	}

	if !utils.IsValidEmail(createUserRequest.Email) {
		utils.WriteErrorResponse(res, http.StatusBadRequest, "Email invalido.")
		return
	}

	hashedPassword, err := utils.HashPassword(createUserRequest.Password)
	if err != nil {
		utils.WriteErrorResponse(res, http.StatusBadRequest, "Senha invalida.")
		return
	}

	createUserRequest.Password = hashedPassword

	ctx := req.Context()

	createUserResponse, err := c.service.CreateUser(ctx, &createUserRequest)
	
	if errors.Is(err, ErrEmailAlreadyExists) {

		utils.WriteErrorResponse(res, http.StatusConflict, "Email ja utilizado anteriormente.")
		return

	}

	if err != nil {
		utils.WriteErrorResponse(res, http.StatusInternalServerError, err.Error())
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