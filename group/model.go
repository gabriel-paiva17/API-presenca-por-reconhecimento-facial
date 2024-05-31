package group

import (
    "myproject/member"
)

type Group struct {
    ID        string           `json:"id" bson:"_id"`
    Name      string           `json:"name" bson:"name"`
    CreatedAt string           `json:"createdAt" bson:"createdAt"`
    Members   []member.Member  `json:"members" bson:"members"`
}