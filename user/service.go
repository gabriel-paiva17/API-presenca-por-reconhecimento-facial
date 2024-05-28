package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	
	newUser := &User{
		ID: 		  uuid.New().String(),
		Username:     req.Username,
		Email:        req.Email,
		Password:     req.Password,
		RegisteredAt: time.Now().Format(time.RFC3339),
	}

	if err := s.repo.CreateUser(ctx, newUser); err != nil {
		return nil, err
	}

	// Resposta com os dados do usuário criado
	response := &CreateUserResponse{
		ID:           newUser.ID,
		Username:     newUser.Username,
		Email:        newUser.Email,
		RegisteredAt: newUser.RegisteredAt,
	}

	return response, nil
}