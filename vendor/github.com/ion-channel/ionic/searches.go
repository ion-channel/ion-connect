package ionic

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ion-channel/ionic/products"
)

const (
	searchEndpoint = "v1/search"
)

// GetSearch takes a query and optionally a resource type to perform
// a productidentifier search against the Ion API, assembling a slice of Ionic
// products.ProductSearchResponse objects
func (ic *IonClient) GetSearch(query, resource, token string) ([]products.Product, error) {
	params := &url.Values{}
	params.Set("q", query)
	if resource != "" {
		params.Set("r", resource)
	}

	b, err := ic.Get(searchEndpoint, token, params, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get search: %v", err.Error())
	}

	var results []products.Product
	err = json.Unmarshal(b, &results)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal search results: %v", err.Error())
	}

	return results, nil
}

func (ic *IonClient) GetRawSearch(query, resource, token string) (*json.RawMessage, error) {
	params := &url.Values{}
	params.Set("q", query)
	if resource != "" {
		params.Set("r", resource)
	}

	b, err := ic.Get(searchEndpoint, token, params, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get search: %v", err.Error())
	}

	var results json.RawMessage
	err = json.Unmarshal(b, &results)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal search results: %v", err.Error())
	}

	return &results, nil
}
