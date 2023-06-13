package config

import (
	// "crypto/tls"
	// "fmt"

	"crypto/tls"
	"fmt"
	"time"

	"github.com/coreos/etcd/pkg/transport"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

var (
	EtcdDialTimeout = 5 * time.Second
)

type EtcdClient struct {
	Endpoints []string `toml:"endpoints" json:"endpoints"`
	Username  string   `toml:"username" json:"username"`
	Password  string   `toml:"password" json:"password"`

	CertFile  string `toml:"cert_file" json:"cert-file"`
	KeyFile   string `toml:"key_file" json:"key-file"`
	TrustedCA string `toml:"trusted_ca" json:"trusted-ca"`
	Insecure  bool   `toml:"insecure" json:"insecure"`
}

func (e *EtcdClient) TLSConfig() (*tls.Config, error) {
	tlsInfo := transport.TLSInfo{
		CertFile:           e.CertFile,
		KeyFile:            e.KeyFile,
		TrustedCAFile:      e.TrustedCA,
		InsecureSkipVerify: e.Insecure,
	}
	cfg, err := tlsInfo.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to generate tls client config: %w", err)
	}
	return cfg, nil
}

func (e *EtcdClient) ClientConfig() (clientv3.Config, error) {
	if len(e.Endpoints) == 0 {
		return clientv3.Config{}, fmt.Errorf("missing endpoints")
	}

	tlsCfg, err := e.TLSConfig()
	if err != nil {
		return clientv3.Config{}, fmt.Errorf("failed to generate tls client config: %w", err)
	}

	clientCfg := clientv3.Config{
		TLS:         tlsCfg,
		DialTimeout: EtcdDialTimeout,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},

		Endpoints: e.Endpoints,
		Username:  e.Username,
		Password:  e.Password,
	}
	return clientCfg, nil
}
