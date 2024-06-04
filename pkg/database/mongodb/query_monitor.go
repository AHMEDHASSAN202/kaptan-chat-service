package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/event"
)

var (
	monitor = &event.CommandMonitor{
		Started: func(ctx context.Context, event *event.CommandStartedEvent) {
			fmt.Println("Command started:", event.CommandName, event.Command, event.DatabaseName)
			fmt.Println("****************************************")
		},
		Succeeded: func(ctx context.Context, event *event.CommandSucceededEvent) {
			fmt.Println("Command succeeded", convertNanoToMSec(event.DurationNanos))
			fmt.Println("****************************************")
		},
		Failed: func(ctx context.Context, event *event.CommandFailedEvent) {
			fmt.Println("Command failed:", event.Failure)
			fmt.Println("****************************************")
		},
	}
)

func convertNanoToMSec(nanoValue int64) string {
	Mseconds := float64(nanoValue) / 1e6
	return fmt.Sprintf("it takes => %f ms", Mseconds)
}
