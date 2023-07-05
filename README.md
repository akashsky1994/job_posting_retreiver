# Job Posting Aggregator
Fetches jobs from multiple sources into single csv




# Terraform Setup

```
cd terraform_deployment
terraform plan -out=terraform.tfplan -var "do_token=${DO_PAT}" -var "pvt_key=${DO_PVT_KEY}"
terraform show terraform.tfplan
terraform apply terraform.tfplan
```

Terraform Delete
```
terraform plan -destroy -out=terraform_del.tfplan -var "do_token=${DO_PAT}" -var "pvt_key=${DO_PVT_KEY}"
terraform apply terraform_del.tfplan
```
