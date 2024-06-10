package group

import (
	"errors"
)

type Group struct {
	ID        string          `json:"id" bson:"_id"`
	Name      string          `json:"name" bson:"name"`
	CreatedAt string          `json:"createdAt" bson:"createdAt"`
	Members   []Member 		  `json:"members" bson:"members"`
	CreatedBy string          `json:"createdBy" bson:"createdBy"`
}

// GET /grupos

type GetAllGroupsResponse struct {
	Groups []GroupByName `json:"groups"`
}

type GroupByName struct {
	Name string `json:"name"`
}

// POST /grupos

type CreateGroupRequest struct {
	Name      string `json:"name"`
	CreatedBy string `json:"createdBy"`
}

type CreateGroupResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	CreatedBy string `json:"createdBy"`
}

// POST /grupos/{nome-do-grupo}/detalhes/adicionar

type Member struct {
	ID         string `json:"id" bson:"_id"`
	Name       string `json:"name" bson:"name"`
	Face       string `json:"face" bson:"face"`
	Attendance int    `json:"attendance" bson:"attendance"`
	AddedAt    string `json:"addedAt" bson:"addedAt"`
}

type AddMemberRequest struct {

	Name	string `json:"name"`
	Face 	string `json:"face"`

}

// erros de grupo

var ErrNameAlreadyExists = errors.New("name already used by you")
var ErrGroupNotFound = errors.New("group not found")
var ErrFaceAlreadyUsed = errors.New("member with the same face already exists in the group")
var ErrNoFaces = errors.New("faces were not captured in the image")
var ErrMoreThanOneFace = errors.New("mais de uma face capturada, tente ficar em um fundo neutro")