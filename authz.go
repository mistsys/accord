package accord

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
)

// Authz is a bare minimum Authorization interface that just checks
// if any of the requested principals are actually valid for the user
// The service should only grant the access to the principals the user
// should have access to. If there is no overlap with the existing principals
// the result is empty and with an error for more details for the service to
// log for administrators or to send to the users
// Any other authorization backends can be added by implementing this interface
type Authz interface {
	Authorized(user string, principals []string) ([]string, error)
}

// GrantAll is the authz module that grants everyone everything
type GrantAll struct{}

func (g GrantAll) Authorized(user string, principals []string) ([]string, error) {
	return principals, nil
}

type SimpleAuth struct {
	Principals []string `json:"principals",yaml:"principals"`
	// these users can get the root-everywhere if they request for it
	// allowing to sign in with full sudo
	AdminUsers []string `json:"admin_users",yaml:"admin_users"`

	// Users -> Principals map
	AccessMap map[string][]string `json:"access_map",yaml:"access_map"`
}

func NewSimpleAuthFromFile(filePath string) (*SimpleAuth, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot read file %s", filePath)
	}
	return NewSimpleAuthFromBuffer(content)
}

func NewSimpleAuthFromBuffer(content []byte) (*SimpleAuth, error) {
	s := &SimpleAuth{}
	err := json.Unmarshal(content, s)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse json for simple auth")
	}
	return s, nil
}

func (s SimpleAuth) IsAdmin(user string) bool {
	for _, adminUser := range s.AdminUsers {
		if user == adminUser {
			return true
		}
	}
	return false
}

func (s SimpleAuth) validPrincipals(principals []string) (bool, error) {
	// Yes this is O(n^2) but the data size is going to be very small
	// at most 100 elements, or at least if this is becoming a bottleneck
	// you want something like LDAP and not hacking this simple implementation
	for _, p := range principals {
		if !contains(p, s.Principals) {
			return false, fmt.Errorf("Principal %s is unknown", p)
		}
	}
	return true, nil
}

func contains(needle string, haystack []string) bool {
	for _, thread := range haystack {
		if needle == thread {
			return true
		}
	}
	return false
}

func (s SimpleAuth) Authorized(user string, principals []string) ([]string, error) {
	if _, err := s.validPrincipals(principals); err != nil {
		return nil, errors.Wrapf(err, "Invalid principals")
	}

	if s.IsAdmin(user) {
		// if admin, grant the principals if they exist
		return principals, nil
	}

	grantedAccess, ok := s.AccessMap[user]
	if !ok {
		return nil, fmt.Errorf("User not granted any access yet, talk to your administrator")
	}

	// grant user their own principal
	grantedPrincipals := []string{}

	for _, p := range principals {
		if contains(p, grantedAccess) {
			grantedPrincipals = append(grantedPrincipals, p)
		}
	}
	return grantedPrincipals, nil
}
