package kubeoptions

import (
	"context"
	"net/http"
	"strings"

	v1 "github.com/selectel/mks-go/pkg/v1"
)

// ListFeatureGates gets a list of available feature gates by Kubernetes versions.
func ListFeatureGates(ctx context.Context, client *v1.ServiceClient) ([]*View, *v1.ResponseResult, error) {
	url := strings.Join([]string{client.Endpoint, v1.ResourceURLFeatureGates}, "/")
	responseResult, err := client.DoRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	// Extract available admission-controllers from the response body.
	var result struct {
		FGList []*View `json:"feature_gates"`
	}
	err = responseResult.ExtractResult(&result)
	if err != nil {
		return nil, responseResult, err
	}

	return result.FGList, responseResult, nil
}

// ListAdmissionControllers gets a list of available admission controllers by Kubernetes versions.
func ListAdmissionControllers(ctx context.Context, client *v1.ServiceClient) ([]*View, *v1.ResponseResult, error) {
	url := strings.Join([]string{client.Endpoint, v1.ResourceURLAdmissionControllers}, "/")
	responseResult, err := client.DoRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	// Extract available admission-controllers from the response body.
	var result struct {
		ACList []*View `json:"admission_controllers"`
	}
	err = responseResult.ExtractResult(&result)
	if err != nil {
		return nil, responseResult, err
	}

	return result.ACList, responseResult, nil
}
