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

func resourceDBaaSUserV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDBaaSUserV1Create,
		ReadContext:   resourceDBaaSUserV1Read,
		UpdateContext: resourceDBaaSUserV1Update,
		DeleteContext: resourceDBaaSUserV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDBaaSUserV1ImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"datastore_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDBaaSUserV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	datastoreID := d.Get("datastore_id").(string)

	selMutexKV.Lock(datastoreID)
	defer selMutexKV.Unlock(datastoreID)

	dbaasClient, diagErr := getDBaaSClient(d, meta)
	if diagErr != nil {
		return diagErr
	}

	userCreateOpts := dbaas.UserCreateOpts{
		DatastoreID: datastoreID,
		Name:        d.Get("name").(string),
		Password:    d.Get("password").(string),
	}

	log.Print(msgCreate(objectUser, userCreateOpts))
	user, err := dbaasClient.CreateUser(ctx, userCreateOpts)
	if err != nil {
		return diag.FromErr(errCreatingObject(objectUser, err))
	}

	log.Printf("[DEBUG] waiting for user %s to become 'ACTIVE'", user.ID)
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForDBaaSUserV1ActiveState(ctx, dbaasClient, user.ID, timeout)
	if err != nil {
		return diag.FromErr(errCreatingObject(objectUser, err))
	}

	d.SetId(user.ID)

	return resourceDBaaSUserV1Read(ctx, d, meta)
}

func resourceDBaaSUserV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	dbaasClient, diagErr := getDBaaSClient(d, meta)
	if diagErr != nil {
		return diagErr
	}

	log.Print(msgGet(objectUser, d.Id()))
	user, err := dbaasClient.User(ctx, d.Id())
	if err != nil {
		return diag.FromErr(errGettingObject(objectUser, d.Id(), err))
	}
	d.Set("datastore_id", user.DatastoreID)
	d.Set("name", user.Name)
	d.Set("status", user.Status)

	return nil
}

func resourceDBaaSUserV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	datastoreID := d.Get("datastore_id").(string)

	selMutexKV.Lock(datastoreID)
	defer selMutexKV.Unlock(datastoreID)

	dbaasClient, diagErr := getDBaaSClient(d, meta)
	if diagErr != nil {
		return diagErr
	}

	if d.HasChange("password") {
		password := d.Get("password").(string)
		updateOpts := dbaas.UserUpdateOpts{
			Password: password,
		}

		log.Print(msgUpdate(objectUser, d.Id(), updateOpts))
		_, err := dbaasClient.UpdateUser(ctx, d.Id(), updateOpts)
		if err != nil {
			return diag.FromErr(errUpdatingObject(objectUser, d.Id(), err))
		}

		log.Printf("[DEBUG] waiting for user %s to become 'ACTIVE'", d.Id())
		timeout := d.Timeout(schema.TimeoutCreate)
		err = waitForDBaaSUserV1ActiveState(ctx, dbaasClient, d.Id(), timeout)
		if err != nil {
			return diag.FromErr(errUpdatingObject(objectUser, d.Id(), err))
		}
	}

	return resourceDBaaSUserV1Read(ctx, d, meta)
}

func resourceDBaaSUserV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	datastoreID := d.Get("datastore_id").(string)

	selMutexKV.Lock(datastoreID)
	defer selMutexKV.Unlock(datastoreID)

	dbaasClient, diagErr := getDBaaSClient(d, meta)
	if diagErr != nil {
		return diagErr
	}

	log.Print(msgDelete(objectUser, d.Id()))
	err := dbaasClient.DeleteUser(ctx, d.Id())
	if err != nil {
		return diag.FromErr(errDeletingObject(objectUser, d.Id(), err))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{strconv.Itoa(http.StatusOK)},
		Target:     []string{strconv.Itoa(http.StatusNotFound)},
		Refresh:    dbaasUserV1DeleteStateRefreshFunc(ctx, dbaasClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	log.Printf("[DEBUG] waiting for user %s to become deleted", d.Id())
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error waiting for the user %s to become deleted: %s", d.Id(), err))
	}

	return nil
}

func resourceDBaaSUserV1ImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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
