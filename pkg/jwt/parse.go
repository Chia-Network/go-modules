package jwt

import (
	"net/http"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/pkg/errors"
)

// rawJWTFromAuthHeader returns a raw jwt token from the auth header
func rawJWTFromAuthHeader(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return ""
	}

	auth = strings.TrimPrefix(auth, "bearer ")
	return strings.TrimPrefix(auth, "Bearer ")
}

// ParseTokenFromRequestWithJWKS Parses a jwt.Token from an Auth header, validating with JWKS Keys from an endpoint
func ParseTokenFromRequestWithJWKS(r *http.Request, jwksEndpoint string) (jwt.Token, error) {
	keyset, err := jwk.Fetch(r.Context(), jwksEndpoint)
	if err != nil {
		return nil, errors.Wrap(err, "ParseTokenFromHeaderWithJWKS():jwk.Fetch()")
	}

	// Disable "infer algorithm" so this always uses the algorithm specified at the jwks endpoint or fails
	// (avoid something along the lines of https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries/)
	token, err := jwt.ParseString(rawJWTFromAuthHeader(r), jwt.WithKeySet(keyset), jwt.InferAlgorithmFromKey(false))
	if err != nil {
		return nil, errors.Wrap(err, "ParseTokenFromHeaderWithJWKS():jwt.ParseString()")
	}

	// ensures that we verify essential claims, like expiration, not before, etc
	err = jwt.Validate(token, jwt.WithAcceptableSkew(60*time.Second))
	if err != nil {
		return nil, errors.Wrap(err, "ParseTokenFromHeaderWithJWKS():jwt.Verify()")
	}

	return token, nil
}
