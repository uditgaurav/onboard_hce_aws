package register

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// RegisterInfra is a function to register infrastructure details using the Harness API.
func RegisterInfra(params InfraParameters) {
	// The API endpoint URL
	url := fmt.Sprintf("https://app.harness.io/gateway/chaos/manager/api/query?accountIdentifier=%s", params.AccountId)

	// GraphQL mutation query to register infrastructure
	query := `mutation($identifiers: IdentifiersRequest!, $request: RegisterInfraRequest!) {
  	registerInfra(identifiers: $identifiers, request: $request) {
    	token
    	infraID
    	name
    	manifest
  	}
	}`

	// Set up the request variables
	variables := Variables{
		Identifiers: Identifiers{
			OrgIdentifier:     params.Organisation,
			AccountIdentifier: params.AccountId,
			ProjectIdentifier: params.Project,
		},
		Request: Request{
			Name:             params.Infra.Name,
			EnvironmentID:    params.Infra.EnvironmentID,
			Description:      params.Infra.Description,
			InfraNamespace:   params.Infra.Namespace,
			ServiceAccount:   params.Infra.ServiceAccount,
			InfraScope:       params.InfraScope,
			InfraNsExists:    params.InfraNsExists,
			InfraSaExists:    params.Infra.InfraSaExists,
			InstallationType: params.Infra.InstallationType,
			SkipSsl:          params.Infra.SkipSsl,
		},
	}

	// Create the payload for the API call
	payload := Payload{
		Query:     query,
		Variables: variables,
	}

	// Serialize the payload to JSON
	body, _ := json.Marshal(payload)

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", params.ApiKey)
	req.Header.Set("Type", "ApiKey")

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error on response:\n", err)
		return
	}
	defer resp.Body.Close()

	// Read and print the response data
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))
}
