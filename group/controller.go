package group

import (
	"encoding/json"
	"errors"
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

	encodeErr := json.NewEncoder(res).Encode(getGroupsResponse)
	if encodeErr != nil {

		utils.WriteErrorResponse(res, http.StatusInternalServerError, "Erro Interno do Server")
		return
	}
}

// POST /grupos

func (c *GroupController) CreateGroupHandler(res http.ResponseWriter, req *http.Request) {

	var createGroupReq CreateGroupRequest
	if err := json.NewDecoder(req.Body).Decode(&createGroupReq); err != nil {
		utils.WriteErrorResponse(res, http.StatusBadRequest, "Request Body Invalido")
		return
	}

	userID, _ := utils.GetAuthenticatedUserId(req)

	createGroupReq.CreatedBy = userID

	ctx := req.Context()
	group, err := c.service.CreateGroup(ctx, &createGroupReq)

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
	encodeErr := json.NewEncoder(res).Encode(map[string]interface{}{
		"message": "Grupo criado com sucesso.",
		"group":   response,
	})

	if encodeErr != nil {
		utils.WriteErrorResponse(res, http.StatusInternalServerError, "Erro ao codificar resposta.")
		return
	}

}

// GET grupos/{nome-do-grupo}

func (c *GroupController) GetGroupDetails(res http.ResponseWriter, req *http.Request) {
	userId, _ := utils.GetAuthenticatedUserId(req)

	vars := mux.Vars(req)
	groupName := vars["nome-do-grupo"]

	group, err := c.service.GetGroupByName(groupName, userId, req.Context())
	if err != nil {
		utils.WriteErrorResponse(res, http.StatusNotFound, "Grupo do usuário não foi encontrado.")
		return
	}

	encodeErr := json.NewEncoder(res).Encode(group)
	if encodeErr != nil {
		utils.WriteErrorResponse(res, http.StatusInternalServerError, "Erro ao codificar resposta.")
		return
	}
}
