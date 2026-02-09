package api

import "net/http"

// GetEncryptionMeta fetches the user's encryption metadata.
func (c *Client) GetEncryptionMeta() (*EncryptionMeta, error) {
	var result EncryptionMeta
	if err := c.doAuthenticatedRequest(http.MethodGet, PathEncryption, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateEncryptionMeta updates the user's encryption metadata (salt, verifier, public_key).
func (c *Client) UpdateEncryptionMeta(data map[string]string) error {
	return c.doAuthenticatedRequest(http.MethodPatch, PathEncryption, data, nil)
}
