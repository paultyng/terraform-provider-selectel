package selectel

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/selectel/go-selvpcclient/selvpcclient"
	"github.com/selectel/go-selvpcclient/selvpcclient/resell/v2/subnets"
	"github.com/terraform-providers/terraform-provider-selectel/selectel/internal/hashcode"
)

func getPrefixLengthFromCIDR(cidr string) (int, error) {
	cidrParts := strings.Split(cidr, "/")
	if len(cidrParts) != 2 {
		return 0, fmt.Errorf("got invalid CIDR: %s", cidr)
	}

	prefixLenght, err := strconv.Atoi(cidrParts[1])
	if err != nil {
		return 0, fmt.Errorf("error reading prefix length from '%s': %s", cidrParts[1], err)
	}

	return prefixLenght, nil
}

func getIPVersionFromPrefixLength(prefixLength int) string {
	// Any subnet with prefix length larger than 32 is a IPv6 protocol subnet
	// and Selectel doesn't provide any IPv6 subnets with smaller prefix lengths.
	if prefixLength > 32 {
		return string(selvpcclient.IPv6)
	}

	return string(selvpcclient.IPv4)
}

// subnetsMapsFromStructs converts the provided subnets.Subnet to
// the slice of maps correspondingly to the resource's schema.
func subnetsMapsFromStructs(subnetsStructs []subnets.Subnet) []map[string]interface{} {
	associatedSubnets := make([]map[string]interface{}, len(subnetsStructs))

	for i, subnet := range subnetsStructs {
		associatedSubnets[i] = map[string]interface{}{
			"network_id":      subnet.NetworkID,
			"subnet_id":       subnet.SubnetID,
			"region":          subnet.Region,
			"cidr":            subnet.CIDR,
			"vlan_id":         subnet.VLANID,
			"project_id":      subnet.ProjectID,
			"vtep_ip_address": subnet.VTEPIPAddress,
		}
	}

	return associatedSubnets
}

// hashSubnets is a hash function to use with the "subnet" set.
func hashSubnets(v interface{}) int {
	m := v.(map[string]interface{})
	return hashcode.String(fmt.Sprintf("%s-", m["network_id"].(string)))
}
