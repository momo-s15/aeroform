# GCP bootstrap (library)

This package implements **GCP-specific** bootstrap steps: workload identity federation for GitHub Actions and a **GCS** bucket for Terraform remote state.

Use from the CLI with Pro Mode config:

```bash
aeroform bootstrap --repo your-org/your-repo
```

Output is rendered by `bootstrap/bootstrap.go` using sections from `bootstrap/gcp/bootstrap.go`.

## What you get

Typical sections include:

- Workload identity pool + provider for GitHub
- Service account IAM bindings for deployment
- GCS bucket creation and `backend "gcs"` HCL fragment

Execute the printed `gcloud` commands with a project owner or equivalent. See [checklist.md](checklist.md).
