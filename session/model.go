package session

import "myproject/group"

type Session struct {

	ID			string
	Name		string
	StartedAt	string
	EndedAt		string
	GroupName 	string
	CreatedBy	string
	Members		[]group.Member

}