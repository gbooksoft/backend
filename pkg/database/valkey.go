package database

import (
	"context"
	"fmt"
	"github.com/valkey-io/valkey-go"
	"os"
	"time"
)

var valkeyClient *valkey.Client

func SetUpValkeyConnection(ctx context.Context) (*valkey.Client, error) {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{"localhost:" + os.Getenv("REDIS_PORT")},
	})
	if err != nil {
		return nil, fmt.Errorf("unable to connect to valkey: %w", err)
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	err = client.Do(timeoutCtx, client.B().Ping().Build()).Error()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to Redis: %w", err)
	}
	valkeyClient = &client

	return valkeyClient, nil
}

func GetValkeyClient() *valkey.Client {
	return valkeyClient
}
