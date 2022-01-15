package selectel

import (
	"context"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/selectel/go-selvpcclient/selvpcclient/resell/v2/users"
)

func resourceVPCUserV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPCUserV2Create,
		ReadContext:   resourceVPCUserV2Read,
		UpdateContext: resourceVPCUserV2Update,
		DeleteContext: resourceVPCUserV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: false,
			},
		},
	}
}

func resourceVPCUserV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	resellV2Client := config.resellV2Client()

	opts := users.UserOpts{
		Name:     d.Get("name").(string),
		Password: d.Get("password").(string),
	}

	log.Print(msgCreate(objectUser, opts))
	user, _, err := users.Create(ctx, resellV2Client, opts)
	if err != nil {
		return diag.FromErr(errCreatingObject(objectUser, err))
	}

	d.SetId(user.ID)

	return resourceVPCUserV2Read(ctx, d, meta)
}

func resourceVPCUserV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	resellV2Client := config.resellV2Client()

	log.Print(msgGet(objectUser, d.Id()))
	user, response, err := users.Get(ctx, resellV2Client, d.Id())
	if err != nil {
		if response != nil {
			if response.StatusCode == http.StatusNotFound {
				d.SetId("")
				return nil
			}
		}

		return diag.FromErr(errGettingObject(objectUser, d.Id(), err))
	}

	d.Set("name", user.Name)
	d.Set("enabled", user.Enabled)

	return nil
}

func resourceVPCUserV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	resellV2Client := config.resellV2Client()

	enabled := d.Get("enabled").(bool)
	opts := users.UserOpts{
		Name:     d.Get("name").(string),
		Password: d.Get("password").(string),
		Enabled:  &enabled,
	}

	log.Print(msgUpdate(objectUser, d.Id(), opts))
	_, _, err := users.Update(ctx, resellV2Client, d.Id(), opts)
	if err != nil {
		return diag.FromErr(errUpdatingObject(objectUser, d.Id(), err))
	}

	return resourceVPCUserV2Read(ctx, d, meta)
}

func resourceVPCUserV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	resellV2Client := config.resellV2Client()

	log.Print(msgDelete(objectUser, d.Id()))
	response, err := users.Delete(ctx, resellV2Client, d.Id())
	if err != nil {
		if response != nil {
			if response.StatusCode == http.StatusNotFound {
				d.SetId("")
				return nil
			}
		}

		return diag.FromErr(errDeletingObject(objectUser, d.Id(), err))
	}

	return nil
}
