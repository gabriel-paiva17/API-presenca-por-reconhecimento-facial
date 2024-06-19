package api

import (
	"encoding/json"
	"errors"
	"myproject/cv"
	"myproject/utils"
	"net/http"

	"github.com/gorilla/mux"
)

type GroupController struct {
	service *GroupService
}

func NewGroupController(service *GroupService) *GroupController {
	return &GroupController{
		service: service,
	}
}

// GET /grupos

func (c *GroupController) GetAllGroupsByUserID(res http.ResponseWriter, req *http.Request) {

	userID, _ := utils.GetAuthenticatedUserId(req)

	groupsByName, err := c.service.GetGroups(userID, req.Context())

	if err != nil {

		utils.WriteErrorResponse(res, http.StatusInternalServerError, err.Error())
		return

	}

	if len(groupsByName) == 0 {

		utils.WriteErrorResponse(res, http.StatusNotFound, "Nenhum grupo foi encontrado.")
		return

	}

	getGroupsResponse := GetAllGroupsResponse{Groups: groupsByName}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	json.NewEncoder(res).Encode(getGroupsResponse)
	
}

// POST /grupos

func (c *GroupController) CreateGroup(res http.ResponseWriter, req *http.Request) {

	var createGroupReq CreateGroupRequest
	if err := json.NewDecoder(req.Body).Decode(&createGroupReq); err != nil {
		utils.WriteErrorResponse(res, http.StatusBadRequest, "Request Body Invalido")
		return
	}

	userID, _ := utils.GetAuthenticatedUserId(req)

	createGroupReq.CreatedBy = userID

	ctx := req.Context()
	group, err := c.service.CreateGroup(ctx, &createGroupReq)

	// garantindo que nao seja criado um grupo para um usuario que nao existe mais
	if errors.Is(err, ErrUserNotFound) {

		utils.WriteErrorResponse(res, http.StatusNotFound, err.Error())
		return

	}

	if errors.Is(err, ErrNameAlreadyExists) {

		utils.WriteErrorResponse(res, http.StatusConflict, "Nome ja utilizado por voce.")
		return

	}

	if err != nil {
		utils.WriteErrorResponse(res, http.StatusInternalServerError, err.Error())
		return
	}

	response := CreateGroupResponse{
		ID:        group.ID,
		Name:      group.Name,
		CreatedAt: group.CreatedAt,
		CreatedBy: group.CreatedBy,
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(map[string]interface{}{
		"message": "Grupo criado com sucesso.",
		"group":   response,
	})

}

// GET grupos/{nome-do-grupo}/detalhes

func (c *GroupController) GetGroupDetails(res http.ResponseWriter, req *http.Request) {
	userId, _ := utils.GetAuthenticatedUserId(req)

	vars := mux.Vars(req)
	groupName := vars["nome-do-grupo"]

	group, err := c.service.GetGroupByName(groupName, userId, req.Context())
	if err != nil {
		utils.WriteErrorResponse(res, http.StatusNotFound, "Grupo do usuário não foi encontrado.")
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(group)
	
}


// POST grupos/{nome-do-grupo}/detalhes/adicionar

func (c *GroupController) AddMemberToGroup(res http.ResponseWriter, req *http.Request) {

    userId, _ := utils.GetAuthenticatedUserId(req)

	vars := mux.Vars(req)
	groupName := vars["nome-do-grupo"]

	var addMemberReq AddMemberRequest
	if err := json.NewDecoder(req.Body).Decode(&addMemberReq); err != nil {
		utils.WriteErrorResponse(res, http.StatusBadRequest, "Request Body Invalido")
		return
	}

    if !cv.IsBase64JPG(addMemberReq.Face) {

        utils.WriteErrorResponse(res, http.StatusUnsupportedMediaType, "Face deve ser uma imagem .jpg na base 64")
        return

	}

	addedMember, err := c.service.AddMemberToGroup(req.Context(), groupName, userId, &addMemberReq)

    if errors.Is(err, cv.ErrNoFaces) || errors.Is(err, cv.ErrMoreThanOneFace) {

        utils.WriteErrorResponse(res, http.StatusBadRequest, err.Error())
        return

	}

	if errors.Is(err, ErrFaceAlreadyUsed) || errors.Is(err, ErrNameAlreadyExists) {

        utils.WriteErrorResponse(res, http.StatusConflict, err.Error())
        return

	}
   
    if errors.Is(err, ErrGroupNotFound) {

        utils.WriteErrorResponse(res, http.StatusNotFound, err.Error())
        return

	}

    if err != nil {

		utils.WriteErrorResponse(res, http.StatusInternalServerError, err.Error())
        return

	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(map[string]interface{}{
		"message": "Membro adicionado com sucesso.",
		"member":   addedMember,
	})

}	


// DELETE /grupos/{nome-do-grupo}/deletar

func (c *GroupController) DeleteOneGroup(res http.ResponseWriter, req *http.Request) {

	userId, _ := utils.GetAuthenticatedUserId(req)

	vars := mux.Vars(req)
	groupName := vars["nome-do-grupo"]

	err := c.service.DeleteOneGroup(req.Context(), groupName, userId)

	if errors.Is(err, ErrGroupNotFound) {

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
		"message": "Grupo deletado coom sucesso.",
	})

}