package selectel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/selectel/mks-go/pkg/v1/kubeoptions"
)

type availableKubeOptionsSearchFilter struct {
	kubeVersion string
}

func dataSourceFeatureGateTypeV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFeatureGateTypeV1Read,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					ru1Region,
					ru2Region,
					ru3Region,
					ru7Region,
					ru8Region,
					ru9Region,
				}, false),
			},
			"filter": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kube_version_minor": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"available_feature_gates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kube_version_minor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"names": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceFeatureGateTypeV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mksClient, diagErr := getMKSClient(ctx, d, meta)
	if diagErr != nil {
		return diagErr
	}

	featureGates, _, err := kubeoptions.ListFeatureGates(ctx, mksClient)
	if diagErr != nil {
		return diag.FromErr(errGettingObjects(objectAvailableFeatureGates, err))
	}

	filterSet := d.Get("filter").(*schema.Set)
	if filterSet.Len() == 0 {
		return setAllAvailableFeatureGates(d, featureGates)
	}

	filterMap := filterSet.List()[0].(map[string]interface{})
	kubeVersionVal, ok := filterMap["kube_version_minor"]
	if !ok {
		return diag.Errorf("kube_version_minor is not set: %v", kubeVersionVal)
	}
	kubeVersion := kubeVersionVal.(string)

	if kubeVersion == "" {
		return diag.Errorf("kubernetes version is not set")
	}
	kubeMinorVersion, err := kubeVersionTrimToMinor(kubeVersion)
	if err != nil {
		return diag.FromErr(err)
	}

	availableFeatureGates, err := filterKubeOptionsByKubeVersion(featureGates, kubeMinorVersion)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(kubeVersion)
	flatFG := flattenFeatureGatesFromSlice(kubeVersion, availableFeatureGates)
	if err := d.Set("available_feature_gates", flatFG); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setAllAvailableFeatureGates(d *schema.ResourceData, availableFeatureGates []*kubeoptions.View) diag.Diagnostics {
	d.SetId("mks_available_feature_gates")
	flatFG := flattenFeatureGates(availableFeatureGates)
	if err := d.Set("available_feature_gates", flatFG); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func filterKubeOptionsByKubeVersion(options []*kubeoptions.View, version string) ([]string, error) {
	if version == "" {
		return nil, fmt.Errorf("kubernetes version is not set")
	}

	for _, view := range options {
		if view.KubeVersion == version {
			return view.Names, nil
		}
	}

	return nil, fmt.Errorf("available kubernetes options for kubernetes version %q is not found", version)
}

func dataSourceAdmissionControllersTypeV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAdmissionControllersTypeV1Read,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					ru1Region,
					ru2Region,
					ru3Region,
					ru7Region,
					ru8Region,
					ru9Region,
				}, false),
			},
			"filter": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kube_version_minor": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"available_admission_controllers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kube_version_minor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"names": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAdmissionControllersTypeV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mksClient, diagErr := getMKSClient(ctx, d, meta)
	if diagErr != nil {
		return diagErr
	}

	admissionControllers, _, err := kubeoptions.ListAdmissionControllers(ctx, mksClient)
	if diagErr != nil {
		return diag.FromErr(errGettingObjects(objectAvailableAdmissionControllers, err))
	}

	filterSet := d.Get("filter").(*schema.Set)
	if filterSet.Len() == 0 {
		return setAllAvailableAdmissionControllers(d, admissionControllers)
	}

	filterMap := filterSet.List()[0].(map[string]interface{})
	kubeVersionVal, ok := filterMap["kube_version_minor"]
	if !ok {
		return diag.Errorf("kube_version_minor is not set: %v", kubeVersionVal)
	}
	kubeVersion := kubeVersionVal.(string)

	if kubeVersion == "" {
		return diag.Errorf("kubernetes version is not set")
	}
	kubeMinorVersion, err := kubeVersionTrimToMinor(kubeVersion)
	if err != nil {
		return diag.FromErr(err)
	}

	availableAdmissionControllers, err := filterKubeOptionsByKubeVersion(admissionControllers, kubeMinorVersion)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(kubeVersion)
	flatFG := flattenAdmissionControllersFromSlice(kubeVersion, availableAdmissionControllers)
	if err := d.Set("available_admission_controllers", flatFG); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setAllAvailableAdmissionControllers(d *schema.ResourceData, availableAdmissionControllers []*kubeoptions.View) diag.Diagnostics {
	d.SetId("mks_available_admission_controllers")
	flatFG := flattenAdmissionControllers(availableAdmissionControllers)
	if err := d.Set("available_admission_controllers", flatFG); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
