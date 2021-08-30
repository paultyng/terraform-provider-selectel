package cluster

import (
	"encoding/json"
	"time"
)

// Status represents custom type for various cluster statuses.
type Status string

const (
	StatusActive                             Status = "ACTIVE"
	StatusPendingCreate                      Status = "PENDING_CREATE"
	StatusPendingUpdate                      Status = "PENDING_UPDATE"
	StatusPendingUpgrade                     Status = "PENDING_UPGRADE"
	StatusPendingRotateCerts                 Status = "PENDING_ROTATE_CERTS"
	StatusPendingDelete                      Status = "PENDING_DELETE"
	StatusPendingResize                      Status = "PENDING_RESIZE"
	StatusPendingNodeReinstall               Status = "PENDING_NODE_REINSTALL"
	StatusPendingUpgradePatchVersion         Status = "PENDING_UPGRADE_PATCH_VERSION"
	StatusPendingUpgradeMinorVersion         Status = "PENDING_UPGRADE_MINOR_VERSION"
	StatusPendingUpdateNodegroup             Status = "PENDING_UPDATE_NODEGROUP"
	StatusPendingUpgradeMastersConfiguration Status = "PENDING_UPGRADE_MASTERS_CONFIGURATION"
	StatusPendingUpgradeClusterConfiguration Status = "PENDING_UPGRADE_CLUSTER_CONFIGURATION"
	StatusMaintenance                        Status = "MAINTENANCE"
	StatusError                              Status = "ERROR"
	StatusUnknown                            Status = "UNKNOWN"
)

// View represents an unmarshalled cluster body from an API response.
type View struct {
	// ID is the identifier of the cluster.
	ID string `json:"id"`

	// CreatedAt is the timestamp in UTC timezone of when the cluster has been created.
	CreatedAt *time.Time `json:"created_at"`

	// UpdatedAt is the timestamp in UTC timezone of when the cluster has been updated.
	UpdatedAt *time.Time `json:"updated_at"`

	// Name represents the name of the cluster.
	Name string `json:"name"`

	// Status represents current status of the cluster.
	Status Status `json:"-"`

	// ProjectID contains reference to the project of the cluster.
	ProjectID string `json:"project_id"`

	// NetworkID contains reference to the network of the cluster.
	NetworkID string `json:"network_id"`

	// SubnetID contains reference to the subnet of the cluster.
	SubnetID string `json:"subnet_id"`

	// KubeAPIIP represents the IP of the Kubernetes API.
	KubeAPIIP string `json:"kube_api_ip"`

	// KubeVersion represents the current Kubernetes version of the cluster.
	KubeVersion string `json:"kube_version"`

	// Region represents the region of where the cluster is located.
	Region string `json:"region"`

	// AdditionalSoftware represents information about additional software installed in the cluster.
	AdditionalSoftware map[string]interface{} `json:"additional_software"`

	// PKITreeUpdatedAt represents the timestamp in UTC timezone of when the PKI-tree of the cluster
	// has been updated.
	PKITreeUpdatedAt *time.Time `json:"pki_tree_updated_at"`

	// MaintenanceWindowStart represents UTC time in "hh:mm:ss" format of when the cluster will start its
	// maintenance tasks.
	MaintenanceWindowStart string `json:"maintenance_window_start"`

	// MaintenanceWindowEnd represents UTC time in "hh:mm:ss" format of when the cluster will end its
	// maintenance tasks.
	MaintenanceWindowEnd string `json:"maintenance_window_end"`

	// MaintenanceLastStart is the timestamp in UTC timezone of the last cluster maintenance start.
	MaintenanceLastStart *time.Time `json:"maintenance_last_start"`

	// EnableAutorepair reflects if worker nodes are allowed to be reinstalled automatically
	// in case of their unavailability or unhealthiness.
	EnableAutorepair bool `json:"enable_autorepair"`

	// EnablePatchVersionAutoUpgrade specifies if Kubernetes patch version of the cluster is allowed to be upgraded
	// automatically.
	EnablePatchVersionAutoUpgrade bool `json:"enable_patch_version_auto_upgrade"`

	// Zonal specifies that cluster has only a single master and that
	// control-plane is not in highly available mode.
	Zonal bool `json:"zonal"`

	// KubernetesOptions represents additional k8s options such as pod security policy,
	// feature gates (Alpha stage only) and admission controllers.
	KubernetesOptions *KubernetesOptions `json:"kubernetes_options,omitempty"`
}

func (result *View) UnmarshalJSON(b []byte) error {
	type tmp View
	var s struct {
		tmp
		Status Status `json:"status"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*result = View(s.tmp)

	// Check cluster status.
	switch s.Status {
	case StatusActive:
		result.Status = StatusActive
	case StatusPendingCreate:
		result.Status = StatusPendingCreate
	case StatusPendingUpdate:
		result.Status = StatusPendingUpdate
	case StatusPendingUpgrade:
		result.Status = StatusPendingUpgrade
	case StatusPendingRotateCerts:
		result.Status = StatusPendingRotateCerts
	case StatusPendingDelete:
		result.Status = StatusPendingDelete
	case StatusPendingResize:
		result.Status = StatusPendingResize
	case StatusPendingNodeReinstall:
		result.Status = StatusPendingNodeReinstall
	case StatusPendingUpgradePatchVersion:
		result.Status = StatusPendingUpgradePatchVersion
	case StatusPendingUpgradeMinorVersion:
		result.Status = StatusPendingUpgradeMinorVersion
	case StatusPendingUpdateNodegroup:
		result.Status = StatusPendingUpdateNodegroup
	case StatusPendingUpgradeMastersConfiguration:
		result.Status = StatusPendingUpgradeMastersConfiguration
	case StatusPendingUpgradeClusterConfiguration:
		result.Status = StatusPendingUpgradeClusterConfiguration
	case StatusMaintenance:
		result.Status = StatusMaintenance
	case StatusError:
		result.Status = StatusError
	default:
		result.Status = StatusUnknown
	}

	return err
}

// KubernetesOptions represents additional k8s options such as pod security policy,
// feature gates (Alpha stage only) and admission controllers.
type KubernetesOptions struct {
	// EnablePodSecurityPolicy indicates if PodSecurityPolicy admission controller
	// must be turned on/off.
	EnablePodSecurityPolicy bool `json:"enable_pod_security_policy"`

	// FeatureGates represents feature gates that should be enabled.
	FeatureGates []string `json:"feature_gates"`

	// AdmissionControllers represents admission controllers that should be enabled.
	AdmissionControllers []string `json:"admission_controllers"`
}
