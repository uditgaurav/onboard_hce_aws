# onboard_hce_aws


1. Create chaos infra.
 - The chaos infra could be a namespace namespace or cluster scope chaos infra.
 - If it is a cluster based chaos infra then ensure you have access to create ClusterRole and CuusterRoleBinding with namespace.

2. Configure AWS account with chaos role.
 - Create AWS Role in the target account
 - Add the given provider URL to the target role.
 - Prepare the policy based on the flag if provided and create the role with the given policy attached to it

3. Annotate the experiment service account with role ARN.
 - Derive the experiment service account and annotate it with experiment service account.

Check out [userguide.md](./docs/UserGuide.md) for more details.

