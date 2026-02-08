package localconfig

import "backend/port"

// ConfigService wraps LocalConfig and provides convenient access methods.
type ConfigService struct {
	config LocalConfig
	logger domain.Logger
}

// NewConfigService creates a new config service by loading configuration.
func NewConfigService(log domain.Logger) (*ConfigService, error) {
	config, err := GetConfig(log)
	if err != nil {
		return nil, err
	}

	return &ConfigService{
		config: config,
		logger: log,
	}, nil
}

// Get returns the entire configuration.
func (s *ConfigService) Get() LocalConfig {
	return s.config
}

// GetConnectionString returns the database connection string.
func (s *ConfigService) GetConnectionString() string {
	return s.config.Database.URL
}

// GetServicePort returns the HTTP server port.
func (s *ConfigService) GetServicePort() int {
	return s.config.Service.Port()
}

// GetServiceName returns the service name.
func (s *ConfigService) GetServiceName() string {
	return s.config.Service.Name
}

// GetDocsPath returns the docs path.
func (s *ConfigService) GetDocsPath() string {
	return s.config.Service.DocsPath
}

// GetMigrationsPath returns the migrations path.
func (s *ConfigService) GetMigrationsPath() string {
	return s.config.Service.MigrationsPath
}

// GetSkipMigrations returns whether to skip migrations.
func (s *ConfigService) GetSkipMigrations() bool {
	return s.config.Service.SkipMigrations
}
