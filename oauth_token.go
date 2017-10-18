package accord

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/mistsys/accord/protocol"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// This package is used for oauth2 3-legged authentication and conversion work
// Convert Oauth token to protobuf version
func OAuthTokenPb(token *oauth2.Token) (*protocol.OauthToken, error) {
	expiryTime, err := ptypes.TimestampProto(token.Expiry)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid timestamp")
	}
	// TODO: do validation with crappy inputs
	pbToken := &protocol.OauthToken{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       expiryTime,
	}
	return pbToken, err
}

// Convert from Protobuf Oauth2 Token to *oauth2.Token
func OAuth2Token(tokenPb *protocol.OauthToken) (*oauth2.Token, error) {
	expiryTime, err := ptypes.Timestamp(tokenPb.Expiry)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid timestamp")
	}
	token := &oauth2.Token{
		AccessToken:  tokenPb.AccessToken,
		TokenType:    tokenPb.TokenType,
		RefreshToken: tokenPb.RefreshToken,
		Expiry:       expiryTime,
	}
	return token, err
}
