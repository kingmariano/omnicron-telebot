package omnicron

import (
	"github.com/Jeffail/gabs/v2"
)

// using the gabs library for dynamic JSON handling
type GabsContainer = gabs.Container

func unmarshalJSONResponse(body []byte) (*GabsContainer, error) {
	container, err := gabs.ParseJSON(body)
	if err != nil {
		return nil, err
	}
	return container, nil
}
