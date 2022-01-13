package utils

import (
	"encoding/json"
	"net/http"
	"os"
	"storeapi/logger"
)

//GetEnv gets enviroment variables
func GetEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

//WriteResponse writes response with given status and body
func WriteResponse(resp http.ResponseWriter, statusCode int, value interface{}) {
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, DELETE, POST")
	resp.WriteHeader(statusCode)
	if err := json.NewEncoder(resp).Encode(value); err != nil {
		logger.Error.Println(err)
		return
	}

}
