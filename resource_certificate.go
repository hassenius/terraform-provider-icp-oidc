package main

import (
      "github.com/hashicorp/terraform/helper/schema"
      "crypto/tls"
      "errors"
      "net/http"
      "fmt"
      "encoding/json"
      "io/ioutil"
      "math/rand"
      "time"
)

func init() {
	rand.Seed(time.Now().Unix())

}

func resourceCertificate() *schema.Resource {
        return &schema.Resource{
                Create: resourceCertificateCreate,
                Read:   resourceCertificateRead,
                Update: resourceCertificateUpdate,
                Delete: resourceCertificateDelete,

                Schema: map[string]*schema.Schema{
                        "host": &schema.Schema{
                  				Type:        schema.TypeString,
                  				Required:    true,
                  				DefaultFunc: schema.EnvDefaultFunc("ICP_OICD_HOST", nil),
                  				Description: "Host to contact the OIDC service",
                  			},
                        "token": &schema.Schema{
                  				Type:        schema.TypeString,
                  				Required:    true,
                  				DefaultFunc: schema.EnvDefaultFunc("ICP_OICD_TOKEN", nil),
                  				Description: "Token to authenticate the OIDC service",
                  			},
                        "request_path": &schema.Schema{
                  				Type:        schema.TypeString,
                  				Optional:    true,
                  				DefaultFunc: schema.EnvDefaultFunc("ICP_OICD_PATH", "/idmgmt/identity/api/v1/certificates"),
                  				Description: "Host to contact the OIDC service",
                  			},
                        "port": &schema.Schema{
                  				Type:        schema.TypeInt,
                  				Optional:    true,
                  				DefaultFunc: schema.EnvDefaultFunc("ICP_OICD_PORT", 8443),
                  				Description: "Port of the OIDC service.",
                  			},
                  			"insecure": &schema.Schema{
                  				Type:        schema.TypeBool,
                  				Optional:    true,
                  				DefaultFunc: schema.EnvDefaultFunc("ICP_OICD_INSECURE", false),
                  				Description: "Skip SSL host validation when using self signed cert",
                  			},
                        "refresh_expired": &schema.Schema{
                  				Type:        schema.TypeBool,
                  				Optional:    true,
                  				DefaultFunc: schema.EnvDefaultFunc("ICP_OICD_AUTOREFRESH", false),
                  				Description: "Automatically refresh expired tokens. Not implemented yet",
                  			},
                        "user_id": &schema.Schema{
                  				Type:        schema.TypeString,
                          Computed:    true,
                          ForceNew:    false,
                  				Description: "Access user retrieved from OIDC Server.",
                  			},
                        "certificate": &schema.Schema{
                          Type:        schema.TypeString,
                          Computed:    true,
                          ForceNew:    false,
                          Description: "Access certificate retrieved from OIDC service.",
                        },
                        "key": &schema.Schema{
                          Type:        schema.TypeString,
                          Computed:    true,
                          ForceNew:    false,
                          Description: "Access key retrieved OIDC service.",
                        },
                        "triggers": {
                          Type:     schema.TypeMap,
                          Optional: true,
                          ForceNew: true,
                        },
                },
        }
}


func resourceCertificateCreate(d *schema.ResourceData, m interface{}) error {
        host         := d.Get("host").(string)
        port         := d.Get("port").(int)
        insecure     := d.Get("insecure").(bool)
        token        := "Bearer " + d.Get("token").(string)
        request_path := d.Get("request_path").(string)
        oidcUrl      := fmt.Sprintf("https://%s:%d%s", host, port, request_path)


        type certstrct struct {
          UserId string `json:"userId"`
          Certificate string `json:"certificate"`
          Key string `json:"key"`
        }

        certresp := certstrct{}

        http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: insecure}

        req, _ := http.NewRequest(http.MethodPost, oidcUrl, nil)

      	req.Header.Set("Authorization", token)

      	resp, err := http.DefaultClient.Do(req)

        if err != nil {
          return errors.New(fmt.Sprintf("Error communicating with OIDC server. %s", err.Error()))
        }

        defer resp.Body.Close()
        cert_json, err := ioutil.ReadAll(resp.Body)

        if nil != err {
          return errors.New(fmt.Sprintf("Error reading response body. %s", err.Error()))
        }

        err = json.Unmarshal(cert_json, &certresp)
        if err != nil {
          return errors.New(fmt.Sprintf("Error unmarshaling response json. %s", err.Error()))
        }

        d.Set("user_id", certresp.UserId)
        d.Set("key", certresp.Key)
        d.Set("certificate", certresp.Certificate)
        d.SetId(fmt.Sprintf("%d", rand.Int()))

        return nil
}


func resourceCertificateRead(d *schema.ResourceData, m interface{}) error {

  return nil
}

func resourceCertificateUpdate(d *schema.ResourceData, m interface{}) error {
        return resourceCertificateRead(d, m)
}

func resourceCertificateDelete(d *schema.ResourceData, m interface{}) error {
        d.SetId("")
        return nil
}
