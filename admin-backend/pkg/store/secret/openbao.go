package secret

import (
	"fmt"

	openbao "github.com/openbao/openbao/api/v2"
)

type Client struct {
	client *openbao.Client
}

// NewClient creates a new OpenBao client authenticated with the root token.
func NewClient(addr, token string) (*Client, error) {
	cfg := openbao.DefaultConfig()
	cfg.Address = addr

	client, err := openbao.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create openbao client: %w", err)
	}

	client.SetToken(token)

	return &Client{client: client}, nil
}

// WriteSecret writes a KV-v2 secret at the given path.
// The path should NOT include "secret/data/" prefix.
func (c *Client) WriteSecret(path string, data map[string]interface{}) error {
	fullPath := fmt.Sprintf("secret/data/%s", path)

	payload := map[string]interface{}{
		"data": data,
	}

	_, err := c.client.Logical().Write(fullPath, payload)
	if err != nil {
		return fmt.Errorf("failed to write secret at %s: %w", fullPath, err)
	}

	return nil
}

// CreatePolicy creates or updates a policy with the given name and HCL rules.
func (c *Client) CreatePolicy(name, rules string) error {
	err := c.client.Sys().PutPolicy(name, rules)
	if err != nil {
		return fmt.Errorf("failed to create policy %s: %w", name, err)
	}
	return nil
}

// EnableAppRoleAuth enables the AppRole auth method if not already enabled.
func (c *Client) EnableAppRoleAuth() error {
	auths, err := c.client.Sys().ListAuth()
	if err != nil {
		return fmt.Errorf("failed to list auth methods: %w", err)
	}

	if _, exists := auths["approle/"]; exists {
		return nil // already enabled
	}

	err = c.client.Sys().EnableAuthWithOptions("approle", &openbao.EnableAuthOptions{
		Type: "approle",
	})
	if err != nil {
		return fmt.Errorf("failed to enable approle auth: %w", err)
	}

	return nil
}

// CreateAppRole creates an AppRole with the specified policies.
func (c *Client) CreateAppRole(roleName string, policies []string) error {
	data := map[string]interface{}{
		"token_policies": policies,
		"token_ttl":      "1h",
		"token_max_ttl":  "4h",
	}

	path := fmt.Sprintf("auth/approle/role/%s", roleName)
	_, err := c.client.Logical().Write(path, data)
	if err != nil {
		return fmt.Errorf("failed to create approle %s: %w", roleName, err)
	}

	return nil
}

// GetRoleID retrieves the role_id for the named AppRole.
func (c *Client) GetRoleID(roleName string) (string, error) {
	path := fmt.Sprintf("auth/approle/role/%s/role-id", roleName)

	secret, err := c.client.Logical().Read(path)
	if err != nil {
		return "", fmt.Errorf("failed to read role-id for %s: %w", roleName, err)
	}
	if secret == nil || secret.Data == nil {
		return "", fmt.Errorf("no role-id found for %s", roleName)
	}

	roleID, ok := secret.Data["role_id"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected role-id format for %s", roleName)
	}

	return roleID, nil
}

// GenerateSecretID generates a new secret_id for the named AppRole.
func (c *Client) GenerateSecretID(roleName string) (string, error) {
	path := fmt.Sprintf("auth/approle/role/%s/secret-id", roleName)

	secret, err := c.client.Logical().Write(path, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate secret-id for %s: %w", roleName, err)
	}
	if secret == nil || secret.Data == nil {
		return "", fmt.Errorf("no secret-id returned for %s", roleName)
	}

	secretID, ok := secret.Data["secret_id"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected secret-id format for %s", roleName)
	}

	return secretID, nil
}
