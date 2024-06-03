package utils

import (
	"net/http"
	"github.com/golang-jwt/jwt"
	"strings"
	"os"
	"fmt"
	"time"
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

	authHeader := req.Header.Get("Authorization")
    if authHeader == "" {
        return "", false
    }

    tokenString, _ := strings.CutPrefix(authHeader, "Bearer ")

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