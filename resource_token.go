package main

import (
      "github.com/hashicorp/terraform/helper/schema"
      "crypto/tls"
      "errors"
      "net/http"
      "net/url"
      "fmt"
      "encoding/json"
      "io/ioutil"
      "math/rand"
      "time"
)

func init() {
	rand.Seed(time.Now().Unix())

}

func resourceToken() *schema.Resource {
        return &schema.Resource{
                Create: resourceTokenCreate,
                Read:   resourceTokenRead,
                Update: resourceTokenUpdate,
                Delete: resourceTokenDelete,

                Schema: map[string]*schema.Schema{
                        "host": &schema.Schema{
                  				Type:        schema.TypeString,
                  				Required:    true,
                  				DefaultFunc: schema.EnvDefaultFunc("ICP_OICD_HOST", nil),
                  				Description: "Host to contact the OIDC service",
                  			},
                        "request_path": &schema.Schema{
                  				Type:        schema.TypeString,
                  				Optional:    true,
                  				DefaultFunc: schema.EnvDefaultFunc("ICP_OICD_PATH", "/idprovider/v1/auth/identitytoken"),
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
                  			"username": &schema.Schema{
                  				Type:        schema.TypeString,
                  				Optional:    true,
                  				DefaultFunc: schema.EnvDefaultFunc("ICP_OICD_USERNAME", nil),
                  				Description: "Username to authenticate with the OIDC service.",
                  			},
                  			"password": &schema.Schema{
                  				Type:        schema.TypeString,
                  				Optional:    true,
                  				DefaultFunc: schema.EnvDefaultFunc("ICP_OICD_PASSWORD", nil),
                  				Description: "Password to authenticate with the OIDC service.",
                  			},
                        "token_data": &schema.Schema{
                  				Type:        schema.TypeMap,
                          Optional:    true,
                          ForceNew:    false,
                  				Description: "Response from OIDC Server.",
                  			},
                        "access_token": &schema.Schema{
                          Type:        schema.TypeString,
                          Computed:    true,
                          ForceNew:    false,
                          Description: "Access token retrieved from OIDC service.",
                        },
                        "token_type": &schema.Schema{
                          Type:        schema.TypeString,
                          Computed:    true,
                          ForceNew:    false,
                          Description: "Access token type retrieved OIDC service.",
                        },
                        "id_token": &schema.Schema{
                          Type:        schema.TypeString,
                          Computed:    true,
                          ForceNew:    false,
                          Description: "ID token retrieved from OIDC service.",
                        },
                        "triggers": {
                          Type:     schema.TypeMap,
                          Optional: true,
                          ForceNew: true,
                        },
                },
        }
}


func resourceTokenCreate(d *schema.ResourceData, m interface{}) error {
        host         := d.Get("host").(string)
        port         := d.Get("port").(int)
        insecure     := d.Get("insecure").(bool)
        username     := d.Get("username").(string)
        password     := d.Get("password").(string)
        request_path := d.Get("request_path").(string)
        oidcUrl      := fmt.Sprintf("https://%s:%d%s", host, port, request_path)

        /* # curl -H “Content-Type: application/x-www-form-urlencoded;charset=UTF-8” -d “grant_type=password&username=admin&password=<ICP admin password>&scope=openid” */
        /*  Define the response structure from OIDC */
        type tokenresp struct {
          AccessToken string `json:"access_token"`
          TokenType string `json:"token_type"`
          IdToken string `json:"id_token"`
        }

        token := tokenresp{}

        v := url.Values{}
        v.Set("grant_type", "password")
        v.Set("username", username)
        v.Set("password", password)
        v.Set("scope", "openid")

        http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: insecure}
        resp, err := http.PostForm(oidcUrl, v)

        if err != nil {
          return errors.New(fmt.Sprintf("Error communicating with OIDC server. %s", err.Error()))
        }

        defer resp.Body.Close()
        token_json, err := ioutil.ReadAll(resp.Body)

        if nil != err {
          return errors.New(fmt.Sprintf("Error reading response body. %s", err.Error()))
        }

        err = json.Unmarshal(token_json, &token)
        if err != nil {
          return errors.New(fmt.Sprintf("Error unmarshaling response json. %s", err.Error()))
        }

        // I'll see about getting this to work another day
        // t := make(map[string]string)
        // for k, v := range token_json {
      	// 	t[k] = fmt.Sprintf("%v", v)
      	// }
        d.Set("token_data", token_json)
        d.Set("access_token", token.AccessToken)
        d.Set("token_type", token.TokenType)
        d.Set("id_token", token.IdToken)
        d.SetId(fmt.Sprintf("%d", rand.Int()))

        return nil
}


func resourceTokenRead(d *schema.ResourceData, m interface{}) error {

  return nil
}

func resourceTokenUpdate(d *schema.ResourceData, m interface{}) error {
        return resourceTokenRead(d, m)
}

func resourceTokenDelete(d *schema.ResourceData, m interface{}) error {
        d.SetId("")
        return nil
}
