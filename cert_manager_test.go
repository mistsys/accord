package accord

import (
	"reflect"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

func TestCertSignRequest_valid(t *testing.T) {
	type fields struct {
		PubKey      []byte
		ValidFrom   time.Time
		ValidUntil  time.Time
		Id          string
		Serial      uint64
		Principals  []string
		Permissions ssh.Permissions
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
		err     error
	}{
		{
			name: "check invalid serial",
			fields: fields{
				ValidFrom:  time.Now().Add(10 * time.Second),
				ValidUntil: time.Now().Add(10 * time.Hour),
				Id:         "test",
				Serial:     0,
			},
			wantErr: true,
			err:     ErrInvalidSerial,
		},
		{
			name: "valid time before now",
			fields: fields{
				ValidFrom:  time.Now().Add(-10 * time.Second),
				ValidUntil: time.Now().Add(10 * time.Hour),
				Id:         "test",
				Serial:     1,
			},
			wantErr: true,
			err:     ErrInvalidStartTime,
		},
		{
			name: "valid end time before start time",
			fields: fields{
				ValidFrom:  time.Now().Add(10 * time.Hour),
				ValidUntil: time.Now().Add(10 * time.Second),
				Id:         "test",
				Serial:     2,
			},
			wantErr: true,
			err:     ErrEndBeforeStartTime,
		},
		{
			name: "longer than 90 days cert requests are rejected",
			fields: fields{
				ValidFrom:  time.Now().Add(10 * time.Second),
				ValidUntil: time.Now().Add(100 * 24 * time.Hour),
				Id:         "test",
				Serial:     2,
			},
			wantErr: true,
			err:     ErrValidityTooLong,
		},
		{
			name: "Empty Ids are rejected",
			fields: fields{
				ValidFrom:  time.Now().Add(10 * time.Second),
				ValidUntil: time.Now().Add(10 * 24 * time.Hour),
				Id:         "",
				Serial:     2,
			},
			wantErr: true,
			err:     ErrEmptyID,
		},
		{
			name: "all fields set correctly",
			fields: fields{
				PubKey:     []byte(`ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCiIlKicXvA58W1EdNKpprPlvfeSsPUQV7hvp5NZLzBOOi6PEp5tUfQgwGNuO7TezeHU8SZnR48pAqHDck5mbYdBo6jL/iLeSEt3Ix3ejqjm6FVBAaXj1v6R0Oys+oSS+pLoTvyO9/8v2C7pVgLD/3rq04l69u7RDpXQSxWYQHXqbRkM1TD1LEK65aOL2+WpZJPqYjZ3WqRzTQS5C8TCLQyVhomxSPCV89Ilnh+plU7Woz2YP5JPCIAtxfMYbio1CqolV4RLuqjTWwg+1SKhoDAn8aVsbDOPJqiSQLFc4ukiCwiXSwhP++vhFSk4TO59yQOoD0AS16RHUGEUruxydYOgm8tk5yLpsKezSyxj8Vv7TANUX1BhZlXbHsc5Xc1X/7hHgwvkvAb0vOwDHET4GOyAwzquWLT6tv1ZnH7eSwnIFz5Y+Uji/3UFbgqiUl+gWOYH3ogHazSgYrrDsxrfebnPLoXq8LdHZcJrJdAdeptMIzDMcLHTsKLu+yW+Y+oe59tZadWE2fTfdjdeI5oRXpm2knsgR5I2gHFCwLPnTD9e4LkNKGCmByiJWO3GGBzH2h+ReCRRGvw7noKRhhvLt/g+Dpeos7swhoiMFUqsAQGa2jkKSRVq2QFa8RybbO3HS8kTlLu0U7CnOq62OAMBv+OR0UAZw8SQ3SUawEmSwkYbw== pgautam@Wintermute-3.local`),
				ValidFrom:  time.Now().Add(10 * time.Second),
				ValidUntil: time.Now().Add(10 * 24 * time.Hour),
				Id:         "testacct",
				Principals: []string{"testacct"},
				Serial:     1,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CertSignRequest{
				PubKey:      tt.fields.PubKey,
				ValidFrom:   tt.fields.ValidFrom,
				ValidUntil:  tt.fields.ValidUntil,
				Id:          tt.fields.Id,
				Serial:      tt.fields.Serial,
				Principals:  tt.fields.Principals,
				Permissions: tt.fields.Permissions,
			}
			got, err := r.valid()
			if (err != nil) != tt.wantErr {
				t.Errorf("CertSignRequest.valid() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else {
				if err != tt.err {
					t.Errorf("CertSignRequest.valid() invalid error = %v, wantErr %v", err, tt.err)
					return
				}
			}
			if got != tt.want {
				t.Errorf("CertSignRequest.valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCertManager_SignUserCert(t *testing.T) {
	type fields struct {
		rootCAPath     string
		rootCAPassword string
		userCAPath     string
		userCAPassword string
	}
	type args struct {
		request *CertSignRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &CertManager{
				rootCAPath:     tt.fields.rootCAPath,
				rootCAPassword: tt.fields.rootCAPassword,
				userCAPath:     tt.fields.userCAPath,
				userCAPassword: tt.fields.userCAPassword,
			}
			got, err := m.SignUserCert(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("CertManager.SignUserCert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CertManager.SignUserCert() = %v, want %v", got, tt.want)
			}
		})
	}
}
