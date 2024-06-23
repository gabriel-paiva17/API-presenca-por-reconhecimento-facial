package api

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

// GET /grupos/{nome-do-grupo}/detalhes

type GetGroupDetailsResponse struct {

	Name      string           `json:"name"`
	CreatedAt string           `json:"createdAt"`
	Members   []MemberResponse `json:"members"`
}

type MemberResponse struct {

	Name       		string `json:"name"`
	Face      		string `json:"face"`
	AddedAt    		string `json:"addedAt"`
	TotalAttendance int	   `json:"totalAttendance"`

}


// GET /grupos

type GetAllGroupsResponse struct {
	Groups []GroupByName `json:"groups"`
}

type GroupByName struct {
	Name string `json:"name"`
}

// POST /grupos/criar

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
var ErrMemberNotFound = errors.New("member not found in group")