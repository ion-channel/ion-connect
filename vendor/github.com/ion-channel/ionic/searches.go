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

// GetSearch takes a query to perform
// a productidentifier search against the Ion API, assembling a slice of Ionic
// products.ProductSearchResponse objects
func (ic *IonClient) GetSearch(query, token string) ([]products.Product, error) {
	params := &url.Values{}
	params.Set("q", query)

	b, err := ic.Get(searchEndpoint, token, params, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get productidentifiers search: %v", err.Error())
	}

	var results []products.Product
	err = json.Unmarshal(b, &results)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal product search results: %v", err.Error())
	}

	return results, nil
}
