package client

import (
	"context"
	"testing"

	"github.com/souliot/bolt/server"
)

var (
	srv *server.Server
	cli *Client
)

func init() {
	var err error
	scfg := server.ServerConfig{
		Name:    "first",
		DataDir: "data",
	}
	srv, err = server.NewServer(scfg)
	if err != nil {
		return
	}
	ccfg := Config{
		Srv: srv,
	}
	cli, err = NewWithConfig(ccfg)
	if err != nil {
		return
	}
}

func TestClientPut(t *testing.T) {
	defer srv.Close()
	ctx := context.Background()
	resp, err := cli.Put(ctx, "asd", "asf1")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
func TestClientGet(t *testing.T) {
	defer srv.Close()
	ctx := context.Background()
	resp, err := cli.Get(ctx, "a11y-profile-manager-indicator")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
func TestClientGetWithPrefix(t *testing.T) {
	defer srv.Close()
	ctx := context.Background()
	resp, err := cli.Get(ctx, "as", WithPrefix())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
func TestClientDel(t *testing.T) {
	defer srv.Close()
	ctx := context.Background()
	resp, err := cli.Delete(ctx, "asf")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
func TestClientDelWithPrefix(t *testing.T) {
	defer srv.Close()
	ctx := context.Background()
	resp, err := cli.Delete(ctx, "as", WithPrefix())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
func TestClientCompact(t *testing.T) {
	defer srv.Close()
	ctx := context.Background()
	resp, err := cli.Compact(ctx, 21, WithCompactPhysical())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
