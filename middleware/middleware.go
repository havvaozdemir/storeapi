package middleware

import (
	"net/http"
	"storeapi/logger"
	"storeapi/utils"
)

//CheckCors checks cors service
func CheckCors() http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method == "OPTIONS" {
			utils.WriteResponse(resp, http.StatusOK, "")
			return
		}
		utils.WriteResponse(resp, http.StatusMethodNotAllowed, "")
	})
}

//LogRequest logs all requests
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		logger.Info.Printf("%s - %s %s %s", req.RemoteAddr, req.Proto, req.Method, req.URL.RequestURI())
		next.ServeHTTP(resp, req)
	})
}
