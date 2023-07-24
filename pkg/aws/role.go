package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/litmuschaos/litmus-go/pkg/cloud/aws/common"
	"github.com/litmuschaos/litmus-go/pkg/log"
	"github.com/pkg/errors"
	"github.com/uditgaurav/onboard_hce_aws/pkg/types"
)

// CreateRoleWithTrustRelationsip will create the role or use a existing role with added OIDC provider
func CreateRoleWithTrustRelationsip(policyARN string, params types.OnboardingParameters) error {

	// Load session from shared config
	sess := common.GetAWSSession(params.Region)

	// Create IAM service client
	svc := iam.New(sess)

	audience := "sts.amazonaws.com"

	// 1. Add provider to a new role with a given role name
	switch strings.TrimSpace(params.RoleName) {
	case "":
		newRoleName := "HCERole-" + params.Infra.Namespace
		log.Infof("[Info]: Creating a new role with role name '%v'", newRoleName)
		if err := addProviderToNewRole(svc, newRoleName, policyARN, params.ProviderARN, audience, params); err != nil {
			return aws.ErrMissingEndpoint
		}
	default:
		log.Infof("[Info]: Using a existing role with roleARN '%v' for adding provider", params.RoleName)
		if err := addProviderToExistingRole(svc, params.RoleName, params.ProviderARN, audience); err != nil {
			return aws.ErrMissingEndpoint
		}
	}
	log.Info("[Info]: The role is created successfully with provider")
	return nil
}

// addProviderToNewRole will add the OIDC provider to a new role
func addProviderToNewRole(svc *iam.IAM, roleName, policyARN, provider, audience string, params types.OnboardingParameters) error {

	_, err := svc.CreateRole(&iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(fmt.Sprintf(`{
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": {
                        "Federated": "%s"
                    },
                    "Action": "sts:AssumeRoleWithWebIdentity",
                    "Condition": {
                        "StringEquals": {
                            "%s:aud": "%s"
                        }
                    }
                }
            ]
        }`, provider, provider, audience)),
		Path:     aws.String("/"),
		RoleName: aws.String(roleName),
	})

	if err != nil {
		return errors.Errorf("Error creating role: %v", err)
	}

	// Attach policy to the newly created role
	_, err = svc.AttachRolePolicy(&iam.AttachRolePolicyInput{
		PolicyArn: aws.String(policyARN),
		RoleName:  aws.String(roleName),
	})

	if err != nil {
		return errors.Errorf("Error attaching policy", err)
	}
	return nil
}

// addProviderToExistingRole will add the OIDC provider to an existing role
func addProviderToExistingRole(svc *iam.IAM, roleName, provider, audience string) error {

	_, err := svc.UpdateAssumeRolePolicy(&iam.UpdateAssumeRolePolicyInput{
		RoleName: aws.String(roleName),
		PolicyDocument: aws.String(fmt.Sprintf(`{
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": {
                        "Federated": "%s"
                    },
                    "Action": "sts:AssumeRoleWithWebIdentity",
                    "Condition": {
                        "StringEquals": {
                            "%s:aud": "%s"
                        }
                    }
                }
            ]
        }`, provider, provider, audience)),
	})

	if err != nil {
		return errors.Errorf("Error updating role", err)
	}
	return nil
}

// GetRoleARN will return the roleARN for given roleName
func GetRoleARN(region, roleName string) (string, error) {

	// Load session from shared config
	sess := common.GetAWSSession(region)
	svc := iam.New(sess)

	input := &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	}

	result, err := svc.GetRole(input)
	if err != nil {
		return "", err
	}

	return *result.Role.Arn, nil
}
