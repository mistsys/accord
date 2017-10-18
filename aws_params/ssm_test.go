package aws_params

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

func Test_client_GetSecureString(t *testing.T) {
	type fields struct {
		cfg *Config
		ssm ssmiface.SSMAPI
	}
	type args struct {
		path string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "available value is returned",
			fields: fields{
				cfg: nil,
				ssm: NewMockedSSMAPI(map[string]string{
					"/params/0": "Much secure, very strong",
				}, nil),
			},
			args: args{path: "/params/0"},
			want: "Much secure, very strong",
		},
		{
			name: "unavailable params return errors",
			fields: fields{
				cfg: nil,
				ssm: NewMockedSSMAPI(nil, nil),
			},
			args:    args{path: "/params/0"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &client{
				cfg: tt.fields.cfg,
				ssm: tt.fields.ssm,
			}
			got, err := c.GetSecureString(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.GetSecureString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("client.GetSecureString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	type args struct {
		cfg *Config
	}
	tests := []struct {
		name    string
		args    args
		want    Client
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
