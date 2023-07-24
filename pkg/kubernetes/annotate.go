package kubernetes

import (
	"context"

	"github.com/pkg/errors"
	"github.com/uditgaurav/onboard_hce_aws/pkg/aws"
	"github.com/uditgaurav/onboard_hce_aws/pkg/clients"
	"github.com/uditgaurav/onboard_hce_aws/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AnnotateServiceAccount will annotate the given experiment service account with aws roleARN
func AnnotateServiceAccount(params types.OnboardingParameters, clients clients.ClientSets) error {

	roleARN, err := aws.GetRoleARN(params.Region, params.RoleName)
	if err != nil {
		return errors.Errorf("failed to retrive roleARN from given role name '%v', err: %v", params.RoleName, err)
	}
	sa, err := clients.KubeClient.CoreV1().ServiceAccounts(params.Infra.Namespace).Get(context.Background(), params.ExperimentServiceAccountName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if sa.Annotations == nil {
		sa.Annotations = make(map[string]string)
	}

	sa.Annotations["eks.amazonaws.com/role-arn"] = roleARN

	_, err = clients.KubeClient.CoreV1().ServiceAccounts(params.Infra.Namespace).Update(context.Background(), sa, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}
