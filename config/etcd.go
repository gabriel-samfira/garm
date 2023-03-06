package config

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"garm/pkg/certificates"
	netUtil "garm/util/net"

	"github.com/BurntSushi/toml"
	zxcvbn "github.com/nbutton23/zxcvbn-go"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	"go.etcd.io/etcd/client/pkg/v3/types"
	"go.etcd.io/etcd/server/v3/embed"
)

// NewEtcdConfig reads in the cluster config file and returns a new EmbededEtcd{} type.
func NewEtcdConfig(cfgFile string) (*EmbededEtcd, error) {
	var config EmbededEtcd
	if _, err := toml.DecodeFile(cfgFile, &config); err != nil {
		return nil, errors.Wrap(err, "decoding toml")
	}

	if err := config.Validate(); err != nil {
		return nil, errors.Wrap(err, "validating config")
	}
	return &config, nil
}

type EmbededEtcd struct {
	Name                     string      `toml:"name"`
	ClusterToken             string      `toml:"cluster_token"`
	PeerListenAddresses      []string    `toml:"peer_listen_addresses"`
	PeerAdvertiseAddresses   []string    `toml:"peer_advertise_addresses"`
	ClientListenAddresses    []string    `toml:"client_listen_addresses"`
	ClientAdvertiseAddresses []string    `toml:"client_advertise_addresses"`
	PeerTLSInfo              EtcdTLSInfo `toml:"peer_tls"`
	ServerTLSInfo            EtcdTLSInfo `toml:"server_tls"`
}

func (e *EmbededEtcd) GetPeerListenAddresses() (types.URLs, error) {
	if len(e.PeerListenAddresses) == 0 {
		return types.NewURLs([]string{"https://0.0.0.0:2380"})
	}

	localAddresses, err := netUtil.GetLocalIPAddresses()
	if err != nil {
		return nil, fmt.Errorf("getting local IP addresses: %w", err)
	}

	withSchema := make([]string, len(e.PeerListenAddresses))
	for idx, val := range e.PeerListenAddresses {
		// Check if a schema was used. Add https if not. We check for http here as well
		// because we want to give a meaningful error message later.
		if !strings.HasPrefix(val, "http://") && !strings.HasPrefix(val, "https://") {
			val = "https://" + val
		}

		parsed, err := url.Parse(val)
		if err != nil {
			return nil, fmt.Errorf("parsing URL: %w", err)
		}
		if parsed.Scheme != "https" {
			return nil, fmt.Errorf("embeded etcd must be over https")
		}

		if parsed.Port() == "" {
			val += ":2380"
		}

		// We want to validate that a local IP address was used here. This is the *listen*
		// address after all.
		addr := parsed.Hostname()
		ipAddr := net.ParseIP(addr)
		if ipAddr == nil {
			return nil, fmt.Errorf("failed to parse IP address: %s", addr)
		}

		// Unless we plan on creating a cluster on the same host, it makes no sense
		// to listen for peer connections on loopback.
		if ipAddr.IsLoopback() {
			return nil, fmt.Errorf("loopback address may not be used here")
		}

		found := false
		for _, localAddr := range localAddresses {
			if localAddr.Equal(ipAddr) {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("address %s not found locally", ipAddr)
		}

		withSchema[idx] = val
	}

	urls, err := types.NewURLs(withSchema)
	if err != nil {
		return nil, errors.Wrap(err, "fetching urls")
	}
	return urls, nil
}

func (e *EmbededEtcd) GetClientListenAddresses() (types.URLs, error) {
	if len(e.ClientListenAddresses) == 0 {
		return types.NewURLs([]string{"https://0.0.0.0:2379"})
	}

	localAddresses, err := netUtil.GetLocalIPAddresses()
	if err != nil {
		return nil, fmt.Errorf("getting local IP addresses: %w", err)
	}

	withSchema := make([]string, len(e.ClientListenAddresses))
	for idx, val := range e.ClientListenAddresses {
		// Check if a schema was used. Add https if not. We check for http here as well
		// because we want to give a meaningful error message later.
		if !strings.HasPrefix(val, "http://") && !strings.HasPrefix(val, "https://") {
			val = "https://" + val
		}

		parsed, err := url.Parse(val)
		if err != nil {
			return nil, fmt.Errorf("parsing URL: %w", err)
		}
		if parsed.Scheme != "https" {
			return nil, fmt.Errorf("embeded etcd must be over https")
		}

		if parsed.Port() == "" {
			val += ":2379"
		}

		// We want to validate that a local IP address was used here. This is the *listen*
		// address after all.
		addr := parsed.Hostname()
		ipAddr := net.ParseIP(addr)
		if ipAddr == nil {
			return nil, fmt.Errorf("failed to parse IP address: %s", addr)
		}

		// For client *listen* addresses, we allow loopback. Anything else, needs to exist
		// on the system.
		if !ipAddr.IsLoopback() {
			found := false
			for _, localAddr := range localAddresses {
				if localAddr.Equal(ipAddr) {
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("address %s not found locally", ipAddr)
			}
		}

		withSchema[idx] = val
	}

	urls, err := types.NewURLs(withSchema)
	if err != nil {
		return nil, errors.Wrap(err, "fetching urls")
	}
	return urls, nil
}

func (e *EmbededEtcd) validateAdvertiseAddress(advertiseAddr, defaultPort string) (string, error) {
	// Check if a schema was used. Add https if not. We check for http here as well
	// because we want to give a meaningful error message later.
	if !strings.HasPrefix(advertiseAddr, "http://") && !strings.HasPrefix(advertiseAddr, "https://") {
		advertiseAddr = "https://" + advertiseAddr
	}

	parsed, err := url.Parse(advertiseAddr)
	if err != nil {
		return "", fmt.Errorf("parsing URL: %w", err)
	}
	if parsed.Scheme != "https" {
		return "", fmt.Errorf("embeded etcd must be over https")
	}

	if parsed.Port() == "" {
		advertiseAddr = fmt.Sprintf("%s:%s", advertiseAddr, defaultPort)
	}

	// We want to validate that a local IP address was used here. Advertising loopback to a
	// remote client is of no help.
	addr := parsed.Hostname()
	if addr == "localhost" {
		return "", fmt.Errorf("loopback address may not be used")
	}

	ipAddr := net.ParseIP(addr)
	if ipAddr != nil {
		if ipAddr.IsLoopback() {
			return "", fmt.Errorf("loopback address may not be used")
		}
	}
	return advertiseAddr, nil
}

func (e *EmbededEtcd) GetPeerAdvertiseAddresses() (types.URLs, error) {
	if len(e.PeerAdvertiseAddresses) == 0 {
		return types.NewURLs([]string{"https://localhost:2380"})
	}

	withSchema := make([]string, len(e.ClientAdvertiseAddresses))
	for idx, val := range e.ClientAdvertiseAddresses {
		var err error
		val, err = e.validateAdvertiseAddress(val, "2380")
		if err != nil {
			return nil, fmt.Errorf("validating address: %w", err)
		}

		withSchema[idx] = val
	}

	urls, err := types.NewURLs(withSchema)
	if err != nil {
		return nil, errors.Wrap(err, "fetching urls")
	}
	return urls, nil
}

func (e *EmbededEtcd) GetClientAdvertiseAddresses() (types.URLs, error) {
	// This is most likely a single node deployment.
	if len(e.ClientAdvertiseAddresses) == 0 {
		return types.NewURLs([]string{"https://localhost:2379"})
	}

	withSchema := make([]string, len(e.ClientAdvertiseAddresses))
	for idx, val := range e.ClientAdvertiseAddresses {
		var err error
		val, err = e.validateAdvertiseAddress(val, "2379")
		if err != nil {
			return nil, fmt.Errorf("validating address: %w", err)
		}

		withSchema[idx] = val
	}

	urls, err := types.NewURLs(withSchema)
	if err != nil {
		return nil, errors.Wrap(err, "fetching urls")
	}
	return urls, nil
}

func (e *EmbededEtcd) Validate() error {
	if e.Name == "" {
		return fmt.Errorf("missing peer name")
	}

	if e.ClusterToken == "" {
		return fmt.Errorf("missing cluster token")
	}

	tokenStenght := zxcvbn.PasswordStrength(e.ClusterToken, nil)
	if tokenStenght.Score < 4 {
		return fmt.Errorf("cluster token is too weak")
	}

	if _, err := e.GetPeerListenAddresses(); err != nil {
		return fmt.Errorf("parsing peer listen address: %w", err)
	}

	if _, err := e.GetPeerAdvertiseAddresses(); err != nil {
		return fmt.Errorf("parsing peer advertise address: %w", err)
	}

	if _, err := e.GetClientListenAddresses(); err != nil {
		return fmt.Errorf("parsing client listen address: %w", err)
	}

	if _, err := e.GetClientAdvertiseAddresses(); err != nil {
		return fmt.Errorf("parsing client advertise address: %w", err)
	}

	return nil
}

func (e *EmbededEtcd) ServerConfig(baseDir string) (*embed.Config, error) {
	dataDir := filepath.Join(baseDir, "etcd")

	if _, err := os.Stat(dataDir); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("accessing data dir: %w", err)
		}
		if err := os.MkdirAll(dataDir, 0700); err != nil {
			return nil, fmt.Errorf("creating data dir: %w", err)
		}
	}

	var err error
	cfg := embed.NewConfig()

	cfg.PeerTLSInfo, err = e.PeerTLSInfo.TLSInfo()
	if err != nil {
		return nil, fmt.Errorf("fetching peer TLS info: %w", err)
	}

	cfg.ClientTLSInfo, err = e.ServerTLSInfo.TLSInfo()
	if err != nil {
		return nil, fmt.Errorf("fetching client TLS info: %w", err)
	}

	cfg.Name = e.Name
	cfg.InitialClusterToken = e.ClusterToken
	cfg.Dir = dataDir
	cfg.AutoCompactionMode = "periodic"
	cfg.AutoCompactionRetention = "5h"
	// TODO(gabriel-samfira): Configure metrics? Listen on localhost and aggregate
	// with the usual garm metrics?

	cfg.LCUrls, err = e.GetClientListenAddresses()
	if err != nil {
		return nil, fmt.Errorf("parsing client listen addresses: %w", err)
	}

	cfg.LPUrls, err = e.GetPeerListenAddresses()
	if err != nil {
		return nil, fmt.Errorf("parsing peer listen addresses: %w", err)
	}

	cfg.ACUrls, err = e.GetClientAdvertiseAddresses()
	if err != nil {
		return nil, fmt.Errorf("parsing client advertise addresses: %w", err)
	}
	cfg.APUrls, err = e.GetPeerAdvertiseAddresses()
	if err != nil {
		return nil, fmt.Errorf("parsing peer advertise addresses: %w", err)
	}

	var initialCluster string
	for _, val := range cfg.APUrls {
		initialCluster += fmt.Sprintf(",%s=%s", e.Name, val.String())
	}

	if len(initialCluster) > 0 {
		cfg.InitialCluster = initialCluster[1:]
	}

	return cfg, nil
}

type EtcdTLSInfo struct {
	CertPath string `toml:"cert"`
	KeyPath  string `toml:"key"`
	CACert   string `toml:"ca_cert"`
}

func (e *EtcdTLSInfo) TLSInfo() (transport.TLSInfo, error) {
	if _, err := tls.LoadX509KeyPair(e.CertPath, e.KeyPath); err != nil {
		return transport.TLSInfo{}, fmt.Errorf("failed to load certificate: %w", err)
	}

	if _, err := certificates.LoadCertFromFile(e.CACert); err != nil {
		return transport.TLSInfo{}, fmt.Errorf("failed to load CA cert: %w", err)
	}

	return transport.TLSInfo{
		CertFile:       e.CertPath,
		KeyFile:        e.KeyPath,
		TrustedCAFile:  e.CACert,
		ClientCertAuth: true,
	}, nil
}
