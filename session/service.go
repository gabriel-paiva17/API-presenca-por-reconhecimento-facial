package session

import (
	"context"
	"myproject/group"
	"time"

	"github.com/google/uuid"
)

type SessionService struct {
	repo *SessionRepository
	groupRepo *group.GroupRepository
}

func NewSessionService(repo *SessionRepository) *SessionService {
	return &SessionService{repo: repo}
}

func (s *SessionService) StartNewSession(ctx context.Context, req *StartSessionRequest) (*Session, error) {

	foundedGroup, found := s.groupRepo.FindOneByNameAndCreator(ctx, req.GroupName, req.CreatedBy)
	if !found {

		return nil, ErrGroupNotFound

	}

	for _, m := range foundedGroup.Members {

		
		m.Attendance = 0

	}

	newSession := &Session{
		ID: uuid.New().String(),
		Name: req.Name,
		MaxAttendance: req.MaxAttendance,
		StartedAt: time.Now().Format(time.RFC3339),
		GroupName: req.GroupName,
		CreatedBy: req.CreatedBy,
		Members: foundedGroup.Members,
	}

	err := s.repo.StartNewSession(ctx, newSession)

	if err != nil {
		
		return nil, err

	}

	return newSession, nil

}