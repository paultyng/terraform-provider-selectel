package selectel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-selectel/selectel/internal/mutexkv"
)

const (
	objectFloatingIP                    = "floating IP"
	objectKeypair                       = "keypair"
	objectLicense                       = "license"
	objectProject                       = "project"
	objectProjectQuotas                 = "quotas for project"
	objectRole                          = "role"
	objectSubnet                        = "subnet"
	objectToken                         = "token"
	objectUser                          = "user"
	objectVRRPSubnet                    = "VRRP subnet"
	objectCluster                       = "cluster"
	objectNodegroup                     = "nodegroup"
	objectDomain                        = "domain"
	objectRecord                        = "record"
	objectDatastore                     = "datastore"
	objectDatabase                      = "database"
	objectGrant                         = "grant"
	objectExtension                     = "extension"
	objectDatastoreTypes                = "datastore-types"
	objectAvailableExtensions           = "available-extensions"
	objectFlavors                       = "flavors"
	objectConfigurationParameters       = "configuration-parameters"
	objectAvailableFeatureGates         = "feature-gates"
	objectAvailableAdmissionControllers = "admission-controllers"
)

// This is a global MutexKV for use within this plugin.
var selMutexKV = mutexkv.NewMutexKV()

// Provider returns the Selectel terraform provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SEL_TOKEN", nil),
				Description: "Token to authorize with the Selectel API.",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SEL_ENDPOINT", nil),
				Description: "Base endpoint to work with the Selectel API.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SEL_PROJECT_ID", nil),
				Description: "VPC project ID to import resources that need the project scope auth token.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SEL_REGION", nil),
				Description: "VPC region to import resources associated with the specific region.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"selectel_domains_domain_v1":                      dataSourceDomainsDomainV1(),
			"selectel_dbaas_datastore_type_v1":                dataSourceDBaaSDatastoreTypeV1(),
			"selectel_dbaas_available_extension_v1":           dataSourceDBaaSAvailableExtensionV1(),
			"selectel_dbaas_flavor_v1":                        dataSourceDBaaSFlavorV1(),
			"selectel_dbaas_configuration_parameter_v1":       dataSourceDBaaSConfigurationParameterV1(),
			"selectel_mks_available_feature_gates_v1":         dataSourceFeatureGateTypeV1(),
			"selectel_mks_available_admission_controllers_v1": dataSourceAdmissionControllersTypeV1(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"selectel_vpc_floatingip_v2":         resourceVPCFloatingIPV2(),
			"selectel_vpc_keypair_v2":            resourceVPCKeypairV2(),
			"selectel_vpc_license_v2":            resourceVPCLicenseV2(),
			"selectel_vpc_project_v2":            resourceVPCProjectV2(),
			"selectel_vpc_role_v2":               resourceVPCRoleV2(),
			"selectel_vpc_subnet_v2":             resourceVPCSubnetV2(),
			"selectel_vpc_token_v2":              resourceVPCTokenV2(),
			"selectel_vpc_user_v2":               resourceVPCUserV2(),
			"selectel_vpc_vrrp_subnet_v2":        resourceVPCVRRPSubnetV2(),
			"selectel_vpc_crossregion_subnet_v2": resourceVPCCrossRegionSubnetV2(),
			"selectel_mks_cluster_v1":            resourceMKSClusterV1(),
			"selectel_mks_nodegroup_v1":          resourceMKSNodegroupV1(),
			"selectel_domains_domain_v1":         resourceDomainsDomainV1(),
			"selectel_domains_record_v1":         resourceDomainsRecordV1(),
			"selectel_dbaas_datastore_v1":        resourceDBaaSDatastoreV1(),
			"selectel_dbaas_user_v1":             resourceDBaaSUserV1(),
			"selectel_dbaas_database_v1":         resourceDBaaSDatabaseV1(),
			"selectel_dbaas_grant_v1":            resourceDBaaSGrantV1(),
			"selectel_dbaas_extension_v1":        resourceDBaaSExtensionV1(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
		Token:    d.Get("token").(string),
		Endpoint: d.Get("endpoint").(string),
	}
	if v, ok := d.GetOk("project_id"); ok {
		config.ProjectID = v.(string)
	}
	if v, ok := d.GetOk("region"); ok {
		config.Region = v.(string)
	}
	if err := config.Validate(); err != nil {
		return nil, diag.FromErr(err)
	}

	return &config, nil
}
