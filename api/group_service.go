package api

import (
	"context"
	"myproject/cv"
	"time"

	"github.com/google/uuid"
)

type GroupService struct {
	groupRepo *GroupRepository
	userRepo *UserRepository
	sessionRepo *SessionRepository
}

func NewGroupService(groupRepo *GroupRepository, userRepo *UserRepository, sessionRepo *SessionRepository) *GroupService {
	return &GroupService{groupRepo:  groupRepo, userRepo: userRepo, sessionRepo: sessionRepo}
}

// GET /grupos

func (s *GroupService) GetGroups(userID string, ctx context.Context) ([]GroupByName, error) {

	groups, err := s.groupRepo.FindAllGroupsByUserID(userID, ctx)

	if err != nil {

		return nil, err

	}

	var groupsByName []GroupByName

	for _, g := range groups {

		groupsByName = append(groupsByName, GroupByName{Name: g.Name})

	}

	return groupsByName, nil

}

// POST /grupos/criar

func (s *GroupService) CreateGroup(ctx context.Context, req *CreateGroupRequest) (*Group, error) {

	_, found := s.userRepo.FindOneByID(ctx, req.CreatedBy)

	if !found {

		return nil, ErrUserNotFound

	}

	group := &Group{
		ID:        uuid.New().String(),
		Name:      req.Name,
		CreatedAt: time.Now().Format(time.RFC3339),
		Members:   []Member{},
		CreatedBy: req.CreatedBy,
	}

	err := s.groupRepo.CreateGroup(ctx, group)
	if err != nil {
		return nil, err
	}

	return group, nil

}

// GET /grupos/{nome-do-grupo}/detalhes

func (s *GroupService) GetGroupDetails(groupName, userId string, ctx context.Context) (*GetGroupDetailsResponse, error) {

	group, ok := s.groupRepo.FindOneByNameAndCreator(ctx, groupName, userId)

	if !ok {

		return nil, ErrGroupNotFound

	}

	totalAttendance, err := s.sessionRepo.CalculateTotalAttendance(ctx, groupName, userId)

	if err != nil {

		return nil , err

	}

	responseMembers := []MemberResponse{}
	
	for i := range group.Members {

		mr := MemberResponse{

			Name: group.Members[i].Name,
			Face: group.Members[i].Face,
			AddedAt: group.Members[i].AddedAt,
			TotalAttendance: totalAttendance[group.Members[i].ID],

		}		

		responseMembers = append(responseMembers, mr)

	}

	response := &GetGroupDetailsResponse {

		Name: group.Name,
		CreatedAt: group.CreatedAt,
		Members: responseMembers,
	}

	return response, nil

}

// POST /grupos/{nome-do-grupo}/detalhes/adicionar

func (s *GroupService) AddMemberToGroup(ctx context.Context, groupName, userID string, req *AddMemberRequest) (*Member, error) {

	err := cv.CheckOnlyOneFace(req.Face)

	if err != nil {

		return nil, err

	}
 
	newMember := &Member{
		ID: uuid.New().String(),
		Name: req.Name,
		Face: req.Face,
		AddedAt: time.Now().Format(time.RFC3339), 
	}

	addedMember, err := s.groupRepo.AddMemberToGroup(ctx, groupName, userID, newMember) 

    if err != nil {

        return nil, err

	}

	return addedMember, nil

}

// DELETE /grupos/{nome-do-grupo}/deletar

func (s *GroupService) DeleteOneGroup(ctx context.Context, groupName, createdBy string) error {

	_, found := s.groupRepo.FindOneByNameAndCreator(ctx, groupName, createdBy)
	if !found {

		return ErrGroupNotFound

	}

	err := s.sessionRepo.DeleteAllActiveSessionsOfAGroup(ctx, groupName, createdBy)
	if err != nil {

		return err

	}

	err = s.sessionRepo.DeleteAllEndedSessionsOfAGroup(ctx, groupName, createdBy)
	if err != nil {

		return err

	}

	return s.groupRepo.DeleteOneGroup(ctx, groupName, createdBy)

}