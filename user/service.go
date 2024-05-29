package user

import (
	"context"
	"myproject/utils"
	"time"
	"fmt"

	"github.com/google/uuid"
	"github.com/dgrijalva/jwt-go"
)

type UserService struct {
	repo UserRepository
	jwtKey   []byte
}

func NewUserService(repo UserRepository, jwtKey string) *UserService {
	return &UserService{
		repo:   repo,
		jwtKey: []byte(jwtKey),
	}
}

/////////////////////////
// POST /auth/register //
/////////////////////////

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

//////////////////////
// POST /auth/login //
//////////////////////

func (s *UserService) LoginUser(ctx context.Context, req *LoginRequest) (string, error) {

	dbUser, found := s.repo.FindOneByEmail(ctx, req.Email)

	if !found {
        return "", fmt.Errorf("usuario nao existe")
    }

	if !utils.IsHashEqualPassword(dbUser.Password, req.Password) {

		return "", fmt.Errorf("senha invalida")


	}

	// configurando token de autenticacao
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "userId": dbUser.ID,
        "exp":    time.Now().Add(time.Hour * 72).Unix(),
    })


	// adicionando chave secreta a assinatura
    tokenString, err := token.SignedString(s.jwtKey)
    if err != nil {
        return "", ErrGeneratingToken
    }

	return tokenString, nil

}