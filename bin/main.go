package main

import (
	"crypto/sha1"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/litmuschaos/litmus-go/pkg/log"
)

func main() {

	providerURL := "https://oidc.eks.us-east-2.amazonaws.com/id/F50A21747502C9909AEC6992E84CC870"
	clientID := "sts.amazonaws.com"
	region := "us-east-2"

	thumbprint, err := getThumbprint(providerURL)
	if err != nil {
		fmt.Println(err)
	}

	log.Info("[Info]: The thumbprint is created successfully")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	if err != nil {
		fmt.Println("Error creating session: ", err)
		return
	}

	svc := iam.New(sess)

	// Replace "ProviderURL" with your OIDC provider URL.
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
		fmt.Println("Error creating OIDC provider: ", err)
		return
	}

	fmt.Println("OIDC provider created: ", result)
}

func getThumbprint(urlStr string) (string, error) {
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println(err)
	}

	conn, err := tls.Dial("tcp", parsedUrl.Host+":443", &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	cert := conn.ConnectionState().PeerCertificates[0]
	digest := sha1.Sum(cert.Raw)
	return strings.ToUpper(hex.EncodeToString(digest[:])), nil
}
