package selectel

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/selectel/dbaas-go"
)

func resourceDBaaSPostgreSQLLogicalReplicationSlotV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDBaaSPostgreSQLLogicalReplicationSlotV1Create,
		ReadContext:   resourceDBaaSPostgreSQLLogicalReplicationSlotV1Read,
		DeleteContext: resourceDBaaSPostgreSQLLogicalReplicationSlotV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDBaaSPostgreSQLLogicalReplicationSlotV1ImportState,
		},
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
			"datastore_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"database_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDBaaSPostgreSQLLogicalReplicationSlotV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	datastoreID := d.Get("datastore_id").(string)

	selMutexKV.Lock(datastoreID)
	defer selMutexKV.Unlock(datastoreID)

	databaseID := d.Get("database_id").(string)

	selMutexKV.Lock(databaseID)
	defer selMutexKV.Unlock(databaseID)

	dbaasClient, diagErr := getDBaaSClient(d, meta)
	if diagErr != nil {
		return diagErr
	}

	slotCreateOpts := dbaas.LogicalReplicationSlotCreateOpts{
		Name:        d.Get("name").(string),
		DatastoreID: datastoreID,
		DatabaseID:  databaseID,
	}

	log.Print(msgCreate(objectLogicalReplicationSlot, slotCreateOpts))
	slot, err := dbaasClient.CreateLogicalReplicationSlot(ctx, slotCreateOpts)
	if err != nil {
		return diag.FromErr(errCreatingObject(objectLogicalReplicationSlot, err))
	}

	log.Printf("[DEBUG] waiting for slot %s to become 'ACTIVE'", slot.ID)
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForDBaaSLogicalReplicationSlotV1ActiveState(ctx, dbaasClient, slot.ID, timeout)
	if err != nil {
		return diag.FromErr(errCreatingObject(objectLogicalReplicationSlot, err))
	}

	d.SetId(slot.ID)

	return resourceDBaaSPostgreSQLLogicalReplicationSlotV1Read(ctx, d, meta)
}

func resourceDBaaSPostgreSQLLogicalReplicationSlotV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	dbaasClient, diagErr := getDBaaSClient(d, meta)
	if diagErr != nil {
		return diagErr
	}

	log.Print(msgGet(objectLogicalReplicationSlot, d.Id()))
	slot, err := dbaasClient.LogicalReplicationSlot(ctx, d.Id())
	if err != nil {
		return diag.FromErr(errGettingObject(objectLogicalReplicationSlot, d.Id(), err))
	}

	d.Set("name", slot.Name)
	d.Set("datastore_id", slot.DatastoreID)
	d.Set("database_id", slot.DatabaseID)

	return nil
}

func resourceDBaaSPostgreSQLLogicalReplicationSlotV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	datastoreID := d.Get("datastore_id").(string)

	selMutexKV.Lock(datastoreID)
	defer selMutexKV.Unlock(datastoreID)

	databaseID := d.Get("database_id").(string)

	selMutexKV.Lock(databaseID)
	defer selMutexKV.Unlock(databaseID)

	dbaasClient, diagErr := getDBaaSClient(d, meta)
	if diagErr != nil {
		return diagErr
	}

	log.Print(msgDelete(objectLogicalReplicationSlot, d.Id()))
	err := dbaasClient.DeleteLogicalReplicationSlot(ctx, d.Id())
	if err != nil {
		return diag.FromErr(errDeletingObject(objectLogicalReplicationSlot, d.Id(), err))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{strconv.Itoa(http.StatusOK)},
		Target:     []string{strconv.Itoa(http.StatusNotFound)},
		Refresh:    dbaasLogicalReplicationSlotV1DeleteStateRefreshFunc(ctx, dbaasClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	log.Printf("[DEBUG] waiting for slot %s to become deleted", d.Id())
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error waiting for the slot %s to become deleted: %s", d.Id(), err))
	}

	return nil
}

func resourceDBaaSPostgreSQLLogicalReplicationSlotV1ImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if config.ProjectID == "" {
		return nil, errors.New("SEL_PROJECT_ID must be set for the resource import")
	}
	if config.Region == "" {
		return nil, errors.New("SEL_REGION must be set for the resource import")
	}

	d.Set("project_id", config.ProjectID)
	d.Set("region", config.Region)

	return []*schema.ResourceData{d}, nil
}
