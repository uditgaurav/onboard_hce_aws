package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uditgaurav/onboard_hce_aws/pkg/register"
)

var params register.InfraParameters

var rootCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new Harness Chaos infrastructure",
	Long:  `A CLI utility to register a new Harness Chaos infrastructure using the given name and namespace.`,
	Run: func(cmd *cobra.Command, args []string) {
		register.RegisterInfra(params)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&params.ApiKey, "api-key", "", "API Key for Harness (required)")
	rootCmd.PersistentFlags().StringVar(&params.AccountId, "account-id", "", "Account ID for Harness (required)")
	rootCmd.PersistentFlags().StringVar(&params.Infra.Name, "name", "", "Name of the Harness Chaos infrastructure (required)")
	rootCmd.PersistentFlags().StringVar(&params.Infra.Namespace, "namespace", "", "Namespace for the Harness Chaos infrastructure (required)")
	rootCmd.PersistentFlags().StringVar(&params.Organisation, "organisation", "", "Organisation Identifier (required)")
	rootCmd.PersistentFlags().StringVar(&params.Project, "project", "", "Project Identifier (required)")
	rootCmd.PersistentFlags().StringVar(&params.InfraScope, "infra-scope", "", "Infrastructure Scope (required)")
	rootCmd.PersistentFlags().BoolVar(&params.InfraNsExists, "infra-ns-exists", false, "Does infrastructure namespace exist (required)")

	rootCmd.MarkFlagRequired("api-key")
	rootCmd.MarkFlagRequired("account-id")
	rootCmd.MarkFlagRequired("name")
	rootCmd.MarkFlagRequired("namespace")
	rootCmd.MarkFlagRequired("organisation")
	rootCmd.MarkFlagRequired("project")
	rootCmd.MarkFlagRequired("infra-scope")
	rootCmd.MarkFlagRequired("infra-ns-exists")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
