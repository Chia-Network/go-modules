package githubjwt

import (
	"net/http"

	"github.com/chia-network/go-modules/pkg/jwt"
)

// ParseTokenFromRequestWithJWKS Wrapper around the JWT version of this function, that is specific to Github tokens
// Sets the proper JWKS endpoint, and returns the github specific token data after validation
func ParseTokenFromRequestWithJWKS(r *http.Request) (*GithubJWT, error) {
	token, err := jwt.ParseTokenFromRequestWithJWKS(r, "https://token.actions.githubusercontent.com/.well-known/jwks")
	if err != nil {
		return nil, err
	}

	return TransformTokenToGithubClaims(token)
}
