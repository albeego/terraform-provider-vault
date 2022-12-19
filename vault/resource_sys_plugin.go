package vault

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-provider-vault/internal/provider"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func sysPluginResource() *schema.Resource {
	return &schema.Resource{
		Create: sysPluginWrite,
		Update: sysPluginWrite,
		Delete: sysPluginDelete,
		Read:   sysPluginRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The name of the plugin to register",
			},
			"type": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auth": {
							Type:     schema.TypeString,
							Required: true,
						},
						"database": {
							Type:     schema.TypeString,
							Required: true,
						},
						"secret": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"sha256value": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "This is the SHA256 sum of the plugin's binary. Before a plugin is run it's SHA will be checked against this value, if they do not match the plugin can not be run.",
			},
		},
	}
}

func sysPluginWrite(d *schema.ResourceData, meta interface{}) error {
	client, e := provider.GetClient(d, meta)
	if e != nil {
		return e
	}

	log.Printf("[DEBUG] Registering plugin")

	name := d.Get("name").(string)
	pluginType, e := consts.ParsePluginType(d.Get("type").(string))
	if e != nil {
		return e
	}
	sha256value := d.Get("sha256value").(string)

	registerPluginInput := api.RegisterPluginInput{
		Name:    name,
		Type:    pluginType,
		Command: name,
		SHA256:  sha256value,
	}

	e = client.Sys().RegisterPlugin(&registerPluginInput)
	if e != nil {
		return e
	}

	d.SetId(uuid.New().String())

	return nil
}

func sysPluginDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func sysPluginRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}
