package group

import (
	"net/http"
)

type GroupController struct {
	service *GroupService
}

func NewGroupController(service *GroupService) *GroupController {
	return &GroupController{
		service: service,
	}
}

func (c *GroupController) CreateGroupHandler(res http.ResponseWriter, req *http.Request) {}