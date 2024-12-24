package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vamshisiddarth/kubeversion/internal/download"
	"github.com/vamshisiddarth/kubeversion/internal/version"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "kubeversion",
		Short: "Kubectl version manager",
		Long:  `Kubeversion is a version manager for kubectl that allows you to switch between different versions easily.`,
	}

	var useCmd = &cobra.Command{
		Use:   "use [version]",
		Short: "Switch to specified kubectl version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return version.SwitchVersion(args[0])
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all available kubectl versions",
		RunE: func(cmd *cobra.Command, args []string) error {
			return version.ListVersions()
		},
	}

	var installCmd = &cobra.Command{
		Use:   "install [version]",
		Short: "Install specified kubectl version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return download.InstallVersion(args[0])
		},
	}

	rootCmd.AddCommand(useCmd, listCmd, installCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
