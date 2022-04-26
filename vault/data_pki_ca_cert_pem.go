package vault

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/vault/api"
)

func caCert() *schema.Resource {
	return &schema.Resource{
		Read: readCaCert,

		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Qualifying path from which the CA certificate will be read.",
			},

			"pem": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "PEM encoded CA Certificate.",
			},
		},
	}
}

func readCaCert(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	path := d.Get("path").(string) + "/ca/pem"

	log.Printf("[DEBUG] Reading %s from Vault", path)

	r := client.NewRequest("GET", "/v1/"+path)
	resp, err := client.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading from Vault: %s", err)
	}
	d.Set("pem", string(body))
	log.Printf("[DEBUG] PEM set as %s from Vault", d.Get("pem").(string))

	return nil
}
