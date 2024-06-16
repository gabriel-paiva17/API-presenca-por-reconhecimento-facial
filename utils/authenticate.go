package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func Authenticate(next func(res http.ResponseWriter, req *http.Request)) func(res http.ResponseWriter, req *http.Request) {

	return func(res http.ResponseWriter, req *http.Request) {
		// Verificar a autenticação aqui
		_, authenticated := GetAuthenticatedUserId(req)

		if !authenticated {

			WriteErrorResponse(res, http.StatusUnauthorized, "Requisicao necessita autenticacao.")
			return

		}

		// Se autenticado, chamar a função original
		next(res, req)

	}
}

func GetAuthenticatedUserId(req *http.Request) (string, bool) {

	jwtKey := []byte(os.Getenv("SECRET_KEY"))

	cookie, err := req.Cookie("auth-token")
    if err != nil {
        return "", false
    }

    tokenString := cookie.Value

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metodo de assinatura nao corresponde")
		}

		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return "", false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", false
	}

	userID, ok := claims["userId"].(string)
	if !ok || userID == "" {
		return "", false
	}

	exp, ok := claims["exp"].(float64)
	if !ok || time.Unix(int64(exp), 0).Before(time.Now()) {
		return "", false
	}

	return userID, true

}

// NOT WORKING IN FRONTEND.
// function that only guarantees the user authentication
func CheckAuthenthentication() func(res http.ResponseWriter, req *http.Request) {

	return Authenticate(func(res http.ResponseWriter, req *http.Request) {

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(map[string]interface{}{
			"message": "Usuário autenticado.",
		})
		
	})
}
