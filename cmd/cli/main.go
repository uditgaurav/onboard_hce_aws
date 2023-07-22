package main

import (
	"log"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/uditgaurav/onboard_hce_aws/pkg/register"
)

var params register.InfraParameters

var rootCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new Harness Chaos infrastructure",
	Long:  `A CLI utility to register a new Harness Chaos infrastructure using the given name and namespace.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := register.RegisterInfra(params); err != nil {
			log.Fatalf("fail to register chaos infra, err: %v", err)
		}
	},
}

func init() {
	rootCmd.Flags().StringVar(&params.ApiKey, "api-key", "", "API Key for Harness (required)")
	rootCmd.Flags().StringVar(&params.AccountId, "account-id", "", "Account ID for Harness (required)")
	rootCmd.Flags().StringVar(&params.Infra.Name, "infra-name", "", "Name of the Harness Chaos infrastructure (required)")
	rootCmd.Flags().StringVar(&params.Project, "project", "", "Project Identifier (required)")

	// Default value for infra-environment-id and infra-platform-name is calculated in RegisterInfra based on infra-name

	rootCmd.Flags().StringVar(&params.Infra.Namespace, "infra-namespace", "hce", "Namespace for the Harness Chaos infrastructure")
	rootCmd.Flags().StringVar(&params.Organisation, "organisation", "default", "Organisation Identifier")
	rootCmd.Flags().StringVar(&params.InfraScope, "infra-scope", "namespace", "Infrastructure Scope")
	rootCmd.Flags().BoolVar(&params.InfraNsExists, "infra-ns-exists", true, "Does infrastructure namespace exist")
	rootCmd.Flags().StringVar(&params.Infra.Description, "infra-description", "Infra for Harness Chaos Testing", "Infra Description")
	rootCmd.Flags().StringVar(&params.Infra.ServiceAccount, "infra-service-account", "hce", "Infra Service Account")
	rootCmd.Flags().BoolVar(&params.Infra.InfraSaExists, "is-infra-sa-exists", false, "Does infrastructure service account exist")
	rootCmd.Flags().StringVar(&params.Infra.EnvironmentID, "infra-environment-id", "", "Infra Environment ID")
	rootCmd.Flags().StringVar(&params.Infra.PlatformName, "infra-platform-name", "", "Infra Platform Name")
	rootCmd.Flags().BoolVar(&params.Infra.SkipSsl, "infra-skip-ssl", false, "Skip SSL for Infra")

}

func main() {
	// Now, mark the necessary flags as required
	if err := rootCmd.MarkFlagRequired("api-key"); err != nil {
		logrus.Fatalf("Error marking 'api-key' as required: %v", err)
	}
	if err := rootCmd.MarkFlagRequired("account-id"); err != nil {
		logrus.Fatalf("Error marking 'account-id' as required: %v", err)
	}
	if err := rootCmd.MarkFlagRequired("infra-name"); err != nil {
		logrus.Fatalf("Error marking 'infra-name' as required: %v", err)
	}
	if err := rootCmd.MarkFlagRequired("project"); err != nil {
		logrus.Fatalf("Error marking 'project' as required: %v", err)
	}
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
	}
}
