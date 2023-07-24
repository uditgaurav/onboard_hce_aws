package register

import (
	"bytes"
	ejson "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"context"

	"github.com/litmuschaos/litmus-go/pkg/log"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/uditgaurav/onboard_hce_aws/pkg/clients"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/runtime/serializer/recognizer"
)

// RegisterInfra is a function to register infrastructure details using the Harness API.
func RegisterInfra(params InfraParameters) error {

	log.InfoWithValues("[Info]: Creating the chaos infra with following details:", logrus.Fields{
		"ChaosInfra Name":             params.Infra.Name,
		"ChaosInfra Namespace":        params.Infra.Namespace,
		"ChoasInfra Scope":            params.InfraScope,
		"ChaosInfra Service Acccount": params.Infra.ServiceAccount,
		"Environment":                 params.Infra.Name + "-env",
	})

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

	// If the user didn't provide infra-environment-id, then set it to infra-name with a '-env' suffix.
	if params.Infra.EnvironmentID == "" {
		params.Infra.EnvironmentID = params.Infra.Name + "-env"
	}

	// If the user didn't provide infra-platform-name, then set it to infra-name with a '-platform' suffix.
	if params.Infra.PlatformName == "" {
		params.Infra.PlatformName = params.Infra.Name + "-platform"
	}

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
			PlatformName:     params.Infra.PlatformName,
			InfraNamespace:   params.Infra.Namespace,
			ServiceAccount:   params.Infra.ServiceAccount,
			InfraScope:       params.InfraScope,
			InfraNsExists:    params.InfraNsExists,
			InfraSaExists:    params.Infra.InfraSaExists,
			InstallationType: "MANIFEST",
			SkipSsl:          params.Infra.SkipSsl,
		},
	}

	// Create the payload for the API call
	payload := Payload{
		Query:     query,
		Variables: variables,
	}

	// Serialize the payload to JSON
	body, err := ejson.Marshal(payload)
	if err != nil {
		return errors.Errorf("Error serializing payload to JSON: %v", err)
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return errors.Errorf("Error creating request: %v", err)
	}

	// Set the required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", params.ApiKey)
	req.Header.Set("Type", "ApiKey")

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Errorf("Error on response: %v", err)

	}
	defer resp.Body.Close()

	// Read the response data
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Errorf("Error reading response data: %v", err)
	}

	// Parse the response data into the Response struct
	var responseData Response
	err = ejson.Unmarshal(data, &responseData)
	if err != nil {
		return errors.Errorf("Error parsing JSON response: %v", err)
	}

	if responseData.Data.RegisterInfra.Manifest != "" {
		log.Info("[Info]: Chaos Infra Manifest prepared")
	} else {
		return errors.Errorf("[Info]: The prepared chaos infra manifest is empty")
	}

	if err := applyChaosManifest(responseData.Data.RegisterInfra.Token, responseData.Data.RegisterInfra.Manifest, responseData.Data.RegisterInfra.InfraID, params); err != nil {
		return errors.Errorf("Failed to create chaos infra manifest: ", err)
	}
	return nil
}

// applyChaosManifest will create the chaosYAML manifest created while registring infra
func applyChaosManifest(token, manifest, infraID string, params InfraParameters) error {
	clients := clients.ClientSets{}

	//Getting kubeConfig and Generate ClientSets
	if err := clients.GenerateClientSetFromKubeConfig(); err != nil {
		return fmt.Errorf("Unable to Get the kubeconfig, err: %v", err)
	}

	// Get the dynamic client for unstructured objects
	dynamicClient := clients.DynamicClient

	// Remove first line of the manifest if it contains '---'
	if strings.HasPrefix(manifest, "---") {
		lines := strings.Split(manifest, "\n")
		if len(lines) > 1 {
			manifest = strings.Join(lines[1:], "\n")
		} else {
			manifest = ""
		}
	}

	// Split the manifest into individual resources
	manifests := strings.Split(manifest, "---")

	log.Info("[Info]: Creating the manifest to install chaos infra")
	for _, m := range manifests {
		// Decode the YAML manifest into an unstructured object
		obj := &unstructured.Unstructured{}

		// Create a new scheme
		s := runtime.NewScheme()

		// Create the recogniser decoder
		d := recognizer.NewDecoder(
			json.NewSerializer(json.DefaultMetaFactory, s, s, false),
			json.NewYAMLSerializer(json.DefaultMetaFactory, s, s),
		)

		// Decode YAML to unstructured object
		if _, _, err := d.Decode([]byte(m), nil, obj); err != nil {
			return fmt.Errorf("Error decoding manifest: %v", err)
		}

		// Get the GVR from the object to create a dynamic client for that GVR
		gvk := obj.GroupVersionKind()
		gvr, _ := meta.UnsafeGuessKindToResource(gvk)

		// Create the object using the dynamic client
		_, err := dynamicClient.Resource(gvr).Namespace(params.Infra.Namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("Error applying manifest: %v", err)
		}
	}

	log.Info("[Info]: Successfully applied chaos infra manifest to Kubernetes cluster")
	if err := waitForChaosInfra(infraID, params); err != nil {
		return errors.Errorf("failed to get the chaos infra in Connected state, err: $v", err)
	}
	return nil
}

// getChaosInfraState fetches the current state of the chaos infrastructure
func getChaosInfraState(infraID string, params InfraParameters) (bool, error) {

	// The API endpoint URL
	url := fmt.Sprintf("https://app.harness.io/gateway/chaos/manager/api/query?accountIdentifier=%s", params.AccountId)

	// Define the GraphQL query and variables
	query := `query GetInfra($infraID: String!, $identifiers: IdentifiersRequest!) {
			getInfra(infraID: $infraID, identifiers: $identifiers) {
				infraID
				name
				description
				tags
				environmentID
				platformName
				isActive
				isInfraConfirmed
				isRemoved
				updatedAt
				createdAt
				noOfSchedules
				noOfWorkflows
				token
				infraNamespace
				serviceAccount
				infraScope
				infraNsExists
				infraSaExists
				installationType
				k8sConnectorID
				lastWorkflowTimestamp
				startTime
				version
				createdBy {
					userID
					username
					email
				}
				updatedBy {
					userID
					username
					email
				}
			}
		}`
	variables := map[string]interface{}{
		"identifiers": map[string]string{
			"orgIdentifier":     params.Organisation,
			"accountIdentifier": params.AccountId,
			"projectIdentifier": params.Project,
		},
		"infraID": infraID,
	}

	// Create the request body
	reqBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}
	reqBodyBytes, err := ejson.Marshal(reqBody)
	if err != nil {
		return false, errors.Errorf("error creating request body: %v", err)
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, strings.NewReader(string(reqBodyBytes)))
	if err != nil {
		return false, errors.Errorf("error creating request: %v", err)
	}

	// Set the required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", params.ApiKey)

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, errors.Errorf("error on response: %v", err)
	}
	defer resp.Body.Close()

	// Read the response data
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, errors.Errorf("error reading response data: %v", err)
	}

	// Parse the response data into the Response struct
	var responseData struct {
		Data struct {
			GetInfra struct {
				IsActive bool `json:"isActive"`
			} `json:"getInfra"`
		} `json:"data"`
	}

	err = ejson.Unmarshal(data, &responseData)
	if err != nil {
		return false, errors.Errorf("error parsing JSON response: %v", err)
	}
	return responseData.Data.GetInfra.IsActive, nil
}

// waitForChaosInfra will wait for the chaos infra to get in active state for the given timeout.
func waitForChaosInfra(infraID string, params InfraParameters) error {
	timeout := time.After(180 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return errors.New("timeout reached while waiting for infra")
		case <-ticker.C:
			result, err := getChaosInfraState(infraID, params)
			if err != nil {
				return err
			}
			if result {
				log.Info("[Info]: The infra is now activated!")
				return nil
			}
			log.Info("[Info]: The infra is not activated yet")
		}
	}
}
