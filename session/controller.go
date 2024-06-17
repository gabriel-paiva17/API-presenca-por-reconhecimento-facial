package session

import (
	"encoding/json"
	"errors"
	"myproject/cv"
	"myproject/utils"
	"net/http"

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
	json.NewEncoder(res).Encode(map[string]interface{}{
		"message": "Sessão iniciada com sucesso.",
		"session": newSession,
	})

}

// PUT /grupos/{nome-do-grupo}/sessoes/{nome-da-sessao}
func (c *SessionController) ValidateFace(res http.ResponseWriter, req *http.Request) {

	var validateFaceRequest ValidateFaceRequest
	if err := json.NewDecoder(req.Body).Decode(&validateFaceRequest); err != nil {
		utils.WriteErrorResponse(res, http.StatusBadRequest, "Request Body Invalido")
		return
	}

	if !cv.IsBase64JPG(validateFaceRequest.Face) {

		utils.WriteErrorResponse(res, http.StatusUnsupportedMediaType, "Face deve ser uma imagem .jpg na base 64")
        return

	}
	
	vars := mux.Vars(req)
	groupName := vars["nome-do-grupo"]
	sessionName := vars["nome-da-sessao"]

	userId, _ := utils.GetAuthenticatedUserId(req)

	validateFaceRequest.GroupName = groupName
	validateFaceRequest.SessionName = sessionName
	validateFaceRequest.CreatedBy = userId

	err := c.service.ValidateFace(req.Context(), &validateFaceRequest)

	if errors.Is(err, cv.ErrNoFaces) || errors.Is(err, cv.ErrMoreThanOneFace) {

		utils.WriteErrorResponse(res, http.StatusBadRequest, err.Error())
		return

	} 

	if errors.Is(err, ErrSessionNotFound) {

		utils.WriteErrorResponse(res, http.StatusNotFound, err.Error())
        return

	}

	if errors.Is(err, ErrSessionHasEnded) {

		utils.WriteErrorResponse(res, http.StatusConflict, err.Error())
		return

	}

	if errors.Is(err, ErrFaceDoesntMatch) {

		utils.WriteErrorResponse(res, http.StatusUnprocessableEntity, err.Error())
		return

	}

	if err != nil {

		utils.WriteErrorResponse(res, http.StatusInternalServerError, err.Error())
        return

	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(map[string]interface{}{
		"message": "Face validada.",
	})

}

// POST /grupos/{nome-do-grupo}/sessoes/{nome-da-sessao}/encerrar
func (c *SessionController) EndSession(res http.ResponseWriter, req *http.Request) {

	userId, _ := utils.GetAuthenticatedUserId(req)

	vars := mux.Vars(req)
	groupName := vars["nome-do-grupo"]
	sessionName := vars["nome-da-sessao"]

	endSessionRequest := EndSessionRequest{
		SessionName: sessionName,
    	GroupName: groupName,
		CreatedBy: userId, 
	}

	err := c.service.EndSession(req.Context(), &endSessionRequest)

	if errors.Is(err, ErrSessionNotFound) {

		utils.WriteErrorResponse(res, http.StatusNotFound, err.Error())
        return

	}

	if errors.Is(err, ErrSessionHasEnded) {
		utils.WriteErrorResponse(res, http.StatusForbidden, err.Error())
        return
	}

	if err != nil {
		utils.WriteErrorResponse(res, http.StatusInternalServerError, err.Error())
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(map[string]interface{}{
		"message": "Sessão encerrada.",
	})

}

// GET /grupos/{nome-do-grupo}/sessoes/finalizadas

func (c *SessionController) GetActiveSessions(res http.ResponseWriter, req *http.Request) {

	userId, _ := utils.GetAuthenticatedUserId(req)

	vars := mux.Vars(req)
	groupName := vars["nome-do-grupo"]

	response, err := c.service.sessionRepo.FindAllActiveSessions(req.Context(), groupName, userId)

	if err != nil {

		utils.WriteErrorResponse(res, http.StatusInternalServerError, err.Error())
		return

	}

	if len(response.Sessions) == 0 {

		utils.WriteErrorResponse(res, http.StatusNotFound, "Nenhuma sessão foi encontrada.")
		return

	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(response)

}

// GET /grupos/{nome-do-grupo}/sessoes/em-andamento

func (c *SessionController) GetEndedSessions(res http.ResponseWriter, req *http.Request) {

	userId, _ := utils.GetAuthenticatedUserId(req)

	vars := mux.Vars(req)
	groupName := vars["nome-do-grupo"]

	response, err := c.service.sessionRepo.FindAllEndedSessions(req.Context(), groupName, userId)

	if err != nil {

		utils.WriteErrorResponse(res, http.StatusInternalServerError, err.Error())
		return

	}

	if len(response.Sessions) == 0 {

		utils.WriteErrorResponse(res, http.StatusNotFound, "Nenhuma sessão foi encontrada.")
		return

	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(response)

}