package client

import (
	"context"
	"sync"

	"github.com/souliot/bolt/server"
)

type Client struct {
	KV
	Lease
	Watcher
	cfg Config
	mu  *sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc
}

// New creates a new etcdv3 client from a given configuration.
func NewWithConfig(cfg Config) (cli *Client, err error) {
	return newClient(&cfg)
}

func New(name, dir string) (cli *Client, err error) {
	scfg := server.ServerConfig{
		Name:    name,
		DataDir: dir,
	}
	srv, err := server.NewServer(scfg)
	if err != nil {
		return
	}
	ccfg := Config{
		Srv: srv,
	}
	return newClient(&ccfg)
}

func (c *Client) Close() error {
	c.cancel()
	// TODO
	// c.Watcher.Close()
	// c.Lease.Close()
	return c.ctx.Err()
}

func (c *Client) RawServer() *server.Server {
	return c.cfg.Srv
}

func (c *Client) CloseServer() {
	c.cfg.Srv.Close()
}

func (c *Client) Ctx() context.Context { return c.ctx }

func newClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		cfg = &Config{}
	}

	// use a temporary skeleton client to bootstrap first connection
	baseCtx := context.TODO()
	if cfg.Context != nil {
		baseCtx = cfg.Context
	}

	ctx, cancel := context.WithCancel(baseCtx)
	client := &Client{
		cfg:    *cfg,
		ctx:    ctx,
		cancel: cancel,
		mu:     new(sync.RWMutex),
	}

	client.KV = NewKV(cfg.Srv)
	// TODO
	// client.Lease = NewLease(client)
	// client.Watcher = NewWatcher(client)

	return client, nil
}

func toErr(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	if ctx.Err() != nil {
		err = ctx.Err()
	}

	return err
}
