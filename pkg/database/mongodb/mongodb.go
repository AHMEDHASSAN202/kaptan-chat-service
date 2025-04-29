package mongodb

import (
	"context"
	"fmt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kaptan/pkg/config"
	"os"
	"time"
)

const timeout = 10 * time.Second

// NewClient established connection to a mongoDb instance using provided URI and auth credentials.
func NewClient(c *config.MongoConfig) (*mongo.Client, *mongo.Database, error) {
	opts := options.Client().ApplyURI(c.MongoConnection)
	logEnabled := os.Getenv("ENABLE_LOG_DB")
	if logEnabled == "true" {
		opts.SetMonitor(monitor)
	}
	err := mgm.SetDefaultConfig(nil, c.MongoDbName, opts)
	if err != nil {
		return nil, nil, err
	}

	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, nil, err
	}

	db := client.Database(c.MongoDbName)
	return client, db, nil
}

var (
	IndexType = struct {
		Asc     int
		Desc    int
		Text    string
		Spatial string
	}{
		Asc:     1,
		Desc:    -1,
		Text:    "text",
		Spatial: "2dsphere",
	}
)

func CreateIndex(collectionConnection *mongo.Collection, unique bool, fields ...bson.E) bool {
	// 1. Lets define the keys for the index we want to create
	var keys bson.D
	for _, field := range fields {
		keys = append(keys, field)
	}
	mod := mongo.IndexModel{
		Keys:    keys, // index in ascending chat or -1 for descending chat
		Options: options.Index().SetUnique(unique),
	}
	// 2. Create the context for this operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 3. Connect to the database and access the collection
	//collectionConnection= reflect.TypeOf(collectionConnection)
	collection := collectionConnection //.(mongo.Collection)

	// 4. Create a single index
	_, err := collection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		// 5. Something went wrong, we log it and return false
		fmt.Println(err.Error())
		return false
	}

	// 6. All went well, we return true
	return true
}

func CreateIndexWithTTL(collectionConnection *mongo.Collection, unique bool, field bson.E, expireAfterSeconds int32) bool {
	mod := mongo.IndexModel{
		Keys:    bson.D{field}, // index in ascending chat or -1 for descending chat
		Options: options.Index().SetUnique(unique).SetExpireAfterSeconds(expireAfterSeconds),
	}
	// 2. Create the context for this operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 3. Connect to the database and access the collection
	//collectionConnection= reflect.TypeOf(collectionConnection)
	collection := collectionConnection //.(mongo.Collection)

	// 4. Create a single index
	_, err := collection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		// 5. Something went wrong, we log it and return false
		fmt.Println(err.Error())
		return false
	}

	// 6. All went well, we return true
	return true
}
