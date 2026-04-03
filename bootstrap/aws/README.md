# AWS bootstrap (library)

This package implements **AWS-specific** bootstrap content: shell-oriented steps for GitHub OIDC, IAM roles, and remote **S3 + DynamoDB** Terraform state.

End users do **not** run code from this folder directly. From a project with `config.yaml`, run:

```bash
aeroform bootstrap --repo your-org/your-repo
```

The CLI loads config, calls `bootstrap.Run`, and prints numbered sections (see `bootstrap/bootstrap.go`). AWS commands are produced from `bootstrap/aws/bootstrap.go` (and related helpers).

## What you get

Typical sections include:

- GitHub OIDC identity provider (if needed)
- IAM role trust policy for `repo:your-org/your-repo`
- S3 bucket + DynamoDB table for state locking (naming derived from account/region)

Review and run the printed commands in an authenticated AWS shell. Adjust names to match your org’s standards.

See [checklist.md](checklist.md) for a short verification list.
