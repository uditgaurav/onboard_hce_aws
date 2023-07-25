package main

import (
	"os"

	"github.com/litmuschaos/litmus-go/pkg/log"
	"github.com/spf13/cobra"
	"github.com/uditgaurav/onboard_hce_aws/execute"
	"github.com/uditgaurav/onboard_hce_aws/pkg/types"
)

var params types.OnboardingParameters

var rootCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new Harness Chaos infrastructure with AWS",
	Long:  `A CLI utility to register a new Harness Chaos infrastructure with AWS account.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := os.Setenv("KUBECONFIG", params.KubeConfigPath); err != nil {
			log.Fatalf("Failed to set KUBECONFIG environment variable, err: %v", err)
		}

		if err := execute.Execute(params); err != nil {
			log.Fatalf("fail to register chaos infra with aws, err: %v", err)
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
	rootCmd.Flags().IntVar(&params.Timeout, "timeout", 180, "Timeout For Infra setup")
	rootCmd.Flags().IntVar(&params.Delay, "delay", 2, "Delay between checking the status of Infra")

	// Flags for aws setup
	rootCmd.Flags().StringVar(&params.ProviderUrl, "provider-url", "", "Provider URL")
	rootCmd.Flags().StringVar(&params.RoleName, "role-name", "", "Role Name")
	rootCmd.Flags().StringVar(&params.Resources, "resources", "", "Resources")
	rootCmd.Flags().StringVar(&params.Region, "region", "", "Target AWS Region")
	rootCmd.Flags().StringVar(&params.ExperimentServiceAccountName, "service-account", "litmus-admin", "Experiment Service Account Name")
	rootCmd.Flags().StringVar(&params.KubeConfigPath, "kubeconfig-path", "", "Path to the kubeconfig file")
	rootCmd.Flags().StringVar(&params.Actions, "actions", "all", "Actions that are performed by this cli. (Default all)")

}

func main() {
	// Now, mark the necessary flags as required
	if err := rootCmd.MarkFlagRequired("api-key"); err != nil {
		log.Fatalf("Error marking 'api-key' as required: %v", err)
	}
	if err := rootCmd.MarkFlagRequired("account-id"); err != nil {
		log.Fatalf("Error marking 'account-id' as required: %v", err)
	}
	if err := rootCmd.MarkFlagRequired("infra-name"); err != nil {
		log.Fatalf("Error marking 'infra-name' as required: %v", err)
	}
	if err := rootCmd.MarkFlagRequired("project"); err != nil {
		log.Fatalf("Error marking 'project' as required: %v", err)
	}
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
