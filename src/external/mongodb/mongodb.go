package mongodb

import (
	"context"
	"log"
	"nudam-ctrader-api/helpers/configs_helper"
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

type MongoDbData struct {
	SymbolId    int64     `bson:"symbolId" json:"symbolId"`
	SymbolName  string    `bson:"symbolName" json:"symbolName"`
	PipPosition int32     `bson:"pipPosition" json:"pipPosition"`
	StepVolume  int64     `bson:"stepVolume" json:"stepVolume"`
	LotSize     int64     `bson:"lotSize" json:"lotSize"`
	Prices      PriceData `bson:"prices" json:"prices"`
	Ema         []Ema     `bson:"ema" json:"ema"`
}

type PriceData struct {
	Bid uint64 `bson:"bid" json:"bid"`
	Ask uint64 `bson:"ask" json:"ask"`
}

type Ema struct {
	Period int64            `bson:"period" json:"period"`
	Values map[string]int64 `bson:"values" json:"values"`
}

// Return mongodb client.
func GetMongoClient() (*mongo.Client, context.Context, error) {
	var err error
	mongoClientOnce.Do(func() {
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		clientOptions := options.Client().ApplyURI(configs_helper.MongoDbConfig.Uri).SetServerAPIOptions(serverAPI)
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

// Save interface to MongoDb based on collection in configs.
func SaveToMongo(doc interface{}, filter bson.M) error {
	client, ctx, err := GetMongoClient()
	if err != nil {
		return err
	}

	coll := client.Database(configs_helper.MongoDbConfig.DatabaseName).Collection(configs_helper.MongoDbConfig.Collection)

	opts := options.Replace().SetUpsert(true)
	_, err = coll.ReplaceOne(ctx, filter, doc, opts)
	if err != nil {
		return err
	}

	return nil
}

// Update doc in MongoDb based on collection in configs and with specific filter and update bson.M.
func UpdateMongo(filter, update bson.M) error {
	client, ctx, err := GetMongoClient()
	if err != nil {
		return err
	}

	coll := client.Database(configs_helper.MongoDbConfig.DatabaseName).Collection(configs_helper.MongoDbConfig.Collection)

	opts := options.Update().SetUpsert(true)
	_, err = coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

// Take symbolId from mongodb based on symbolName.
func FindSymbolId(symbolName string) (int64, error) {
	client, ctx, err := GetMongoClient()
	if err != nil {
		return 0, err
	}

	coll := client.Database(configs_helper.MongoDbConfig.DatabaseName).Collection(configs_helper.MongoDbConfig.Collection)

	var result ctrader.SymbolList
	if err = coll.FindOne(ctx, bson.M{"symbolName": symbolName}).Decode(&result); err != nil {
		return 0, err
	}

	return result.SymbolId, nil
}

// Take emas from mongodb based on symbolName.
func FindEmas(symbolName string) ([]Ema, error) {
	client, ctx, err := GetMongoClient()
	if err != nil {
		return nil, err
	}

	coll := client.Database(configs_helper.MongoDbConfig.DatabaseName).Collection(configs_helper.MongoDbConfig.Collection)

	var result MongoDbData
	if err = coll.FindOne(ctx, bson.M{"symbolName": symbolName}).Decode(&result); err != nil {
		return nil, err
	}

	return result.Ema, nil
}

// Take bid and ask price based on symbolName.
func FindBidAsk(symbolName string) (PriceData, error) {
	client, ctx, err := GetMongoClient()
	if err != nil {
		return PriceData{}, err
	}

	coll := client.Database(configs_helper.MongoDbConfig.DatabaseName).Collection(configs_helper.MongoDbConfig.Collection)

	var result MongoDbData
	if err = coll.FindOne(ctx, bson.M{"symbolName": symbolName}).Decode(&result); err != nil {
		return PriceData{}, err
	}

	return result.Prices, nil
}
