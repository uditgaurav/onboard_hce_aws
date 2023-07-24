package aws

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/litmuschaos/litmus-go/pkg/cloud/aws/common"
	"github.com/litmuschaos/litmus-go/pkg/log"
	"github.com/pkg/errors"
	"github.com/uditgaurav/onboard_hce_aws/pkg/types"
)

// PreparePolicyAndCreateRole will prepare a policy JSON based on the target resource provided
func PreparePolicyAndCreateRole(params types.OnboardingParameters) error {

	policyName := "HCEChaosPolicy-" + params.Infra.Namespace

	log.Info("[Info]: Preparing policy for the role")
	resources := strings.Split(params.Resources, ",")

	// Prepare combined policy
	combinedPolicy := Policy{Version: "2012-10-17"}
	actionSet := make(map[string]bool)
	for _, resource := range resources {
		switch resource {
		case "ec2":
			for _, stmt := range ec2Policy.Statement {
				for _, action := range stmt.Action {
					actionSet[action] = true
				}
			}
		case "lambda":
			for _, stmt := range lambdaPolicy.Statement {
				for _, action := range stmt.Action {
					actionSet[action] = true
				}
			}
		case "aws-access-restrict":
			for _, stmt := range awsAccessRestrictPolicy.Statement {
				for _, action := range stmt.Action {
					actionSet[action] = true
				}
			}
		case "az":
			for _, stmt := range azPolicy.Statement {
				for _, action := range stmt.Action {
					actionSet[action] = true
				}
			}
		case "ebs":
			for _, stmt := range ebsPolicy.Statement {
				for _, action := range stmt.Action {
					actionSet[action] = true
				}
			}
		case "ec2-state":
			for _, stmt := range ec2StatePolicy.Statement {
				for _, action := range stmt.Action {
					actionSet[action] = true
				}
			}
		case "ecs-ec2":
			for _, stmt := range ecsEc2Policy.Statement {
				for _, action := range stmt.Action {
					actionSet[action] = true
				}
			}
		case "ecs-fargate":
			for _, stmt := range ecsFargatePolicy.Statement {
				for _, action := range stmt.Action {
					actionSet[action] = true
				}
			}
		case "ecs-state":
			for _, stmt := range ecsStatePolicy.Statement {
				for _, action := range stmt.Action {
					actionSet[action] = true
				}
			}
		case "lambda-permission":
			for _, stmt := range lambdaPermissionPolicy.Statement {
				for _, action := range stmt.Action {
					actionSet[action] = true
				}
			}
		case "rds":
			for _, stmt := range rdsPolicy.Statement {
				for _, action := range stmt.Action {
					actionSet[action] = true
				}
			}
		case "windows":
			for _, stmt := range windowsPolicy.Statement {
				for _, action := range stmt.Action {
					actionSet[action] = true
				}
			}
		case "all":
			for _, stmt := range allPolicy.Statement {
				for _, action := range stmt.Action {
					actionSet[action] = true
				}
			}
		default:
			return errors.Errorf("unknown resource type: %v", resource)
		}
	}

	for action := range actionSet {
		combinedPolicy.Statement = append(combinedPolicy.Statement, Statement{Effect: "Allow", Action: []string{action}, Resource: "*"})
	}

	log.Infof("[Info]: Prepared policy: %+v", combinedPolicy)
	policyARN, err := createPolicy(combinedPolicy, policyName, params.Region)
	if err != nil {
		return err
	}

	log.Infof("[Info]: The policy is successfully created")
	log.Infof("[Info]: Creating AWS role")
	if err := CreateRoleWithTrustRelationsip(policyARN, params); err != nil {
		return errors.Errorf("failed to create role, err: %v", err)
	}
	return nil
}

// createPolicy will create the given policy
func createPolicy(policy Policy, policyName, region string) (string, error) {

	// Load session from shared config
	sess := common.GetAWSSession(region)

	// Create IAM service client
	svc := iam.New(sess)

	// Define policy document from your struct (assuming json.Marshal does the right thing here)
	policyDoc, err := json.Marshal(policy)
	if err != nil {
		return "", errors.Errorf("Error marshaling policy document: %v", err)
	}

	// Create policy
	resp, err := svc.CreatePolicy(&iam.CreatePolicyInput{
		PolicyDocument: aws.String(string(policyDoc)),
		PolicyName:     aws.String(policyName),
	})
	if err != nil {
		return "", errors.Errorf("Error marshaling policy document: %v", err)
	}

	log.Info("[Indo]: Policy successfully created.")
	return *resp.Policy.Arn, nil
}
