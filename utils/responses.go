package utils

import (
    "net/http"
    "encoding/json"
)

type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
}

var errorMap = map[int]string{
    http.StatusBadRequest:          "Bad Request",
    http.StatusUnauthorized:        "Unauthorized",
    http.StatusNotFound:            "Not Found",
    http.StatusInternalServerError: "Internal Server Error",
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