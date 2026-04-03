# Azure bootstrap (library)

This package implements **Azure-specific** bootstrap steps: federated identity (GitHub Actions OIDC / workload identity) and **Blob storage** backend snippets for Terraform state.

Use from the CLI with Pro Mode config:

```bash
aeroform bootstrap --repo your-org/your-repo
```

Output is rendered by `bootstrap/bootstrap.go` using sections from `bootstrap/azure/bootstrap.go`.

## What you get

Typical sections include:

- App registration / federated credential hints for GitHub
- Storage account + container recommendations for remote state
- Backend `azurerm` HCL fragments suitable to paste into your root module

Run generated commands in `az` with an account that can create those resources. See [checklist.md](checklist.md).
