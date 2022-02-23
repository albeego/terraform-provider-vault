package vault

import (
	"fmt"
	"github.com/google/uuid"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/vault/api"
)

func sysUnsealResource() *schema.Resource {
	return &schema.Resource{
		Create: sysUnsealWrite,
		Update: sysUnsealWrite,
		Delete: sysUnsealDelete,
		Read:   sysUnsealRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"keys": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Unseal keys",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sealed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Seal status of the vault cluster",
			},
			"threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Total number of keys required to unseal",
			},
			"number_of_shares": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Total number of keys vault cluster was initialized with",
			},
			"progress": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Number of keys left to submit before vault cluster is unsealed",
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Vault version number of cluster",
			},
		},
	}
}

func sysUnsealWrite(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	ikeys := d.Get("keys").([]interface{})
	keys := make([]string, 0, len(ikeys))
	for _, ikey := range ikeys {
		keys = append(keys, ikey.(string))
	}

	log.Printf("[DEBUG] Unsealing vault")

	result, err := client.Sys().Unseal(keys[0])
	if err != nil {
		return fmt.Errorf("error unsealing vault: %s", err)
	}

	if result.T > len(keys)-1 {
		return fmt.Errorf("can't unseal vault, need: %d keys but only have %d", result.T+1, len(keys))
	}
	for i := 1; i < len(keys); i++ {
		result, err = client.Sys().Unseal(keys[0])
		if err != nil {
			return fmt.Errorf("error unsealing vault: %s", err)
		}
		if result.Progress == 0 {
			break
		}
	}

	d.SetId(uuid.New().String())
	d.Set("sealed", result.Sealed)
	d.Set("threshold", result.T)
	d.Set("number_of_shares", result.N)
	d.Set("progress", result.Progress)
	d.Set("version", result.Version)

	return nil
}

func sysUnsealDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func sysUnsealRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}
