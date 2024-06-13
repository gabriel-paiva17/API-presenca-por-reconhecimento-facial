package session

import (
	"errors"
)

type Session struct {
	ID            string         `json:"id" bson:"_id"`
	Name          string         `json:"name" bson:"name"`
	MaxAttendance int            `json:"maxAttendance" bson:"maxAttendance"`
	StartedAt     string         `json:"startedAt" bson:"startedAt"`
	EndedAt       string         `json:"endedAt,omitempty" bson:"endedAt,omitempty"`
	GroupName     string         `json:"groupName" bson:"groupName"`
	CreatedBy     string         `json:"createdBy" bson:"createdBy"`
	Members       []SessionMember `json:"members" bson:"members"`
}

type SessionMember struct {
	ID        			string 		`json:"id" bson:"_id"`
	Name      			string 		`json:"name" bson:"name"`
	Face       			string 		`json:"face" bson:"face"`
	Attendance 			int    		`json:"attendance" bson:"attendance"`
	WasFaceValidated	bool		`json:"wasFaceValidated" bson:"wasFaceValidated"`
}

// POST /grupos/{nome-do-grupo}/sessoes/iniciar

type StartSessionRequest struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	MaxAttendance int    `json:"maxAttendance"`
	StartedAt     string `json:"startedAt"`
	GroupName     string `json:"groupName"`
	CreatedBy     string `json:"createdBy"`
}

type ValidateFaceRequest struct {

	Face 		  string `json:"face"`
	SessionName   string `json:"sessionName"`
	GroupName     string `json:"groupName"`
	CreatedBy     string `json:"createdBy"`

}

var ErrSessionAlreadyExists = errors.New("essa sessao ja existe, ou está em andamento")
var ErrGroupNotFound = errors.New("grupo nao encontrado")
var ErrSessionNotFound = errors.New("sessao nao encontrada")
var ErrFaceDoesntMatch = errors.New("face enviada nao corresponde a de nenhum membro do grupo")