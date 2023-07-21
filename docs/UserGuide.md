# Register Harness Chaos Infrastructure

A command-line interface to register a new Harness Chaos infrastructure using specific parameters including a given name, namespace, API key, and account ID.

## Usage

```code
$ ./cli register --api-key your_api_key --account-id your_account_id --infra-name your_infra_name --project your_project --optional-flags
```

### Mandatory Flags

| Flag         | Description                                          | Example                   |
|--------------|------------------------------------------------------|---------------------------|
| `--api-key` | API Key for Harness (required) | `--api-key abc123` |
| `--account-id` | Account ID for Harness (required) | `--account-id def456` |
| `--infra-name` | Name of the Harness Chaos infrastructure (required) | `--infra-name infra-name` |
| `--project` | Project Identifier (required) | `--project project_id` |

### Optional Flags

| Flag                           | Description                                                                                       | Default                                   | Example                                      |
|--------------------------------|---------------------------------------------------------------------------------------------------|-------------------------------------------|----------------------------------------------|
| `--infra-namespace`            | Namespace for the Harness Chaos infrastructure                                                    | "hce"                                     | `--infra-namespace custom-namespace`         |
| `--organisation`               | Organisation Identifier                                                                           | "default"                                 | `--organisation organisation_id`             |
| `--infra-scope`                | Infrastructure Scope                                                                              | "namespace"                               | `--infra-scope cluster`                      |
| `--infra-ns-exists`            | Does infrastructure namespace exist                                                               | true                                      | `--infra-ns-exists false`                    |
| `--infra-description`          | Infra Description                                                                                 | "Infra for Harness Chaos Testing"         | `--infra-description "custom description"`   |
| `--infra-service-account`      | Infra Service Account                                                                             | "hce"                                     | `--infra-service-account custom-account`     |
| `--is-infra-sa-exists`         | Does infrastructure service account exist                                                         | false                                     | `--is-infra-sa-exists true`                  |
| `--infra-environment-id`       | Infra Environment ID                                                                              | ""                                        | `--infra-environment-id environment_id`      |
| `--infra-platform-name`        | Infra Platform Name                                                                               | ""                                        | `--infra-platform-name platform_name`        |
| `--infra-skip-ssl`             | Skip SSL for Infra                                                                                | false                                     | `--infra-skip-ssl true`                      |

## Description

This CLI utility is used to register a new chaos infrastructure in a Harness SaaS environment. It uses the provided API key and account ID to authenticate with the Harness API and create a new chaos infrastructure with the given name and namespace.

The utility makes a POST request to the `https://app.harness.io/gateway/api/graphql?accountId=<account_id>` endpoint with a JSON payload containing the name and namespace for the new infrastructure. The `x-api-key` HTTP header is used for authentication.

