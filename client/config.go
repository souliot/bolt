package client

import (
	"context"

	"github.com/souliot/bolt/server"
)

type Config struct {
	Context context.Context
	Srv     *server.Server
}
