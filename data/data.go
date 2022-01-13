package data

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"storeapi/config"
	types "storeapi/errors"
	"storeapi/logger"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//Repository interface for data
type Repository interface {
	Add(key, value string)
	Get(key string) (string, bool)
	FetchRecords(request RequestRecord) (ResponseRecord, error)
}

//data properties
type data struct {
	store   map[string]string
	mtx     sync.Mutex
	mongoDB *mongo.Client
}

//New creates Data, starts and loads data from file
func New(store map[string]string, db *mongo.Client) Repository {
	d := &data{store: store, mongoDB: db}
	d.startDataService()
	d.loadExistingDataFromFile()
	return d
}

//RequestData gets jsonbody from post request
type RequestData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

//RequestRecord definition
type RequestRecord struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	MinCount  int    `json:"minCount"`
	MaxCount  int    `json:"maxCount"`
}

//Record definition
type Record struct {
	Key        string    `json:"key" bson:"key"`
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
	TotalCount int       `json:"totalCount" bson:"totalCount"`
}

//ResponseRecord definition
type ResponseRecord struct {
	Code    int       `json:"code"`
	Msg     string    `json:"msg"`
	Records []*Record `json:"records"`
}

//startDataService starts data service using duration period
func (d *data) startDataService() {
	quit := make(chan bool)
	duration, err := strconv.Atoi(config.Duration)
	if err != nil {
		logger.Error.Println(err)
		duration = 10
	}
	ticker := time.NewTicker(time.Duration(duration) * time.Minute)

	go d.saveStoreDataToFile(ticker, quit)
}

//loadExistingDataFromFile loads existing json file and sets datastore
func (d *data) loadExistingDataFromFile() {
	storeFile, err := os.Open(config.ExportFileName)
	if err != nil {
		return
	}
	d.mtx.Lock()
	defer d.mtx.Unlock()
	byteValue, _ := ioutil.ReadAll(storeFile)
	err = json.Unmarshal(byteValue, &d.store)
	if err != nil {
		logger.Error.Println(err)
	}
}

//saveStoreDataToFile saves stores data to file from ticker
func (d *data) saveStoreDataToFile(ticker *time.Ticker, quit chan bool) {
	for {
		select {
		case <-ticker.C:
			_, err := os.OpenFile(config.ExportFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

			if err != nil {
				logger.Error.Println(err)
			}
			dataBytes, err := json.Marshal(d.store)
			if err != nil {
				logger.Error.Println(err)
			}
			err = ioutil.WriteFile(config.ExportFileName, dataBytes, 0777)
			if err != nil {
				logger.Error.Println(err)
			} else {
				logger.Info.Println("Sucsessfully added store data to file")
			}

		case <-quit:
			ticker.Stop()
			return
		}
	}
}

//Add adds data to store
func (d *data) Add(key, value string) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	d.store[key] = value
}

//Get gets value of given key from store
func (d *data) Get(key string) (string, bool) {
	value, ok := d.store[key]
	return value, ok
}

//FetchRecords gets all movies from datastore with given filters
func (d *data) FetchRecords(reqRecord RequestRecord) (ResponseRecord, error) {
	result := ResponseRecord{}

	ctx := context.Background()
	startDate, err := time.Parse("2006-01-02", reqRecord.StartDate)
	if err != nil {
		return result, errors.New(types.DateFormat)
	}
	endDate, err := time.Parse("2006-01-02", reqRecord.EndDate)
	if err != nil {
		return result, errors.New(types.DateFormat)
	}
	matchStage := []bson.M{
		{
			"$match": bson.M{
				"createdAt": bson.M{
					"$gt": startDate,
					"$lt": endDate}},
		},
		{"$project": bson.M{
			"_id":        0,
			"key":        "$key",
			"createdAt":  "$createdAt",
			"totalCount": bson.M{"$sum": "$counts"}},
		},
		{
			"$match": bson.M{
				"totalCount": bson.M{
					"$gt": reqRecord.MinCount,
					"$lt": reqRecord.MaxCount}},
		},
		{
			"$sort": bson.M{
				"totalCount": 1,
			},
		},
	}

	coll := d.mongoDB.Database("getir-case-study").Collection("records")
	cur, err := coll.Aggregate(ctx, matchStage)
	if err != nil {
		logger.Error.Println(err)
		return result, err
	}
	defer cur.Close(ctx)
	resArr := []*Record{}
	for cur.Next(ctx) {
		record := &Record{}
		if err = cur.Decode(record); err != nil {
			logger.Error.Println(err)
			return result, err
		}
		resArr = append(resArr, record)
	}
	result.Code = 0
	result.Msg = "Success"
	result.Records = resArr
	return result, nil
}
