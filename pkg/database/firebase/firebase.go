package firebase

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"google.golang.org/api/option"
	"samm/pkg/config"
	"samm/pkg/logger"
)

func NewFirebaseClient(logger logger.ILogger, c *config.FirebaseConfig) *db.Client {
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: c.DatabaseURL,
	}
	// Fetch the service account key JSON file contents
	opt := option.WithCredentialsFile("pkg/database/firebase/katha-dev-firebase-adminsdk-n9vyu-2e9e6ef932.json")

	// Initialize the app with a service account, granting admin privileges
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		logger.Fatalf("error firebase initializing app: %v\n", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		logger.Fatalf("error firebase connecting app: %v\n", err)
		return nil
	}
	return client
}
