package group

import (
    "myproject/member"
	"errors"
)

type Group struct {
    ID        string           `json:"id" bson:"_id"`
    Name      string           `json:"name" bson:"name"`
    CreatedAt string           `json:"createdAt" bson:"createdAt"`
    Members   []member.Member  `json:"members" bson:"members"`
    CreatedBy string           `json:"createdBy" bson:"createdBy"`
}

// POST /group

type CreateGroupRequest struct {
    Name string         `json:"name"`
    CreatedBy string    `json:"createdBy"` 
}

type CreateGroupResponse struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    CreatedAt string    `json:"createdAt"`
    CreatedBy string    `json:"createdBy"`
}

// erros de grupo

var ErrNameAlreadyExists = errors.New("name already used by you")