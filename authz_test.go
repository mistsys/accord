package accord

import (
	"reflect"
	"testing"
)

func TestGrantAll_Authorized(t *testing.T) {
	type args struct {
		user       string
		principals []string
	}
	tests := []struct {
		name    string
		g       GrantAll
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Grants all principals requested",
			g:    GrantAll{},
			args: args{user: "test", principals: []string{"root-everywhere", "db-server"}},
			want: []string{"root-everywhere", "db-server"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := GrantAll{}
			got, err := g.Authorized(tt.args.user, tt.args.principals)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrantAll.Authorized() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GrantAll.Authorized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSimpleAuthFromFile(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    *SimpleAuth
		wantErr bool
	}{
		{
			name: "Initializes from the file",
			args: args{
				filePath: "test_assets/authz.json",
			},
			want: &SimpleAuth{
				Principals: []string{"root-everywhere", "zones-db"},
				AdminUsers: []string{"user1@ex.ample.com"},
				AccessMap: map[string][]string{
					"user1@ex.ample.com": []string{"root-everywhere"},
					"user2@ex.ample.com": []string{"zones-db"},
				},
			},
		},
		{
			name: "Invalid Files throw error",
			args: args{
				filePath: "test_assets/authz_file_does_not_exist.json",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSimpleAuthFromFile(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSimpleAuthFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSimpleAuthFromFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSimpleAuthFromBuffer(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *SimpleAuth
		wantErr bool
	}{
		{
			name: "simple auth initializes from valid json",
			args: args{
				buf: []byte(`{
"principals": ["root-everywhere", "zones-db"],
"admin_users": ["user1@ex.ample.com", "user2@ex.ample.com"],
"access_map": {
    "user1@ex.ample.com": ["root-everywhere"],
    "user2@ex.ample.com": ["zones-db"]
}

}`,
				),
			},
			want: &SimpleAuth{
				Principals: []string{"root-everywhere", "zones-db"},
				AdminUsers: []string{"user1@ex.ample.com", "user2@ex.ample.com"},
				AccessMap: map[string][]string{
					"user1@ex.ample.com": []string{"root-everywhere"},
					"user2@ex.ample.com": []string{"zones-db"},
				},
			},
		},
		{
			name: "simple auth fails gracefully on invalid json",
			args: args{
				buf: []byte(`{
"principals": ["root-everywhere", "zones-db"],
"admin_users": ["user1@ex.ample.com", "user2@ex.ample.com"],

}`,
				),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSimpleAuthFromBuffer(tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSimpleAuthFromBuffer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSimpleAuthFromBuffer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleAuth_IsAdmin(t *testing.T) {
	type fields struct {
		Principals []string
		AdminUsers []string
		AccessMap  map[string][]string
	}
	testAuthFields := fields{
		Principals: []string{"root-everywhere", "zones-db"},
		AdminUsers: []string{"user1@ex.ample.com"},
		AccessMap: map[string][]string{
			"user1@ex.ample.com": []string{"root-everywhere"},
			"user2@ex.ample.com": []string{"zones-db"},
		},
	}
	type args struct {
		user string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "admin check",
			fields: testAuthFields,
			args: args{
				user: "user1@ex.ample.com",
			},
			want: true,
		},
		{
			name:   "admin check",
			fields: testAuthFields,
			args: args{
				user: "user2@ex.ample.com",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SimpleAuth{
				Principals: tt.fields.Principals,
				AdminUsers: tt.fields.AdminUsers,
				AccessMap:  tt.fields.AccessMap,
			}
			if got := s.IsAdmin(tt.args.user); got != tt.want {
				t.Errorf("SimpleAuth.IsAdmin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleAuth_validPrincipals(t *testing.T) {
	type fields struct {
		Principals []string
		AdminUsers []string
		AccessMap  map[string][]string
	}
	testAuthFields := fields{
		Principals: []string{"root-everywhere", "zones-db"},
		AdminUsers: []string{"user1@ex.ample.com"},
		AccessMap: map[string][]string{
			"user1@ex.ample.com": []string{"root-everywhere"},
			"user2@ex.ample.com": []string{"zones-db"},
		},
	}
	type args struct {
		principals []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:   "valid principals are passed through",
			fields: testAuthFields,
			args: args{
				principals: []string{"root-everywhere"},
			},
			want: true,
		},
		{
			name:   "valid principals are passed through",
			fields: testAuthFields,
			args: args{
				principals: []string{"my-house"},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SimpleAuth{
				Principals: tt.fields.Principals,
				AdminUsers: tt.fields.AdminUsers,
				AccessMap:  tt.fields.AccessMap,
			}
			got, err := s.validPrincipals(tt.args.principals)
			if (err != nil) != tt.wantErr {
				t.Errorf("SimpleAuth.validPrincipals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SimpleAuth.validPrincipals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleAuth_Authorized(t *testing.T) {
	type fields struct {
		Principals []string
		AdminUsers []string
		AccessMap  map[string][]string
	}
	testAuthFields := fields{
		Principals: []string{"root-everywhere", "zones-db", "zones-willywonka"},
		AdminUsers: []string{"user1@ex.ample.com"},
		AccessMap: map[string][]string{
			"user1@ex.ample.com": []string{"root-everywhere"},
			"user2@ex.ample.com": []string{"zones-db"},
		},
	}
	type args struct {
		user       string
		principals []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:   "valid principals are granted access for admin",
			fields: testAuthFields,
			args: args{
				user:       "user1@ex.ample.com",
				principals: []string{"root-everywhere"},
			},
			want: []string{"root-everywhere"},
		},
		{
			name:   "valid principals are granted access for non-admin",
			fields: testAuthFields,
			args: args{
				user:       "user2@ex.ample.com",
				principals: []string{"zones-db"},
			},
			want: []string{"zones-db"},
		},
		{
			name:   "unknown users shouldn't be allowed access to any principals",
			fields: testAuthFields,
			args: args{
				user:       "user3@ex.ample.com",
				principals: []string{"zones-willywonka"},
			},
			wantErr: true,
		},
		{
			name:   "Invalid Principals throw error",
			fields: testAuthFields,
			args: args{
				user:       "user3@ex.ample.com",
				principals: []string{"zones-chocolatefactory"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SimpleAuth{
				Principals: tt.fields.Principals,
				AdminUsers: tt.fields.AdminUsers,
				AccessMap:  tt.fields.AccessMap,
			}
			got, err := s.Authorized(tt.args.user, tt.args.principals)
			if (err != nil) != tt.wantErr {
				t.Errorf("SimpleAuth.Authorized() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SimpleAuth.Authorized() = %v, want %v", got, tt.want)
			}
		})
	}
}
