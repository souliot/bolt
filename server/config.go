package server

import (
	"path/filepath"
	"time"
)

const (
	// DefaultQuotaBytes is the number of bytes the backend Size may
	// consume before exceeding the space quota.
	DefaultQuotaBytes = int64(2 * 1024 * 1024 * 1024) // 2GB
)

type ServerConfig struct {
	Name               string
	DataDir            string
	BootstrapTimeout   time.Duration
	AutoCompactionMode string
	QuotaBackendBytes  int64
}

func (c *ServerConfig) SnapDir() string     { return filepath.Join(c.DataDir, "snap") }
func (c *ServerConfig) backendPath() string { return filepath.Join(c.SnapDir(), "db") }

// ReqTimeout returns timeout for request to finish.
func (c *ServerConfig) ReqTimeout() time.Duration {
	// 5s for queue waiting, computation and disk IO delay
	// + 2 * election timeout for possible leader election
	return 5*time.Second + 2*time.Duration(1000)*time.Millisecond
}

func (c *ServerConfig) bootstrapTimeout() time.Duration {
	if c.BootstrapTimeout != 0 {
		return c.BootstrapTimeout
	}
	return time.Second
}
