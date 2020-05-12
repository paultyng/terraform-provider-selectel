package selectel

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	v1 "github.com/selectel/domains-go/pkg/v1"
	"github.com/selectel/domains-go/pkg/v1/domain"
)

const domainsEndpoint = "https://api.selectel.ru/domains/v1"

func resourceDomainsDomainV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceDomainsDomainV1Create,
		Read:   resourceDomainsDomainV1Read,
		Update: resourceDomainsDomainV1Update,
		Delete: resourceDomainsDomainV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"create_date": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"change_date": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceDomainsDomainV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx := context.Background()
	client := v1.NewDomainsClientV1(config.Token, domainsEndpoint)

	createOpts := &domain.CreateOpts{
		Name: d.Get("name").(string),
	}

	log.Print(msgCreate(objectDomain, createOpts))
	domainObj, _, err := domain.Create(ctx, client, createOpts)
	if err != nil {
		return errCreatingObject(objectDomain, err)
	}

	d.SetId(strconv.Itoa(domainObj.ID))

	return resourceDomainsDomainV1Read(d, meta)
}

func resourceDomainsDomainV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx := context.Background()
	client := v1.NewDomainsClientV1(config.Token, domainsEndpoint)

	log.Print(msgGet(objectDomain, d.Id()))

	domainID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("failed to parse domain ID: %w", err)
	}

	domainObj, resp, err := domain.GetByID(ctx, client, domainID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return errGettingObject(objectDomain, d.Id(), err)
	}

	d.Set("name", domainObj.Name)
	d.Set("create_date", domainObj.CreateDate)
	d.Set("change_date", domainObj.ChangeDate)
	d.Set("user_id", domainObj.UserID)

	return nil
}

func resourceDomainsDomainV1Update(d *schema.ResourceData, meta interface{}) error {
	return resourceDomainsDomainV1Read(d, meta)
}

func resourceDomainsDomainV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx := context.Background()
	client := v1.NewDomainsClientV1(config.Token, domainsEndpoint)

	log.Print(msgDelete(objectDomain, d.Id()))

	domainID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("failed to parse domain ID: %w", err)
	}

	_, err = domain.Delete(ctx, client, domainID)
	if err != nil {
		return errDeletingObject(objectDomain, d.Id(), err)
	}

	return nil
}
