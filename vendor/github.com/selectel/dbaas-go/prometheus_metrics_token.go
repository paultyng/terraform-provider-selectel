package dbaas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// PrometheusMetricTokenCreateOpts represents options for the prometheus metrics token Create request.
type PrometheusMetricTokenCreateOpts struct {
	Name string `json:"name"`
}

// PrometheusMetricTokenUpdateOpts represents options for the prometheus metrics token Update request.
type PrometheusMetricTokenUpdateOpts struct {
	Name string `json:"name"`
}

// PrometheusMetricToken is the API response for the prometheus metrics tokens.
type PrometheusMetricToken struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	ProjectID string `json:"project_id"`
	Name      string `json:"name"`
	Value     string `json:"value"`
}

// PrometheusMetricToken returns a token based on the ID.
func (api *API) PrometheusMetricToken(
	ctx context.Context,
	prometheusMetricTokenID string,
) (PrometheusMetricToken, error) {
	uri := fmt.Sprintf("/prometheus-metrics-tokens/%s", prometheusMetricTokenID)

	resp, err := api.makeRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return PrometheusMetricToken{}, err
	}

	var result PrometheusMetricToken
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return PrometheusMetricToken{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result, nil
}

// PrometheusMetricTokens returns all tokens.
func (api *API) PrometheusMetricTokens(ctx context.Context) ([]PrometheusMetricToken, error) {
	uri := "/prometheus-metrics-tokens"

	resp, err := api.makeRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return []PrometheusMetricToken{}, err
	}

	var result struct {
		PrometheusMetricTokens []PrometheusMetricToken `json:"prometheus-metrics-tokens"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return []PrometheusMetricToken{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.PrometheusMetricTokens, nil
}

// CreatePrometheusMetricToken creates a new token.
func (api *API) CreatePrometheusMetricToken(
	ctx context.Context,
	opts PrometheusMetricTokenCreateOpts,
) (PrometheusMetricToken, error) {
	uri := "/prometheus-metrics-tokens"
	createPrometheusMetricTokensOpts := struct {
		PrometheusMetricToken PrometheusMetricTokenCreateOpts `json:"prometheus-metrics-token"`
	}{
		PrometheusMetricToken: opts,
	}
	requestBody, err := json.Marshal(createPrometheusMetricTokensOpts)
	if err != nil {
		return PrometheusMetricToken{}, fmt.Errorf("Error marshalling params to JSON, %w", err)
	}

	resp, err := api.makeRequest(ctx, http.MethodPost, uri, requestBody)
	if err != nil {
		return PrometheusMetricToken{}, err
	}

	var result struct {
		PrometheusMetricToken PrometheusMetricToken `json:"prometheus-metrics-token"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return PrometheusMetricToken{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.PrometheusMetricToken, nil
}

// DeletePrometheusMetricToken deletes an existing token.
func (api *API) DeletePrometheusMetricToken(ctx context.Context, prometheusMetricTokenID string) error {
	uri := fmt.Sprintf("/prometheus-metrics-tokens/%s", prometheusMetricTokenID)

	_, err := api.makeRequest(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return err
	}

	return nil
}

// UpdatePrometheusMetricToken updates an existing token.
func (api *API) UpdatePrometheusMetricToken(
	ctx context.Context,
	prometheusMetricTokenID string,
	opts PrometheusMetricTokenUpdateOpts,
) (PrometheusMetricToken, error) {
	uri := fmt.Sprintf("/prometheus-metrics-tokens/%s", prometheusMetricTokenID)
	updatePrometheusMetricTokensOpts := struct {
		PrometheusMetricToken PrometheusMetricTokenUpdateOpts `json:"prometheus-metrics-token"`
	}{
		PrometheusMetricToken: opts,
	}
	requestBody, err := json.Marshal(updatePrometheusMetricTokensOpts)
	if err != nil {
		return PrometheusMetricToken{}, fmt.Errorf("Error marshalling params to JSON, %w", err)
	}

	resp, err := api.makeRequest(ctx, http.MethodPut, uri, requestBody)
	if err != nil {
		return PrometheusMetricToken{}, err
	}

	var result PrometheusMetricToken
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return PrometheusMetricToken{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result, nil
}
