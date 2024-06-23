package api

import (
	"context"
	"fmt"
	"myproject/utils"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo    *UserRepository
	groupRepo   *GroupRepository
	sessionRepo *SessionRepository
	jwtKey      []byte
}

func NewUserService(userRepo *UserRepository, groupRepo *GroupRepository, sessionRepo *SessionRepository, jwtKey string) *UserService {
	return &UserService{
		userRepo:    userRepo,
		groupRepo:   groupRepo,
		sessionRepo: sessionRepo,
		jwtKey: []byte(jwtKey),
	}
}

/////////////////////////
// POST /auth/register //
/////////////////////////

func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {

	newUser := &User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		Email:        req.Email,
		Password:     req.Password,
		RegisteredAt: time.Now().Format(time.RFC3339),
	}

	if err := s.userRepo.CreateUser(ctx, newUser); err != nil {
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

//////////////////////
// POST /auth/login //
//////////////////////

func (s *UserService) LoginUser(ctx context.Context, req *LoginRequest, res http.ResponseWriter) error {

	dbUser, found := s.userRepo.FindOneByEmail(ctx, req.Email)

	if !found {
		return fmt.Errorf("usuario nao existe")
	}

	if !utils.IsHashEqualPassword(dbUser.Password, req.Password) {

		return fmt.Errorf("senha invalida")

	}

	// configurando token de autenticacao
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": dbUser.ID,
		"exp":    time.Now().Add(time.Hour * 12).Unix(),
	})

	// adicionando chave secreta a assinatura
	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		return ErrGeneratingToken
	}

	cookie := &http.Cookie{
        Name:     "auth-token",
        Value:    tokenString,
		Path: "/",
        Expires:  time.Now().Add(12 * time.Hour),
        HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
    }
    http.SetCookie(res, cookie)

	return nil

}

// DELETE /auth/delete

func (s *UserService) DeleteUser(ctx context.Context, userId string) error {

	err := s.userRepo.DeleteUser(ctx, userId)

	if err != nil {

		return err

	}
	
	err = s.groupRepo.DeleteAllGroupsFromUser(ctx, userId)

	if err != nil {

		return err

	}

	err = s.sessionRepo.DeleteAllSessionsFromUser(ctx, userId)

	return err

} 