package firebase

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/db"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"kaptan/pkg/config"
	"kaptan/pkg/logger"
	"kaptan/pkg/utils"
	"strings"
	"time"
)

func NewFirebaseClient(logger logger.ILogger, c *config.FirebaseConfig) (*db.Client, *auth.Client) {
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: c.DatabaseURL,
	}
	// Fetch the service account key JSON file contents
	opt := option.WithCredentialsFile("")

	// Initialize the app with a service account, granting admin privileges
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		logger.Fatalf("error firebase initializing app: %v\n", err)
	}
	auth, err := app.Auth(ctx)
	if err != nil {
		logger.Fatalf("error firebase auth client: %v\n", err)
		return nil, nil
	}
	client, err := app.Database(ctx)
	if err != nil {
		logger.Fatalf("error firebase database connecting app: %v\n", err)
		return nil, nil
	}

	//todo: add to cronjob remove all expired data
	removeExpiredNodesFromDB(ctx, client)
	return client, auth
}

func removeExpiredNodesFromDB(ctx context.Context, client *db.Client) {
	ref := client.NewRef("chats/ttl")

	var items map[string]interface{}
	if err := ref.Get(ctx, &items); err != nil {
		logrus.Fatalf("error firebase database connecting app: %v\n", err)
	}

	for key, value := range items {
		doc := parseStringToMap(key)
		if t, err := time.Parse(utils.DefaultDateTimeFormat, value.(string)); err == nil && t.Before(time.Now()) {
			deleteExpiredChats(key, doc, client, ctx)
		}
	}
}

func deleteExpiredChats(key string, doc map[string]interface{}, client *db.Client, ctx context.Context) {
	err := client.NewRef("orders/raw").Child(doc["orderId"].(string)).Delete(ctx)
	if err != nil {
		logrus.Error("raw: sync to firebase order => ", err)
	}

	err = client.NewRef("orders/users").Child(doc["userId"].(string)).Child(doc["orderId"].(string)).Delete(ctx)
	if err != nil {
		logrus.Error("users: sync to firebase order => ", err)
	}

	err = client.NewRef("orders/ttl").Child(key).Delete(ctx)
	if err != nil {
		logrus.Error("users: sync to firebase order => ", err)
	}

	for _, kitchenId := range doc["kitchenId"].([]string) {
		err = client.NewRef("orders/kitchens").Child(kitchenId).Child(doc["orderId"].(string)).Delete(ctx)
		if err != nil {
			logrus.Error("kitchens: sync to firebase order => ", err)
		}
	}
}

func parseStringToMap(input string) map[string]interface{} {
	parts := strings.Split(input, " ")
	result := make(map[string]interface{})
	kitchenIds := []string{}

	for i := 0; i < len(parts); i++ {
		switch parts[i] {
		case "^orderId":
			if i+1 < len(parts) {
				result["orderId"] = parts[i+1]
				i++
			}
		case "^userId":
			if i+1 < len(parts) {
				result["userId"] = parts[i+1]
				i++
			}
		case "^kitchenId":
			if i+1 < len(parts) {
				kitchenIds = append(kitchenIds, parts[i+1])
				i++
			}
		}
	}

	result["kitchenId"] = kitchenIds
	return result
}
