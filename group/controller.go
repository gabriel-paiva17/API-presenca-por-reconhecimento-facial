package group

import (
	"encoding/json"
	"myproject/utils"
	"net/http"
	"errors"
)

type GroupController struct {
	service *GroupService
}

func NewGroupController(service *GroupService) *GroupController {
	return &GroupController{
		service: service,
	}
}

func (c *GroupController) CreateGroupHandler(res http.ResponseWriter, req *http.Request) {

	var createGroupReq CreateGroupRequest
	if err := json.NewDecoder(req.Body).Decode(&createGroupReq); err != nil {
		http.Error(res, "Invalid request body", http.StatusBadRequest)
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
	json.NewEncoder(res).Encode(map[string]interface{}{
		"message": "Grupo criado com sucesso.",
		"group":    response,
	})

}