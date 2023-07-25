package aws

import (
	"crypto/sha1"
	"crypto/tls"
	"encoding/hex"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/litmuschaos/litmus-go/pkg/cloud/aws/common"
	"github.com/litmuschaos/litmus-go/pkg/log"
	"github.com/pkg/errors"
	hce_types "github.com/uditgaurav/onboard_hce_aws/pkg/types"
)

// ConnectOIDCProvider will connect the provided OIDC provider in the AWS account
func ConnectOIDCProvider(onboardingParams hce_types.OnboardingParameters) (string, error) {

	clientID := "sts.amazonaws.com"
	var providerArn string

	thumbprint, err := getThumbprint(onboardingParams.ProviderUrl)
	if err != nil {
		return "", err
	}

	log.Info("[Info]: The thumbprint is created successfully")

	// Load session from shared config
	sess := common.GetAWSSession(onboardingParams.Region)

	svc := iam.New(sess)
	params := &iam.CreateOpenIDConnectProviderInput{
		Url: aws.String(onboardingParams.ProviderUrl),
		ThumbprintList: []*string{
			aws.String(thumbprint),
		},
		ClientIDList: []*string{
			aws.String(clientID),
		},
	}

	result, err := svc.CreateOpenIDConnectProvider(params)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Warnf("[Warning]: %v", err)
			arn, err := getProviderArn(onboardingParams.ProviderUrl, onboardingParams.Region)
			if err != nil {
				return "", errors.Errorf("Error getting OIDC provider ARN: %v", err)
			}
			log.Infof("[Info]: The providerARN for the given URL is: %v", arn)
			providerArn = arn
		} else {
			return "", errors.Errorf("Error creating OIDC provider: %v", err)
		}
	} else {
		log.Infof("[Info]: OIDC provider created with ARN: %v", *result.OpenIDConnectProviderArn)
		providerArn = *result.OpenIDConnectProviderArn
	}
	return providerArn, nil
}

// getThumbprint will create the thumbprint for the given provider URL
func getThumbprint(urlStr string) (string, error) {
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	conn, err := tls.Dial("tcp", parsedUrl.Host+":443", &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return "", err
	}
	defer conn.Close()

	cert := conn.ConnectionState().PeerCertificates[0]
	digest := sha1.Sum(cert.Raw)
	return strings.ToUpper(hex.EncodeToString(digest[:])), nil
}

func getProviderArn(identityProviderUrl, region string) (string, error) {

	// Load session from shared config
	sess := common.GetAWSSession(region)
	svc := iam.New(sess)

	input := &iam.ListOpenIDConnectProvidersInput{}
	result, err := svc.ListOpenIDConnectProviders(input)
	if err != nil {
		return "", errors.Errorf("Failed to list providers, err: %v", err)
	}

	for _, provider := range result.OpenIDConnectProviderList {
		providerDetails, err := svc.GetOpenIDConnectProvider(&iam.GetOpenIDConnectProviderInput{
			OpenIDConnectProviderArn: provider.Arn,
		})
		if err != nil {
			log.Infof("Failed to get provider details for %s, %v", *provider.Arn, err)
			continue
		}
		if strings.Contains(identityProviderUrl, *providerDetails.Url) {
			return *provider.Arn, nil
		}
	}

	return "", errors.Errorf("no provider found with the given URL: %s", identityProviderUrl)
}
