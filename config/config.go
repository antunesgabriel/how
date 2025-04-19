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
	ProviderOpenAI   Provider = "openai"
	ProviderGemini   Provider = "gemini"
	ProviderClaude   Provider = "claude"
	ProviderDeepseek Provider = "deepseek"
	ProviderOllama   Provider = "ollama"
)

// OpenAIChatModelConfig contains the configuration options for the OpenAI model
type OpenAIChatModelConfig struct {
	// APIKey is your authentication key
	// Use OpenAI API key or Azure API key depending on the service
	// Required
	APIKey string `yaml:"api_key"`

	// Timeout specifies the maximum duration to wait for API responses in milliseconds
	// If HTTPClient is set, Timeout will not be used.
	// Optional. Default: 30000 (30 seconds)
	Timeout int `yaml:"timeout,omitempty"`

	// The following three fields are only required when using Azure OpenAI Service, otherwise they can be ignored.
	// For more details, see: https://learn.microsoft.com/en-us/azure/ai-services/openai/

	// ByAzure indicates whether to use Azure OpenAI Service
	// Required for Azure
	ByAzure bool `yaml:"by_azure,omitempty"`

	// BaseURL is the Azure OpenAI endpoint URL
	// Format: https://{YOUR_RESOURCE_NAME}.openai.azure.com. YOUR_RESOURCE_NAME is the name of your resource that you have created on Azure.
	// Required for Azure
	BaseURL string `yaml:"base_url,omitempty"`

	// APIVersion specifies the Azure OpenAI API version
	// Required for Azure
	APIVersion string `yaml:"api_version,omitempty"`

	// The following fields correspond to OpenAI's chat completion API parameters
	// Ref: https://platform.openai.com/docs/api-reference/chat/create

	// Model specifies the ID of the model to use
	// Required
	Model string `yaml:"model"` // use gpt-4o as a default

	// MaxTokens limits the maximum number of tokens that can be generated in the chat completion
	// Optional. Default: model's maximum
	MaxTokens *int `yaml:"max_tokens,omitempty"`
}

// GeminiConfig contains the configuration options for the Gemini model
type GeminiConfig struct {
	// Model specifies which Gemini model to use
	// Examples: "gemini-pro", "gemini-pro-vision", "gemini-1.5-flash"
	// Required
	Model string `yaml:"model"`

	// MaxTokens limits the maximum number of tokens in the response
	// Optional
	MaxTokens *int `yaml:"max_tokens,omitempty"`

	// APIKey is your Gemini API key
	// Required
	APIKey string `yaml:"api_key"`
}

// ClaudeConfig contains the configuration options for the Claude model
type ClaudeConfig struct {
	// ByBedrock indicates whether to use Bedrock Service
	// Required for Bedrock
	ByBedrock bool `yaml:"by_bedrock,omitempty"`

	// AccessKey is your Bedrock API Access key
	// Obtain from: https://docs.aws.amazon.com/bedrock/latest/userguide/getting-started.html
	// Required for Bedrock
	AccessKey string `yaml:"access_key,omitempty"`

	// SecretAccessKey is your Bedrock API Secret Access key
	// Obtain from: https://docs.aws.amazon.com/bedrock/latest/userguide/getting-started.html
	// Required for Bedrock
	SecretAccessKey string `yaml:"secret_access_key,omitempty"`

	// SessionToken is your Bedrock API Session Token
	// Obtain from: https://docs.aws.amazon.com/bedrock/latest/userguide/getting-started.html
	// Optional for Bedrock
	SessionToken string `yaml:"session_token,omitempty"`

	// Region is your Bedrock API region
	// Obtain from: https://docs.aws.amazon.com/bedrock/latest/userguide/getting-started.html
	// Required for Bedrock
	Region string `yaml:"region,omitempty"`

	// BaseURL is the custom API endpoint URL
	// Use this to specify a different API endpoint, e.g., for proxies or enterprise setups
	// Optional. Example: "https://custom-claude-api.example.com"
	BaseURL *string `yaml:"base_url,omitempty"`

	// APIKey is your Anthropic API key
	// Obtain from: https://console.anthropic.com/account/keys
	// Required
	APIKey string `yaml:"api_key"`

	// Model specifies which Claude model to use
	// Required
	Model string `yaml:"model"`

	// MaxTokens limits the maximum number of tokens in the response
	// Range: 1 to model's context length
	// Required. Example: 2000 for a medium-length response
	MaxTokens int `yaml:"max_tokens"`
}

// DeepseekChatModelConfig contains the configuration options for the Deepseek model
type DeepseekChatModelConfig struct {
	// APIKey is your authentication key
	// Required
	APIKey string `yaml:"api_key"`

	// Timeout specifies the maximum duration to wait for API responses in milliseconds
	// Optional. Default: 60000 (1 minute)
	Timeout int `yaml:"timeout,omitempty"`

	// BaseURL is your custom deepseek endpoint url
	// Optional. Default: https://api.deepseek.com/
	BaseURL string `yaml:"base_url,omitempty"`

	// Model specifies the ID of the model to use
	// Models available:
	// - deepseek-chat: DeepSeekChat is the official model for chat completions
	// - deepseek-coder: DeepSeekCoder has been combined with DeepSeekChat, but you can still use it
	// - deepseek-reasoner: DeepSeekReasoner is the official model for reasoning completions
	// - DeepSeek-R1: Azure model for DeepSeek R1
	// - deepseek/deepseek-r1: OpenRouter model for DeepSeek R1
	// - deepseek/deepseek-r1-distill-llama-70b: DeepSeek R1 Distill Llama 70B
	// - deepseek/deepseek-r1-distill-llama-8b: DeepSeek R1 Distill Llama 8B
	// - deepseek/deepseek-r1-distill-qwen-14b: DeepSeek R1 Distill Qwen 14B
	// - deepseek/deepseek-r1-distill-qwen-1.5b: DeepSeek R1 Distill Qwen 1.5B
	// - deepseek/deepseek-r1-distill-qwen-32b: DeepSeek R1 Distill Qwen 32B
	// Required
	Model string `yaml:"model"`

	// MaxTokens limits the maximum number of tokens that can be generated in the chat completion
	// Range: [1, 8192].
	// Optional. Default: 4096
	MaxTokens int `yaml:"max_tokens,omitempty"`
}

// OllamaChatModelConfig contains the configuration options for the Ollama model
type OllamaChatModelConfig struct {
	// BaseURL is the URL of your Ollama server
	// Required
	BaseURL string `yaml:"base_url"`

	// Timeout specifies the maximum duration to wait for API responses in milliseconds
	// Optional. Default: 30000 (30 seconds)
	Timeout int `yaml:"timeout,omitempty"`

	// Model specifies which Ollama model to use
	// Required
	Model string `yaml:"model"`
}

type Config struct {
	DefaultProvider Provider                 `yaml:"default_provider"`
	OpenAI          *OpenAIChatModelConfig   `yaml:"openai,omitempty"`
	Gemini          *GeminiConfig            `yaml:"gemini,omitempty"`
	Claude          *ClaudeConfig            `yaml:"claude,omitempty"`
	Deepseek        *DeepseekChatModelConfig `yaml:"deepseek,omitempty"`
	Ollama          *OllamaChatModelConfig   `yaml:"ollama,omitempty"`
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
		if c.OpenAI.Model == "" {
			return errors.New("openai.model is required")
		}
	case ProviderGemini:
		if c.Gemini == nil {
			return errors.New("gemini configuration is required when default_provider is gemini")
		}
		if c.Gemini.APIKey == "" {
			return errors.New("gemini.api_key is required")
		}
		if c.Gemini.Model == "" {
			return errors.New("gemini.model is required")
		}
	case ProviderClaude:
		if c.Claude == nil {
			return errors.New("claude configuration is required when default_provider is claude")
		}
		if c.Claude.APIKey == "" {
			return errors.New("claude.api_key is required")
		}
		if c.Claude.Model == "" {
			return errors.New("claude.model is required")
		}
	case ProviderDeepseek:
		if c.Deepseek == nil {
			return errors.New("deepseek configuration is required when default_provider is deepseek")
		}
		if c.Deepseek.APIKey == "" {
			return errors.New("deepseek.api_key is required")
		}
		if c.Deepseek.Model == "" {
			return errors.New("deepseek.model is required")
		}
	case ProviderOllama:
		if c.Ollama == nil {
			return errors.New("ollama configuration is required when default_provider is ollama")
		}
		if c.Ollama.BaseURL == "" {
			return errors.New("ollama.base_url is required")
		}
		if c.Ollama.Model == "" {
			return errors.New("ollama.model is required")
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
		OpenAI: &OpenAIChatModelConfig{
			APIKey: "your-openai-api-key",
			Model:  "gpt-4o",
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

	// Default timeout values in milliseconds
	defaultTimeout := 30000  // 30 seconds
	deepseekTimeout := 60000 // 1 minute

	// Create example config with helpful comments
	exampleConfig := Config{
		DefaultProvider: ProviderOpenAI,
		OpenAI: &OpenAIChatModelConfig{
			APIKey:    "your-openai-api-key",
			BaseURL:   "https://api.openai.com/v1", // Optional: Default OpenAI API endpoint
			Model:     "gpt-4o",                    // Default model
			Timeout:   defaultTimeout,              // Optional: 30 seconds timeout
			ByAzure:   false,                       // Set to true if using Azure OpenAI
			MaxTokens: nil,                         // Optional: Use model's default max tokens
		},
		Gemini: &GeminiConfig{
			APIKey:    "your-gemini-api-key",
			Model:     "gemini-pro", // Other options: gemini-pro-vision, gemini-1.5-flash
			MaxTokens: nil,          // Optional: Use model's default max tokens
		},
		Claude: &ClaudeConfig{
			APIKey:    "your-anthropic-api-key",
			Model:     "claude-3-opus-20240229",
			MaxTokens: 2000, // Required: Example for medium-length responses
			// Bedrock configuration (optional, only if using AWS Bedrock)
			ByBedrock:       false,
			AccessKey:       "",
			SecretAccessKey: "",
			Region:          "",
		},
		Deepseek: &DeepseekChatModelConfig{
			APIKey:    "your-deepseek-api-key",
			BaseURL:   "https://api.deepseek.com/", // Optional: Default Deepseek API endpoint
			Model:     "deepseek-coder",            // Default model for coding tasks
			Timeout:   deepseekTimeout,             // Optional: 1 minute timeout
			MaxTokens: 4096,                        // Optional: Default is 4096
		},
		Ollama: &OllamaChatModelConfig{
			BaseURL: "http://localhost:11434", // Default Ollama server URL
			Model:   "llama3",                 // Choose your locally available model
			Timeout: defaultTimeout,           // Optional: 30 seconds timeout
		},
	}

	data, err := yaml.Marshal(&exampleConfig)
	if err != nil {
		return fmt.Errorf("error creating example config: %w", err)
	}

	// Add a header comment to the YAML file
	yamlWithComments := "# How AI Configuration File\n" +
		"# This file configures the AI providers for the How AI CLI tool\n" +
		"# Generated with 'how init' command\n" +
		"# \n" +
		"# Available providers: openai, gemini, claude, deepseek, ollama\n" +
		"# Set default_provider to one of these values\n" +
		"# \n" +
		"# Timeout values are in milliseconds (1000ms = 1 second)\n" +
		"# \n\n" +
		string(data)

	err = os.WriteFile(configPath, []byte(yamlWithComments), 0600)
	if err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}
