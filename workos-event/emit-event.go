package main

import (
	"context"
	"time"
	"fmt"
	"os"
	"github.com/workos/workos-go/pkg/auditlogs"
)

func main() {
	auditlogs.SetAPIKey(os.Getenv("WORKOS_API_KEY"))

	err := auditlogs.CreateEvent(context.Background(), auditlogs.CreateEventOpts{
		OrganizationID: "org_01GSG5CS0E174RT22N38GEKY2W",
		Event: auditlogs.Event{
			Action:     "user.signed_in",
			OccurredAt: time.Now(),
			Actor: auditlogs.Actor{
				ID:   "user_01GBNJC3MX9ZZJW1FSTF4C5938",
				Type: "user",
			},
			Targets: []auditlogs.Target{
				{ID: "team_01GBNJD4MKHVKJGEWK42JNMBGS", Type: "team"},
			},
			Context: auditlogs.Context{
				Location:  "123.123.123.123",
				UserAgent: "Chrome/104.0.0.0",
			},
		},
	})
	if err != nil {
		fmt.Printf("error")
		fmt.Printf("error: %s", err)
	}
}

