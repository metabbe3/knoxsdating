package helpers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success    bool        `json:"success"`
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Error      interface{} `json:"error"`
}

func GenerateResponse(success bool, statusCode int, message string, data interface{}, err interface{}) Response {
	return Response{
		Success:    success,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Error:      err,
	}
}

func SendJSONResponse(w http.ResponseWriter, statusCode int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
