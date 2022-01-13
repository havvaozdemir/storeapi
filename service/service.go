package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"storeapi/data"
	types "storeapi/errors"
	"storeapi/utils"

	"go.mongodb.org/mongo-driver/mongo"
)

//Manager describes handler service interface
type Manager interface {
	SetStore(resp http.ResponseWriter, req *http.Request)
	GetStore(resp http.ResponseWriter, req *http.Request)
	FetchRecords(resp http.ResponseWriter, req *http.Request)
}

//service describes properties for api
type service struct {
	data data.Repository
}

//New creates new service
func New(store map[string]string, db *mongo.Client) Manager {
	return &service{data: data.New(store, db)}
}

// swagger:route POST /in-memory body addStore
// Sets key value from store
// responses:
// 201: json
// 400: KeyRequired, ValueRequired, Reading body errors,

//SetStore sets the store with request json key-value body
func (h *service) SetStore(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	post := data.RequestData{}
	if err := json.NewDecoder(req.Body).Decode(&post); err != nil {
		utils.WriteResponse(resp, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if post.Key == "" {
		utils.WriteResponse(resp, http.StatusBadRequest, types.KeyRequired)
		return
	}
	if post.Value == "" {
		utils.WriteResponse(resp, http.StatusBadRequest, types.ValueRequired)
		return
	}
	h.data.Add(post.Key, post.Value)
	result := data.RequestData{}
	result.Value = post.Value
	result.Key = post.Key
	utils.WriteResponse(resp, http.StatusCreated, result)
}

// swagger:route GET /in-memory key queryparam
// Return get the value of key from store
// responses:
// 200: json
// 404: KeyNotFound
// 400: KeyRequired

//GetStore return get the value of key from store
func (h *service) GetStore(resp http.ResponseWriter, req *http.Request) {
	var key string
	if keyParam, ok := req.URL.Query()["key"]; ok {
		key = keyParam[0]
	}
	if key == "" {
		utils.WriteResponse(resp, http.StatusBadRequest, types.KeyRequired)
		return
	}

	if value, ok := h.data.Get(key); ok {
		result := data.RequestData{}
		result.Value = value
		result.Key = key
		utils.WriteResponse(resp, http.StatusOK, result)
		return
	}

	utils.WriteResponse(resp, http.StatusNotFound, fmt.Sprintf(types.KeyNotFound, key))
}

// swagger:route POST /records with requestbody
// Return fetch the records from mongodb using body params
// responses:
// 200: json
// 400: StatusBadRequest

//FetchRecords return get the value of key from store
func (h *service) FetchRecords(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	post := data.RequestRecord{}
	if err := json.NewDecoder(req.Body).Decode(&post); err != nil {
		utils.WriteResponse(resp, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if post.StartDate == "" {
		utils.WriteResponse(resp, http.StatusBadRequest, types.StartDateRequired)
		return
	}
	if post.EndDate == "" {
		utils.WriteResponse(resp, http.StatusBadRequest, types.EndDateRequired)
		return
	}
	if post.MaxCount <= 0 {
		utils.WriteResponse(resp, http.StatusBadRequest, types.MaxCountRequired)
		return
	}
	if post.MinCount <= 0 {
		utils.WriteResponse(resp, http.StatusBadRequest, types.MinCountRequired)
		return
	}
	records, err := h.data.FetchRecords(post)
	if err != nil {
		utils.WriteResponse(resp, http.StatusBadRequest, err)
		return
	}
	utils.WriteResponse(resp, http.StatusOK, records)
}
