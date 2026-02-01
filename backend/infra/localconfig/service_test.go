package localconfig

import (
	"testing"

	"backend/infra/logger"
	"github.com/stretchr/testify/assert"
)

func TestConfigService_NewConfigService(t *testing.T) {
	log := logger.NewNoop()

	// Test that NewConfigService doesn't panic
	assert.NotPanics(t, func() {
		_, err := NewConfigService(log)
		// This will likely fail due to missing env, but shouldn't panic
		_ = err
	})
}

func TestConfigService_Get(t *testing.T) {
	// Test with minimal config
	config := LocalConfig{
		Service: Service{
			Port: 8080,
			Name: "test-service",
		},
		Database: Database{
			URL: "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable",
		},
	}

	service := &ConfigService{
		config: config,
		logger: logger.NewNoop(),
	}

	result := service.Get()
	assert.Equal(t, config, result, "Get should return the stored config")
}

func TestConfigService_GetConnectionString(t *testing.T) {
	tests := []struct {
		name     string
		config   Database
		expected string
	}{
		{
			name: "basic connection string",
			config: Database{
				URL: "postgres://user:pass@localhost:5432/dbname?sslmode=disable",
			},
			expected: "postgres://user:pass@localhost:5432/dbname?sslmode=disable",
		},
		{
			name: "connection with SSL require",
			config: Database{
				URL: "postgres://myuser:mypass@db.example.com:5432/mydb?sslmode=require",
			},
			expected: "postgres://myuser:mypass@db.example.com:5432/mydb?sslmode=require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &ConfigService{
				config: LocalConfig{Database: tt.config},
				logger: logger.NewNoop(),
			}

			result := service.GetConnectionString()
			assert.Equal(t, tt.expected, result, "GetConnectionString should return correct connection string")
		})
	}
}

func TestConfigService_GetServicePort(t *testing.T) {
	config := LocalConfig{
		Service: Service{
			Port: 9090,
			Name: "test-service",
		},
	}

	service := &ConfigService{
		config: config,
		logger: logger.NewNoop(),
	}

	result := service.GetServicePort()
	assert.Equal(t, 9090, result, "GetServicePort should return the correct port")
}

func TestConfigService_GetServiceName(t *testing.T) {
	config := LocalConfig{
		Service: Service{
			Port: 8080,
			Name: "my-api",
		},
	}

	service := &ConfigService{
		config: config,
		logger: logger.NewNoop(),
	}

	result := service.GetServiceName()
	assert.Equal(t, "my-api", result, "GetServiceName should return the correct name")
}
