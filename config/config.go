package config

import "storeapi/utils"

var (
	//ExportFileName definition
	ExportFileName = utils.GetEnv("EXPORT_FILE_NAME", "output.json")
	//APIPort definition
	APIPort = utils.GetEnv("API_PORT", ":8080")
	//Duration period of store save to json file
	Duration = utils.GetEnv("DURATION", "10")
	//MongoDBURI for connecting db
	MongoDBURI = utils.GetEnv("MONGODB_URI", "mongodb+srv://challengeUser:WUMglwNBaydH8Yvu@challenge-xzwqd.mongodb.net/getir-case-study?retryWrites=true")
)
