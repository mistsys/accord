package cloud_metadata

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
)

type Cloud string

const (
	AWS     Cloud = "aws"
	GCE           = "gce"
	Unknown       = "unknown"
)

func urlAlive(url string, headers map[string]string) (bool, error) {
	timeout := time.Duration(1 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	req, _ := http.NewRequest("GET", url, nil)
	if headers != nil {
		for name, value := range headers {
			req.Header.Set(name, value)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, errors.Wrapf(err, "Failed to check the metadata url %s", url)
	}
	if resp.StatusCode == 200 {
		return true, nil
	}
	return false, fmt.Errorf("URL reachable but invalid code: %d", resp.StatusCode)
}

// Do some checks to return the name of cloud env or error
// This isn't the most secure way to do it. In theory you should check the Document
// with the signature against the AWS's public certificate
// this is okay for short term deployments
func CloudService() (Cloud, error) {

	awsMetadataUrl := "http://169.254.169.254/latest/meta-data"
	// Based on https://cloud.google.com/compute/docs/storing-retrieving-metadata
	gceMetadataUrl := "http://metadata.google.internal/computeMetadata/v1"
	alive, awsErr := urlAlive(awsMetadataUrl, nil)
	if alive {
		return AWS, nil
	}
	// GCE requires all requests to have the header Metadata-Flavor: Google
	alive, gceErr := urlAlive(gceMetadataUrl, map[string]string{"Metadata-Flavor": "Google"})
	if awsErr != nil && gceErr != nil {
		return Unknown, fmt.Errorf("All Clouds checked and exhausted, not a known cloud instance: aws: %s gce: %s", awsErr, gceErr)
	}
	if alive {
		return GCE, nil
	}

	return Unknown, nil
}

// some additional data along with the Instance Document from the Metadata
type AWSInstanceInfo struct {
	IdentityDocument ec2metadata.EC2InstanceIdentityDocument `json:"identity_document"`
	IAMInfo          ec2metadata.EC2IAMInfo                  `json:"iam_info"`
	SecurityGroups   []string                                `json:"security_groups"`
	PublicIPv4       string                                  `json:"public_ipv4"`
	PublicHostname   string                                  `json:"public_hostname"`
	MAC              string                                  `json:"mac"`
	VPCId            string                                  `json:"vpc_id"`
	SubnetId         string                                  `json:"subnet_id"`
	SSHKey           string                                  `json:"ssh_key"`
}

// TODO: possibly
func GetAWSInstanceInfo() (instanceInfo *AWSInstanceInfo, err error) {
	instanceInfo = &AWSInstanceInfo{}
	sess := session.Must(session.NewSession())
	meta := ec2metadata.New(sess)
	identityDocument, err := meta.GetInstanceIdentityDocument()
	if err != nil {
		err = errors.Wrapf(err, "Failed to get Instance Document")
		return
	}
	instanceInfo.IdentityDocument = identityDocument
	iamInfo, err := meta.IAMInfo()
	if err != nil {
		err = errors.Wrapf(err, "Failed to get iam role")
		return
	}
	instanceInfo.IAMInfo = iamInfo
	secGroupStrings, err := meta.GetMetadata("security-groups")
	if err != nil {
		err = errors.Wrapf(err, "Failed to get security groups metadata")
		return
	}
	if secGroupStrings != "" {
		instanceInfo.SecurityGroups = strings.Split(secGroupStrings, "\n")
	}
	publicIPV4, err := meta.GetMetadata("public-ipv4")
	if err != nil {
		err = errors.Wrapf(err, "Failed to get ipv4 metadata")
		return
	}
	instanceInfo.PublicIPv4 = publicIPV4
	publicHostname, err := meta.GetMetadata("public-hostname")
	if err != nil {
		err = errors.Wrapf(err, "Failed to get public hostname")
		return
	}
	instanceInfo.PublicHostname = publicHostname
	mac, err := meta.GetMetadata("mac")
	if err != nil {
		err = errors.Wrapf(err, "Failed to get the mac")
		return
	}
	instanceInfo.MAC = mac
	vpcId, err := meta.GetMetadata("network/interfaces/macs/" + mac + "/vpc-id")
	if err != nil {
		err = errors.Wrapf(err, "Failed to get the vpc-id")
		return
	}
	instanceInfo.VPCId = vpcId
	subnetId, err := meta.GetMetadata("network/interfaces/macs/" + mac + "/subnet-id")
	if err != nil {
		err = errors.Wrapf(err, "Failed to get the subnet-id")
		return
	}
	instanceInfo.SubnetId = subnetId
	sshKey, err := meta.GetMetadata("public-keys/0/openssh-key")
	if err != nil {
		err = errors.Wrapf(err, "Failed to get openssh public key for the server")
		return
	}
	instanceInfo.SSHKey = sshKey
	return
}
