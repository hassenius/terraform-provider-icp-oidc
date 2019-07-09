package main

import (
        "github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
        return &schema.Provider{
                ResourcesMap: map[string]*schema.Resource{
                  "icp_oidc_token": resourceToken(),
                  "icp_oidc_certificate": resourceCertificate(),
                },
        }
}
