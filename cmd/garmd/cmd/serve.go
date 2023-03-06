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
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"srv"},
	Short:   "Run the GARM main loop",
	Long:    `Start all workers, watchers and api server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// ctx, stop := signal.NotifyContext(context.Background(), signals...)
		// defer stop()

		// etcdCfg, err := config.EmbededEtcd.ServerConfig(garmConfig.Default.ConfigDir)
		// if err != nil {
		// 	return err
		// }

		// e, err := embed.StartEtcd(etcdCfg)
		// if err != nil {
		// 	return err
		// }
		// // defer e.Close()

		// select {
		// case <-e.Server.ReadyNotify():
		// 	log.Printf("Server is ready!")
		// case <-time.After(60 * time.Second):
		// 	e.Server.Stop() // trigger a shutdown
		// 	srvErr := <-e.Err()
		// 	return fmt.Errorf("timed out waiting for server to be ready: %s", srvErr)
		// }

		// <-ctx.Done()
		// e.Server.Stop()
		// <-e.Server.StopNotify()
		// e.Close()
		// etcdErr := <-e.Err()
		// if etcdErr != nil {
		// 	return etcdErr
		// }
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// // Create new CA
	// certificatesCmd.AddCommand(generateRootCACmd)

	// // Create new peer certificate
	// certificatesCmd.AddCommand(generatePeerCertCmd)
	// generatePeerCertCmd.Flags().StringSliceVar(&certificateSANs, "sans", []string{}, "SANs the peer will be known by (fqdns, IP addresses)")
	// generatePeerCertCmd.Flags().StringVar(&peerName, "name", "", "A descriptive name of the peer. This will be included in the certificates.")
	// generatePeerCertCmd.MarkFlagRequired("sans")
	// generatePeerCertCmd.MarkFlagRequired("name")
}
