package mongodb

import (
	"context"
	"fmt"
	"log"
	"nudam-ctrader-api/helpers/configs_helper"
	"nudam-ctrader-api/logger"
	"nudam-ctrader-api/types/ctrader"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient     *mongo.Client
	mongoClientCtx  context.Context
	mongoClientOnce sync.Once
)

func GetMongoClient() (*mongo.Client, context.Context, error) {
	var err error
	mongoClientOnce.Do(func() {
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		clientOptions := options.Client().ApplyURI(configs_helper.MongoDb.Uri).SetServerAPIOptions(serverAPI)
		mongoClientCtx := context.TODO()
		mongoClient, err = mongo.Connect(mongoClientCtx, clientOptions)
		if err != nil {
			log.Panic(err)
			return
		}

		if err := mongoClient.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
			log.Panic(err)
			return
		}
	})

	return mongoClient, mongoClientCtx, err
}

// Saves interface to MongoDb based on collection in configs.
func SaveToMongo(doc interface{}, payloadType int) error {
	client, ctx, err := GetMongoClient()
	if err != nil {
		return err
	}

	coll := client.Database(configs_helper.MongoDb.DatabaseName).Collection(configs_helper.MongoDb.CollectionName)

	filter := bson.M{"payloadtype": payloadType}
	opts := options.Replace().SetUpsert(true)

	result, err := coll.ReplaceOne(ctx, filter, doc, opts)
	if err != nil {
		return err
	}

	if result.UpsertedID != nil {
		logger.LogMessage(fmt.Sprintf("new doc added: %v", result.UpsertedID))
	} else {
		logger.LogMessage("doc updated")
	}

	return nil
}

// Takes symbolIds from mongodb based on currency pairs in constants.json.
func FindSymbolIds(symbolNames []string, payloadType int) ([]int64, error) {
	client, ctx, err := GetMongoClient()
	if err != nil {
		return nil, err
	}

	coll := client.Database(configs_helper.MongoDb.DatabaseName).Collection(configs_helper.MongoDb.CollectionName)
	filter := bson.M{"payloadtype": payloadType}

	var result ctrader.Message[ctrader.ProtoOASymbolsListRes]
	if err = coll.FindOne(ctx, filter).Decode(&result); err != nil {
		return nil, err
	}

	symbolIds := make([]int64, 0, len(symbolNames))
	for _, symbolName := range symbolNames {
		found := false
		for _, symbol := range result.Payload.Symbol {
			if *symbol.SymbolName == symbolName {
				symbolIds = append(symbolIds, symbol.SymbolId)
				found = true
				break
			}
		}
		if !found {
			return nil, err
		}
	}

	return symbolIds, nil
}
