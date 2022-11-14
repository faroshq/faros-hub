package main

import (
	"context"

	"github.com/faroshq/faros-hub/pkg/config"
	"github.com/faroshq/faros-hub/pkg/server"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		panic(err)
	}
}

func run(ctx context.Context) error {
	c, err := config.LoadAPI()
	if err != nil {
		return err
	}

	server, err := server.New(c)
	if err != nil {
		return err
	}
	return server.Run(ctx)
}
