# Job Posting Aggregator
Fetches jobs from multiple sources into single csv

### Running the project

Standalone
```
go run . --env ./config/.env.prod
```

In Docker
```
docker compose --env-file ./config/.env.prod build
docker compose --env-file ./config/.env.prod up -d
```




# Terraform Setup

```
cd terraform_deployment
terraform plan -out=terraform.tfplan -var "do_token=${DO_PAT}" -var "public_ip=${JR_PUBLIC_IP}"
terraform show terraform.tfplan
terraform apply terraform.tfplan
```

Terraform Delete
```
terraform plan -destroy -out=terraform_del.tfplan -var "do_token=${DO_PAT}" -var "pvt_key=${DO_PVT_KEY}"
terraform apply terraform_del.tfplan
```
