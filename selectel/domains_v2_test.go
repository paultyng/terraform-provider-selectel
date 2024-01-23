package selectel

import (
	"context"
	"strconv"
	"testing"

	domainsV2 "github.com/selectel/domains-go/pkg/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockedDNSv2Client struct {
	mock.Mock
	domainsV2.Client
}

func (client *mockedDNSv2Client) ListZones(ctx context.Context, opts *map[string]string) (domainsV2.Listable[domainsV2.Zone], error) {
	args := client.Called(ctx, opts)
	zones := args.Get(0).(domainsV2.Listable[domainsV2.Zone])
	err := args.Error(1)

	return zones, err
}

func (client *mockedDNSv2Client) ListRRSets(ctx context.Context, zoneID string, opts *map[string]string) (domainsV2.Listable[domainsV2.RRSet], error) {
	args := client.Called(ctx, zoneID, opts)
	rrsets := args.Get(0).(domainsV2.Listable[domainsV2.RRSet])
	err := args.Error(1)

	return rrsets, err
}

func TestGetZoneByName_whenNeededZoneInResponseWithOffset(t *testing.T) {
	nameForSearch := "test.xyz."
	correctIDForSearch := "mocked-uuid-2"

	mDNSClient := new(mockedDNSv2Client)
	ctx := context.Background()
	nextOffset := 3
	opts1 := &map[string]string{
		"filter": nameForSearch,
		"limit":  "1000",
		"offset": "0",
	}
	opts2 := &map[string]string{
		"filter": nameForSearch,
		"limit":  "1000",
		"offset": strconv.Itoa(nextOffset),
	}
	incorrectNameForSearch := "a." + nameForSearch
	incorrectIDForSearch := "mocked-uuid-1"
	zonesWithNextOffset := domainsV2.Listable[domainsV2.Zone](domainsV2.List[domainsV2.Zone]{
		Count:      1,
		NextOffset: nextOffset,
		Items: []*domainsV2.Zone{
			{
				UUID: incorrectIDForSearch,
				Name: incorrectNameForSearch,
			},
		},
	})
	mDNSClient.On("ListZones", ctx, opts1).Return(zonesWithNextOffset, nil)
	zonesWithoutNextOffset := domainsV2.Listable[domainsV2.Zone](domainsV2.List[domainsV2.Zone]{
		Count:      1,
		NextOffset: 0,
		Items: []*domainsV2.Zone{
			{
				UUID: correctIDForSearch,
				Name: nameForSearch,
			},
		},
	})
	mDNSClient.On("ListZones", ctx, opts2).Return(zonesWithoutNextOffset, nil)

	zone, err := getZoneByName(ctx, mDNSClient, nameForSearch)

	assert.NoError(t, err)

	assert.NotNil(t, zone)
	assert.Equal(t, correctIDForSearch, zone.UUID)
	assert.Equal(t, nameForSearch, zone.Name)
}

func TestGetRrsetByNameAndType_whenNeededRrrsetInResponseWithOffset(t *testing.T) {
	rrsetNameForSearch := "test.xyz."
	rrsetTypeForSearch := "A"
	correctIDForSearch := "mocked-uuid-2"
	mockedZoneID := "mopcked-zone-id"
	mDNSClient := new(mockedDNSv2Client)
	ctx := context.Background()
	nextOffset := 3
	opts1 := &map[string]string{
		"name":        rrsetNameForSearch,
		"rrset_types": rrsetTypeForSearch,
		"limit":       "1000",
		"offset":      "0",
	}
	opts2 := &map[string]string{
		"name":        rrsetNameForSearch,
		"rrset_types": rrsetTypeForSearch,
		"limit":       "1000",
		"offset":      strconv.Itoa(nextOffset),
	}
	incorrectNameForSearch := "a." + rrsetNameForSearch
	incorrectIDForSearch := "mocked-uuid-1"
	rrsetWithNextOffset := domainsV2.Listable[domainsV2.RRSet](domainsV2.List[domainsV2.RRSet]{
		Count:      1,
		NextOffset: nextOffset,
		Items: []*domainsV2.RRSet{
			{
				UUID: incorrectIDForSearch,
				Name: incorrectNameForSearch,
				Type: domainsV2.RecordType(rrsetTypeForSearch),
			},
		},
	})
	mDNSClient.On("ListRRSets", ctx, mockedZoneID, opts1).Return(rrsetWithNextOffset, nil)
	rrsetsWithoutNextOffset := domainsV2.Listable[domainsV2.RRSet](domainsV2.List[domainsV2.RRSet]{
		Count:      1,
		NextOffset: 0,
		Items: []*domainsV2.RRSet{
			{
				UUID: correctIDForSearch,
				Name: rrsetNameForSearch,
				Type: domainsV2.RecordType(rrsetTypeForSearch),
			},
		},
	})
	mDNSClient.On("ListRRSets", ctx, mockedZoneID, opts2).Return(rrsetsWithoutNextOffset, nil)

	rrset, err := getRrsetByNameAndType(ctx, mDNSClient, mockedZoneID, rrsetNameForSearch, rrsetTypeForSearch)

	assert.NoError(t, err)

	assert.NotNil(t, rrset)
	assert.Equal(t, correctIDForSearch, rrset.UUID)
	assert.Equal(t, rrsetNameForSearch, rrset.Name)
	assert.Equal(t, rrsetTypeForSearch, string(rrset.Type))
}
