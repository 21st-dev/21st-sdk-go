package agents

import "context"

// TokensResource manages short-lived authentication tokens.
type TokensResource struct {
	client *AgentClient
}

func (t *TokensResource) Create(ctx context.Context, req *CreateTokenRequest) (*Token, error) {
	if req == nil {
		req = &CreateTokenRequest{}
	}

	expiresIn := req.ExpiresIn
	if expiresIn == "" {
		expiresIn = "1h"
	}

	var agents []string
	if req.Agent != "" {
		agents = []string{req.Agent}
	}

	body := struct {
		Agents    []string `json:"agents,omitempty"`
		UserID    string   `json:"userId,omitempty"`
		ExpiresIn string   `json:"expiresIn"`
	}{
		Agents:    agents,
		UserID:    req.UserID,
		ExpiresIn: expiresIn,
	}

	var out Token
	err := t.client.do(ctx, "POST", "/v1/tokens", body, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
