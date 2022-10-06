package githubjwt

import (
	"context"
	"encoding/json"
	"time"

	"github.com/lestrrat-go/jwx/jwt"
	"github.com/pkg/errors"
)

// GithubJWT are the fields defined in the Github JWT Token
type GithubJWT struct {
	Actor                string    `json:"actor"`
	ActorID              string    `json:"actor_id"`
	Aud                  []string  `json:"aud"`
	BaseRef              string    `json:"base_ref"`
	Enterprise           string    `json:"enterprise"`
	EventName            string    `json:"event_name"`
	Exp                  time.Time `json:"exp"`
	HeadRef              string    `json:"head_ref"`
	Iat                  time.Time `json:"iat"`
	Iss                  string    `json:"iss"`
	JobWorkflowRef       string    `json:"job_workflow_ref"`
	Jti                  string    `json:"jti"`
	Nbf                  time.Time `json:"nbf"`
	Ref                  string    `json:"ref"`
	RefType              string    `json:"ref_type"`
	Repository           string    `json:"repository"`
	RepositoryID         string    `json:"repository_id"`
	RepositoryOwner      string    `json:"repository_owner"`
	RepositoryOwnerID    string    `json:"repository_owner_id"`
	RepositoryVisibility string    `json:"repository_visibility"`
	RunAttempt           string    `json:"run_attempt"`
	RunID                string    `json:"run_id"`
	RunNumber            string    `json:"run_number"`
	Sha                  string    `json:"sha"`
	Sub                  string    `json:"sub"`
	Workflow             string    `json:"workflow"`
}

// TransformTokenToGithubClaims takes a jwt.Token and transforms it to the structured claims provided
func TransformTokenToGithubClaims(token jwt.Token) (*GithubJWT, error) {
	asMap, err := token.AsMap(context.TODO())
	if err != nil {
		return nil, err
	}

	jsonToken, err := json.Marshal(asMap)
	if err != nil {
		return nil, errors.Wrap(err, "TransformTokenToClaims():json.Marshal(token)")
	}

	githubToken := &GithubJWT{}

	err = json.Unmarshal(jsonToken, githubToken)
	if err != nil {
		return nil, errors.Wrap(err, "TransformTokenToClaims():json.Unmarshal(jsonToken, claims)")
	}

	return githubToken, nil
}
