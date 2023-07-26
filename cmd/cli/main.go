package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/litmuschaos/litmus-go/pkg/log"
	"github.com/spf13/cobra"
	"github.com/uditgaurav/onboard_hce_aws/execute"
	"github.com/uditgaurav/onboard_hce_aws/pkg/types"
)

var params types.OnboardingParameters
var configFile string

var rootCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new Harness Chaos infrastructure with AWS",
	Long:  `A CLI utility to register a new Harness Chaos infrastructure with AWS account.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if config file flag is provided
		if configFile != "" {
			// Read from the config file
			configBytes, err := ioutil.ReadFile(configFile)
			if err != nil {
				log.Fatalf("Unable to read config file: %v", err)
			}

			// Unmarshal JSON into params slice
			var paramsSlice []types.OnboardingParameters
			err = json.Unmarshal(configBytes, &paramsSlice)
			if err != nil {
				log.Fatalf("Unable to parse config JSON: %v", err)
			}

			// Iterate over each item in the paramsSlice
			for _, params := range paramsSlice {
				// Call the function using the params from the JSON
				registerInfra(params)
			}
		} else {
			// Call the function using the params from the flags
			registerInfra(params)
		}
	},
}

func registerInfra(params types.OnboardingParameters) {
	if err := os.Setenv("KUBECONFIG", params.KubeConfigPath); err != nil {
		log.Fatalf("Failed to set KUBECONFIG environment variable, err: %v", err)
	}

	if err := execute.Execute(params); err != nil {
		log.Fatalf("fail to register chaos infra with aws, err: %v", err)
	}
}

func init() {
	rootCmd.Flags().StringVar(&params.ApiKey, "api-key", "", "API Key for Harness")
	rootCmd.Flags().StringVar(&params.AccountId, "account-id", "", "Account ID for Harness")
	rootCmd.Flags().StringVar(&params.Infra.Name, "infra-name", "", "Name of the Harness Chaos infrastructure")
	rootCmd.Flags().StringVar(&params.Project, "project", "", "Project Identifier")

	// Default value for infra-environment-id and infra-platform-name is calculated in RegisterInfra based on infra-name

	rootCmd.Flags().StringVar(&params.Infra.Namespace, "infra-namespace", "hce", "Namespace for the Harness Chaos infrastructure")
	rootCmd.Flags().StringVar(&params.Organisation, "organisation", "default", "Organisation Identifier")
	rootCmd.Flags().StringVar(&params.Infra.InfraScope, "infra-scope", "namespace", "Infrastructure Scope")
	rootCmd.Flags().BoolVar(&params.Infra.InfraNsExists, "infra-ns-exists", true, "Does infrastructure namespace exist")
	rootCmd.Flags().StringVar(&params.Infra.InfraDescription, "infra-description", "Infra for Harness Chaos Testing", "Infra Description")
	rootCmd.Flags().StringVar(&params.Environment.EnvironmentDescription, "env-description", "Environment for Harness Chaos Testing", "Environment Description")
	rootCmd.Flags().StringVar(&params.Environment.EnvironmentType, "env-type", "PreProduction", "Specify the environment type whether Production or PreProduction")

	rootCmd.Flags().StringVar(&params.Infra.ServiceAccount, "infra-service-account", "hce", "Infra Service Account")
	rootCmd.Flags().BoolVar(&params.Infra.InfraSaExists, "is-infra-sa-exists", false, "Does infrastructure service account exist")
	rootCmd.Flags().StringVar(&params.Environment.EnvironmentName, "environment-name", "", "Environment Name")
	rootCmd.Flags().StringVar(&params.Infra.PlatformName, "infra-platform-name", "", "Infra Platform Name")
	rootCmd.Flags().BoolVar(&params.Infra.SkipSsl, "infra-skip-ssl", false, "Skip SSL for Infra")
	rootCmd.Flags().IntVar(&params.Timeout, "timeout", 180, "Timeout For Infra setup")
	rootCmd.Flags().IntVar(&params.Delay, "delay", 2, "Delay between checking the status of Infra")

	// Flags for aws setup
	rootCmd.Flags().StringVar(&params.ProviderUrl, "provider-url", "", "Provider URL")
	rootCmd.Flags().StringVar(&params.RoleName, "role-name", "", "Role Name")
	rootCmd.Flags().StringVar(&params.Resources, "resources", "all", "Resources")
	rootCmd.Flags().StringVar(&params.Region, "region", "", "Target AWS Region")
	rootCmd.Flags().StringVar(&params.ExperimentServiceAccountName, "service-account", "litmus-admin", "Experiment Service Account Name")
	rootCmd.Flags().StringVar(&params.KubeConfigPath, "kubeconfig-path", "", "Path to the kubeconfig file")
	rootCmd.Flags().StringVar(&params.Actions, "actions", "all", "Actions that are performed by this cli. (Default all)")
	rootCmd.Flags().StringVar(&params.AWSCredentialFile, "aws-credential-file", "", "Path To The AWS Credential File (default $HOME/.aws/credentials)")
	rootCmd.Flags().StringVar(&params.AWSProfile, "aws-profile", "default", "Provide the AWS profile (Default 'default')")
	rootCmd.Flags().StringVar(&configFile, "config", "", "Config file containing parameters")

	if params.AWSCredentialFile == "" {
		params.AWSCredentialFile = fmt.Sprintf("%s/.aws/credentials", os.Getenv("HOME"))
	}

	if err := os.Setenv("AWS_SHARED_CREDENTIALS_FILE", params.AWSCredentialFile); err != nil {
		log.Fatalf("Failed to set AWS_SHARED_CREDENTIALS_FILE environment variable, err: %v", err)
	}

	// Set the AWS_PROFILE environment variable
	if err := os.Setenv("AWS_PROFILE", params.AWSProfile); err != nil {
		log.Fatalf("Failed to set AWS_PROFILE environment variable, err: %v", err)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
