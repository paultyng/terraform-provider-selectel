package cluster

import "github.com/selectel/mks-go/pkg/v1/nodegroup"

// CreateOpts represents options for the cluster Create request.
type CreateOpts struct {
	// Name represent the needed name of the cluster.
	// It shouldn't contain more than 32 symbols and can contain latin letters
	// with numbers and hyphens and start with a letter or a number.
	Name string `json:"name,omitempty"`

	// NetworkID contains a reference to the network of the cluster.
	// It can be set in cases where network is pre-created.
	NetworkID string `json:"network_id,omitempty"`

	// SubnetID contains a reference to the subnet of the cluster.
	// It can be set in cases where subnet is pre-created.
	SubnetID string `json:"subnet_id,omitempty"`

	// KubeVersion represents the needed Kubernetes version of the cluster.
	// It should be in x.y.z format.
	KubeVersion string `json:"kube_version,omitempty"`

	// Region represents the needed region.
	Region string `json:"region,omitempty"`

	// Nodegroups contains groups of nodes with their parameters.
	Nodegroups []*nodegroup.CreateOpts `json:"nodegroups,omitempty"`

	// AdditionalSoftware represents parameters of additional software that can be installed
	// in the Kubernetes cluster.
	AdditionalSoftware map[string]interface{} `json:"additional_software,omitempty"`

	// MaintenanceWindowStart represents UTC time of when the cluster will start its maintenance tasks.
	// It should be in hh:mm:ss format if provided.
	MaintenanceWindowStart string `json:"maintenance_window_start,omitempty"`

	// EnableAutorepair reflects if worker nodes are allowed to be reinstalled automatically
	// in case of their unavailability or unhealthiness. Enabled by default.
	EnableAutorepair *bool `json:"enable_autorepair,omitempty"`

	// EnablePatchVersionAutoUpgrade specifies if Kubernetes patch version of the cluster is allowed to be upgraded
	// automatically. Enabled by default.
	EnablePatchVersionAutoUpgrade *bool `json:"enable_patch_version_auto_upgrade,omitempty"`

	// Zonal specifies that only a single zonal master will be created.
	// It is needed if highly available control-plane is not required.
	Zonal *bool `json:"zonal,omitempty"`

	// KubernetesOptions represents additional k8s options such as pod security policy,
	// feature gates (Alpha stage only) and admission controllers.
	KubernetesOptions *KubernetesOptions `json:"kubernetes_options,omitempty"`
}

// UpdateOpts represents options for the cluster Update request.
type UpdateOpts struct {
	// MaintenanceWindowStart represents UTC time of when the cluster will start its maintenance tasks.
	// It should be in hh:mm:ss format if provided.
	MaintenanceWindowStart string `json:"maintenance_window_start,omitempty"`

	// EnableAutorepair reflects if worker nodes are allowed to be reinstalled automatically
	// in case of their unavailability or unhealthiness. Enabled by default.
	EnableAutorepair *bool `json:"enable_autorepair,omitempty"`

	// EnablePatchVersionAutoUpgrade specifies if Kubernetes patch version of the cluster is allowed to be upgraded
	// automatically. Enabled by default.
	EnablePatchVersionAutoUpgrade *bool `json:"enable_patch_version_auto_upgrade,omitempty"`

	// KubernetesOptions represents additional k8s options such as pod security policy,
	// feature gates (Alpha stage only) and admission controllers.
	KubernetesOptions *KubernetesOptions `json:"kubernetes_options,omitempty"`
}
