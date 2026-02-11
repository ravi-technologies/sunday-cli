package api

import "net/http"

// ListIdentities returns all identities for the authenticated user.
func (c *Client) ListIdentities() ([]Identity, error) {
	var identities []Identity
	if err := c.doAuthenticatedRequest(http.MethodGet, PathIdentities, nil, &identities); err != nil {
		return nil, err
	}
	return identities, nil
}

// BindIdentity exchanges an unbound token pair for one with the given
// identity UUID baked into the JWT claims.
func (c *Client) BindIdentity(identityUUID string) (*BindIdentityResponse, error) {
	req := BindIdentityRequest{
		Identity: identityUUID,
	}
	var resp BindIdentityResponse
	if err := c.doAuthenticatedRequest(http.MethodPost, PathBindIdentity, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
