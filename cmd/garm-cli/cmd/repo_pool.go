/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"garm/config"
	"garm/params"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	poolProvider       string
	poolMaxRunners     uint
	poolMinIdleRunners uint
	poolImage          string
	poolFlavor         string
	poolOSType         string
	poolOSArch         string
	poolTags           string
	poolEnabled        bool
)

// repoPoolCmd represents the pool command
var repoPoolCmd = &cobra.Command{
	Use:          "pool",
	SilenceUsage: true,
	Aliases:      []string{"pools"},
	Short:        "Manage pools",
	Long: `Manage pools for a repository.

Repositories and organizations can define multiple pools with different
characteristics, which in turn will spawn github self hosted runners on
compute instances that reflect those characteristics.

For example, one pool can define a runner with tags "GPU,ML" which will
spin up instances with access to a GPU, on the desired provider.`,
	Run: nil,
}

var poolAddCmd = &cobra.Command{
	Use:          "add",
	Aliases:      []string{"create"},
	Short:        "Add pool",
	Long:         `Add a new pool repository to the manager.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if needsInit {
			return needsInitError
		}

		if len(args) == 0 {
			return fmt.Errorf("requires a repository ID")
		}

		if len(args) > 1 {
			return fmt.Errorf("too many arguments")
		}

		tags := strings.Split(poolTags, ",")
		newPoolParams := params.CreatePoolParams{
			ProviderName:   poolProvider,
			MaxRunners:     poolMaxRunners,
			MinIdleRunners: poolMinIdleRunners,
			Image:          poolImage,
			Flavor:         poolFlavor,
			OSType:         config.OSType(poolOSType),
			OSArch:         config.OSArch(poolOSArch),
			Tags:           tags,
			Enabled:        poolEnabled,
		}
		if err := newPoolParams.Validate(); err != nil {
			return err
		}
		pool, err := cli.CreateRepoPool(args[0], newPoolParams)
		if err != nil {
			return err
		}
		formatOnePool(pool)
		return nil
	},
}

var poolListCmd = &cobra.Command{
	Use:          "list",
	Aliases:      []string{"ls"},
	Short:        "List repository pools",
	Long:         `List all configured pools for a given repository.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if needsInit {
			return needsInitError
		}

		if len(args) == 0 {
			return fmt.Errorf("requires a repository ID")
		}

		if len(args) > 1 {
			return fmt.Errorf("too many arguments")
		}

		pools, err := cli.ListRepoPools(args[0])
		if err != nil {
			return err
		}
		formatPools(pools)
		return nil
	},
}

var poolShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show details for one pool",
	Long:  `Displays detailed information about a single pool.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if needsInit {
			return needsInitError
		}

		if len(args) < 2 || len(args) > 2 {
			return fmt.Errorf("command requires repoID and poolID")
		}

		pool, err := cli.GetRepoPool(args[0], args[1])
		if err != nil {
			return err
		}

		formatOnePool(pool)
		return nil
	},
}

var poolDeleteCmd = &cobra.Command{
	Use:          "delete",
	Aliases:      []string{"remove", "rm", "del"},
	Short:        "Removes one pool",
	Long:         `Delete one repository pool from the manager.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if needsInit {
			return needsInitError
		}
		if len(args) < 2 || len(args) > 2 {
			return fmt.Errorf("command requires repoID and poolID")
		}

		if err := cli.DeleteRepoPool(args[0], args[1]); err != nil {
			return err
		}
		return nil
	},
}

var poolUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update one pool",
	Long: `Updates pool characteristics.

This command updates the pool characteristics. Runners already created prior to updating
the pool, will not be recreated. IF they no longer suit your needs, you will need to
explicitly remove them using the runner delete command.
	`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if needsInit {
			return needsInitError
		}

		if len(args) < 2 || len(args) > 2 {
			return fmt.Errorf("command requires repoID and poolID")
		}

		poolUpdateParams := params.UpdatePoolParams{}

		if cmd.Flags().Changed("image") {
			poolUpdateParams.Image = poolImage
		}

		if cmd.Flags().Changed("flavor") {
			poolUpdateParams.Flavor = poolFlavor
		}

		if cmd.Flags().Changed("tags") {
			poolUpdateParams.Tags = strings.Split(poolTags, ",")
		}

		if cmd.Flags().Changed("os-type") {
			poolUpdateParams.OSType = config.OSType(poolOSType)
		}

		if cmd.Flags().Changed("os-arch") {
			poolUpdateParams.OSArch = config.OSArch(poolOSArch)
		}

		if cmd.Flags().Changed("max-runners") {
			poolUpdateParams.MaxRunners = &poolMaxRunners
		}

		if cmd.Flags().Changed("min-idle-runners") {
			poolUpdateParams.MinIdleRunners = &poolMinIdleRunners
		}

		if cmd.Flags().Changed("enabled") {
			poolUpdateParams.Enabled = &poolEnabled
		}

		pool, err := cli.UpdateRepoPool(args[0], args[1], poolUpdateParams)
		if err != nil {
			return err
		}

		formatOnePool(pool)
		return nil
	},
}

func init() {
	poolAddCmd.Flags().StringVar(&poolProvider, "provider-name", "", "The name of the provider where runners will be created.")
	poolAddCmd.Flags().StringVar(&poolImage, "image", "", "The provider-specific image name to use for runners in this pool.")
	poolAddCmd.Flags().StringVar(&poolFlavor, "flavor", "", "The flavor to use for this runner.")
	poolAddCmd.Flags().StringVar(&poolTags, "tags", "", "A comma separated list of tags to assign to this runner.")
	poolAddCmd.Flags().StringVar(&poolOSType, "os-type", "linux", "Operating system type (windows, linux, etc).")
	poolAddCmd.Flags().StringVar(&poolOSArch, "os-arch", "amd64", "Operating system architecture (amd64, arm, etc).")
	poolAddCmd.Flags().UintVar(&poolMaxRunners, "max-runners", 5, "The maximum number of runner this pool will create.")
	poolAddCmd.Flags().UintVar(&poolMinIdleRunners, "min-idle-runners", 1, "Attempt to maintain a minimum of idle self-hosted runners of this type.")
	poolAddCmd.Flags().BoolVar(&poolEnabled, "enabled", false, "Enable this pool.")
	poolAddCmd.MarkFlagRequired("provider-name")
	poolAddCmd.MarkFlagRequired("image")
	poolAddCmd.MarkFlagRequired("flavor")
	poolAddCmd.MarkFlagRequired("tags")

	poolUpdateCmd.Flags().StringVar(&poolImage, "image", "", "The provider-specific image name to use for runners in this pool.")
	poolUpdateCmd.Flags().StringVar(&poolFlavor, "flavor", "", "The flavor to use for this runner.")
	poolUpdateCmd.Flags().StringVar(&poolTags, "tags", "", "A comma separated list of tags to assign to this runner.")
	poolUpdateCmd.Flags().StringVar(&poolOSType, "os-type", "linux", "Operating system type (windows, linux, etc).")
	poolUpdateCmd.Flags().StringVar(&poolOSArch, "os-arch", "amd64", "Operating system architecture (amd64, arm, etc).")
	poolUpdateCmd.Flags().UintVar(&poolMaxRunners, "max-runners", 5, "The maximum number of runner this pool will create.")
	poolUpdateCmd.Flags().UintVar(&poolMinIdleRunners, "min-idle-runners", 1, "Attempt to maintain a minimum of idle self-hosted runners of this type.")
	poolUpdateCmd.Flags().BoolVar(&poolEnabled, "enabled", false, "Enable this pool.")

	repoPoolCmd.AddCommand(
		poolListCmd,
		poolAddCmd,
		poolShowCmd,
		poolDeleteCmd,
		poolUpdateCmd,
	)

	repositoryCmd.AddCommand(repoPoolCmd)
}

func formatPools(pools []params.Pool) {
	t := table.NewWriter()
	header := table.Row{"ID", "Image", "Flavor", "Tags", "Enabled"}
	t.AppendHeader(header)

	for _, pool := range pools {
		tags := []string{}
		for _, tag := range pool.Tags {
			tags = append(tags, tag.Name)
		}
		t.AppendRow(table.Row{pool.ID, pool.Image, pool.Flavor, strings.Join(tags, " "), pool.Enabled})
		t.AppendSeparator()
	}
	fmt.Println(t.Render())
}

func formatOnePool(pool params.Pool) {
	t := table.NewWriter()
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}

	header := table.Row{"Field", "Value"}

	tags := []string{}
	for _, tag := range pool.Tags {
		tags = append(tags, tag.Name)
	}

	t.AppendHeader(header)
	t.AppendRow(table.Row{"ID", pool.ID})
	t.AppendRow(table.Row{"Provider Name", pool.ProviderName})
	t.AppendRow(table.Row{"Image", pool.Image})
	t.AppendRow(table.Row{"Flavor", pool.Flavor})
	t.AppendRow(table.Row{"OS Type", pool.OSType})
	t.AppendRow(table.Row{"OS Architecture", pool.OSArch})
	t.AppendRow(table.Row{"Max Runners", pool.MaxRunners})
	t.AppendRow(table.Row{"Min Idle Runners", pool.MinIdleRunners})
	t.AppendRow(table.Row{"Tags", strings.Join(tags, ", ")})
	t.AppendRow(table.Row{"Enabled", pool.Enabled})

	if len(pool.Instances) > 0 {
		for _, instance := range pool.Instances {
			t.AppendRow(table.Row{"Instances", fmt.Sprintf("%s (%s)", instance.Name, instance.ID)}, rowConfigAutoMerge)
		}
	}

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
		{Number: 2, AutoMerge: true},
	})
	fmt.Println(t.Render())
}