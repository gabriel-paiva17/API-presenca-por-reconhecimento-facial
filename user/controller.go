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

/////////////////////////
// POST /auth/register //
/////////////////////////

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

//////////////////////
// POST /auth/login //
//////////////////////

func (c *UserController) LoginUserHandler(res http.ResponseWriter, req *http.Request) {

	var loginReq LoginRequest
    if err := json.NewDecoder(req.Body).Decode(&loginReq); err != nil {
        utils.WriteErrorResponse(res, http.StatusBadRequest, "Request Body invalido.")
        return
    }

    token, err := c.service.LoginUser(req.Context(), &loginReq)
    
	if err == ErrGeneratingToken {

		utils.WriteErrorResponse(res, http.StatusInternalServerError, err.Error())
        return
		
	}
	
	if err != nil {
		utils.WriteErrorResponse(res, http.StatusUnauthorized, err.Error())
        return
    }

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Expose-Headers", "Authorization")
    res.Header().Set("Authorization", "Bearer "+token)

    res.WriteHeader(http.StatusOK)

	response := LoginResponse{Message: "Login realizado com sucesso."}
    if err := json.NewEncoder(res).Encode(response); err != nil {
        utils.WriteErrorResponse(res, http.StatusInternalServerError, "Erro ao codificar resposta.")
    }

}