package session

import "myproject/group"

type Session struct {
	ID            string         `json:"id" bson:"_id"`
	Name          string         `json:"name" bson:"name"`
	MaxAttendance int            `json:"maxAttendance" bson:"maxAttendance"`
	StartedAt     string         `json:"startedAt" bson:"startedAt"`
	EndedAt       string         `json:"endedAt,omitempty" bson:"endedAt,omitempty"`
	GroupName     string         `json:"groupName" bson:"groupName"`
	CreatedBy     string         `json:"createdBy" bson:"createdBy"`
	Members       []group.Member `json:"members" bson:"members"`
}

type StartSessionRequest struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	MaxAttendance int    `json:"maxAttendance"`
	StartedAt     string `json:"startedAt"`
	GroupName     string `json:"groupName"`
	CreatedBy     string `json:"createdBy"`
}