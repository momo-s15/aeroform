# LLM Backends

Aeroform uses an LLM to understand your infrastructure prompt and select the best matching template(s). The LLM never generates raw Terraform code — it only selects from pre-validated templates. This is the **constrained generation** approach that keeps deployments safe and predictable.

Three backends are supported. Ollama is the default and recommended for most users.

**Toolchain:** Build Aeroform with **Go 1.26.1+** (see repository `go.mod`). LLM configuration lives under `llm:` in `config.yaml` (Pro) or defaults for Simple Mode.

---

## Ollama (Default)

Ollama runs models locally on your machine. No API keys, no cloud costs, no data leaving your laptop.

### Setup

1. Install Ollama from [ollama.com/download](https://ollama.com/download)

2. Run Aeroform setup to pull the default model:

```bash
aeroform setup
```

This installs `llama3.2` (~2GB). To use a different model:

```bash
aeroform setup --model mistral
```

3. Verify Ollama is running:

```bash
ollama list
```

### Configuration

Ollama is the default backend — no config is needed. To be explicit:

```yaml
llm:
  backend: ollama
  model: llama3.2
  host: http://localhost:11434
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `AEROFORM_LLM_BACKEND` | Override backend | `ollama` |
| `AEROFORM_LLM_MODEL` | Override model | `llama3.2` |
| `AEROFORM_LLM_HOST` | Ollama API endpoint | `http://localhost:11434` |

### Recommended Models

| Model | Size | Speed | Quality |
|-------|------|-------|---------|
| `llama3.2` | 2GB | Fast | Good for most prompts |
| `llama3.1:8b` | 4.7GB | Medium | Better at complex prompts |
| `mistral` | 4.1GB | Medium | Good general-purpose |
| `phi3` | 2.2GB | Fast | Lightweight alternative |

For template selection (which is Aeroform's use case), smaller models work well because the task is classification, not generation.

---

## OpenAI

Use OpenAI's API for template selection. Useful when you don't want to run a local model or need the quality of GPT-4.

### Setup

1. Get an API key from [platform.openai.com/api-keys](https://platform.openai.com/api-keys)

2. Set the key as an environment variable:

```bash
export OPENAI_API_KEY=sk-...
```

Or configure in `config.yaml`:

```yaml
llm:
  backend: openai
  api_key: sk-...
  model: gpt-4o-mini
```

### Configuration

```yaml
llm:
  backend: openai
  api_key: sk-...        # or use OPENAI_API_KEY env var
  model: gpt-4o-mini     # optional, defaults to gpt-4o-mini
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `AEROFORM_LLM_BACKEND` | Set to `openai` | — |
| `OPENAI_API_KEY` | Your API key | — |
| `AEROFORM_LLM_MODEL` | Model to use | `gpt-4o-mini` |

### Cost

OpenAI charges per token. For Aeroform's template selection prompts (typically ~200 tokens in, ~50 tokens out), the cost is negligible — usually under $0.01 per run with `gpt-4o-mini`.

---

## AWS Bedrock

Use Amazon Bedrock for teams that want to keep everything within AWS. Bedrock supports Claude and Titan models.

### Prerequisites

- An AWS account with Bedrock model access enabled
- AWS credentials configured (`aws configure` or environment variables)
- The target model must be enabled in your Bedrock console for your region

### Setup

1. Enable model access in the AWS Bedrock console (region-specific)

2. Configure your IAM role/user with Bedrock permissions:

```json
{
  "Effect": "Allow",
  "Action": ["bedrock:InvokeModel"],
  "Resource": "arn:aws:bedrock:*:*:foundation-model/*"
}
```

3. Configure in `config.yaml`:

```yaml
llm:
  backend: bedrock
  region: us-east-1
  model: anthropic.claude-3-haiku-20240307-v1:0
```

### Configuration

```yaml
llm:
  backend: bedrock
  region: us-east-1                                    # required
  model: anthropic.claude-3-haiku-20240307-v1:0        # optional
  role_arn: arn:aws:iam::123456789012:role/bedrock-role # optional, for cross-account
```

### Supported Models

| Model ID | Provider | Notes |
|----------|----------|-------|
| `anthropic.claude-3-haiku-20240307-v1:0` | Anthropic | Default, fast and cheap |
| `anthropic.claude-3-sonnet-20240229-v1:0` | Anthropic | Higher quality |
| `amazon.titan-text-express-v1` | Amazon | No additional model access needed |
| `amazon.titan-text-lite-v1` | Amazon | Smallest/fastest Titan |

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `AEROFORM_LLM_BACKEND` | Set to `bedrock` | — |
| `AEROFORM_LLM_REGION` | Bedrock region | — |
| `AEROFORM_LLM_MODEL` | Model ID | `anthropic.claude-3-haiku-*` |
| `AWS_REGION` | Fallback region | — |

### Cross-Account Access

If your Bedrock models are in a different account, use `role_arn` to assume a role:

```yaml
llm:
  backend: bedrock
  region: us-east-1
  role_arn: arn:aws:iam::999888777666:role/bedrock-cross-account
```

---

## Backend Comparison

| Feature | Ollama | OpenAI | Bedrock |
|---------|--------|--------|---------|
| **Privacy** | Full (local) | Data sent to OpenAI | Data stays in AWS |
| **Cost** | Free | ~$0.01/run | ~$0.001/run |
| **Latency** | ~1-3s (depends on hardware) | ~1-2s | ~1-2s |
| **Offline** | Yes | No | No |
| **Setup** | Install Ollama + model | API key | AWS credentials + model access |
| **Best for** | Solo devs, privacy-focused | Quick setup, high quality | AWS-native teams |

## How It Works Internally

Regardless of backend, the LLM interaction follows the same pattern:

1. Aeroform builds a prompt containing the user's description and the list of available templates
2. The LLM returns the name(s) of the best-matching template(s)
3. Aeroform validates the response against the known template list
4. If the response is invalid, Aeroform falls back to keyword matching

The LLM is a selection mechanism, not a generation mechanism. It never writes Terraform code. This is why even a small local model works well — the task is template classification, not infrastructure code generation.
