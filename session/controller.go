package session

import (
	"net/http"
	"errors"
	"encoding/json"
	"myproject/utils"

	"github.com/gorilla/mux"

)

type SessionController struct {
	service *SessionService
}

func NewSessionController(service *SessionService) *SessionController {
	return &SessionController{
		service: service,
	}
}

func (c *SessionController) StartNewSession(res http.ResponseWriter, req *http.Request) {

	var startSessionRequest StartSessionRequest
	if err := json.NewDecoder(req.Body).Decode(&startSessionRequest); err != nil {
		utils.WriteErrorResponse(res, http.StatusBadRequest, "Request Body Invalido")
		return
	}

	userID, _ := utils.GetAuthenticatedUserId(req)
	startSessionRequest.CreatedBy = userID

	vars := mux.Vars(req)
	groupName := vars["nome-do-grupo"]	
	startSessionRequest.GroupName = groupName

	newSession, err := c.service.StartNewSession(req.Context(), &startSessionRequest)

	if errors.Is(err, ErrGroupNotFound) {

		utils.WriteErrorResponse(res, http.StatusNotFound, err.Error())
		return

	}

	if errors.Is(err, ErrSessionAlreadyExists) {

		utils.WriteErrorResponse(res, http.StatusConflict, err.Error())
		return

	}

	if err != nil {

		utils.WriteErrorResponse(res, http.StatusInternalServerError, err.Error())
		return

	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	encodeErr := json.NewEncoder(res).Encode(map[string]interface{}{
		"message": "Sessão iniciada com sucesso.",
		"session": newSession,
	})
	if encodeErr != nil {
		utils.WriteErrorResponse(res, http.StatusInternalServerError, "Erro ao codificar resposta.")
	}	

}