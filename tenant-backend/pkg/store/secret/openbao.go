package secret

import (
	"fmt"

	openbao "github.com/openbao/openbao/api/v2"
)

type Client struct {
	client *openbao.Client
}

func NewClient(addr string) (*Client, error) {
	cfg := openbao.DefaultConfig()
	cfg.Address = addr

	client, err := openbao.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create openbao client: %w", err)
	}

	return &Client{client: client}, nil
}

func (c *Client) LoginAppRole(roleID, secretID string) error {
	data := map[string]interface{}{
		"role_id":   roleID,
		"secret_id": secretID,
	}

	resp, err := c.client.Logical().Write("auth/approle/login", data)
	if err != nil {
		return fmt.Errorf("approle login failed: %w", err)
	}
	if resp.Auth == nil {
		return fmt.Errorf("approle login returned no auth info")
	}

	c.client.SetToken(resp.Auth.ClientToken)
	return nil
}

// ReadSecret reads a KV-v2 secret. The path should NOT include "secret/data/" prefix —
// it will be constructed as "secret/data/{path}".
func (c *Client) ReadSecret(path string) (map[string]interface{}, error) {
	fullPath := fmt.Sprintf("secret/data/%s", path)

	secret, err := c.client.Logical().Read(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret at %s: %w", fullPath, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no secret found at %s", fullPath)
	}

	// KV-v2 wraps the actual data under a "data" key
	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected secret data format at %s", fullPath)
	}

	return data, nil
}

// TenantSecrets holds the parsed secrets for a tenant.
type TenantSecrets struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string

	RedisHost     string
	RedisPort     string
	RedisPassword string
}

// ReadTenantSecrets reads and parses secrets for a specific tenant.
func (c *Client) ReadTenantSecrets(env, tenantID string) (*TenantSecrets, error) {
	path := fmt.Sprintf("%s/dp/%s", env, tenantID)

	data, err := c.ReadSecret(path)
	if err != nil {
		return nil, err
	}

	getString := func(key string) (string, error) {
		val, ok := data[key].(string)
		if !ok || val == "" {
			return "", fmt.Errorf("missing or empty secret key: %s", key)
		}
		return val, nil
	}

	secrets := &TenantSecrets{}

	if secrets.DBHost, err = getString("db_host"); err != nil {
		return nil, err
	}
	if secrets.DBPort, err = getString("db_port"); err != nil {
		return nil, err
	}
	if secrets.DBUser, err = getString("db_user"); err != nil {
		return nil, err
	}
	if secrets.DBPassword, err = getString("db_password"); err != nil {
		return nil, err
	}
	if secrets.DBName, err = getString("db_name"); err != nil {
		return nil, err
	}
	if secrets.S3Endpoint, err = getString("s3_endpoint"); err != nil {
		return nil, err
	}
	if secrets.S3AccessKey, err = getString("s3_access_key"); err != nil {
		return nil, err
	}
	if secrets.S3SecretKey, err = getString("s3_secret_key"); err != nil {
		return nil, err
	}
	if secrets.S3Bucket, err = getString("s3_bucket"); err != nil {
		return nil, err
	}
	if secrets.RedisHost, err = getString("redis_host"); err != nil {
		return nil, err
	}
	if secrets.RedisPort, err = getString("redis_port"); err != nil {
		return nil, err
	}
	if secrets.RedisPassword, err = getString("redis_password"); err != nil {
		return nil, err
	}

	return secrets, nil
}
