package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

var errorMap = map[int]string{
	http.StatusBadRequest:          "Bad Request",
	http.StatusUnauthorized:        "Unauthorized",
	http.StatusForbidden:			"Forbidden",
	http.StatusNotFound:            "Not Found",
	http.StatusConflict:            "Conflict",
	http.StatusInternalServerError: "Internal Server Error",
	http.StatusUnsupportedMediaType:"Unsupported Media Type", 
	http.StatusUnprocessableEntity: "Unprocessable Entity",
}

func WriteErrorResponse(res http.ResponseWriter, statusCode int, Message string) {

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)

	errorText, ok := errorMap[statusCode]
	if !ok {
		errorText = "Undefined Error in errorMap"
	}

	errResponse := ErrorResponse{Error: errorText, Message: Message}
	json.NewEncoder(res).Encode(errResponse)
}
