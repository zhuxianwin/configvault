package ssm_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsssm "github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"

	"configvault/internal/provider/ssm"
)

type mockSSMClient struct {
	params []types.Parameter
	err    error
}

func (m *mockSSMClient) GetParametersByPath(_ context.Context, _ *awsssm.GetParametersByPathInput, _ ...func(*awsssm.Options)) (*awsssm.GetParametersByPathOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &awsssm.GetParametersByPathOutput{Parameters: m.params}, nil
}

func TestProvider_GetSecrets_ReturnsKeyValues(t *testing.T) {
	client := &mockSSMClient{
		params: []types.Parameter{
			{Name: aws.String("/myapp/db-host"), Value: aws.String("localhost"), Type: types.ParameterTypeString},
			{Name: aws.String("/myapp/db-pass"), Value: aws.String("s3cr3t"), Type: types.ParameterTypeSecureString},
		},
	}
	p := ssm.New(client, true)

	secrets, err := p.GetSecrets(context.Background(), "/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", secrets["DB_HOST"])
	}
	if secrets["DB_PASS"] != "s3cr3t" {
		t.Errorf("expected DB_PASS=s3cr3t, got %q", secrets["DB_PASS"])
	}
}

func TestProvider_GetSecrets_SkipsSecureStringWhenNoDecrypt(t *testing.T) {
	client := &mockSSMClient{
		params: []types.Parameter{
			{Name: aws.String("/myapp/api-key"), Value: aws.String("hidden"), Type: types.ParameterTypeSecureString},
		},
	}
	p := ssm.New(client, false)

	secrets, err := p.GetSecrets(context.Background(), "myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := secrets["API_KEY"]; ok {
		t.Error("expected API_KEY to be skipped when decrypt=false")
	}
}

func TestProvider_Name(t *testing.T) {
	p := ssm.New(nil, false)
	if p.Name() != "aws-ssm" {
		t.Errorf("expected name aws-ssm, got %q", p.Name())
	}
}
