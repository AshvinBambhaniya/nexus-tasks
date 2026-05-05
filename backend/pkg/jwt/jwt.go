package jwt

import (
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

var (
	issuer = "nexus-tasks.example.com"
)

/**
JWT token creation guideline
---------------------------------
alg: HS256
payload: {
	subject: uid
	issuer: nexus-tasks.example.com
}
We are also using jti for implementation signout all device
We are using symmetric key methodology for sign token
---------------------------------
*/

// ParseToken parse, validate the jwt token
// On valid token it returns the decoded token
func ParseToken(secret string, token string) (jwt.Token, error) {
	key, err := jwk.FromRaw([]byte(secret))
	if err != nil {
		return nil, err
	}

	claims, err := jwt.Parse([]byte(token), jwt.WithKey(jwa.HS256, key), jwt.WithIssuer(issuer))
	return claims, err
}

func CreateToken(secret, sub string, exp time.Time) (string, error) {
	stringToken := ""
	token, err := jwt.NewBuilder().Subject(sub).Expiration(exp).Issuer(issuer).Build()
	if err != nil {
		return stringToken, err
	}
	key, err := jwk.FromRaw([]byte(secret))
	if err != nil {
		return stringToken, err
	}
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, key))
	if err != nil {
		return stringToken, err
	}
	stringToken = string(signed)
	return stringToken, nil
}
