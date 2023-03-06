package etcd

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/server/v3/embed"
)

// init new instance:
// * is this a new cluster
//   - skip to questions about addresses
//   - on which address(es) should we listen on for inbound connections? (default: 0.0.0.0)
//   - on which IP address(es) can peers connect to us? (this can be a floating IP, a local
//     address or a DNS A record)
//
// * are we joining an existing cluster
//
//   - on which address(es) should we listen on for inbound connections? (default: 0.0.0.0)
//
//   - on which IP address(es) can peers connect to us? (this can be a floating IP, a local
//     address or a DNS A record)
//
//   - input join token. This is a JWT which includes all the info needed to connect to existing members
//     and request certificates, peer info, etc.
//
//   - after joining the etcd cluster, the node will need to be unsealed. Unsealing means that the node has
//     temporary access to a key which can decrypt the root key used to encrypt all other keys on the system.
//     The data at rest is encrypted with these keys using AES-GCM. The unencrypted root key is kept in memory
//     until the unit is restarted or explicitly sealed by an operator.
//     Unseal key --> root key --> encryption keys --> data. Besides the unseal key, everything else is saved
//     in the etcd cluster. A new node may join and sync all data, but without the unseal key, it won't have
//     access to any of the data.
//
//   - Input unseal passphrase (random, 32 characters long). This is used to decrypt keys,
//     which in turn will decrypt data in the cluster. Use Shamir secret sharing to split the unseal key
//     into multiple keys.
//
// * Save unseal passphrase in the startup config? If not, member will have to be unsealed using shamir keys. (5 keys total, min 3 to unseal)
// * on which address(es) should we listen on for inbound connections? (default: 0.0.0.0)
// * On which IP address(es) can peers connect to us? (this can be a floating IP, a local address or a DNS A record)
func NewEtcdServer(ctx context.Context, conf *embed.Config) (*etcdServer, error) {
	// srv, err := embed.StartEtcd(conf)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to start etcd server: %w", err)
	// }
	return &etcdServer{
		ctx:    ctx,
		config: conf,
		// srv:    srv,
	}, nil
}

type etcdServer struct {
	ctx    context.Context
	config *embed.Config
	srv    *embed.Etcd
}

func (e *etcdServer) WaitReady() error {
	select {
	case <-e.srv.Server.ReadyNotify():
	case <-time.After(3600 * time.Second):
		e.srv.Server.Stop()
		e.srv.Close()
		srvErr := <-e.srv.Err()
		return fmt.Errorf("timed out waiting for server to be ready: %s", srvErr)
	}
	return nil
}
