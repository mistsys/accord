package aws_params

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type MockedSSMAPI struct {
	ssmiface.SSMAPI
	mSecureStrings map[string]*ssm.Parameter
	mStrings       map[string]*ssm.Parameter
}

func NewMockedSSMAPI(secureStrings map[string]string, insecureStrings map[string]string) *MockedSSMAPI {
	secStrings := make(map[string]*ssm.Parameter)
	insecStrings := make(map[string]*ssm.Parameter)
	for name, str := range secureStrings {
		secStrings[name] = &ssm.Parameter{
			Name:  aws.String(name),
			Type:  aws.String("Secure String"),
			Value: aws.String(str),
		}
	}

	for name, str := range insecureStrings {
		insecStrings[name] = &ssm.Parameter{
			Name:  aws.String(name),
			Type:  aws.String("String"),
			Value: aws.String(str),
		}
	}

	return &MockedSSMAPI{
		mSecureStrings: secStrings,
		mStrings:       insecStrings,
	}
}

// we only need GetParameter but ssmiface requires all of these
// Generated with this command: impl 'c *MockedSSMAPI' github.com/aws/aws-sdk-go/service/ssm/ssmiface.SSMAPI
func (c *MockedSSMAPI) GetParameters(i *ssm.GetParametersInput) (*ssm.GetParametersOutput, error) {
	params := []*ssm.Parameter{}
	if *i.WithDecryption {
		for _, n := range i.Names {
			param, ok := c.mSecureStrings[*n]
			if !ok {
				return nil, fmt.Errorf("No parameter %s", *n)
			}
			params = append(params, param)
		}
	} else {
		for _, n := range i.Names {
			param, ok := c.mSecureStrings[*n]
			if !ok {
				return nil, fmt.Errorf("No parameter %s", *n)
			}
			params = append(params, param)
		}
	}
	return &ssm.GetParametersOutput{Parameters: params}, nil
}
