package register

import (
	"bytes"
	ejson "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"context"

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
		logrus.Error("Error serializing payload to JSON:", err)
		return
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		logrus.Error("Error creating request:", err)
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
		logrus.Error("Error on response:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response data
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("Error reading response data:", err)
		return
	}

	// Parse the response data into the Response struct
	var responseData Response
	err = ejson.Unmarshal(data, &responseData)
	if err != nil {
		logrus.Error("Error parsing JSON response:", err)
		return
	}

	// Print the manifest
	if responseData.Data.RegisterInfra.Manifest != "" {
		logrus.Info("Chaos infra Manifest prepared")
	} else {
		logrus.Error("Chaos infra Manifest is empty")
	}

	if err := applyChaosManifest(responseData.Data.RegisterInfra.Token, params.Infra.Namespace, responseData.Data.RegisterInfra.Manifest); err != nil {
		logrus.Fatal("Failed to create chaos infra manifest: ", err)
	}

}

func applyChaosManifest(token string, namespace string, manifest string) error {
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
		_, err := dynamicClient.Resource(gvr).Namespace(namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("Error applying manifest: %v", err)
		}
	}

	logrus.Info("Successfully applied chaos infra manifest to Kubernetes cluster")
	return nil
}
