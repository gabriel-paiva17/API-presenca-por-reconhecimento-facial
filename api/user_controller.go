package api

import (
	"encoding/json"
	"errors"
	"myproject/utils"
	"net/http"
	"time"
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

func (c *UserController) CreateUser(res http.ResponseWriter, req *http.Request) {

	var createUserRequest CreateUserRequest
	if err := json.NewDecoder(req.Body).Decode(&createUserRequest); err != nil {
		utils.WriteErrorResponse(res, http.StatusBadRequest, "Request Body invalido.")
		return
	}

	if createUserRequest.Username == "" {
		utils.WriteErrorResponse(res, http.StatusBadRequest, "Nome precisa ser preenchido.")
		return
	}

	if len(createUserRequest.Password) < 8 {
		utils.WriteErrorResponse(res, http.StatusBadRequest, "Senha precisa ter no mínimo 8 caracteres.")
		return
	}

	hashedPassword, err := utils.HashPassword(createUserRequest.Password)
	if err != nil {
		utils.WriteErrorResponse(res, http.StatusBadRequest, "Senha invalida.")
		return
	}


	if !utils.IsValidEmail(createUserRequest.Email) {
		utils.WriteErrorResponse(res, http.StatusBadRequest, "Email invalido.")
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

func (c *UserController) LoginUser(res http.ResponseWriter, req *http.Request) {

	var loginReq LoginRequest
	if err := json.NewDecoder(req.Body).Decode(&loginReq); err != nil {
		utils.WriteErrorResponse(res, http.StatusBadRequest, "Request Body invalido.")
		return
	}

	err := c.service.LoginUser(req.Context(), &loginReq, res)

	if err == ErrGeneratingToken {

		utils.WriteErrorResponse(res, http.StatusInternalServerError, err.Error())
		return

	}

	if err != nil {
		utils.WriteErrorResponse(res, http.StatusUnauthorized, err.Error())
		return
	}

	res.Header().Set("Content-Type", "application/json")

	res.WriteHeader(http.StatusOK)

	json.NewEncoder(res).Encode(map[string]interface{}{
		"message": "Login realizado com sucesso.",
	})

}

///////////////////////
// POST /auth/logout //
///////////////////////

func (c *UserController) LogoutUser(res http.ResponseWriter, req *http.Request) {

	var response LogoutResponse
    response.Date = time.Now().Format(time.RFC3339)
    response.Message = "Logout realizado com sucesso!"

    // Clear the authentication cookie
    cookie := &http.Cookie{
        Name:     "auth-token",
        Value:    "",
        Path: 	  "/",
		Expires:  time.Unix(0, 0),
        HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
    }
    http.SetCookie(res, cookie)

    res.Header().Set("Content-Type", "application/json")
    res.WriteHeader(http.StatusOK)
    json.NewEncoder(res).Encode(response)

}

// DELETE /auth/delete

func (c *UserController) DeleteUser(res http.ResponseWriter, req *http.Request) {

	userId, _ := utils.GetAuthenticatedUserId(req)

	err := c.service.DeleteUser(req.Context(), userId)

	if err == ErrUserNotFound {

		utils.WriteErrorResponse(res, http.StatusNotFound, err.Error())
		return

	}

	if err != nil {

		utils.WriteErrorResponse(res, http.StatusInternalServerError, err.Error())
		return

	}

	res.Header().Set("Content-Type", "application/json")

	res.WriteHeader(http.StatusOK)

	json.NewEncoder(res).Encode(map[string]interface{}{
		"message": "Usuário deletado com sucesso.",
	})

}