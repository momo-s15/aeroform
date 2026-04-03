# GCP bootstrap checklist

After running `aeroform bootstrap --repo owner/repo` and applying the steps:

- [ ] Workload identity pool and GitHub provider are configured
- [ ] Service account used by Actions can impersonate or use WIF correctly
- [ ] GCS state bucket exists and backend `gcs` block matches bucket + prefix
- [ ] Required APIs (e.g. IAM, Storage) are enabled on the project
