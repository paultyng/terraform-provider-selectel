package selectel

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/selectel/go-selvpcclient/selvpcclient/resell/v2/tokens"
	v1 "github.com/selectel/mks-go/pkg/v1"
	"github.com/selectel/mks-go/pkg/v1/cluster"
	"github.com/selectel/mks-go/pkg/v1/kubeoptions"
	"github.com/selectel/mks-go/pkg/v1/kubeversion"
	"github.com/selectel/mks-go/pkg/v1/node"
	"github.com/selectel/mks-go/pkg/v1/nodegroup"
)

const (
	ru1MKSClusterV1Endpoint = "https://ru-1.mks.selcloud.ru/v1"
	ru2MKSClusterV1Endpoint = "https://ru-2.mks.selcloud.ru/v1"
	ru3MKSClusterV1Endpoint = "https://ru-3.mks.selcloud.ru/v1"
	ru7MKSClusterV1Endpoint = "https://ru-7.mks.selcloud.ru/v1"
	ru8MKSClusterV1Endpoint = "https://ru-8.mks.selcloud.ru/v1"
	ru9MKSClusterV1Endpoint = "https://ru-9.mks.selcloud.ru/v1"
)

func getMKSClusterV1Endpoint(region string) (endpoint string) {
	switch region {
	case ru1Region:
		endpoint = ru1MKSClusterV1Endpoint
	case ru2Region:
		endpoint = ru2MKSClusterV1Endpoint
	case ru3Region:
		endpoint = ru3MKSClusterV1Endpoint
	case ru7Region:
		endpoint = ru7MKSClusterV1Endpoint
	case ru8Region:
		endpoint = ru8MKSClusterV1Endpoint
	case ru9Region:
		endpoint = ru9MKSClusterV1Endpoint
	}

	return
}

func waitForMKSClusterV1ActiveState(
	ctx context.Context, client *v1.ServiceClient, clusterID string, timeout time.Duration) error {
	pending := []string{
		string(cluster.StatusPendingCreate),
		string(cluster.StatusPendingUpdate),
		string(cluster.StatusPendingUpgradePatchVersion),
		string(cluster.StatusPendingUpgradeMinorVersion),
		string(cluster.StatusPendingResize),
	}
	target := []string{
		string(cluster.StatusActive),
	}

	stateConf := &resource.StateChangeConf{
		Pending:    pending,
		Target:     target,
		Refresh:    mksClusterV1StateRefreshFunc(ctx, client, clusterID),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"error waiting for the cluster %s to become 'ACTIVE': %s",
			clusterID, err)
	}

	return nil
}

func mksClusterV1StateRefreshFunc(
	ctx context.Context, client *v1.ServiceClient, clusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		c, _, err := cluster.Get(ctx, client, clusterID)
		if err != nil {
			return nil, "", err
		}

		return c, string(c.Status), nil
	}
}

func mksClusterV1KubeVersionDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Id() == "" {
		return false
	}

	currentMajor, err := kubeVersionToMajor(old)
	if err != nil {
		log.Printf("[DEBUG] error getting a major part of the current kube version %s: %s", old, err)
		return false
	}
	desiredMajor, err := kubeVersionToMajor(new)
	if err != nil {
		log.Printf("[DEBUG] error getting a major part of the desired kube version %s: %s", new, err)
		return false
	}

	// If the desired major version is newer than current, do not suppress diff.
	if desiredMajor > currentMajor {
		return false
	}

	// If the desired major version is less than current, suppress diff.
	if desiredMajor < currentMajor {
		return true
	}

	currentMinor, err := kubeVersionToMinor(old)
	if err != nil {
		log.Printf("[DEBUG] error getting a minor part of the current kube version %s: %s", old, err)
		return false
	}
	desiredMinor, err := kubeVersionToMinor(new)
	if err != nil {
		log.Printf("[DEBUG] error getting a minor part of the desired kube version %s: %s", new, err)
		return false
	}

	// If the desired minor version is newer than current, do not suppress diff.
	if desiredMinor > currentMinor {
		return false
	}

	// If the desired minor version is less than current, suppress diff.
	if desiredMinor < currentMinor {
		return true
	}

	currentPatch, err := kubeVersionToPatch(old)
	if err != nil {
		log.Printf("[DEBUG] error getting a patch part of the current kube version %s: %s", old, err)
		return false
	}
	desiredPatch, err := kubeVersionToPatch(new)
	if err != nil {
		log.Printf("[DEBUG] error getting a patch part of the desired kube version %s: %s", new, err)
		return false
	}

	// If the desired patch version is less than current, suppress diff.
	if desiredPatch < currentPatch {
		return true
	}

	return false
}

func mksClusterV1GetLatestPatchVersions(ctx context.Context, client *v1.ServiceClient) (map[string]string, error) {
	kubeVersions, _, err := kubeversion.List(ctx, client)
	if err != nil {
		return nil, err
	}

	result := map[string]string{}

	for _, version := range kubeVersions {
		minor, err := kubeVersionTrimToMinor(version.Version)
		if err != nil {
			return nil, err
		}

		current, ok := result[minor]
		if !ok {
			result[minor] = version.Version
			continue
		}

		latest, err := compareTwoKubeVersionsByPatch(version.Version, current)
		if err != nil {
			return nil, err
		}

		result[minor] = latest
	}

	return result, nil
}

func upgradeMKSClusterV1KubeVersion(ctx context.Context, d *schema.ResourceData, client *v1.ServiceClient) error {
	o, n := d.GetChange("kube_version")
	currentVersion := o.(string)
	desiredVersion := n.(string)

	log.Printf("[DEBUG] current kube version: %s", currentVersion)
	log.Printf("[DEBUG] desired kube version: %s", desiredVersion)

	// Compare current and desired major versions.
	currentMajor, err := kubeVersionToMajor(currentVersion)
	if err != nil {
		return fmt.Errorf("error getting a major part of the current version %s: %s", currentVersion, err)
	}
	desiredMajor, err := kubeVersionToMajor(desiredVersion)
	if err != nil {
		return fmt.Errorf("error getting a major part of the desired version %s: %s", desiredVersion, err)
	}
	if desiredMajor != currentMajor {
		return fmt.Errorf("current version %s can't be upgraded to version %s", currentVersion, desiredVersion)
	}

	// Compare current and desired minor versions.
	currentMinor, err := kubeVersionTrimToMinor(currentVersion)
	if err != nil {
		return fmt.Errorf("error getting a minor part of the current version %s: %s", currentVersion, err)
	}
	desiredMinor, err := kubeVersionTrimToMinor(desiredVersion)
	if err != nil {
		return fmt.Errorf("error getting a minor part of the desired version %s: %s", desiredVersion, err)
	}
	if desiredMinor != currentMinor {
		log.Print("[DEBUG] upgrading minor version")

		// Increment minor version.
		currentMinorNew, err := kubeVersionTrimToMinorIncremented(currentVersion)
		if err != nil {
			return fmt.Errorf("error getting incremented minor part of the current version %s: %s", currentVersion, err)
		}

		// Get latest patch versions for every minor version.
		latestPatchVersions, err := mksClusterV1GetLatestPatchVersions(ctx, client)
		if err != nil {
			return fmt.Errorf("error getting latest patch versions: %s", err)
		}

		// Check that we have a Kubernetes version of the current minor version + 1.
		latestVersion, ok := latestPatchVersions[currentMinorNew]
		if !ok {
			return fmt.Errorf("the cluster is already on the latest available minor version: %s", currentMinor)
		}

		log.Printf("[DEBUG] latest kube version: %s", latestVersion)

		// Compare the latest patch version with the desired version.
		if desiredVersion != latestVersion {
			return fmt.Errorf(
				"current version %s can't be upgraded to version %s, the latest available version is: %s",
				currentVersion, desiredVersion, latestVersion)
		}

		_, _, err = cluster.UpgradeMinorVersion(ctx, client, d.Id())
		if err != nil {
			return fmt.Errorf("error upgrading minor version: %s", err)
		}

		log.Printf("[DEBUG] waiting for cluster %s to become 'ACTIVE'", d.Id())
		timeout := d.Timeout(schema.TimeoutUpdate)
		err = waitForMKSClusterV1ActiveState(ctx, client, d.Id(), timeout)
		if err != nil {
			return fmt.Errorf("error waiting for the minor version upgrade: %s", err)
		}

		return nil
	}

	log.Print("[DEBUG] upgrading patch version")

	// Get latest patch versions for every minor version.
	latestPatchVersions, err := mksClusterV1GetLatestPatchVersions(ctx, client)
	if err != nil {
		return fmt.Errorf("error getting latest patch versions: %s", err)
	}

	// Find the latest patch version corresponding to the current minor version.
	latestVersion, ok := latestPatchVersions[currentMinor]
	if !ok {
		return fmt.Errorf("unable to find the latest patch version for the current minor version %s", currentMinor)
	}

	log.Printf("[DEBUG] latest kube version: %s", latestVersion)

	if desiredVersion != latestVersion {
		return fmt.Errorf(
			"current version %s can't be upgraded to version %s, the latest available patch version is: %s",
			currentVersion, desiredVersion, latestVersion)
	}

	_, _, err = cluster.UpgradePatchVersion(ctx, client, d.Id())
	if err != nil {
		return fmt.Errorf("error upgrading patch version: %s", err)
	}

	log.Printf("[DEBUG] waiting for cluster %s to become 'ACTIVE'", d.Id())
	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waitForMKSClusterV1ActiveState(ctx, client, d.Id(), timeout)
	if err != nil {
		return fmt.Errorf("error waiting for the patch version upgrade: %s", err)
	}

	return nil
}

// kubeVersionToMajor returns given Kubernetes version major part.
func kubeVersionToMajor(kubeVersion string) (int, error) {
	// Trim version prefix if needed.
	kubeVersion = strings.TrimPrefix(kubeVersion, "v")

	kubeVersionParts := strings.Split(kubeVersion, ".")
	majorPart := kubeVersionParts[0]
	major, err := strconv.Atoi(majorPart)
	if err != nil {
		return 0, errKubeVersionIsInvalidFmt(kubeVersion, "major part is not an integer number")
	}
	if major < 0 {
		return 0, errKubeVersionIsInvalidFmt(kubeVersion, "major part is a negative number")
	}

	return major, nil
}

// kubeVersionToMinor returns given Kubernetes version minor part.
func kubeVersionToMinor(kubeVersion string) (int, error) {
	// Trim version prefix if needed.
	kubeVersion = strings.TrimPrefix(kubeVersion, "v")

	kubeVersionParts := strings.Split(kubeVersion, ".")
	if len(kubeVersionParts) < 2 {
		return 0, errKubeVersionIsInvalidFmt(kubeVersion, "expected to have major and minor version parts")
	}

	minorPart := kubeVersionParts[1]
	minor, err := strconv.Atoi(minorPart)
	if err != nil {
		return 0, errKubeVersionIsInvalidFmt(kubeVersion, "minor part is not an integer number")
	}
	if minor < 0 {
		return 0, errKubeVersionIsInvalidFmt(kubeVersion, "minor part is a negative number")
	}

	return minor, nil
}

// kubeVersionToPatch returns given Kubernetes version patch part.
func kubeVersionToPatch(kubeVersion string) (int, error) {
	// Trim version prefix if needed.
	kubeVersion = strings.TrimPrefix(kubeVersion, "v")

	kubeVersionParts := strings.Split(kubeVersion, ".")
	if len(kubeVersionParts) < 3 {
		return 0, errKubeVersionIsInvalidFmt(kubeVersion, "expected to have major, minor and patch version parts")
	}

	patchPart := kubeVersionParts[2]
	patch, err := strconv.Atoi(patchPart)
	if err != nil {
		return 0, errKubeVersionIsInvalidFmt(kubeVersion, "patch part is not an integer number")
	}
	if patch < 0 {
		return 0, errKubeVersionIsInvalidFmt(kubeVersion, "patch part is a negative number")
	}

	return patch, nil
}

// compareTwoKubeVersionsByPatch parses two Kubernetes versions, compares their patch versions and returns
// the latest version.
// It doesn't check minor version so it will give bad result in case of different minor versions.
func compareTwoKubeVersionsByPatch(a, b string) (string, error) {
	aPatch, err := kubeVersionToPatch(a)
	if err != nil {
		return "", fmt.Errorf("unable to compare kube versions: %s", err)
	}

	bPatch, err := kubeVersionToPatch(b)
	if err != nil {
		return "", fmt.Errorf("unable to compare kube versions: %s", err)
	}

	if aPatch > bPatch {
		return a, nil
	}

	return b, nil
}

// kubeVersionTrimToMinor returns given Kubernetes version trimmed to minor.
func kubeVersionTrimToMinor(kubeVersion string) (string, error) {
	major, err := kubeVersionToMajor(kubeVersion)
	if err != nil {
		return "", err
	}

	minor, err := kubeVersionToMinor(kubeVersion)
	if err != nil {
		return "", err
	}

	return strings.Join([]string{strconv.Itoa(major), strconv.Itoa(minor)}, "."), nil
}

// KubeVersionTrimToMinorIncremented returns given Kubernetes version trimmed to minor incremented by 1.
func kubeVersionTrimToMinorIncremented(kubeVersion string) (string, error) {
	major, err := kubeVersionToMajor(kubeVersion)
	if err != nil {
		return "", err
	}

	minor, err := kubeVersionToMinor(kubeVersion)
	if err != nil {
		return "", err
	}

	// Increment minor version.
	minor++

	return strings.Join([]string{strconv.Itoa(major), strconv.Itoa(minor)}, "."), nil
}

func mksNodegroupV1ParseID(id string) (string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		return "", "", errParseMKSNodegroupV1ID(id)
	}
	if parts[0] == "" || parts[1] == "" {
		return "", "", errParseMKSNodegroupV1ID(id)
	}

	return parts[0], parts[1], nil
}

func flattenMKSNodegroupV1Nodes(views []*node.View) []map[string]interface{} {
	nodes := make([]map[string]interface{}, len(views))
	for i, view := range views {
		nodes[i] = make(map[string]interface{})
		nodes[i]["id"] = view.ID
		nodes[i]["ip"] = view.IP
		nodes[i]["hostname"] = view.Hostname
	}

	return nodes
}

func flattenMKSNodegroupV1Taints(views []nodegroup.Taint) []interface{} {
	taints := make([]interface{}, len(views))
	for i, view := range views {
		taints[i] = map[string]interface{}{
			"key":    view.Key,
			"value":  view.Value,
			"effect": string(view.Effect),
		}
	}

	return taints
}

func flattenFeatureGates(views []*kubeoptions.View) []interface{} {
	availableFeatureGates := make([]interface{}, len(views))
	for i, fg := range views {
		availableFeatureGates[i] = map[string]interface{}{
			"kube_version_minor": fg.KubeVersion,
			"names":              strings.Join(fg.Names, ","),
		}
	}

	return availableFeatureGates
}

func flattenFeatureGatesFromSlice(kubeVersion string, featureGates []string) []interface{} {
	availableFeatureGates := make([]interface{}, 1)
	availableFeatureGates[0] = map[string]interface{}{
		"kube_version_minor": kubeVersion,
		"names":              strings.Join(featureGates, ","),
	}

	return availableFeatureGates
}

func flattenAdmissionControllers(views []*kubeoptions.View) []interface{} {
	availableAdmissionControllers := make([]interface{}, len(views))
	for i, fg := range views {
		availableAdmissionControllers[i] = map[string]interface{}{
			"kube_version_minor": fg.KubeVersion,
			"names":              strings.Join(fg.Names, ","),
		}
	}

	return availableAdmissionControllers
}

func flattenAdmissionControllersFromSlice(kubeVersion string, admissionControllers []string) []interface{} {
	availableAdmissionControllers := make([]interface{}, 1)
	availableAdmissionControllers[0] = map[string]interface{}{
		"kube_version_minor": kubeVersion,
		"names":              strings.Join(admissionControllers, ","),
	}

	return availableAdmissionControllers
}

func expandMKSNodegroupV1Taints(taints []interface{}) []nodegroup.Taint {
	result := make([]nodegroup.Taint, len(taints))
	for i := range taints {
		taint := nodegroup.Taint{}
		obj := taints[i].(map[string]interface{})
		taint.Key = obj["key"].(string)
		taint.Value = obj["value"].(string)

		switch obj["effect"].(string) {
		case string(nodegroup.NoScheduleEffect):
			taint.Effect = nodegroup.NoScheduleEffect
		case string(nodegroup.NoExecuteEffect):
			taint.Effect = nodegroup.NoExecuteEffect
		case string(nodegroup.PreferNoScheduleEffect):
			taint.Effect = nodegroup.PreferNoScheduleEffect
		}

		result[i] = taint
	}

	return result
}

func expandMKSNodegroupV1Labels(labels map[string]interface{}) map[string]string {
	result := make(map[string]string)

	for k, v := range labels {
		result[k] = v.(string)
	}

	return result
}

func getMKSClient(ctx context.Context, d *schema.ResourceData, meta interface{}) (*v1.ServiceClient, diag.Diagnostics) {
	config := meta.(*Config)
	resellV2Client := config.resellV2Client()
	tokenOpts := tokens.TokenOpts{
		ProjectID: d.Get("project_id").(string),
	}

	token, _, err := tokens.Create(ctx, resellV2Client, tokenOpts)
	if err != nil {
		return nil, diag.FromErr(errCreatingObject(objectToken, err))
	}

	region := d.Get("region").(string)
	endpoint := getMKSClusterV1Endpoint(region)
	mksClient := v1.NewMKSClientV1(token.ID, endpoint)

	return mksClient, nil
}
