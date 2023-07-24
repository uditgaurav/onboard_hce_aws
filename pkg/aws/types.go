package aws

type Statement struct {
	Effect   string
	Action   []string
	Resource string
}

type Policy struct {
	Version   string
	Statement []Statement
}

// This is the predefined policy statement for EC2.
var ec2PolicyStatement = []Statement{
	{
		Effect: "Allow",
		Action: []string{
			"ssm:GetDocument",
			"ssm:DescribeDocument",
			"ssm:GetParameter",
			"ssm:GetParameters",
			"ssm:SendCommand",
			"ssm:CancelCommand",
			"ssm:CreateDocument",
			"ssm:DeleteDocument",
			"ssm:GetCommandInvocation",
			"ssm:UpdateInstanceInformation",
			"ssm:DescribeInstanceInformation",
			"ec2messages:AcknowledgeMessage",
			"ec2messages:DeleteMessage",
			"ec2messages:FailMessage",
			"ec2messages:GetEndpoint",
			"ec2messages:GetMessages",
			"ec2messages:SendReply",
			"ec2:DescribeInstanceStatus",
			"ec2:DescribeInstances",
		},
		Resource: "*",
	},
}

// This is the predefined policy statement for Lambda.
var lambdaPolicyStatement = []Statement{
	{
		Effect: "Allow",
		Action: []string{
			"lambda:ListEventSourceMappings",
			"lambda:DeleteEventSourceMapping",
			"lambda:UpdateEventSourceMapping",
			"lambda:CreateEventSourceMapping",
			"lambda:UpdateFunctionConfiguration",
			"lambda:GetFunctionConcurrency",
			"lambda:GetFunction",
			"lambda:DeleteFunctionConcurrency",
			"lambda:PutFunctionConcurrency",
		},
		Resource: "*",
	},
}

// This is the predefined policy statement for RDS.
var rdsPolicyStatement = []Statement{
	{
		Effect: "Allow",
		Action: []string{
			"ec2:DescribeInstanceStatus",
			"ec2:DescribeInstances",
			"rds:DescribeDBClusters",
			"rds:DescribeDBInstances",
			"rds:DeleteDBInstance",
			"rds:RebootDBInstance",
		},
		Resource: "*",
	},
}

// This is the predefined policy statement for AWS Access Restrict.
var awsAccessRestrictPolicyStatement = []Statement{
	{
		Effect: "Allow",
		Action: []string{
			"ec2:DescribeSecurityGroups",
			"ec2:RevokeSecurityGroupIngress",
			"ec2:AuthorizeSecurityGroupIngress",
			"ec2:RevokeSecurityGroupEgress",
			"ec2:AuthorizeSecurityGroupEgress",
		},
		Resource: "*",
	},
}

// This is the predefined policy statement for AZ.
var azPolicyStatement = []Statement{
	{
		Effect: "Allow",
		Action: []string{
			"ec2:DescribeInstanceStatus",
			"ec2:DescribeInstances",
			"ec2:DescribeSubnets",
			"elasticloadbalancing:DetachLoadBalancerFromSubnets",
			"elasticloadbalancing:AttachLoadBalancerToSubnets",
			"elasticloadbalancing:DescribeLoadBalancers",
			"ec2:CreateNetworkAcl",
			"ec2:CreateNetworkAclEntry",
			"ec2:DescribeNetworkAcls",
			"ec2:ReplaceNetworkAclAssociation",
			"ec2:DeleteNetworkAcl",
		},
		Resource: "*",
	},
}

// This is the predefined policy statement for EBS.
var ebsPolicyStatement = []Statement{
	{
		Effect: "Allow",
		Action: []string{
			"ec2:AttachVolume",
			"ec2:DetachVolume",
			"ec2:DescribeVolumes",
			"ec2:DescribeInstanceStatus",
			"ec2:DescribeInstances",
		},
		Resource: "*",
	},
}

// This is the predefined policy statement for EC2-State.
var ec2StatePolicyStatement = []Statement{
	{
		Effect: "Allow",
		Action: []string{
			"ec2:StartInstances",
			"ec2:StopInstances",
			"ec2:DescribeInstanceStatus",
			"ec2:DescribeInstances",
			"autoscaling:DescribeAutoScalingInstances",
		},
		Resource: "*",
	},
}

// This is the predefined policy statement for ECS-EC2.
var ecsEc2PolicyStatement = []Statement{
	{
		Effect: "Allow",
		Action: []string{
			"ssm:GetDocument",
			"ssm:DescribeDocument",
			"ssm:GetParameter",
			"ssm:GetParameters",
			"ssm:SendCommand",
			"ssm:CancelCommand",
			"ssm:CreateDocument",
			"ssm:DeleteDocument",
			"ssm:GetCommandInvocation",
			"ssm:UpdateInstanceInformation",
			"ssm:DescribeInstanceInformation",
			"ec2messages:AcknowledgeMessage",
			"ec2messages:DeleteMessage",
			"ec2messages:FailMessage",
			"ec2messages:GetEndpoint",
			"ec2messages:GetMessages",
			"ec2messages:SendReply",
			"ec2:DescribeInstanceStatus",
			"ec2:DescribeInstances",
		},
		Resource: "*",
	},
}

// This is the predefined policy statement for ECS-Fargate.
var ecsFargatePolicyStatement = []Statement{
	{
		Effect: "Allow",
		Action: []string{
			"ecs:DescribeTasks",
			"ecs:DescribeServices",
			"ecs:DescribeTaskDefinition",
			"ecs:RegisterTaskDefinition",
			"ecs:UpdateService",
			"ecs:ListTasks",
			"ecs:DeregisterTaskDefinition",
			"iam:PassRole",
		},
		Resource: "*",
	},
}

// This is the predefined policy statement for ECS-State.
var ecsStatePolicyStatement = []Statement{
	{
		Effect: "Allow",
		Action: []string{
			"ecs:ListServices",
			"ecs:ListTasks",
			"ecs:StopTask",
			"ecs:DescribeServices",
			"ecs:DescribeTasks",
			"ecs:ListContainerInstances",
			"ecs:DescribeContainerInstances",
			"ec2:StartInstances",
			"ec2:StopInstances",
			"ec2:DescribeInstanceStatus",
			"ec2:DescribeInstances",
			"autoscaling:DescribeAutoScalingInstances",
		},
		Resource: "*",
	},
}

// This is the predefined policy statement for Lambda-Permission.
var lambdaPermissionPolicyStatement = []Statement{
	{
		Effect: "Allow",
		Action: []string{
			"iam:PassRole",
			"lambda:GetFunction",
			"lambda:UpdateFunctionConfiguration",
			"iam:AttachRolePolicy",
			"iam:DetachRolePolicy",
			"iam:ListAttachedRolePolicies",
			"iam:GetRolePolicy",
		},
		Resource: "*",
	},
}

// This is the predefined policy statement for Windows.
var windowsPolicyStatement = []Statement{
	{
		Effect: "Allow",
		Action: []string{
			"ssm:GetDocument",
			"ssm:DescribeDocument",
			"ssm:GetParameter",
			"ssm:GetParameters",
			"ssm:SendCommand",
			"ssm:CancelCommand",
			"ssm:CreateDocument",
			"ssm:DeleteDocument",
			"ssm:GetCommandInvocation",
			"ssm:UpdateInstanceInformation",
			"ssm:DescribeInstanceInformation",
			"ec2messages:AcknowledgeMessage",
			"ec2messages:DeleteMessage",
			"ec2messages:FailMessage",
			"ec2messages:GetEndpoint",
			"ec2messages:GetMessages",
			"ec2messages:SendReply",
			"ec2:DescribeInstanceStatus",
			"ec2:DescribeInstances",
		},
		Resource: "*",
	},
}

// This is the predefined policy statement for All.
var allPolicyStatement = []Statement{
	{
		Effect: "Allow",
		Action: []string{
			"ec2:StartInstances",
			"ec2:StopInstances",
			"ec2:AttachVolume",
			"ec2:DetachVolume",
			"ec2:DescribeVolumes",
			"ec2:DescribeSubnets",
			"ec2:DescribeInstanceStatus",
			"ec2:DescribeInstances",
			"ec2messages:AcknowledgeMessage",
			"ec2messages:DeleteMessage",
			"ec2messages:FailMessage",
			"ec2messages:GetEndpoint",
			"ec2messages:GetMessages",
			"ec2messages:SendReply",
			"ec2:AuthorizeSecurityGroupEgress",
			"ec2:RevokeSecurityGroupEgress",
			"ec2:RevokeSecurityGroupIngress",
			"ec2:DescribeSecurityGroups",
			"autoscaling:DescribeAutoScalingInstances",
			"ssm:GetDocument",
			"ssm:DescribeDocument",
			"ssm:GetParameter",
			"ssm:GetParameters",
			"ssm:SendCommand",
			"ssm:CancelCommand",
			"ssm:CreateDocument",
			"ssm:DeleteDocument",
			"ssm:GetCommandInvocation",
			"ssm:UpdateInstanceInformation",
			"ssm:DescribeInstanceInformation",
			"ecs:UpdateContainerInstancesState",
			"ecs:RegisterContainerInstance",
			"ecs:ListContainerInstances",
			"ecs:DeregisterContainerInstance",
			"ecs:DescribeContainerInstances",
			"ecs:ListTasks",
			"ecs:DescribeClusters",
			"ecs:ListServices",
			"ecs:StopTask",
			"ecs:DescribeServices",
			"ecs:DescribeTaskDefinition",
			"ecs:RegisterTaskDefinition",
			"ecs:DeregisterTaskDefinition",
			"ecs:UpdateService",
			"ecs:DescribeTasks",
			"elasticloadbalancing:DetachLoadBalancerFromSubnets",
			"elasticloadbalancing:AttachLoadBalancerToSubnets",
			"elasticloadbalancing:DescribeLoadBalancers",
			"lambda:ListEventSourceMappings",
			"lambda:DeleteEventSourceMapping",
			"lambda:UpdateEventSourceMapping",
			"lambda:CreateEventSourceMapping",
			"lambda:UpdateFunctionConfiguration",
			"lambda:GetFunctionConcurrency",
			"lambda:GetFunction",
			"lambda:DeleteFunctionConcurrency",
			"lambda:PutFunctionConcurrency",
			"lambda:DeleteLayerVersion",
			"lambda:GetLayerVersion",
			"lambda:ListLayerVersions",
			"rds:DescribeDBClusters",
			"rds:DescribeDBInstances",
			"rds:DeleteDBInstance",
			"rds:RebootDBInstance",
		},
		Resource: "*",
	},
}

// Initialize policies for each resource type
var ec2Policy = Policy{Version: "2012-10-17", Statement: ec2PolicyStatement}
var rdsPolicy = Policy{Version: "2012-10-17", Statement: rdsPolicyStatement}
var lambdaPolicy = Policy{Version: "2012-10-17", Statement: lambdaPolicyStatement}
var awsAccessRestrictPolicy = Policy{Version: "2012-10-17", Statement: awsAccessRestrictPolicyStatement}
var azPolicy = Policy{Version: "2012-10-17", Statement: azPolicyStatement}
var ebsPolicy = Policy{Version: "2012-10-17", Statement: ebsPolicyStatement}
var ec2StatePolicy = Policy{Version: "2012-10-17", Statement: ec2StatePolicyStatement}
var ecsEc2Policy = Policy{Version: "2012-10-17", Statement: ecsEc2PolicyStatement}
var ecsFargatePolicy = Policy{Version: "2012-10-17", Statement: ecsFargatePolicyStatement}
var ecsStatePolicy = Policy{Version: "2012-10-17", Statement: ecsStatePolicyStatement}
var lambdaPermissionPolicy = Policy{Version: "2012-10-17", Statement: lambdaPermissionPolicyStatement}
var windowsPolicy = Policy{Version: "2012-10-17", Statement: windowsPolicyStatement}
var allPolicy = Policy{Version: "2012-10-17", Statement: allPolicyStatement}
