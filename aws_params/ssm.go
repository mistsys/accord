package aws_params

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/pkg/errors"
)

// Config for getting parameters from the environment
// Region is the AWS Region
// Role is the special role that has permission to read the configs, strings, etc
type Config struct {
	Region string
	// RoleArn is the special string that
	RoleArn string
}

type Client interface {
	GetSecureString(path string) (string, error)
}

type client struct {
	cfg *Config
	ssm ssmiface.SSMAPI
}

func (c *client) GetSecureString(path string) (string, error) {
	params := &ssm.GetParametersInput{
		Names: []*string{
			aws.String(path),
		},
		WithDecryption: aws.Bool(true),
	}
	resp, err := c.ssm.GetParameters(params)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to read secure string from %s", path)
	}
	if len(resp.InvalidParameters) > 0 {
		return "", fmt.Errorf("Invalid parameter %s", path)
	}
	val := resp.Parameters[0].Value
	return *val, nil
}

func NewConfig(region string) *Config {
	cfg := &Config{Region: region}
	return cfg
}

func NewClient(cfg *Config) (Client, error) {
	if cfg == nil {
		// just default
		cfg = NewConfig("us-east-1")
	}

	if cfg.Region == "" {
		return nil, errors.New("params.NewClient requires an AWS Region")
	}

	if cfg.RoleArn == "" {
		return nil, errors.New("params.NewClient requires a RoleARN to assume")
	}

	var awsConfig *aws.Config

	sess, err := session.NewSession()
	if err != nil {
		return nil, errors.Wrapf(err, "params.NewClient failed to create a ssm session")
	}

	// this is just for private environments, and not cloud
	if cfg.RoleArn != "" {
		creds := stscreds.NewCredentials(sess, cfg.RoleArn)
		awsConfig = &aws.Config{Credentials: creds, Region: aws.String(cfg.Region), CredentialsChainVerboseErrors: aws.Bool(true)}
	} else {
		awsConfig = aws.NewConfig().WithRegion(cfg.Region).WithCredentialsChainVerboseErrors(true)
	}

	cl := &client{
		cfg: cfg,
		ssm: ssm.New(sess, awsConfig),
	}
	return cl, nil

}
