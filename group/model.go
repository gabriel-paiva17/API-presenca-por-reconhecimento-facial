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
}

// POST /group

type CreateGroupRequest struct {
    Name string `json:"name"`
}

type CreateGroupResponse struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    CreatedAt string    `json:"createdAt"`
}


// erros de grupo

var ErrNameAlreadyExists = errors.New("name already used")