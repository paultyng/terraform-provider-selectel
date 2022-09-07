package selectel

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/hashicorp/go-retryablehttp"
	domainsV1 "github.com/selectel/domains-go/pkg/v1"
	"github.com/selectel/go-selvpcclient/selvpcclient"
	"github.com/selectel/go-selvpcclient/selvpcclient/resell"
	resellV2 "github.com/selectel/go-selvpcclient/selvpcclient/resell/v2"
	"github.com/selectel/go-selvpcclient/selvpcclient/resell/v2/tokens"
)

const (
	DefaultOSEndpoint = "https://api.selvpc.ru/identity/v3"
)

// Config contains all available configuration options.
type Config struct {
	Token      string
	OSEndpoint string
	Endpoint   string
	ProjectID  string
	DomainName string
	Region     string
	User       string
	Password   string
}

// Validate performs config validation.
func (c *Config) Validate() error {
	if c.Token == "" && !c.isKeystoneCredentials() {
		return errors.New("token or credentials with domain name must be specified")
	}
	if c.Endpoint == "" {
		c.Endpoint = strings.Join([]string{resell.Endpoint, resellV2.APIVersion}, "/")
	}
	if c.Region != "" {
		if err := validateRegion(c.Region); err != nil {
			return err
		}
	}
	if c.OSEndpoint == "" {
		c.OSEndpoint = DefaultOSEndpoint
	}

	return nil
}

// Initialize Selectel resell client.
func (c *Config) resellV2Client() *selvpcclient.ServiceClient {
	return resellV2.NewV2ResellClientWithEndpoint(c.Token, c.Endpoint)
}

// Create Keystone token by Selectel token or Keystone credentials.
func (c *Config) getToken(ctx context.Context, p string, r string) (string, error) {
	if !c.isKeystoneCredentials() {
		return c.getTokenBySelectelToken(ctx, p, r)
	}

	return c.getTokenByCredentials(ctx, p, r)
}

// useSelectelToken will determine which auth type to use:
// either the Selectel token or Keystone credentials.
func (c *Config) isKeystoneCredentials() bool {
	if c.User == "" || c.Password == "" || c.DomainName == "" {
		return false
	}

	return true
}

// Create Keystone token by Selectel token.
func (c *Config) getTokenBySelectelToken(ctx context.Context, p string, r string) (string, error) {
	tokenOpts := tokens.TokenOpts{
		ProjectID: p,
	}
	resellV2Client := c.resellV2Client()
	log.Print(msgCreate(objectToken, tokenOpts))

	token, _, err := tokens.Create(ctx, resellV2Client, tokenOpts)
	if err != nil {
		return "", err
	}

	return token.ID, err
}

// Create Keystone token by Keystone credentials.
func (c *Config) getTokenByCredentials(ctx context.Context, p string, r string) (string, error) {
	providerOpts := gophercloud.AuthOptions{
		AllowReauth:      true,
		IdentityEndpoint: c.OSEndpoint,
		Username:         c.User,
		Password:         c.Password,
		DomainName:       c.DomainName,
		Scope: &gophercloud.AuthScope{
			ProjectID: p,
		},
	}

	newProvider, err := openstack.AuthenticatedClient(providerOpts)
	if err != nil {
		fmt.Println(err)
	}

	var tokenID string
	if newProvider != nil {
		tokenID, err = newProvider.GetAuthResult().ExtractTokenID()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		return "", errors.New("authentication failed")
	}

	return tokenID, err
}

// Initialize Selectel domains client.
func (c *Config) domainsV1Client() *domainsV1.ServiceClient {
	domainsClient := domainsV1.NewDomainsClientV1WithDefaultEndpoint(c.Token)
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = nil // Ignore retyablehttp client logs
	retryClient.RetryWaitMin = domainsV1DefaultRetryWaitMin
	retryClient.RetryWaitMax = domainsV1DefaultRetryWaitMax
	retryClient.RetryMax = domainsV1DefaultRetry
	domainsClient.HTTPClient = retryClient.StandardClient()

	return domainsClient
}
