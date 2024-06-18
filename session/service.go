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

// POST /grupos/{nome-do-grupo}/sessoes/iniciar 

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

	err := cv.CheckOnlyOneFace(req.Face)
	if err != nil {

		return err

	}

	session, found := s.sessionRepo.FindOneSession(ctx, req.GroupName, req.CreatedBy, req.SessionName)
	if !found {

		return ErrSessionNotFound

	}

	if session.EndedAt != "" {

		return ErrSessionHasEnded

	}


	// caso a face seja encontrada, essa face é atribuida ao membro na sessao,
	// e o membro ganha o maximo de presencas da sessao
	faceValidated := false

	for i := range session.Members {
        sameFace, err := cv.CompareFaces(session.Members[i].Face, req.Face)
        if err != nil {
            return err
        }

        if sameFace {
            session.Members[i].Face = req.Face
            session.Members[i].Attendance = session.MaxAttendance
            session.Members[i].WasFaceValidated = true
            faceValidated = true
            break
        }
    }

	if !faceValidated {

		return ErrFaceDoesntMatch

	}

	return s.sessionRepo.UpdateMembers(ctx, session, session.Members)

}

func (s *SessionService) EndSession(ctx context.Context, req *EndSessionRequest) error {

	session, found := s.sessionRepo.FindOneSession(ctx, req.GroupName, req.CreatedBy, req.SessionName)
	if !found {

		return ErrSessionNotFound

	}

	if session.EndedAt != "" {

		return ErrSessionHasEnded

	}

	for i := range session.Members {

		if !session.Members[i].WasFaceValidated {

			session.Members[i].Face = ""

		}

	}

	err := s.sessionRepo.UpdateMembers(ctx, session, session.Members)
	if err != nil {

		return err

	}

	return s.sessionRepo.EndSession(ctx, session)

}

// DELETE /grupos/{nome-do-grupo}/sessoes/{nome-da-sessao}

func (s *SessionService) DeleteOneSession(ctx context.Context, groupName, createdBy, sessionName string) error {
    _, found := s.sessionRepo.FindOneSession(ctx, groupName, createdBy, sessionName)
    if !found {
        return ErrSessionNotFound
    }
  
	return 	s.sessionRepo.DeleteOneSession(ctx, groupName, createdBy, sessionName)

}