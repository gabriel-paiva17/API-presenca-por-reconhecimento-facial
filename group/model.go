package group

import (
    "myproject/member" // Substitua "your_project_name" pelo nome real do seu módulo
)

type Group struct {
    ID        string           `json:"id" bson:"_id"`
    Name      string           `json:"name" bson:"name"`
    CreatedAt string           `json:"createdAt" bson:"createdAt"`
    Members   []member.Member  `json:"members" bson:"members" `
}