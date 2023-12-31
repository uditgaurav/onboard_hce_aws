package clients

import (
	"os"

	"github.com/pkg/errors"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// ClientSets is a collection of clientSets and kubeConfig needed
type ClientSets struct {
	KubeClient    *kubernetes.Clientset
	KubeConfig    *rest.Config
	DynamicClient dynamic.Interface
}

// GenerateClientSetFromKubeConfig will generation both ClientSets (k8s, and Litmus) as well as the KubeConfig
func (clientSets *ClientSets) GenerateClientSetFromKubeConfig() error {

	config, err := getKubeConfig()
	if err != nil {
		return err
	}
	k8sClientSet, err := generateK8sClientSet(config)
	if err != nil {
		return err
	}
	dynamicClientSet, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}
	clientSets.KubeClient = k8sClientSet
	clientSets.KubeConfig = config
	clientSets.DynamicClient = dynamicClientSet
	return nil
}

// getKubeConfig setup the config for access cluster resource
func getKubeConfig() (*rest.Config, error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	// If KUBECONFIG is not specified, look at the default location
	if kubeconfig == "" {
		kubeconfig = os.Getenv("HOME") + "/.kube/config"
	}
	// It uses in-cluster config, if kubeconfig path is not specified
	config, err := buildConfigFromFlags("", kubeconfig)
	return config, err
}

// generateK8sClientSet will generation k8s client
func generateK8sClientSet(config *rest.Config) (*kubernetes.Clientset, error) {
	k8sClientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to generate kubernetes clientSet, err: %v: ", err)
	}
	return k8sClientSet, nil
}

// buildConfigFromFlags is a helper function that builds configs from a master
// url or a kubeconfig filepath, if nothing is provided it falls back to inClusterConfig
func buildConfigFromFlags(masterUrl, kubeconfigPath string) (*restclient.Config, error) {
	if kubeconfigPath == "" && masterUrl == "" {
		kubeconfig, err := restclient.InClusterConfig()
		if err == nil {
			return kubeconfig, nil
		}
	}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: masterUrl}}).ClientConfig()
}
