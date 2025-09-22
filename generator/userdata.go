package generator

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/lefeck/ubuntu-autoinstaller/config"
	"github.com/lefeck/ubuntu-autoinstaller/utils"
)

// UserDataGenerator generates cloud-init user-data content.
type UserDataGenerator struct{}

// NewUserDataGenerator creates a new user-data generator.
func NewUserDataGenerator() *UserDataGenerator {
	return &UserDataGenerator{}
}

// GenerateFromConfig generates user-data from a config struct.
func (gen *UserDataGenerator) GenerateFromConfig(cfg *config.Config) ([]byte, error) {
	// Validate configuration
	if err := gen.validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %v", err)
	}

	// Ensure identity.password is hashed with SHA-512 crypt ($6$...)
	if cfg.Autoinstall.Identity.Password != "" && !utils.IsSHA512Crypt(cfg.Autoinstall.Identity.Password) {
		hashed, err := utils.HashSHA512Crypt(cfg.Autoinstall.Identity.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %v", err)
		}
		cfg.Autoinstall.Identity.Password = hashed
	}

	// Generate user-data YAML content
	userData, err := gen.generateUserData(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to generate user-data: %v", err)
	}

	return userData, nil
}

// GenerateFromTemplate generates user-data from a predefined template (not implemented yet).
func (gen *UserDataGenerator) GenerateFromTemplate(templateName string, data interface{}) ([]byte, error) {
	// 这里可以实现从预定义模板生成的功能
	// 暂时返回错误，表示功能未实现
	return nil, fmt.Errorf("模板生成功能暂未实现")
}

// validateConfig validates the provided config.
func (gen *UserDataGenerator) validateConfig(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("config must not be nil")
	}

	// Delegate validation to config struct
	return cfg.Validate()
}

// generateUserData marshals the config into YAML with #cloud-config header.
func (gen *UserDataGenerator) generateUserData(cfg *config.Config) ([]byte, error) {
	// Serialize config directly to avoid omitempty surprises via interface{}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal YAML: %v", err)
	}

	// Prepend required cloud-config header
	result := append([]byte("#cloud-config\n"), data...)
	return result, nil
}

// ValidateUserData checks that user-data YAML is syntactically valid and contains required fields.
func (gen *UserDataGenerator) ValidateUserData(userData []byte) error {
	if len(userData) == 0 {
		return fmt.Errorf("user-data must not be empty")
	}

	// Parse YAML to validate syntax
	var config map[string]interface{}
	if err := yaml.Unmarshal(userData, &config); err != nil {
		return fmt.Errorf("invalid YAML syntax: %v", err)
	}

	// Ensure required 'autoinstall' field exists
	if _, exists := config["autoinstall"]; !exists {
		return fmt.Errorf("missing required 'autoinstall' field")
	}

	return nil
}

// GenerateDefaultConfig returns a default config encoded as user-data.
func (gen *UserDataGenerator) GenerateDefaultConfig() ([]byte, error) {
	defaultConfig := config.NewDefaultConfig()
	return gen.GenerateFromConfig(defaultConfig)
}

// LoadConfigFromYAML parses a Config from YAML bytes.
func (gen *UserDataGenerator) LoadConfigFromYAML(yamlData []byte) (*config.Config, error) {
	var cfg config.Config
	if err := yaml.Unmarshal(yamlData, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %v", err)
	}
	return &cfg, nil
}

// SaveConfigToYAML serializes a Config to YAML bytes.
func (gen *UserDataGenerator) SaveConfigToYAML(cfg *config.Config) ([]byte, error) {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %v", err)
	}
	return data, nil
}
