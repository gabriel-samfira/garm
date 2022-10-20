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
	"garm/params"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	enterpriseName          string
	enterpriseWebhookSecret string
	enterpriseCreds         string
)

// enterpriseCmd represents the enterprise command
var enterpriseCmd = &cobra.Command{
	Use:          "enterprise",
	Aliases:      []string{"ent"},
	SilenceUsage: true,
	Short:        "Manage enterprise",
	Long: `Add, remove or update enterprise for which we manage
self hosted runners.

This command allows you to define a new enterprise or manage an existing
enterprise for which garm maintains pools of self hosted runners.`,
	Run: nil,
}

var enterpriseAddCmd = &cobra.Command{
	Use:          "add",
	Aliases:      []string{"create"},
	Short:        "Add enterprise",
	Long:         `Add a new enterprise to the manager.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if needsInit {
			return errNeedsInitError
		}

		newEnterpriseReq := params.CreateEnterpriseParams{
			Name:            enterpriseName,
			WebhookSecret:   enterpriseWebhookSecret,
			CredentialsName: enterpriseCreds,
		}
		enterprise, err := cli.CreateEnterprise(newEnterpriseReq)
		if err != nil {
			return err
		}
		formatOneEnterprise(enterprise)
		return nil
	},
}

var enterpriseListCmd = &cobra.Command{
	Use:          "list",
	Aliases:      []string{"ls"},
	Short:        "List enterprises",
	Long:         `List all configured enterprises that are currently managed.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if needsInit {
			return errNeedsInitError
		}

		enterprises, err := cli.ListEnterprises()
		if err != nil {
			return err
		}
		formatEnterprises(enterprises)
		return nil
	},
}

var enterpriseShowCmd = &cobra.Command{
	Use:          "show",
	Short:        "Show details for one enterprise",
	Long:         `Displays detailed information about a single enterprise.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if needsInit {
			return errNeedsInitError
		}
		if len(args) == 0 {
			return fmt.Errorf("requires a enterprise ID")
		}
		if len(args) > 1 {
			return fmt.Errorf("too many arguments")
		}
		enterprise, err := cli.GetEnterprise(args[0])
		if err != nil {
			return err
		}
		formatOneEnterprise(enterprise)
		return nil
	},
}

var enterpriseDeleteCmd = &cobra.Command{
	Use:          "delete",
	Aliases:      []string{"remove", "rm", "del"},
	Short:        "Removes one enterprise",
	Long:         `Delete one enterprise from the manager.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if needsInit {
			return errNeedsInitError
		}
		if len(args) == 0 {
			return fmt.Errorf("requires a enterprise ID")
		}
		if len(args) > 1 {
			return fmt.Errorf("too many arguments")
		}
		if err := cli.DeleteEnterprise(args[0]); err != nil {
			return err
		}
		return nil
	},
}

func init() {

	enterpriseAddCmd.Flags().StringVar(&enterpriseName, "name", "", "The name of the enterprise")
	enterpriseAddCmd.Flags().StringVar(&enterpriseWebhookSecret, "webhook-secret", "", "The webhook secret for this enterprise")
	enterpriseAddCmd.Flags().StringVar(&enterpriseCreds, "credentials", "", "Credentials name. See credentials list.")
	enterpriseAddCmd.MarkFlagRequired("credentials")
	enterpriseAddCmd.MarkFlagRequired("name")

	enterpriseCmd.AddCommand(
		enterpriseListCmd,
		enterpriseAddCmd,
		enterpriseShowCmd,
		enterpriseDeleteCmd,
	)

	rootCmd.AddCommand(enterpriseCmd)
}

func formatEnterprises(enterprises []params.Enterprise) {
	t := table.NewWriter()
	header := table.Row{"ID", "Name", "Credentials name", "Pool mgr running"}
	t.AppendHeader(header)
	for _, val := range enterprises {
		t.AppendRow(table.Row{val.ID, val.Name, val.CredentialsName, val.PoolManagerStatus.IsRunning})
		t.AppendSeparator()
	}
	fmt.Println(t.Render())
}

func formatOneEnterprise(enterprise params.Enterprise) {
	t := table.NewWriter()
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	header := table.Row{"Field", "Value"}
	t.AppendHeader(header)
	t.AppendRow(table.Row{"ID", enterprise.ID})
	t.AppendRow(table.Row{"Name", enterprise.Name})
	t.AppendRow(table.Row{"Credentials", enterprise.CredentialsName})
	t.AppendRow(table.Row{"Pool manager running", enterprise.PoolManagerStatus.IsRunning})
	if !enterprise.PoolManagerStatus.IsRunning {
		t.AppendRow(table.Row{"Failure reason", enterprise.PoolManagerStatus.FailureReason})
	}

	if len(enterprise.Pools) > 0 {
		for _, pool := range enterprise.Pools {
			t.AppendRow(table.Row{"Pools", pool.ID}, rowConfigAutoMerge)
		}
	}
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
		{Number: 2, AutoMerge: true},
	})

	fmt.Println(t.Render())
}
