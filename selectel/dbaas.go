package selectel

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/selectel/dbaas-go"
	waiters "github.com/terraform-providers/terraform-provider-selectel/selectel/waiters/dbaas"
)

const (
	postgreSQLDatastoreType  = "postgresql"
	mySQLDatastoreType       = "mysql"
	mySQLNativeDatastoreType = "mysql_native"
	redisDatastoreType       = "redis"
	kafkaDatastoreType       = "kafka"
	masterRole               = "MASTER"
	replicaRole              = "REPLICA"
)

func getDBaaSClient(d *schema.ResourceData, meta interface{}) (*dbaas.API, diag.Diagnostics) {
	config := meta.(*Config)
	projectID := d.Get("project_id").(string)
	region := d.Get("region").(string)

	selvpcClient, err := config.GetSelVPCClientWithProjectScope(projectID)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("can't get project-scope selvpc client for dbaas: %w", err))
	}

	err = validateRegion(selvpcClient, DBaaS, region)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("can't validate region: %w", err))
	}

	endpoint, err := selvpcClient.Catalog.GetEndpoint(DBaaS, region)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("can't get endpoint to init dbaas client: %w", err))
	}

	client, err := dbaas.NewDBAASClient(selvpcClient.GetXAuthToken(), endpoint.URL)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("can't create dbaas client: %w", err))
	}

	return client, nil
}

func stringChecksum(s string) (string, error) {
	// #nosec G401
	h := md5.New() // nosemgrep: go.lang.security.audit.crypto.use_of_weak_crypto.use-of-md5
	_, err := h.Write([]byte(s))
	if err != nil {
		return "", err
	}
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs), nil
}

func stringListChecksum(s []string) (string, error) {
	sort.Strings(s)
	checksum, err := stringChecksum(strings.Join(s, ""))
	if err != nil {
		return "", err
	}

	return checksum, nil
}

func convertFieldToStringByType(field interface{}) string {
	switch fieldValue := field.(type) {
	case int:
		return strconv.Itoa(fieldValue)
	case float64:
		return strconv.FormatFloat(fieldValue, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(fieldValue), 'f', -1, 32)
	case string:
		return fieldValue
	case bool:
		return strconv.FormatBool(fieldValue)
	default:
		return ""
	}
}

func RandomWithPrefix(name string) string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return fmt.Sprintf("%s_%d", name, n.Int64())
}

func flavorSchema() *schema.Resource {
	return resourceDBaaSDatastoreV1().Schema["flavor"].Elem.(*schema.Resource)
}

func flavorHashSetFunc() schema.SchemaSetFunc {
	return schema.HashResource(flavorSchema())
}

func resourceDBaaSDatastoreV1FlavorFromSet(flavorSet *schema.Set) (*dbaas.Flavor, error) {
	if flavorSet.Len() == 0 {
		return nil, nil
	}
	var resourceVcpusRaw, resourceRAMRaw, resourceDiskRaw interface{}
	var ok bool

	resourceFlavorMap := flavorSet.List()[0].(map[string]interface{})
	if resourceVcpusRaw, ok = resourceFlavorMap["vcpus"]; !ok {
		return &dbaas.Flavor{}, errors.New("flavor.vcpus value isn't provided")
	}
	if resourceRAMRaw, ok = resourceFlavorMap["ram"]; !ok {
		return &dbaas.Flavor{}, errors.New("flavor.ram value isn't provided")
	}
	if resourceDiskRaw, ok = resourceFlavorMap["disk"]; !ok {
		return &dbaas.Flavor{}, errors.New("flavor.disk value isn't provided")
	}

	resourceVcpus := resourceVcpusRaw.(int)
	resourceRAM := resourceRAMRaw.(int)
	resourceDisk := resourceDiskRaw.(int)

	flavor := &dbaas.Flavor{
		Vcpus: resourceVcpus,
		RAM:   resourceRAM,
		Disk:  resourceDisk,
	}

	return flavor, nil
}

func resourceDBaaSDatastoreV1FlavorToSet(flavor dbaas.Flavor) *schema.Set {
	flavorSet := &schema.Set{
		F: flavorHashSetFunc(),
	}

	flavorSet.Add(map[string]interface{}{
		"vcpus": flavor.Vcpus,
		"ram":   flavor.RAM,
		"disk":  flavor.Disk,
	})

	return flavorSet
}

func resourceDBaaSDatastoreV1InstancesToList(instances []dbaas.Instances) []interface{} {
	flattenedInstances := make([]interface{}, len(instances))

	for i, instance := range instances {
		flattenedInstance := map[string]interface{}{
			"role":        instance.Role,
			"floating_ip": instance.FloatingIP,
		}
		flattenedInstances[i] = flattenedInstance
	}

	return flattenedInstances
}

func resourceDBaaSDatastoreV1RestoreOptsFromSet(restoreSet *schema.Set) (*dbaas.Restore, error) {
	if restoreSet.Len() == 0 {
		return nil, nil
	}
	var resourceDatastoreIDRaw, resourceTargetTimeRaw interface{}
	var ok bool

	resourceRestoreMap := restoreSet.List()[0].(map[string]interface{})
	if resourceDatastoreIDRaw, ok = resourceRestoreMap["datastore_id"]; !ok {
		return &dbaas.Restore{}, errors.New("restore.datastore_id value isn't provided")
	}
	if resourceTargetTimeRaw, ok = resourceRestoreMap["target_time"]; !ok {
		return &dbaas.Restore{}, errors.New("restore.target_time value isn't provided")
	}

	resourceDatastoreID := resourceDatastoreIDRaw.(string)
	resourceTargetTime := resourceTargetTimeRaw.(string)

	restore := &dbaas.Restore{
		DatastoreID: resourceDatastoreID,
		TargetTime:  resourceTargetTime,
	}

	return restore, nil
}

func resourceDBaaSDatastoreV1FloatingIPsOptsFromSet(floatingIPsSet *schema.Set) (*dbaas.FloatingIPs, error) {
	if floatingIPsSet.Len() == 0 {
		return nil, nil
	}

	floatingIPMap, ok := floatingIPsSet.List()[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid format for floating IPs data")
	}

	masterCount, ok := floatingIPMap["master"].(int)
	if !ok {
		return nil, fmt.Errorf(
			"invalid or missing %s floating IPs count",
			masterRole,
		)
	}

	replicasCount, ok := floatingIPMap["replica"].(int)
	if !ok {
		return nil, fmt.Errorf(
			"invalid or missing %s floating IPs count",
			replicaRole,
		)
	}

	return &dbaas.FloatingIPs{
		Master:  masterCount,
		Replica: replicasCount,
	}, nil
}

func updateDatastoreName(ctx context.Context, d *schema.ResourceData, client *dbaas.API) error {
	var updateOpts dbaas.DatastoreUpdateOpts
	updateOpts.Name = d.Get("name").(string)

	log.Print(msgUpdate(objectDatastore, d.Id(), updateOpts))
	_, err := client.UpdateDatastore(ctx, d.Id(), updateOpts)
	if err != nil {
		return errUpdatingObject(objectDatastore, d.Id(), err)
	}

	log.Printf("[DEBUG] waiting for datastore %s to become 'ACTIVE'", d.Id())
	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waiters.WaitForDBaaSDatastoreV1ActiveState(ctx, client, d.Id(), timeout)
	if err != nil {
		return errUpdatingObject(objectDatastore, d.Id(), err)
	}

	return nil
}

func updateDatastoreFirewall(ctx context.Context, d *schema.ResourceData, client *dbaas.API) error {
	firewallSet := d.Get("firewall").(*schema.Set)
	firewallOpts, err := resourceDBaaSDatastoreV1FirewallOptsFromSet(firewallSet)
	if err != nil {
		return errParseDatastoreV1Firewall(err)
	}

	log.Print(msgUpdate(objectDatastore, d.Id(), firewallOpts))
	_, err = client.FirewallDatastore(ctx, d.Id(), firewallOpts)
	if err != nil {
		return errUpdatingObject(objectDatastore, d.Id(), err)
	}

	log.Printf("[DEBUG] waiting for datastore %s to become 'ACTIVE'", d.Id())
	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waiters.WaitForDBaaSDatastoreV1ActiveState(ctx, client, d.Id(), timeout)
	if err != nil {
		return errUpdatingObject(objectDatastore, d.Id(), err)
	}

	return nil
}

func updateDatastoreConfig(ctx context.Context, d *schema.ResourceData, client *dbaas.API) error {
	var configOpts dbaas.DatastoreConfigOpts
	datastore, err := client.Datastore(ctx, d.Id())
	if err != nil {
		return err
	}
	config := d.Get("config").(map[string]interface{})

	for param := range datastore.Config {
		if _, ok := config[param]; !ok {
			config[param] = nil
		}
	}

	configOpts.Config = config

	log.Print(msgUpdate(objectDatastore, d.Id(), configOpts))
	_, err = client.ConfigDatastore(ctx, d.Id(), configOpts)
	if err != nil {
		return errUpdatingObject(objectDatastore, d.Id(), err)
	}

	log.Printf("[DEBUG] waiting for datastore %s to become 'ACTIVE'", d.Id())
	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waiters.WaitForDBaaSDatastoreV1ActiveState(ctx, client, d.Id(), timeout)
	if err != nil {
		return errUpdatingObject(objectDatastore, d.Id(), err)
	}

	return nil
}

func updateDatastoreBackups(ctx context.Context, d *schema.ResourceData, client *dbaas.API) error {
	var backupsOpts dbaas.DatastoreBackupsOpts
	backupsOpts.BackupRetentionDays = d.Get("backup_retention_days").(int)

	log.Print(msgUpdate(objectDatastore, d.Id(), backupsOpts))
	_, err := client.BackupsDatastore(ctx, d.Id(), backupsOpts)
	if err != nil {
		return errUpdatingObject(objectDatastore, d.Id(), err)
	}

	log.Printf("[DEBUG] waiting for datastore %s to become 'ACTIVE'", d.Id())
	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waiters.WaitForDBaaSDatastoreV1ActiveState(ctx, client, d.Id(), timeout)
	if err != nil {
		return errUpdatingObject(objectDatastore, d.Id(), err)
	}

	return nil
}

func getResizeOpts(ctx context.Context, d *schema.ResourceData, client *dbaas.API) (dbaas.DatastoreResizeOpts, error) {
	var resizeOpts dbaas.DatastoreResizeOpts
	nodeCount := d.Get("node_count").(int)
	resizeOpts.NodeCount = nodeCount

	if d.HasChange("flavor_id") {
		flavorID := d.Get("flavor_id")
		resizeOpts.FlavorID = flavorID.(string)
	} else if d.HasChange("flavor") {
		flavorRaw := d.Get("flavor")
		flavorSet := flavorRaw.(*schema.Set)
		flavor, err := resourceDBaaSDatastoreV1FlavorFromSet(flavorSet)
		if err != nil {
			return resizeOpts, errParseDatastoreV1Resize(err)
		}
		resizeOpts.Flavor = flavor
	}

	typeID := d.Get("type_id").(string)
	datastoreType, err := client.DatastoreType(ctx, typeID)
	if err != nil {
		return resizeOpts, errors.New("Couldnt get datastore type with id" + typeID)
	}

	if datastoreType.Engine == "redis" {
		resizeOpts.Flavor = nil
	}

	return resizeOpts, nil
}

func resizeDatastore(ctx context.Context, d *schema.ResourceData, client *dbaas.API) error {
	resizeOpts, err := getResizeOpts(ctx, d, client)
	if err != nil {
		return errUpdatingObject(objectDatastore, d.Id(), err)
	}

	log.Print(msgUpdate(objectDatastore, d.Id(), resizeOpts))
	_, err = client.ResizeDatastore(ctx, d.Id(), resizeOpts)
	if err != nil {
		return errUpdatingObject(objectDatastore, d.Id(), err)
	}

	log.Printf("[DEBUG] waiting for datastore %s to become 'ACTIVE'", d.Id())
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waiters.WaitForDBaaSDatastoreV1ActiveState(ctx, client, d.Id(), timeout)
	if err != nil {
		return errUpdatingObject(objectDatastore, d.Id(), err)
	}

	return nil
}

func containDatastoreType(expectedTypes []string, datastoreType string) bool {
	for _, expectedType := range expectedTypes {
		if expectedType == datastoreType {
			return true
		}
	}

	return false
}

func buildDatastoreTypeErrorMessage(expectedDatastoreTypeEngines []string, datastoreTypeEngine string) string {
	var baseMessage string
	if len(expectedDatastoreTypeEngines) > 1 {
		baseMessage = "Provided datastore type must have one of the following engine types: "
	} else {
		baseMessage = "Provided datastore type must have an engine "
	}

	return baseMessage + strings.Join(expectedDatastoreTypeEngines, ", ") + " for this resource. But provided type is " + datastoreTypeEngine
}

func validateDatastoreType(ctx context.Context, expectedDatastoreTypeEngines []string, typeID string, client *dbaas.API) diag.Diagnostics {
	datastoreType, err := client.DatastoreType(ctx, typeID)
	if err != nil {
		return diag.FromErr(errors.New("Couldnt get datastore type with id " + typeID))
	}
	if !containDatastoreType(expectedDatastoreTypeEngines, datastoreType.Engine) {
		return diag.FromErr(errors.New(buildDatastoreTypeErrorMessage(expectedDatastoreTypeEngines, datastoreType.Engine)))
	}

	return nil
}

func getDatastoreMasterInstance(datastore dbaas.Datastore) (masterInstance dbaas.Instances, found bool) {
	for _, instance := range datastore.Instances {
		if instance.Role == masterRole {
			found, masterInstance = true, instance
			break
		}
	}

	return
}

func getDatastoreReplicasInstancesIDsWithFloatings(datastore dbaas.Datastore) (replicasIDs []string) {
	for _, instance := range datastore.Instances {
		if instance.Role == replicaRole && instance.FloatingIP != "" {
			replicasIDs = append(replicasIDs, instance.ID)
		}
	}

	return replicasIDs
}

func getDatastoreReplicasInstancesIDsWithoutFloatings(datastore dbaas.Datastore) (replicasIDs []string) {
	for _, instance := range datastore.Instances {
		if instance.Role == replicaRole && instance.FloatingIP == "" {
			replicasIDs = append(replicasIDs, instance.ID)
		}
	}

	return replicasIDs
}

// Floating IPs

func refreshDatastoreInstancesOutputsDiff(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	if diff.HasChanges("floating_ips") {
		if err := diff.SetNewComputed("instances"); err != nil {
			return err
		}
	}

	return nil
}

func dbaasFloatingIPCreate(ctx context.Context, d *schema.ResourceData, client *dbaas.API, instanceID string) error {
	var floatingIPOpts dbaas.FloatingIPsOpts
	floatingIPOpts.InstanceID = instanceID

	log.Printf("[DEBUG] Creating floating IP for instance %s", instanceID)
	err := client.CreateFloatingIP(ctx, floatingIPOpts)
	if err != nil {
		return fmt.Errorf(
			"error creating Floating IP for the instance %s", instanceID,
		)
	}
	log.Printf("[DEBUG] waiting for datastore %s to become 'ACTIVE'", d.Id())
	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waiters.WaitForDBaaSDatastoreV1ActiveState(ctx, client, d.Id(), timeout)
	if err != nil {
		return errUpdatingObject(objectDatastore, d.Id(), err)
	}

	return nil
}

func dbaasFloatingIPDelete(ctx context.Context, d *schema.ResourceData, client *dbaas.API, instanceID string) error {
	var floatingIPOpts dbaas.FloatingIPsOpts
	floatingIPOpts.InstanceID = instanceID

	log.Printf("[DEBUG] Delete floating IP from instance %s", instanceID)
	err := client.DeleteFloatingIP(ctx, floatingIPOpts)
	if err != nil {
		return fmt.Errorf(
			"error deleting Floating IP from the instance %s", instanceID,
		)
	}

	log.Printf("[DEBUG] waiting for datastore %s to become 'ACTIVE'", d.Id())
	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waiters.WaitForDBaaSDatastoreV1ActiveState(ctx, client, d.Id(), timeout)
	if err != nil {
		return errUpdatingObject(objectDatastore, d.Id(), err)
	}

	return nil
}

func getFloatingIPSchemaDiff(oldSchema *dbaas.FloatingIPs, newSchema *dbaas.FloatingIPs) (int, int) {
	masterDiff := newSchema.Master - oldSchema.Master
	replicasDiff := newSchema.Replica - oldSchema.Replica

	return masterDiff, replicasDiff
}

func getFloatingIPsSchemaFromDatastoreInstances(datastore dbaas.Datastore) (*dbaas.FloatingIPs, error) {
	masterCount, replicasCount := 0, 0
	for _, instance := range datastore.Instances {
		if instance.FloatingIP == "" {
			continue
		}

		switch instance.Role {
		case masterRole:
			masterCount++
			if masterCount > 1 {
				return nil, fmt.Errorf(
					"more than one %s found id the datastore",
					masterRole,
				)
			}
		case replicaRole:
			replicasCount++
		}
	}

	return &dbaas.FloatingIPs{
		Master:  masterCount,
		Replica: replicasCount,
	}, nil
}

func updateDatastoreFloatingIPs(ctx context.Context, d *schema.ResourceData, client *dbaas.API) error {
	_, n := d.GetChange("floating_ips")
	datastore, err := client.Datastore(ctx, d.Id())
	if err != nil {
		return fmt.Errorf(
			"failed to retrieve datastore to update floating IPs: %v",
			err,
		)
	}
	oldSchema, err := getFloatingIPsSchemaFromDatastoreInstances(datastore)
	if err != nil {
		return fmt.Errorf(
			"failed to get old floating IPs schema: %v",
			err,
		)
	}
	newSchema, err := resourceDBaaSDatastoreV1FloatingIPsOptsFromSet(n.(*schema.Set))
	if err != nil {
		return fmt.Errorf(
			"failed to get new floating IPs schema from set: %v",
			err,
		)
	}
	masterDiff, replicasDiff := getFloatingIPSchemaDiff(oldSchema, newSchema)

	if err := manageFloatingIPs(ctx, d, client, datastore, masterDiff, replicasDiff); err != nil {
		return err
	}

	return nil
}

func manageFloatingIPs(ctx context.Context, d *schema.ResourceData, client *dbaas.API, datastore dbaas.Datastore, masterDiff int, replicasDiff int) error {
	masterInstance, ok := getDatastoreMasterInstance(datastore)
	if !ok {
		return fmt.Errorf(
			"%s instance not found in the datastore",
			masterRole,
		)
	}

	var masterFloatingIP int
	if masterInstance.FloatingIP != "" {
		masterFloatingIP = 1
	} else {
		masterFloatingIP = 0
	}
	needed := float64(masterDiff)
	floatingIPCounter := float64(masterFloatingIP) + needed

	// if user tries to attach more than one floating IP to the master instance
	if int(math.Abs(floatingIPCounter)) > 1 {
		return fmt.Errorf(
			"floating IPs count for %s could not be greater than 1",
			masterRole,
		)
	}

	// if user set negative value for floating IP counter
	if floatingIPCounter < 0 {
		return fmt.Errorf(
			"floating IPs count for %s could not be less than 0",
			masterRole,
		)
	}

	if replicasDiff > 0 {
		replicasWithoutFloatingsIDs := getDatastoreReplicasInstancesIDsWithoutFloatings(datastore)
		if len(replicasWithoutFloatingsIDs) < replicasDiff {
			return fmt.Errorf(
				"insufficient replicas without floating IPs: needed %d, found %d",
				replicasDiff,
				len(replicasWithoutFloatingsIDs),
			)
		}
		if err := addReplicasFloatingIPs(ctx, d, client, replicasWithoutFloatingsIDs, replicasDiff); err != nil {
			return fmt.Errorf("failed to adjust replicas floating IPs: %v", err)
		}
	}

	if replicasDiff < 0 {
		replicasWithFloatingsIDs := getDatastoreReplicasInstancesIDsWithFloatings(datastore)
		if err := removeReplicasFloatingIPs(ctx, d, client, replicasWithFloatingsIDs, replicasDiff); err != nil {
			return fmt.Errorf("failed to adjust replicas floating IPs: %v", err)
		}
	}

	masterID := masterInstance.ID

	if masterDiff != 0 {
		if err := adjustMasterFloatingIPs(ctx, d, client, masterID, masterDiff); err != nil {
			return fmt.Errorf("failed to adjust master floating IPs: %v", err)
		}
	}

	return nil
}

func adjustMasterFloatingIPs(ctx context.Context, d *schema.ResourceData, client *dbaas.API, instanceID string, diff int) error {
	if diff > 0 {
		err := dbaasFloatingIPCreate(ctx, d, client, instanceID)
		if err != nil {
			return fmt.Errorf(
				"could not create floating IP for %s instance %s",
				masterRole,
				instanceID,
			)
		}
	}

	if diff < 0 {
		err := dbaasFloatingIPDelete(ctx, d, client, instanceID)
		if err != nil {
			return fmt.Errorf(
				"could not delete floating IP from %s instance %s",
				masterRole,
				instanceID,
			)
		}
	}

	return nil
}

func addReplicasFloatingIPs(ctx context.Context, d *schema.ResourceData, client *dbaas.API, replicasWithoutFIPs []string, diff int) error {
	for i := 0; i < diff; i++ {
		err := dbaasFloatingIPCreate(ctx, d, client, replicasWithoutFIPs[i])
		if err != nil {
			return fmt.Errorf(
				"could not create floating IP for %s instance %s",
				replicaRole,
				replicasWithoutFIPs[i],
			)
		}
	}

	return nil
}

func removeReplicasFloatingIPs(ctx context.Context, d *schema.ResourceData, client *dbaas.API, replicasWithFIPs []string, diff int) error {
	needed := int(math.Abs(float64(diff)))
	for i := 0; i < needed; i++ {
		err := dbaasFloatingIPDelete(ctx, d, client, replicasWithFIPs[i])
		if err != nil {
			return fmt.Errorf(
				"could not delete floating IP from %s instance %s",
				replicaRole,
				replicasWithFIPs[i],
			)
		}
	}

	return nil
}
