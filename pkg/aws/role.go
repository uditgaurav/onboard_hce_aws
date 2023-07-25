package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/litmuschaos/litmus-go/pkg/log"
	"github.com/pkg/errors"
	"github.com/uditgaurav/onboard_hce_aws/pkg/types"
)

// CreateRoleWithTrustRelationsip will create the role or use a existing role with added OIDC provider
func CreateRoleWithTrustRelationsip(policyARN string, params types.OnboardingParameters) error {

	log.Infof("Provider ARN, %v", params.ProviderARN)
	// 1. Add provider to a new role with a given role name
	switch strings.TrimSpace(params.RoleName) {
	case "":
		newRoleName := "HCERole-" + params.Infra.Namespace
		log.Infof("[Info]: Creating a new role with role name '%v'", newRoleName)
		if err := addProviderToNewRole(newRoleName, policyARN, params.ProviderARN, params); err != nil {
			return err
		}
	default:
		log.Infof("[Info]: Using a existing role with roleARN '%v' for adding provider", params.RoleName)
		if err := addProviderToExistingRole(params.RoleName, params.ProviderARN, params.Region, params.Infra.Namespace, params.ExperimentServiceAccountName); err != nil {
			return err
		}
	}
	log.Info("[Info]: The role is updated successfully with provider")
	return nil
}

// addProviderToNewRole will add the OIDC provider to a new role
func addProviderToNewRole(roleName, policyARN, provider string, params types.OnboardingParameters) error {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(params.Region),
	})

	if err != nil {
		return err
	}
	svc := iam.New(sess)

	_, err = svc.CreateRole(&iam.CreateRoleInput{
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
                            "%s:aud": "system:serviceaccount:%s:%s"
                        }
                    }
                }
            ]
        }`, provider, provider, params.Infra.Namespace, params.ExperimentServiceAccountName)),
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
func addProviderToExistingRole(roleName, provider, region, ns, experimentServiceAccount string) error {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	if err != nil {
		return err
	}
	svc := iam.New(sess)

	_, err = svc.UpdateAssumeRolePolicy(&iam.UpdateAssumeRolePolicyInput{
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
                            "%s:aud": "system:serviceaccount:%s:%s"
                        }
                    }
                }
            ]
        }`, provider, provider, ns, experimentServiceAccount)),
	})

	if err != nil {
		return errors.Errorf("Error updating role", err)
	}
	return nil
}

// GetRoleARN will return the roleARN for given roleName
func GetRoleARN(region, roleName string) (string, error) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	if err != nil {
		return "", err
	}
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
