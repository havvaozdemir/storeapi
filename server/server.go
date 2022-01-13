package server

import (
	"context"
	"net/http"
	"storeapi/config"
	"storeapi/logger"
	"storeapi/middleware"
	"storeapi/service"
	"storeapi/utils"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func handler(resp http.ResponseWriter, req *http.Request) {
	utils.WriteResponse(resp, http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
}

//SetupDB db connection
func SetupDB(dbURI string) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(dbURI))
	if err != nil {
		logger.Fatal.Println(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		defer client.Disconnect(ctx)
		logger.Fatal.Printf("error, not sent ping to database, %v", err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.Fatal.Println(err)
	}

	return client
}

//NewServer creates store-api server with mongodb and in-memory db
func NewServer() {
	logger.Info.Println("db setup")
	db := SetupDB(config.MongoDBURI)
	store := make(map[string]string)
	service := service.New(store, db)

	logger.Info.Println("Server starting")

	r := mux.NewRouter()
	r.Use(middleware.LogRequest)
	r.HandleFunc("/", handler)
	r.HandleFunc("/in-memory", service.SetStore).Methods("POST")
	r.HandleFunc("/in-memory", service.GetStore).Methods("GET")
	r.HandleFunc("/records", service.FetchRecords).Methods("POST")
	r.MethodNotAllowedHandler = middleware.CheckCors()
	logger.Info.Printf("Server started %s", config.APIPort)

	srv := &http.Server{
		Handler:      r,
		Addr:         config.APIPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal.Fatalf("Could not listen on %s: %v\n", config.APIPort, err)
	}
	logger.Info.Println("Server stopped")
}
