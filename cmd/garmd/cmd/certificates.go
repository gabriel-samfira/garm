// Copyright 2022 Cloudbase Solutions SRL
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package cmd

import (
	"fmt"
	"garm/pkg/certificates"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	peerName        string
	certificateSANs []string
)

// certificatesCmd represents the certificates command
var certificatesCmd = &cobra.Command{
	Use:     "certificates",
	Aliases: []string{"cert"},
	Short:   "Generate etcd x509 certificates",
	Long:    `Generate x509 certificates that are suitable for the embeded etcd server.`,
	Run:     nil,
}

func getCertsDir() string {
	if garmConfig == nil {
		fmt.Println("garm config not set")
		os.Exit(1)
	}

	return filepath.Join(garmConfig.Default.ConfigDir, "certs")
}

var generateRootCACmd = &cobra.Command{
	Use:   "create-ca",
	Short: "Generate new certificate authority",
	Long:  `Generate a new certificate authority.`,
	// SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		caBaseDir := filepath.Join(getCertsDir(), "ca")
		if _, err := os.Stat(caBaseDir); err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("cannot access %s", caBaseDir)
			}

			if err := os.MkdirAll(caBaseDir, 0700); err != nil {
				return fmt.Errorf("cannot create %s", caBaseDir)
			}
		}

		caCertFile := filepath.Join(caBaseDir, "ca-cert.pem")
		caKeyFile := filepath.Join(caBaseDir, "ca-key.pem")

		if err := certificates.GenerateNewInternalCA(caKeyFile, caCertFile); err != nil {
			return fmt.Errorf("failed to generate a new CA: %s", err)
		}

		return nil
	},
}

var generatePeerCertCmd = &cobra.Command{
	Use:   "create-peer",
	Short: "Generate new certificate for a peer",
	Long: `Generate new certificate for a peer.

This command generates 3 x509 key pairs:

  * client-cert.pem - This can be used by garm to access the embeded etcd as a client.
  * server-cert.pem - This is used by the embeded etcd instance to serve client connections.
  * peer-cert.pem - This is used by the embeded etcd server to secure cluster peer communication.

By default the certificates are saved in /etc/garm/certs/peers/[peer name]/. The value in [DEFAULT]/config_dir
is used as a base for the certs dir.
`,
	// SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		certsDir := getCertsDir()
		caBaseDir := filepath.Join(getCertsDir(), "ca")

		caCert := filepath.Join(caBaseDir, "ca-cert.pem")
		caKey := filepath.Join(caBaseDir, "ca-key.pem")

		if _, err := os.Stat(caCert); err != nil {
			return fmt.Errorf("failed to access CA cert")
		}

		if _, err := os.Stat(caKey); err != nil {
			return fmt.Errorf("failed to access CA key")
		}

		peerCertDir := filepath.Join(certsDir, "peers", peerName)

		peerCertPath := filepath.Join(peerCertDir, "peer-cert.pem")
		peerKeyPath := filepath.Join(peerCertDir, "peer-key.pem")
		serverCertPath := filepath.Join(peerCertDir, "server-cert.pem")
		serverKeyPath := filepath.Join(peerCertDir, "server-key.pem")
		clientCertPath := filepath.Join(peerCertDir, "client-cert.pem")
		clientKeyPath := filepath.Join(peerCertDir, "client-key.pem")

		if _, err := os.Stat(peerCertDir); err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("cannot access %s", peerCertDir)
			}

			if err := os.MkdirAll(peerCertDir, 0700); err != nil {
				return fmt.Errorf("cannot create %s", peerCertDir)
			}
		}

		peerCertTpl, err := certificates.NewServerCertificateTemplate(peerName, certificateSANs, true)
		if err != nil {
			return err
		}
		serverCertTpl, err := certificates.NewServerCertificateTemplate(peerName, certificateSANs, false)
		if err != nil {
			return err
		}
		clientCertTpl, err := certificates.NewCertificateTemplate(peerName)
		if err != nil {
			return err
		}

		if err := certificates.NewCertificatesFromTemplate(peerCertTpl, peerKeyPath, peerCertPath, caCert, caKey); err != nil {
			return err
		}
		if err := certificates.NewCertificatesFromTemplate(serverCertTpl, serverKeyPath, serverCertPath, caCert, caKey); err != nil {
			return err
		}
		if err := certificates.NewCertificatesFromTemplate(clientCertTpl, clientKeyPath, clientCertPath, caCert, caKey); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(certificatesCmd)

	// Create new CA
	certificatesCmd.AddCommand(generateRootCACmd)

	// Create new peer certificate
	certificatesCmd.AddCommand(generatePeerCertCmd)
	generatePeerCertCmd.Flags().StringSliceVar(&certificateSANs, "sans", []string{}, "SANs the peer will be known by (fqdns, IP addresses)")
	generatePeerCertCmd.Flags().StringVar(&peerName, "name", "", "A descriptive name of the peer. This will be included in the certificates.")
	generatePeerCertCmd.MarkFlagRequired("sans")
	generatePeerCertCmd.MarkFlagRequired("name")
}
