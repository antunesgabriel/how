package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Provider string

const (
	ProviderOpenAI      Provider = "openai"
	ProviderGemini      Provider = "gemini"
	ProviderAnthropic   Provider = "anthropic"
	ProviderDeepseek    Provider = "deepseek"
	ProviderOllama      Provider = "ollama"
	ProviderOpenRouter  Provider = "openrouter"
	ProviderAzureOpenAI Provider = "azure_openai"
)

type OpenAIConfig struct {
	APIKey       string `yaml:"api_key"`
	BaseURL      string `yaml:"base_url,omitempty"`
	DefaultModel string `yaml:"default_model"`
}

type GeminiConfig struct {
	APIKey       string `yaml:"api_key"`
	BaseURL      string `yaml:"base_url,omitempty"`
	DefaultModel string `yaml:"default_model"`
}

type AnthropicConfig struct {
	APIKey       string `yaml:"api_key"`
	BaseURL      string `yaml:"base_url,omitempty"`
	DefaultModel string `yaml:"default_model"`
}

type DeepseekConfig struct {
	APIKey       string `yaml:"api_key"`
	BaseURL      string `yaml:"base_url,omitempty"`
	DefaultModel string `yaml:"default_model"`
}

type OllamaConfig struct {
	BaseURL      string `yaml:"base_url"`
	DefaultModel string `yaml:"default_model"`
}

type OpenRouterConfig struct {
	APIKey       string `yaml:"api_key"`
	BaseURL      string `yaml:"base_url,omitempty"`
	DefaultModel string `yaml:"default_model"`
}

type AzureOpenAIConfig struct {
	APIKey       string `yaml:"api_key"`
	Endpoint     string `yaml:"endpoint"`
	DefaultModel string `yaml:"default_model"`
	APIVersion   string `yaml:"api_version"`
}

type Config struct {
	DefaultProvider Provider           `yaml:"default_provider"`
	OpenAI          *OpenAIConfig      `yaml:"openai,omitempty"`
	Gemini          *GeminiConfig      `yaml:"gemini,omitempty"`
	Anthropic       *AnthropicConfig   `yaml:"anthropic,omitempty"`
	Deepseek        *DeepseekConfig    `yaml:"deepseek,omitempty"`
	Ollama          *OllamaConfig      `yaml:"ollama,omitempty"`
	OpenRouter      *OpenRouterConfig  `yaml:"openrouter,omitempty"`
	AzureOpenAI     *AzureOpenAIConfig `yaml:"azure_openai,omitempty"`
}

func ConfigFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, ".how_ai", "config.yaml")
}

func ConfigDirPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, ".how_ai")
}

func Load() (*Config, error) {
	configPath := ConfigFilePath()
	if configPath == "" {
		return nil, errors.New("could not determine user home directory")
	}

	_, err := os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found at %s, please create it or run 'how init' to create a default config", configPath)
		}
		return nil, fmt.Errorf("error checking config file: %w", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	err = config.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

func (c *Config) Validate() error {
	if c.DefaultProvider == "" {
		return errors.New("default_provider is required")
	}

	switch c.DefaultProvider {
	case ProviderOpenAI:
		if c.OpenAI == nil {
			return errors.New("openai configuration is required when default_provider is openai")
		}
		if c.OpenAI.APIKey == "" {
			return errors.New("openai.api_key is required")
		}
		if c.OpenAI.DefaultModel == "" {
			return errors.New("openai.default_model is required")
		}
	case ProviderGemini:
		if c.Gemini == nil {
			return errors.New("gemini configuration is required when default_provider is gemini")
		}
		if c.Gemini.APIKey == "" {
			return errors.New("gemini.api_key is required")
		}
		if c.Gemini.DefaultModel == "" {
			return errors.New("gemini.default_model is required")
		}
	case ProviderAnthropic:
		if c.Anthropic == nil {
			return errors.New("anthropic configuration is required when default_provider is anthropic")
		}
		if c.Anthropic.APIKey == "" {
			return errors.New("anthropic.api_key is required")
		}
		if c.Anthropic.DefaultModel == "" {
			return errors.New("anthropic.default_model is required")
		}
	case ProviderDeepseek:
		if c.Deepseek == nil {
			return errors.New("deepseek configuration is required when default_provider is deepseek")
		}
		if c.Deepseek.APIKey == "" {
			return errors.New("deepseek.api_key is required")
		}
		if c.Deepseek.DefaultModel == "" {
			return errors.New("deepseek.default_model is required")
		}
	case ProviderOllama:
		if c.Ollama == nil {
			return errors.New("ollama configuration is required when default_provider is ollama")
		}
		if c.Ollama.BaseURL == "" {
			return errors.New("ollama.base_url is required")
		}
		if c.Ollama.DefaultModel == "" {
			return errors.New("ollama.default_model is required")
		}
	case ProviderOpenRouter:
		if c.OpenRouter == nil {
			return errors.New("openrouter configuration is required when default_provider is openrouter")
		}
		if c.OpenRouter.APIKey == "" {
			return errors.New("openrouter.api_key is required")
		}
		if c.OpenRouter.DefaultModel == "" {
			return errors.New("openrouter.default_model is required")
		}
	case ProviderAzureOpenAI:
		if c.AzureOpenAI == nil {
			return errors.New("azure_openai configuration is required when default_provider is azure_openai")
		}
		if c.AzureOpenAI.APIKey == "" {
			return errors.New("azure_openai.api_key is required")
		}
		if c.AzureOpenAI.Endpoint == "" {
			return errors.New("azure_openai.endpoint is required")
		}
		if c.AzureOpenAI.DefaultModel == "" {
			return errors.New("azure_openai.default_model is required")
		}
		if c.AzureOpenAI.APIVersion == "" {
			return errors.New("azure_openai.api_version is required")
		}
	default:
		return fmt.Errorf("unsupported provider: %s", c.DefaultProvider)
	}

	return nil
}

func EnsureConfigDirExists() error {
	configDir := ConfigDirPath()
	if configDir == "" {
		return errors.New("could not determine user home directory")
	}

	_, err := os.Stat(configDir)
	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf("error checking config directory: %w", err)
	}

	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	return nil
}

func CreateDefaultConfig() error {
	configPath := ConfigFilePath()
	if configPath == "" {
		return errors.New("could not determine user home directory")
	}

	err := EnsureConfigDirExists()
	if err != nil {
		return err
	}

	_, err = os.Stat(configPath)
	if err == nil {
		return fmt.Errorf("config file already exists at %s", configPath)
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("error checking config file: %w", err)
	}

	defaultConfig := Config{
		DefaultProvider: ProviderOpenAI,
		OpenAI: &OpenAIConfig{
			APIKey:       "your-openai-api-key",
			DefaultModel: "gpt-3.5-turbo",
		},
	}

	data, err := yaml.Marshal(&defaultConfig)
	if err != nil {
		return fmt.Errorf("error creating default config: %w", err)
	}

	err = os.WriteFile(configPath, data, 0600)
	if err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

func SaveConfig(config *Config) error {
	configPath := ConfigFilePath()
	if configPath == "" {
		return errors.New("could not determine user home directory")
	}

	err := EnsureConfigDirExists()
	if err != nil {
		return err
	}

	err = config.Validate()
	if err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	err = os.WriteFile(configPath, data, 0600)
	if err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

func CreateExampleConfig() error {
	configPath := ConfigFilePath()
	if configPath == "" {
		return errors.New("could not determine user home directory")
	}

	err := EnsureConfigDirExists()
	if err != nil {
		return err
	}

	_, err = os.Stat(configPath)
	if err == nil {
		return fmt.Errorf("config file already exists at %s", configPath)
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("error checking config file: %w", err)
	}

	exampleConfig := Config{
		DefaultProvider: ProviderOpenAI,
		OpenAI: &OpenAIConfig{
			APIKey:       "your-openai-api-key",
			BaseURL:      "https://api.openai.com/v1",
			DefaultModel: "gpt-3.5-turbo",
		},
		Gemini: &GeminiConfig{
			APIKey:       "your-gemini-api-key",
			DefaultModel: "gemini-pro",
		},
		Anthropic: &AnthropicConfig{
			APIKey:       "your-anthropic-api-key",
			DefaultModel: "claude-3-opus-20240229",
		},
		Deepseek: &DeepseekConfig{
			APIKey:       "your-deepseek-api-key",
			DefaultModel: "deepseek-coder",
		},
		Ollama: &OllamaConfig{
			BaseURL:      "http://localhost:11434",
			DefaultModel: "llama3",
		},
		OpenRouter: &OpenRouterConfig{
			APIKey:       "your-openrouter-api-key",
			DefaultModel: "openai/gpt-4-turbo",
		},
		AzureOpenAI: &AzureOpenAIConfig{
			APIKey:       "your-azure-openai-api-key",
			Endpoint:     "https://your-resource-name.openai.azure.com",
			DefaultModel: "gpt-4",
			APIVersion:   "2023-05-15",
		},
	}

	data, err := yaml.Marshal(&exampleConfig)
	if err != nil {
		return fmt.Errorf("error creating example config: %w", err)
	}

	err = os.WriteFile(configPath, data, 0600)
	if err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}
