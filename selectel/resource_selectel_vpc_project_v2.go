package selectel

import (
	"context"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/selectel/go-selvpcclient/selvpcclient/resell/v2/projects"
	"github.com/selectel/go-selvpcclient/selvpcclient/resell/v2/quotas"
)

func resourceVPCProjectV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPCProjectV2Create,
		ReadContext:   resourceVPCProjectV2Read,
		UpdateContext: resourceVPCProjectV2Update,
		DeleteContext: resourceVPCProjectV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"custom_url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"theme": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"auto_quotas": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"quotas": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},
						"resource_quotas": {
							Type:     schema.TypeSet,
							Required: true,
							ForceNew: false,
							Set:      hashResourceQuotas,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeInt,
										Required: true,
										ForceNew: false,
									},
									"region": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: false,
									},
									"zone": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: false,
									},
								},
							},
						},
					},
				},
			},
			"all_quotas": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_quotas": {
							Type:     schema.TypeSet,
							Computed: true,
							Set:      hashResourceQuotas,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"zone": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"used": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceVPCProjectV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	resellV2Client := config.resellV2Client()

	var opts projects.CreateOpts
	quotaSet := d.Get("quotas").(*schema.Set)
	if quotaSet.Len() != 0 {
		quotasOpts, err := resourceVPCProjectV2QuotasOptsFromSet(quotaSet)
		if err != nil {
			return diag.FromErr(errParseProjectV2Quotas(err))
		}
		opts.Quotas = quotasOpts
	}
	opts.Name = d.Get("name").(string)
	opts.AutoQuotas = d.Get("auto_quotas").(bool)

	log.Print(msgCreate(objectProject, opts))
	project, _, err := projects.Create(ctx, resellV2Client, opts)
	if err != nil {
		return diag.FromErr(errCreatingObject(objectProject, err))
	}

	d.SetId(project.ID)

	return resourceVPCProjectV2Read(ctx, d, meta)
}

func resourceVPCProjectV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	resellV2Client := config.resellV2Client()

	log.Print(msgGet(objectProject, d.Id()))
	project, response, err := projects.Get(ctx, resellV2Client, d.Id())
	if err != nil {
		if response != nil {
			if response.StatusCode == http.StatusNotFound {
				d.SetId("")
				return nil
			}
		}

		return diag.FromErr(errGettingObject(objectProject, d.Id(), err))
	}

	projectCustomURL, err := resourceVPCProjectV2URLWithoutSchema(project.CustomURL)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("custom_url", projectCustomURL)
	d.Set("name", project.Name)
	d.Set("url", project.URL)
	d.Set("enabled", project.Enabled)
	if err := d.Set("theme", flattenVPCProjectV2Theme(project.Theme)); err != nil {
		log.Print(errSettingComplexAttr("theme", err))
	}

	// Set all quotas. This can be different from what the user specified since
	// the project will have all available resource quotas automatically applied.
	allQuotas := resourceVPCProjectV2QuotasToSet(project.Quotas)
	if err := d.Set("all_quotas", allQuotas); err != nil {
		log.Print(errSettingComplexAttr("all_quotas", err))
	}

	return nil
}

func resourceVPCProjectV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	resellV2Client := config.resellV2Client()

	var hasChange, projectChange, quotaChange bool
	var projectOpts projects.UpdateOpts
	var projectQuotasOpts quotas.UpdateProjectQuotasOpts

	if d.HasChange("name") {
		hasChange, projectChange = true, true
		projectOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("custom_url") {
		hasChange, projectChange = true, true
		customURL := d.Get("custom_url").(string)
		projectOpts.CustomURL = &customURL
	}
	if d.HasChange("theme") {
		hasChange, projectChange = true, true
		themeMap := d.Get("theme").(map[string]interface{})
		updateThemeOpts := resourceProjectV2UpdateThemeOptsFromMap(themeMap)
		projectOpts.Theme = updateThemeOpts
	}
	if d.HasChange("quotas") {
		hasChange, quotaChange = true, true
		quotaSet := d.Get("quotas").(*schema.Set)
		quotasOpts, err := resourceVPCProjectV2QuotasOptsFromSet(quotaSet)
		if err != nil {
			return diag.FromErr(errParseProjectV2Quotas(err))
		}
		projectQuotasOpts.QuotasOpts = quotasOpts
	}

	if hasChange {
		// Update project options if needed.
		if projectChange {
			log.Print(msgUpdate(objectProject, d.Id(), projectOpts))
			_, _, err := projects.Update(ctx, resellV2Client, d.Id(), projectOpts)
			if err != nil {
				return diag.FromErr(errUpdatingObject(objectProject, d.Id(), err))
			}
		}
		// Update project quotas if needed.
		if quotaChange {
			log.Print(msgUpdate(objectProjectQuotas, d.Id(), projectQuotasOpts))
			_, _, err := quotas.UpdateProjectQuotas(ctx, resellV2Client, d.Id(), projectQuotasOpts)
			if err != nil {
				return diag.FromErr(errUpdatingObject(objectProjectQuotas, d.Id(), err))
			}
		}
	}

	return resourceVPCProjectV2Read(ctx, d, meta)
}

func resourceVPCProjectV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	resellV2Client := config.resellV2Client()

	log.Print(msgDelete(objectProject, d.Id()))
	response, err := projects.Delete(ctx, resellV2Client, d.Id())
	if err != nil {
		if response != nil {
			if response.StatusCode == http.StatusNotFound {
				d.SetId("")
				return nil
			}
		}

		return diag.FromErr(errDeletingObject(objectProject, d.Id(), err))
	}

	return nil
}
