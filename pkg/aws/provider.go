package main

import (
	"crypto/sha1"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/litmuschaos/litmus-go/pkg/cloud/aws/common"
	"github.com/litmuschaos/litmus-go/pkg/log"
	"github.com/pkg/errors"
)

// ConnectOIDCProvider will connect the provided OIDC provider in the AWS account
func ConnectOIDCProvider(providerURL, clientID, region string) error {

	thumbprint, err := getThumbprint(providerURL)
	if err != nil {
		fmt.Println(err)
	}

	log.Info("[Info]: The thumbprint is created successfully")

	// Load session from shared config
	sess := common.GetAWSSession(region)

	svc := iam.New(sess)
	params := &iam.CreateOpenIDConnectProviderInput{
		Url: aws.String(providerURL),
		ThumbprintList: []*string{
			aws.String(thumbprint),
		},
		ClientIDList: []*string{
			aws.String(clientID),
		},
	}

	result, err := svc.CreateOpenIDConnectProvider(params)
	if err != nil {
		return errors.Errorf("Error creating OIDC provider: %v", err)
	}
	log.Infof("[Info]: OIDC provider created with ARN: %v", result.OpenIDConnectProviderArn)
	return nil
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
