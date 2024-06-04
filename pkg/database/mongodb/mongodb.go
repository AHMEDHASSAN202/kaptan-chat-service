package mongodb

import (
	"context"
	"example.com/fxdemo/pkg/config"
	"fmt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func CreateIndex(collectionConnection *mongo.Collection, unique bool, fields ...string) bool {
	// 1. Lets define the keys for the index we want to create
	var keys bson.D
	for _, field := range fields {
		keys = append(keys, bson.E{field, 1})
	}
	mod := mongo.IndexModel{
		Keys:    keys, // index in ascending order or -1 for descending order
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
