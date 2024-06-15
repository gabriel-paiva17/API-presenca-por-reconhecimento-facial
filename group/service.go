package group

import (
	"context"
	"errors"
	"myproject/cv"
	"myproject/user"
	"time"

	"github.com/google/uuid"
)

type GroupService struct {
	groupRepo *GroupRepository
	userRepo *user.UserRepository
}

func NewGroupService(groupRepo *GroupRepository, userRepo *user.UserRepository) *GroupService {
	return &GroupService{groupRepo:  groupRepo, userRepo: userRepo}
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

// POST /grupos

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

func (s *GroupService) GetGroupByName(groupName, userId string, ctx context.Context) (*Group, error) {

	group, ok := s.groupRepo.FindOneByNameAndCreator(ctx, groupName, userId)

	if !ok {

		return nil, errors.New("not found")

	}

	return group, nil

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
		Attendance: 0,
		AddedAt: time.Now().Format(time.RFC3339), 
	}

	addedMember, err := s.groupRepo.AddMemberToGroup(ctx, groupName, userID, newMember) 

    if err != nil {

        return nil, err

	}

	return addedMember, nil

}