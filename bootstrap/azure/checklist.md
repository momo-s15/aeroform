# Azure bootstrap checklist

After running `aeroform bootstrap --repo owner/repo` and applying the steps:

- [ ] Federated credential / app registration allows your GitHub repo subject
- [ ] Storage account and container for Terraform state exist
- [ ] Backend `azurerm` configuration matches resource group, account, and container
- [ ] Pipeline identity can read/write state and deploy resources
