package execute

import (
	"github.com/pkg/errors"
	"github.com/uditgaurav/onboard_hce_aws/pkg/aws"
	"github.com/uditgaurav/onboard_hce_aws/pkg/clients"
	"github.com/uditgaurav/onboard_hce_aws/pkg/kubernetes"
	"github.com/uditgaurav/onboard_hce_aws/pkg/register"
	"github.com/uditgaurav/onboard_hce_aws/pkg/types"
)

// Execute is the main function which is responsible for the whole onboarding process. It takes OnboardingParameters as an argument.
func Execute(params types.OnboardingParameters) error {
	// RegisterInfra is a function to register the ChaosInfra.
	if err := register.RegisterInfra(params); err != nil {
		return errors.Errorf("failed to register ChaosInfra, err: %v", err)
	}

	// ConnectOIDCProvider is a function to connect to the OIDC provider.
	if err := aws.ConnectOIDCProvider(params); err != nil {
		return errors.Errorf("failed to connect OIDC provider, err: %v", err)
	}

	// PreparePolicyAndCreateRole and CreateRoleWithTrustRelationsip are functions to create IAM roles in AWS.
	if params.RoleName == "" {
		if err := aws.PreparePolicyAndCreateRole(params); err != nil {
			return errors.Errorf("failed to create policy and role, err: %v", err)
		}
	} else {
		if err := aws.CreateRoleWithTrustRelationsip("", params); err != nil {
			return errors.Errorf("failed to create role, err: %v", err)
		}
	}

	// AnnotateServiceAccount is a function to annotate a Kubernetes ServiceAccount with the ARN of the created IAM role.
	if err := kubernetes.AnnotateServiceAccount(params, clients.ClientSets{}); err != nil {
		return errors.Errorf("failed to annotate experiment service account with role arn, err: %v", err)
	}

	return nil
}
