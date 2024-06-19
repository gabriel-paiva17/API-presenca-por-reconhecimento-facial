package handler

import "errors"

//////////////////
// Modelo do BD //
//////////////////

type User struct {
	ID           string `json:"id" bson:"_id"`
	Username     string `json:"username" bson:"username"`
	Email        string `json:"email" bson:"email"`
	Password     string `json:"password" bson:"password"`
	RegisteredAt string `json:"registeredAt" bson:"registeredAt"`
}

/////////////////////////
// POST /auth/register //
/////////////////////////

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	RegisteredAt string `json:"registeredAt"`
}

//////////////////////
// POST /auth/login //
//////////////////////

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

///////////////////////
// POST /auth/logout //
///////////////////////

type LogoutResponse struct {
	Message string `json:"message"`
	Date    string `json:"date"`
}

//////////////////////
// Erros de usuario //
//////////////////////

var ErrEmailAlreadyExists = errors.New("email already used")
var ErrGeneratingToken = errors.New("cannot generate token")
