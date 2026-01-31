package api

import (
	"context"
	"time"

	"accio/internal/models"
)

type Client interface {
	GetStacks(ctx context.Context) ([]models.Stack, error)
	SubmitIntent(ctx context.Context, intent models.InfraIntent) error
	GetStatus(ctx context.Context) (map[string]string, error)
}

type MockClient struct{}

func NewMockClient() Client {
	return &MockClient{}
}

func (m *MockClient) GetStacks(ctx context.Context) ([]models.Stack, error) {
	// Simulate latency
	time.Sleep(500 * time.Millisecond)

	return []models.Stack{
		{
			Name:   "prod-db",
			Cloud:  "aws",
			Region: "us-east-1",
			Status: models.StatusReady,
			Resources: []models.Resource{
				{Name: "rds-instance", Kind: "RDSInstance", Status: "Available"},
			},
		},
		{
			Name:   "staging-cluster",
			Cloud:  "gcp",
			Region: "us-central1",
			Status: models.StatusReconciling,
			Resources: []models.Resource{
				{Name: "gke-cluster", Kind: "GKECluster", Status: "Provisioning"},
			},
		},
	}, nil
}

func (m *MockClient) SubmitIntent(ctx context.Context, intent models.InfraIntent) error {
	time.Sleep(1 * time.Second)
	return nil
}

func (m *MockClient) GetStatus(ctx context.Context) (map[string]string, error) {
	return map[string]string{
		"Platform": "Healthy",
		"API":      "Connected",
		"Auth":     "Authenticated",
	}, nil
}
