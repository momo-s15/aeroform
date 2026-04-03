# Getting Started - Pro

Pro Mode is for teams who want explicit config control and cloud-provider specific planning.

## Prepare Config

Copy `config.yaml.example` to `config.yaml` and fill in the cloud-specific fields.

## First Run

```bash
go run . bootstrap
go run . plan "resilient kubernetes cluster"
go run . generate "resilient kubernetes cluster"
```

## What Happens

1. Aeroform loads the config file.
2. Aeroform selects the cloud provider implementation.
3. Aeroform builds a Pro plan and security scan output.
4. Aeroform keeps deployment decisions explicit.

