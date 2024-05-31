package utils

import "net/http"

func Authenticate(next func(res http.ResponseWriter, req *http.Request)) func(res http.ResponseWriter, req *http.Request) {
	
	return func(res http.ResponseWriter, req *http.Request) {
		// Verificar a autenticação aqui
		authenticated := isAuthenticated(req)

		if !authenticated {
			// Se autenticado, chamar a função original
			WriteErrorResponse(res, http.StatusUnauthorized, "Acesso necessita autenticacao.")
		}

		next(res, req)		

	}
}

func isAuthenticated(req *http.Request) bool {

	return false

}