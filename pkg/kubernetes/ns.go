package kubernetes

import (
	"context"

	"github.com/litmuschaos/litmus-go/pkg/log"
	"github.com/uditgaurav/onboard_hce_aws/pkg/clients"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateNS will create a namespace using client-go
func CreateNS(namespaceName string, clients clients.ClientSets) error {

	namespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespaceName,
		},
	}

	_, err := clients.KubeClient.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	if err != nil {
		if err.Error() == "namespace already exists" {
			log.Infof("[Info]: The namespace %v already exist", namespaceName)
			return nil
		} else {
			return err
		}
	}

	return nil
}
