# Register Harness Chaos Infrastructure

A command-line interface to register a new Harness Chaos infrastructure using a given name, namespace, API key, and account ID.

## Usage

```code
$ ./cli register --api-key your_api_key --account-id your_account_id --name your_name --namespace your_namespace
```

## Flags

| Flag         | Description                                          | Example                   |
|--------------|------------------------------------------------------|---------------------------|
| `--api-key` | API Key for Harness (required) | `--api-key abc123` |
| `--account-id` | Account ID for Harness (required) | `--account-id def456` |
| `--name` | Name of the Harness Chaos infrastructure (required) | `--name infra-name` |
| `--namespace` | Namespace for the Harness Chaos infrastructure (required) | `--namespace infra-namespace` |

## Description

This CLI utility is used to register a new chaos infrastructure in a Harness SaaS environment. It uses the provided API key and account ID to authenticate with the Harness API and create a new chaos infrastructure with the given name and namespace.

The utility makes a POST request to the `https://app.harness.io/gateway/api/graphql?accountId=<account_id>` endpoint with a JSON payload containing the name and namespace for the new infrastructure. The `x-api-key` HTTP header is used for authentication.

