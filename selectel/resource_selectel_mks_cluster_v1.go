package selectel

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/selectel/go-selvpcclient/v4/selvpcclient/quotamanager/quotas"
	"github.com/selectel/mks-go/pkg/v1/cluster"
)

func resourceMKSClusterV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMKSClusterV1Create,
		ReadContext:   resourceMKSClusterV1Read,
		UpdateContext: resourceMKSClusterV1Update,
		DeleteContext: resourceMKSClusterV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMKSClusterV1ImportState,
		},
		CustomizeDiff: customdiff.All(
			customdiff.ComputedIf(
				"maintenance_window_end",
				func(_ context.Context, d *schema.ResourceDiff, _ interface{}) bool {
					return d.HasChange("maintenance_window_start")
				}),
		),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(_, oldVersion, newVersion string, _ *schema.ResourceData) bool {
					return strings.EqualFold(oldVersion, newVersion)
				},
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"kube_version": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         false,
				DiffSuppressFunc: mksClusterV1KubeVersionDiffSuppressFunc,
				StateFunc: func(v interface{}) string {
					return strings.TrimPrefix(v.(string), "v")
				},
			},
			"enable_autorepair": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: false,
			},
			"enable_patch_version_auto_upgrade": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"enable_pod_security_policy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: false,
			},
			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"maintenance_window_start": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: false,
			},
			"zonal": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"maintenance_window_end": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kube_api_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"feature_gates": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"admission_controllers": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"private_kube_api": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"enable_audit_logs": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: false,
			},
			"oidc": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				MaxItems: 1,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"provider_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"issuer_url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"client_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"username_claim": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
							DiffSuppressFunc: func(_, oldVersion, newVersion string, _ *schema.ResourceData) bool {
								// Ignore diff on default value from API.
								return oldVersion == "sub" && newVersion == ""
							},
						},
						"groups_claim": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
							DiffSuppressFunc: func(_, oldVersion, newVersion string, _ *schema.ResourceData) bool {
								// Ignore diff on default value from API.
								return oldVersion == "groups" && newVersion == ""
							},
						},
					},
				},
			},
		},
	}
}

func resourceMKSClusterV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mksClient, diagErr := getMKSClient(d, meta)
	if diagErr != nil {
		return diagErr
	}

	config := meta.(*Config)
	projectID := d.Get("project_id").(string)
	region := d.Get("region").(string)

	selvpcClient, err := config.GetSelVPCClientWithProjectScope(projectID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("can't get project-scope selvpc for cluster object: %w", err))
	}
	err = validateRegion(selvpcClient, MKS, region)
	if err != nil {
		return diag.FromErr(fmt.Errorf("can't validate region: %w", err))
	}

	// Prepare cluster create options.
	enableAutorepair := d.Get("enable_autorepair").(bool)
	enablePodSecurityPolicy := d.Get("enable_pod_security_policy").(bool)
	zonal := d.Get("zonal").(bool)
	privateKubeAPI := d.Get("private_kube_api").(bool)
	enableAuditLogs := d.Get("enable_audit_logs").(bool)

	enablePatchVersionAutoUpgrade := !zonal // true by default only for regional clusters
	if v, ok := d.GetOk("enable_patch_version_auto_upgrade"); ok {
		enablePatchVersionAutoUpgrade = v.(bool)
	}

	// Check if "enable_patch_version_auto_upgrade" and "zonal" arguments are both not set to true.
	if enablePatchVersionAutoUpgrade && zonal {
		return diag.FromErr(errors.New("\"enable_patch_version_auto_upgrade\" argument should be explicitly " +
			"set to false in case of zonal cluster"))
	}

	featureGates, err := getSetAsStrings(d, featureGatesKey)
	if err != nil {
		return diag.FromErr(errCreatingObject(objectCluster, err))
	}

	admissionControllers, err := getSetAsStrings(d, admissionControllersKey)
	if err != nil {
		return diag.FromErr(errCreatingObject(objectCluster, err))
	}

	oidc, err := expandAndValidateMKSClusterV1OIDC(d)
	if err != nil {
		return diag.FromErr(err)
	}

	createOpts := &cluster.CreateOpts{
		Name:                          d.Get("name").(string),
		NetworkID:                     d.Get("network_id").(string),
		SubnetID:                      d.Get("subnet_id").(string),
		KubeVersion:                   d.Get("kube_version").(string),
		MaintenanceWindowStart:        d.Get("maintenance_window_start").(string),
		EnableAutorepair:              &enableAutorepair,
		EnablePatchVersionAutoUpgrade: &enablePatchVersionAutoUpgrade,
		Region:                        region,
		KubernetesOptions: &cluster.KubernetesOptions{
			EnablePodSecurityPolicy: enablePodSecurityPolicy,
			FeatureGates:            featureGates,
			AdmissionControllers:    admissionControllers,
			AuditLogs: cluster.AuditLogs{
				Enabled: enableAuditLogs,
			},
			OIDC: oidc,
		},
		Zonal:          &zonal,
		PrivateKubeAPI: &privateKubeAPI,
	}

	projectQuotas, _, err := quotas.GetProjectQuotas(selvpcClient, projectID, region)
	if err != nil {
		return diag.FromErr(errGettingObject(objectProjectQuotas, projectID, err))
	}

	if err := checkQuotasForCluster(projectQuotas, zonal); err != nil {
		return diag.FromErr(errCreatingObject(objectCluster, err))
	}

	log.Print(msgCreate(objectCluster, createOpts))
	newCluster, _, err := cluster.Create(ctx, mksClient, createOpts)
	if err != nil {
		return diag.FromErr(errCreatingObject(objectCluster, err))
	}

	log.Printf("[DEBUG] waiting for cluster %s to become 'ACTIVE'", newCluster.ID)
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForMKSClusterV1ActiveState(ctx, mksClient, newCluster.ID, timeout)
	if err != nil {
		return diag.FromErr(errCreatingObject(objectCluster, err))
	}

	d.SetId(newCluster.ID)

	return resourceMKSClusterV1Read(ctx, d, meta)
}

func resourceMKSClusterV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mksClient, diagErr := getMKSClient(d, meta)
	if diagErr != nil {
		return diagErr
	}

	log.Print(msgGet(objectCluster, d.Id()))
	mksCluster, response, err := cluster.Get(ctx, mksClient, d.Id())
	if err != nil {
		if response != nil {
			if response.StatusCode == http.StatusNotFound {
				d.SetId("")
				return nil
			}
		}

		return diag.FromErr(errGettingObject(objectCluster, d.Id(), err))
	}

	d.Set("name", mksCluster.Name)
	d.Set("status", mksCluster.Status)
	d.Set("project_id", mksCluster.ProjectID)
	d.Set("network_id", mksCluster.NetworkID)
	d.Set("subnet_id", mksCluster.SubnetID)
	d.Set("kube_api_ip", mksCluster.KubeAPIIP)
	d.Set("kube_version", mksCluster.KubeVersion)
	d.Set("region", mksCluster.Region)
	d.Set("maintenance_window_start", mksCluster.MaintenanceWindowStart)
	d.Set("maintenance_window_end", mksCluster.MaintenanceWindowEnd)
	d.Set("enable_autorepair", mksCluster.EnableAutorepair)
	d.Set("enable_patch_version_auto_upgrade", mksCluster.EnablePatchVersionAutoUpgrade)
	d.Set("enable_pod_security_policy", mksCluster.KubernetesOptions.EnablePodSecurityPolicy)
	d.Set("zonal", mksCluster.Zonal)
	d.Set("private_kube_api", mksCluster.PrivateKubeAPI)
	d.Set("enable_audit_logs", mksCluster.KubernetesOptions.AuditLogs.Enabled)
	d.Set("oidc", flattenMKSClusterV1OIDC(mksCluster))

	return nil
}

func resourceMKSClusterV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mksClient, diagErr := getMKSClient(d, meta)
	if diagErr != nil {
		return diagErr
	}

	if d.HasChange("kube_version") {
		if err := upgradeMKSClusterV1KubeVersion(ctx, d, mksClient); err != nil {
			return diag.FromErr(errUpdatingObject(objectCluster, d.Id(), err))
		}
	}

	var updateOpts cluster.UpdateOpts
	if d.HasChange("maintenance_window_start") {
		updateOpts.MaintenanceWindowStart = d.Get("maintenance_window_start").(string)
	}
	if d.HasChange("enable_autorepair") {
		v := d.Get("enable_autorepair").(bool)
		updateOpts.EnableAutorepair = &v
	}
	if d.HasChange("enable_patch_version_auto_upgrade") {
		v := d.Get("enable_patch_version_auto_upgrade").(bool)
		updateOpts.EnablePatchVersionAutoUpgrade = &v
	}

	kubeOptions := new(cluster.KubernetesOptions)
	if d.HasChange("enable_pod_security_policy") {
		v := d.Get("enable_pod_security_policy").(bool)
		kubeOptions.EnablePodSecurityPolicy = v
	}
	if d.HasChange(featureGatesKey) {
		v, err := getSetAsStrings(d, featureGatesKey)
		if err != nil {
			return diag.FromErr(errCreatingObject(objectCluster, err))
		}
		kubeOptions.FeatureGates = v
	}
	if d.HasChange(admissionControllersKey) {
		v, err := getSetAsStrings(d, admissionControllersKey)
		if err != nil {
			return diag.FromErr(errCreatingObject(objectCluster, err))
		}
		kubeOptions.AdmissionControllers = v
	}
	if d.HasChange("enable_audit_logs") {
		v := d.Get("enable_audit_logs").(bool)
		kubeOptions.AuditLogs.Enabled = v
	}
	if d.HasChange("oidc") {
		oidc, err := expandAndValidateMKSClusterV1OIDC(d)
		if err != nil {
			return diag.FromErr(err)
		}
		kubeOptions.OIDC = oidc
	}

	updateOpts.KubernetesOptions = kubeOptions

	if updateOpts != (cluster.UpdateOpts{}) {
		log.Print(msgUpdate(objectCluster, d.Id(), updateOpts))
		_, _, err := cluster.Update(ctx, mksClient, d.Id(), &updateOpts)
		if err != nil {
			return diag.FromErr(errUpdatingObject(objectCluster, d.Id(), err))
		}

		log.Printf("[DEBUG] waiting for cluster %s to become 'ACTIVE'", d.Id())
		timeout := d.Timeout(schema.TimeoutUpdate)
		err = waitForMKSClusterV1ActiveState(ctx, mksClient, d.Id(), timeout)
		if err != nil {
			return diag.FromErr(errUpdatingObject(objectCluster, d.Id(), err))
		}
	}

	return resourceMKSClusterV1Read(ctx, d, meta)
}

func resourceMKSClusterV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mksClient, diagErr := getMKSClient(d, meta)
	if diagErr != nil {
		return diagErr
	}

	log.Print(msgDelete(objectCluster, d.Id()))
	_, err := cluster.Delete(ctx, mksClient, d.Id())
	if err != nil {
		return diag.FromErr(errDeletingObject(objectCluster, d.Id(), err))
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{strconv.Itoa(http.StatusOK)},
		Target:  []string{strconv.Itoa(http.StatusNotFound)},
		Refresh: func() (result interface{}, state string, err error) {
			result, response, err := cluster.Get(ctx, mksClient, d.Id())
			if err != nil {
				if response != nil {
					return result, strconv.Itoa(response.StatusCode), nil
				}

				return nil, "", err
			}

			return result, strconv.Itoa(response.StatusCode), err
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	log.Printf("[DEBUG] waiting for cluster %s to become deleted", d.Id())
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error waiting for the cluster %s to become deleted: %s", d.Id(), err))
	}

	return nil
}

func resourceMKSClusterV1ImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if config.ProjectID == "" {
		return nil, errors.New("INFRA_PROJECT_ID must be set for the resource import")
	}
	if config.Region == "" {
		return nil, errors.New("INFRA_REGION must be set for the resource import")
	}

	d.Set("project_id", config.ProjectID)
	d.Set("region", config.Region)

	return []*schema.ResourceData{d}, nil
}
