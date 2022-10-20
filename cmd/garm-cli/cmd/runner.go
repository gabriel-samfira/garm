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
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	runnerRepository   string
	runnerOrganization string
	runnerEnterprise   string
	runnerAll          bool
	forceRemove        bool
)

// runnerCmd represents the runner command
var runnerCmd = &cobra.Command{
	Use:          "runner",
	Aliases:      []string{"run"},
	SilenceUsage: true,
	Short:        "List runners in a pool",
	Long: `Given a pool ID, of either a repository or an organization,
list all instances.`,
	Run: nil,
}

var runnerListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List runners",
	Long: `List runners of pools, repositories, orgs or all of the above.
	
This command expects to get either a pool ID as a positional parameter, or it expects
that one of the supported switches be used to fetch runners of --repo, --org or --all

Example:

	List runners from one pool:
	garm-cli runner list e87e70bd-3d0d-4b25-be9a-86b85e114bcb

	List runners from one repo:
	garm-cli runner list --repo=05e7eac6-4705-486d-89c9-0170bbb576af

	List runners from one org:
	garm-cli runner list --org=5493e51f-3170-4ce3-9f05-3fe690fc6ec6

	List runners from one enterprise:
	garm-cli runner list --enterprise=a966188b-0e05-4edc-9b82-bc81a1fd38ed

	List all runners from all pools belonging to all repos and orgs:
	garm-cli runner list --all

`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if needsInit {
			return errNeedsInitError
		}

		var instances []params.Instance
		var err error

		switch len(args) {
		case 1:
			if cmd.Flags().Changed("repo") ||
				cmd.Flags().Changed("org") ||
				cmd.Flags().Changed("enterprise") ||
				cmd.Flags().Changed("all") {

				return fmt.Errorf("specifying a pool ID and any of [all org repo enterprise] are mutually exclusive")
			}
			instances, err = cli.ListPoolInstances(args[0])
		case 0:
			if cmd.Flags().Changed("repo") {
				instances, err = cli.ListRepoInstances(runnerRepository)
			} else if cmd.Flags().Changed("org") {
				instances, err = cli.ListOrgInstances(runnerOrganization)
			} else if cmd.Flags().Changed("enterprise") {
				instances, err = cli.ListEnterpriseInstances(runnerEnterprise)
			} else if cmd.Flags().Changed("all") {
				instances, err = cli.ListAllInstances()
			} else {
				cmd.Help()
				os.Exit(0)
			}
		default:
			cmd.Help()
			os.Exit(0)
		}

		if err != nil {
			return err
		}
		formatInstances(instances)
		return nil
	},
}

var runnerShowCmd = &cobra.Command{
	Use:          "show",
	Short:        "Show details for a runner",
	Long:         `Displays a detailed view of a single runner.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if needsInit {
			return errNeedsInitError
		}

		if len(args) == 0 {
			return fmt.Errorf("requires a runner name")
		}

		if len(args) > 1 {
			return fmt.Errorf("too many arguments")
		}

		instance, err := cli.GetInstanceByName(args[0])
		if err != nil {
			return err
		}
		formatSingleInstance(instance)
		return nil
	},
}

var runnerDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Remove a runner",
	Aliases: []string{"remove", "rm", "del"},
	Long: `Remove a runner.

This command deletes an existing runner. If it registered in Github
and we recorded an agent ID for it, we will attempt to remove it from
Github first, then mark the runner as pending_delete so it will be
cleaned up by the provider.

NOTE: An active runner cannot be removed from Github. You will have
to either cancel the workflow or wait for it to finish.
`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if needsInit {
			return errNeedsInitError
		}

		if len(args) == 0 {
			return fmt.Errorf("requires a runner name")
		}

		if !forceRemove {
			return fmt.Errorf("use --force-remove-runner=true to remove a runner")
		}

		if err := cli.DeleteRunner(args[0]); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	runnerListCmd.Flags().StringVarP(&runnerRepository, "repo", "r", "", "List all runners from all pools within this repository.")
	runnerListCmd.Flags().StringVarP(&runnerOrganization, "org", "o", "", "List all runners from all pools withing this organization.")
	runnerListCmd.Flags().StringVarP(&runnerEnterprise, "enterprise", "e", "", "List all runners from all pools withing this enterprise.")
	runnerListCmd.Flags().BoolVarP(&runnerAll, "all", "a", false, "List all runners, regardless of org or repo.")
	runnerListCmd.MarkFlagsMutuallyExclusive("repo", "org", "enterprise", "all")

	runnerDeleteCmd.Flags().BoolVarP(&forceRemove, "force-remove-runner", "f", false, "Confirm you want to delete a runner")
	runnerDeleteCmd.MarkFlagsMutuallyExclusive("force-remove-runner")

	runnerCmd.AddCommand(
		runnerListCmd,
		runnerShowCmd,
		runnerDeleteCmd,
	)

	rootCmd.AddCommand(runnerCmd)
}

func formatInstances(param []params.Instance) {
	t := table.NewWriter()
	header := table.Row{"Name", "Status", "Runner Status", "Pool ID"}
	t.AppendHeader(header)

	for _, inst := range param {
		t.AppendRow(table.Row{inst.Name, inst.Status, inst.RunnerStatus, inst.PoolID})
		t.AppendSeparator()
	}
	fmt.Println(t.Render())
}

func formatSingleInstance(instance params.Instance) {
	t := table.NewWriter()

	header := table.Row{"Field", "Value"}

	t.AppendHeader(header)
	t.AppendRow(table.Row{"ID", instance.ID}, table.RowConfig{AutoMerge: false})
	t.AppendRow(table.Row{"Provider ID", instance.ProviderID}, table.RowConfig{AutoMerge: false})
	t.AppendRow(table.Row{"Name", instance.Name}, table.RowConfig{AutoMerge: false})
	t.AppendRow(table.Row{"OS Type", instance.OSType}, table.RowConfig{AutoMerge: false})
	t.AppendRow(table.Row{"OS Architecture", instance.OSArch}, table.RowConfig{AutoMerge: false})
	t.AppendRow(table.Row{"OS Name", instance.OSName}, table.RowConfig{AutoMerge: false})
	t.AppendRow(table.Row{"OS Version", instance.OSVersion}, table.RowConfig{AutoMerge: false})
	t.AppendRow(table.Row{"Status", instance.Status}, table.RowConfig{AutoMerge: false})
	t.AppendRow(table.Row{"Runner Status", instance.RunnerStatus}, table.RowConfig{AutoMerge: false})
	t.AppendRow(table.Row{"Pool ID", instance.PoolID}, table.RowConfig{AutoMerge: false})

	if len(instance.Addresses) > 0 {
		for _, addr := range instance.Addresses {
			t.AppendRow(table.Row{"Addresses", addr.Address}, table.RowConfig{AutoMerge: true})
		}
	}

	if len(instance.ProviderFault) > 0 {
		t.AppendRow(table.Row{"Provider Fault", string(instance.ProviderFault)}, table.RowConfig{AutoMerge: true})
	}

	if len(instance.StatusMessages) > 0 {
		for _, msg := range instance.StatusMessages {
			t.AppendRow(table.Row{"Status Updates", fmt.Sprintf("%s: %s", msg.CreatedAt.Format("2006-01-02T15:04:05"), msg.Message)}, table.RowConfig{AutoMerge: true})
		}
	}

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
		{Number: 2, AutoMerge: false, WidthMax: 100},
	})
	fmt.Println(t.Render())
}
