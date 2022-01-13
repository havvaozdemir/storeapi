// Package Classification Store API
//
// RESTfull API that includes endpoints for fetching data from the provided MongoBD and getting, setting key-value pairs by in-memory store.
//
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
// Contact: Havva Ozdemir <havvaozdemir34@gmail.com>
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package main

import (
	"storeapi/logger"
	"storeapi/server"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			logger.Fatal.Printf("Failed: (%v)", r)
		}
	}()
	server.NewServer()
}
