# AWS bootstrap checklist

After running `aeroform bootstrap --repo owner/repo` and applying the commands:

- [ ] GitHub OIDC provider exists in IAM (or was already present)
- [ ] Deploy role trusts `token.actions.githubusercontent.com` for your repository
- [ ] S3 state bucket and DynamoDB lock table exist in the expected region
- [ ] Backend block in Terraform matches bucket, key, region, and `dynamodb_table`
- [ ] CI role can `sts:AssumeRoleWithWebIdentity` from GitHub Actions
