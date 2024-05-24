package user

import (
    "net/http"
    "time"
    "encoding/json"
    "github.com/google/uuid" // Para gerar IDs únicos
    "main/utils"     // Importe o pacote de utilitários para hash de senha e validação de e-mail

)

func (u *User) CreateUser(res http.ResponseWriter, req *http.Request) {
    // Decodificar o JSON da requisição para a estrutura CreateUserBody
    var createUserBody CreateUserBody
    if err := json.NewDecoder(req.Body).Decode(&createUserBody); err != nil {
        http.Error(res, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validar o e-mail
    if !utils.IsValidEmail(createUserBody.Email) {
        http.Error(res, "Invalid email format", http.StatusBadRequest)
        return
    }

    // Hash da senha (suponha que você tenha uma função HashPassword em utils)
    hashedPassword, err := utils.HashPassword(createUserBody.Password)
    if err != nil {
        http.Error(res, "Error hashing password", http.StatusInternalServerError)
        return
    }

    // Gerar ID único para o usuário
    userID := uuid.New().String()

    // Criar a instância do usuário com os dados fornecidos na requisição
    newUser := User{
        ID:          userID,
        Username:    createUserBody.Username,
        Email:       createUserBody.Email,
        Password:    hashedPassword,
        RegisteredAt: time.Now().Format(time.RFC3339), // Timestamp atual
    }

    // Serializar o usuário criado como resposta
    createUserResponse := CreateUserResponse{
        ID:          newUser.ID,
        Username:    newUser.Username,
        Email:       newUser.Email,
        RegisteredAt: newUser.RegisteredAt,
    }

    // Serializar a resposta como JSON e enviar
    res.Header().Set("Content-Type", "application/json")
    res.WriteHeader(http.StatusCreated)
    json.NewEncoder(res).Encode(map[string]interface{}{
        "message": "Usuário criado com sucesso.",
        "user": createUserResponse,
    })
}
