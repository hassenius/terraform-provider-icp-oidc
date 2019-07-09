# Terraform provider ICP OIDC

A really quick and dirty hack to get token and certificates from IBM Cloud Private OIDC

## Example use
```hcl
resource "icp_oidc_token" "aws_token" {
  host = "${var.console}"
  username      = "admin"
  password = "${var.password}"
  insecure  = true
}

resource "icp_oidc_certificate" "aws_cert" {
  host = "${var.console}"
  token  = "${icp_oidc_token.aws_token.access_token}"
  insecure  = true
}

provider "kubernetes" {
  host      = "https://${var.console_host}:8001"
  # client_certificate = "${icp_oidc_certificate.aws_cert.certificate}"
  # client_key = "${icp_oidc_certificate.aws_cert.key}"
  insecure  = "true"
  token    = "${icp_oidc_token.aws_token.id_token}"
  # alias     = "awskube"
}
```

## Reference
https://medium.com/@raghavendra.1729/leverage-certificate-authentication-in-ibm-cloud-private-f45159f04d5e


## Use
go build -o terraform-provider-icp_oidc

cp terraform-provider-icp_oidc <path_to_terraform>/terraform.d/plugins/darwin_amd64/terraform-provider-icp
