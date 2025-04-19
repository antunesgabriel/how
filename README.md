# How AI CLI

A terminal-based AI assistant powered by various LLM providers.

## Configuration

How AI CLI supports both global and local configurations, allowing you to have different settings for different projects.

### Global Configuration

Global configuration applies to all projects and is stored in `~/.how/config.yaml`.

1. Create a default global configuration:

```bash
how init
```

This creates a minimal configuration file with OpenAI as the default provider.

2. Edit the configuration file to add your API keys and customize settings:

```bash
# Open with your preferred editor
nano ~/.how/config.yaml
# or
code ~/.how/config.yaml
```

### Local Configuration

Local configuration applies only to the current directory and is stored in `./.how/config.yaml`. Local configuration takes precedence over global configuration.

1. Create a default local configuration:

```bash
how init --local
```

2. Edit the local configuration file:

```bash
# Open with your preferred editor
nano ./.how/config.yaml
# or
code ./.how/config.yaml
```

### Supported Providers

The following AI providers are supported:

#### OpenAI

```yaml
default_provider: openai
openai:
  api_key: "your-openai-api-key"
  model: "gpt-4o"  # or any other OpenAI model
  timeout: 30000   # optional, in milliseconds (30 seconds)
  # Azure OpenAI specific settings (optional)
  by_azure: false
  base_url: ""     # required for Azure: https://{YOUR_RESOURCE_NAME}.openai.azure.com
  api_version: ""  # required for Azure
```

#### Claude (Anthropic)

```yaml
default_provider: claude
claude:
  api_key: "your-anthropic-api-key"
  model: "claude-3-opus-20240229"  # or any other Claude model
  max_tokens: 2000  # required
  # AWS Bedrock specific settings (optional)
  by_bedrock: false
  access_key: ""
  secret_access_key: ""
  region: ""
```

#### Gemini (Google)

```yaml
default_provider: gemini
gemini:
  api_key: "your-gemini-api-key"
  model: "gemini-pro"  # or gemini-pro-vision, gemini-1.5-flash, etc.
```

#### Deepseek

```yaml
default_provider: deepseek
deepseek:
  api_key: "your-deepseek-api-key"
  model: "deepseek-coder"  # or other Deepseek models
  timeout: 60000  # optional, in milliseconds (1 minute)
  base_url: "https://api.deepseek.com/"  # optional
```

#### Ollama (Local Models)

```yaml
default_provider: ollama
ollama:
  base_url: "http://localhost:11434"  # URL to your Ollama server
  model: "llama3"  # or any other model you have installed
  timeout: 30000  # optional, in milliseconds (30 seconds)
```

### Configuration Priority

When running the How AI CLI, it will:

1. First check for a local configuration file (`./.how/config.yaml`)
2. If no local configuration is found, use the global configuration (`~/.how/config.yaml`)
3. If no configuration is found at all, prompt you to create one with `how init`

This allows you to have different settings for different projects while maintaining a global default.
