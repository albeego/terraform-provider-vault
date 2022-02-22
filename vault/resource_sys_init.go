package vault

import (
	"fmt"
	"github.com/google/uuid"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/vault/api"
)

func sysInitResource() *schema.Resource {
	return &schema.Resource{
		Create: sysInitWrite,
		Update: sysInitWrite,
		Delete: sysInitDelete,
		Read:   sysInitRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"secret_shares": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Number of key shares to split the generated master key into. This is the number of \"unseal keys\" to generate",
			},
			"secret_threshold": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Number of key shares required to reconstruct the master key. This must be less than or equal to -key-shares",
			},
			"stored_shares": {
				Type:        schema.TypeInt,
				Required:    false,
				Computed:    true,
				Optional:    true,
				Description: "Number of unseal keys to store on an HSM. This must be equal to -key-shares",
			},
			"pgp_keys": {
				Type:        schema.TypeList,
				Required:    false,
				Computed:    true,
				Optional:    true,
				Description: "Comma-separated list of paths to files on disk containing public PGP keys OR a comma-separated list of Keybase usernames using the format keybase:<username>. When supplied, the generated unseal keys will be encrypted and base64-encoded in the order specified in this list. The number of entries must match -key-shares, unless -stored-shares are used",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"recovery_shares": {
				Type:        schema.TypeInt,
				Required:    false,
				Computed:    true,
				Optional:    true,
				Description: "Number of key shares to split the recovery key into. This is only used Auto Unseal seals (HSM, KMS and Transit seals)",
			},
			"recovery_threshold": {
				Type:        schema.TypeInt,
				Required:    false,
				Computed:    true,
				Optional:    true,
				Description: "Number of key shares required to reconstruct the recovery key. This is only used Auto Unseal seals (HSM, KMS and Transit seals)",
			},
			"recovery_pgp_keys": {
				Type:        schema.TypeList,
				Required:    false,
				Computed:    true,
				Optional:    true,
				Description: "Behaves like -pgp-keys, but for the recovery key shares. This is only used with Auto Unseal seals (HSM, KMS and Transit seals)",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"root_token_pgp_key": {
				Type:        schema.TypeString,
				Required:    false,
				Computed:    true,
				Optional:    true,
				Description: "Path to a file on disk containing a binary or base64-encoded public PGP key. This can also be specified as a Keybase username using the format keybase:<username>. When supplied, the generated root token will be encrypted and base64-encoded with the given public key",
			},
			"keys": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Key shares the generated master key is split into. These are the \"unseal keys\"",
			},
			"keys_base64": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Key shares in base64",
			},
			"recovery_keys": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Key shares the recovery key is split into. These are only used in Auto Unseal seals (HSM, KMS and Transit seals)",
			},
			"recovery_keys_base64": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Recovery key shares in base64",
			},
			"root_token": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The generated root token",
			},
		},
	}
}

func sysInitWrite(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	secretShares := d.Get("secret_shares").(int)
	secretThreshold := d.Get("secret_threshold").(int)

	initRequest := api.InitRequest{
		SecretShares:    secretShares,
		SecretThreshold: secretThreshold,
	}

	if storedShares, hasStoredShares := d.GetOk("stored_shares"); hasStoredShares {
		initRequest.StoredShares = storedShares.(int)
	}
	if pgpKeys, hasPgpKeys := d.GetOk("pgp_keys"); hasPgpKeys {
		initRequest.PGPKeys = pgpKeys.([]string)
	}
	if recoveryShares, hasRecoveryShares := d.GetOk("recovery_shares"); hasRecoveryShares {
		initRequest.RecoveryShares = recoveryShares.(int)
	}
	if recoveryThreshold, hasRecoveryThreshold := d.GetOk("recovery_threshold"); hasRecoveryThreshold {
		initRequest.RecoveryThreshold = recoveryThreshold.(int)
	}
	if recoveryPgpKeys, hasRecoveryPgpKeys := d.GetOk("recovery_pgp_keys"); hasRecoveryPgpKeys {
		initRequest.RecoveryPGPKeys = recoveryPgpKeys.([]string)
	}
	if rootTokenPgpKey, hasRootTokenPgpKey := d.GetOk("recovery_shares"); hasRootTokenPgpKey {
		initRequest.RootTokenPGPKey = rootTokenPgpKey.(string)
	}

	log.Printf("[DEBUG] Initializing vault with %d secret shares and a %d secret threshold", secretShares, secretThreshold)

	result, err := client.Sys().Init(&initRequest)

	if err != nil {
		return fmt.Errorf("error initializing vault: %s", err)
	}

	d.SetId(uuid.New().String())
	d.Set("keys", result.Keys)
	d.Set("keys_base64", result.KeysB64)
	d.Set("recovery_keys", result.RecoveryKeys)
	d.Set("recovery_keys_base64", result.RecoveryKeysB64)
	d.Set("root_token", result.RootToken)

	return nil
}

func sysInitDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func sysInitRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}
