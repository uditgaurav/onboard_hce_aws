package execute

import (
	"github.com/pkg/errors"
	"github.com/uditgaurav/onboard_hce_aws/pkg/aws"
	"github.com/uditgaurav/onboard_hce_aws/pkg/clients"
	"github.com/uditgaurav/onboard_hce_aws/pkg/kubernetes"
	"github.com/uditgaurav/onboard_hce_aws/pkg/register"
	"github.com/uditgaurav/onboard_hce_aws/pkg/types"
)

func Execute(params types.OnboardingParameters) error {
	// Create a new ClientSets
	clients := &clients.ClientSets{}

	// Initialize KubeClient
	if err := clients.GenerateClientSetFromKubeConfig(); err != nil {
		return errors.Errorf("Failed to initialize KubeClient: %v", err)
	}

	switch params.Actions {

	case "all":
		if err := register.RegisterInfra(params); err != nil {
			return errors.Errorf("failed to register ChaosInfra, err: %v", err)
		}
		providerARN, err := aws.ConnectOIDCProvider(params)
		if err != nil {
			return errors.Errorf("failed to connect OIDC provider, err: %v", err)
		}
		params.ProviderARN = providerARN
		if params.RoleName == "" {
			if err := aws.PreparePolicyAndCreateRole(params); err != nil {
				return errors.Errorf("failed to create policy and role, err: %v", err)
			}
		} else {
			if err := aws.CreateRoleWithTrustRelationsip("", params); err != nil {
				return errors.Errorf("failed to create role, err: %v", err)
			}
		}
		if err := kubernetes.AnnotateServiceAccount(params, *clients); err != nil {
			return errors.Errorf("failed to annotate experiment service account with role arn, err: %v", err)
		}

	case "only_install":

		if err := register.RegisterInfra(params); err != nil {
			return errors.Errorf("failed to register ChaosInfra, err: %v", err)
		}

	case "install_with_provider":

		if err := register.RegisterInfra(params); err != nil {
			return errors.Errorf("failed to register ChaosInfra, err: %v", err)
		}
		providerARN, err := aws.ConnectOIDCProvider(params)
		if err != nil {
			return errors.Errorf("failed to connect OIDC provider, err: %v", err)
		}
		params.ProviderARN = providerARN
		if params.RoleName == "" {
			if err := aws.PreparePolicyAndCreateRole(params); err != nil {
				return errors.Errorf("failed to create policy and role, err: %v", err)
			}
		} else {
			if err := aws.CreateRoleWithTrustRelationsip("", params); err != nil {
				return errors.Errorf("failed to create role, err: %v", err)
			}
		}

	case "only_provider":
		providerARN, err := aws.ConnectOIDCProvider(params)
		if err != nil {
			return errors.Errorf("failed to connect OIDC provider, err: %v", err)
		}
		params.ProviderARN = providerARN
		if params.RoleName == "" {
			if err := aws.PreparePolicyAndCreateRole(params); err != nil {
				return errors.Errorf("failed to create policy and role, err: %v", err)
			}
		} else {
			if err := aws.CreateRoleWithTrustRelationsip("", params); err != nil {
				return errors.Errorf("failed to create role, err: %v", err)
			}
		}

	case "only_annotate":
		if err := kubernetes.AnnotateServiceAccount(params, *clients); err != nil {
			return errors.Errorf("failed to annotate experiment service account with role arn, err: %v", err)
		}

	default:
		return errors.Errorf("invalid action: %s", params.Actions)

	}
	return nil
}
