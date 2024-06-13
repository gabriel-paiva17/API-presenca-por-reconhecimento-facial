package session

import (
	"context"
	"myproject/cv"
	"myproject/group"
	"time"

	"github.com/google/uuid"
)

type SessionService struct {
	sessionRepo *SessionRepository
	groupRepo *group.GroupRepository
}

func NewSessionService(sessionRepo *SessionRepository, groupRepo *group.GroupRepository) *SessionService {
	return &SessionService{sessionRepo: sessionRepo, groupRepo: groupRepo}
}

func (s *SessionService) StartNewSession(ctx context.Context, req *StartSessionRequest) (*Session, error) {

	foundedGroup, found := s.groupRepo.FindOneByNameAndCreator(ctx, req.GroupName, req.CreatedBy)
	if !found {

		return nil, ErrGroupNotFound

	}

	var sessionMembers []SessionMember
	
	for _, m := range foundedGroup.Members {

		sm := SessionMember{
			ID: m.ID,
			Name: m.Name,
			Face: m.Face,
			Attendance: 0,
			WasFaceValidated: false,
		}
		
		sessionMembers = append(sessionMembers, sm)

	}

	newSession := &Session{
		ID: uuid.New().String(),
		Name: req.Name,
		MaxAttendance: req.MaxAttendance,
		StartedAt: time.Now().Format(time.RFC3339),
		GroupName: req.GroupName,
		CreatedBy: req.CreatedBy,
		Members: sessionMembers,
	}

	err := s.sessionRepo.StartNewSession(ctx, newSession)

	if err != nil {
		
		return nil, err

	}

	return newSession, nil

}

// PUT /grupos/{nome-do-grupo}/sessoes/{nome-da-sessao}/validar-face

func (s *SessionService) ValidateFace(ctx context.Context, req *ValidateFaceRequest) error {

	/* mover para controller posteriormente
	if !cv.IsBase64JPG(req.Face) {

		return 

	}
	*/

	err := cv.CheckOnlyOneFace(req.Face)
	if err != nil {

		return err

	}

	session, found := s.sessionRepo.FindOneSession(ctx, req.GroupName, req.CreatedBy, req.SessionName)
	if !found {

		return ErrSessionNotFound

	}


	// caso a face seja encontrada, essa face é atribuida ao membro na sessao,
	// e o membro ganha o maximo de presencas da sessao
	faceValidated := false

	for _, m := range session.Members {

		sameFace, err  := cv.CompareFaces(m.Face, req.Face)

		if err != nil {

			return err

		}
		
		if sameFace {

			m.Face = req.Face
			m.Attendance = session.MaxAttendance
			m.WasFaceValidated = true
			faceValidated = true
			break
		}

	}

	if !faceValidated {

		return ErrFaceDoesntMatch

	}

	err = s.sessionRepo.UpdateMembers(ctx, session, session.Members)

	return err

}
