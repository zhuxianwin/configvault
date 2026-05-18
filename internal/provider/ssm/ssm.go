package ssm

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

// ssmClient is a subset of the SSM API used by Provider.
type ssmClient interface {
	GetParametersByPath(ctx context.Context, params *ssm.GetParametersByPathInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersByPathOutput, error)
}

// Provider retrieves secrets from AWS SSM Parameter Store.
type Provider struct {
	client  ssmClient
	decrypt bool
}

// New creates a new SSM Provider.
func New(client ssmClient, decrypt bool) *Provider {
	return &Provider{client: client, decrypt: decrypt}
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return "aws-ssm"
}

// normalizePath ensures the path starts with a leading slash.
func normalizePath(path string) string {
	if !strings.HasPrefix(path, "/") {
		return "/" + path
	}
	return path
}

// paramNameToKey converts an SSM parameter name to an environment-variable-style key
// by stripping the path prefix, uppercasing, and replacing hyphens with underscores.
func paramNameToKey(name, pathPrefix string) string {
	key := strings.TrimPrefix(name, pathPrefix+"/")
	return strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
}

// GetSecrets fetches all parameters under the given path prefix.
// Parameters of type SecureString are skipped when decryption is disabled.
func (p *Provider) GetSecrets(ctx context.Context, path string) (map[string]string, error) {
	path = normalizePath(path)

	secrets := make(map[string]string)
	var nextToken *string

	for {
		out, err := p.client.GetParametersByPath(ctx, &ssm.GetParametersByPathInput{
			Path:           aws.String(path),
			WithDecryption: aws.Bool(p.decrypt),
			Recursive:      aws.Bool(false),
			NextToken:      nextToken,
		})
		if err != nil {
			return nil, fmt.Errorf("ssm: GetParametersByPath: %w", err)
		}

		for _, param := range out.Parameters {
			if param.Type == types.ParameterTypeSecureString && !p.decrypt {
				continue
			}
			key := paramNameToKey(aws.ToString(param.Name), path)
			secrets[key] = aws.ToString(param.Value)
		}

		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}

	return secrets, nil
}
