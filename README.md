# Onboard HCE on AWS Using CLI Utility

Onboard HCE AWS is a command-line interface (CLI) utility tool that simplifies the process of onboarding on Harness Chaos Engineering (HCE). The users can simply install the chaos infrastructure with just one click using the CLI and the appropriate flags.

This tool aims to not only streamline the onboarding process but also assists in integrating with your AWS account to seamlessly run W experiments. We understand that configuring or authenticating Harness Chaos Engineering with AWS involves multiple steps that could be time-consuming. This utility has been developed to simplify those tasks.

One notable feature is that the policies required and the exact permissions needed to run a certain experiment can also be defined in the CLI as a flag.

## Key Features

- Simplifies the HCE onboarding process
- Facilitates the integration of your AWS account to run AWS chaos experiments
- Defines required policies and permissions as CLI flags and has the ability to create the Roles in AWS.
- Supports different flags that has the capability to customise the role.
- The CLI is also compatible with different Windows and Linux versions

## Documentation

For a detailed understanding of how to use this CLI utility and learn about the different flags it supports, please refer to our [User Docs](./docs/UserGuide.md).

---

1. Create chaos infra.
 - The chaos infra could be a namespace namespace or cluster scope chaos infra. [Done]
 - If it is a cluster based chaos infra then ensure you have access to create ClusterRole and CuusterRoleBinding with namespace. [Done]

2. Configure AWS account with chaos role.
 - Create AWS Role in the target account
 - Add the given provider URL to the target role.
 - Prepare the policy based on the flag if provided and create the role with the given policy attached to it

3. Annotate the experiment service account with role ARN.
 - Derive the experiment service account and annotate it with experiment service account.

Check out [userguide.md](./docs/UserGuide.md) for more details.
