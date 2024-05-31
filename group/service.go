package group

import (
	"context"
	"myproject/member"
	"time"

	"github.com/google/uuid"
)

type GroupService struct {
	repo *GroupRepository
}

func NewGroupService(repo *GroupRepository) *GroupService {
	return &GroupService{repo: repo}
}

func (s *GroupService) CreateGroup(ctx context.Context, req *CreateGroupRequest) (*Group, error) {
	group := &Group{
		ID:        uuid.New().String(),
		Name:      req.Name,
		CreatedAt: time.Now().Format(time.RFC3339),
		Members:   []member.Member{},
	}

	err := s.repo.CreateGroup(ctx, group)
	if err != nil {
		return nil, err
	}

	return group, nil

}
