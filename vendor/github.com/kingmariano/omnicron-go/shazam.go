package omnicron

import (
	"context"
	"os"
)

type ShazamParams struct {
	File *os.File `form:"file"`
}

func (c *Client) Shazam(ctx context.Context, params ShazamParams) (*GabsContainer, error) {
	if params.File == nil {
		return nil, ErrNoFileProvided
	}
	body, err := c.newFormWithFilePostRequest(ctx, "/shazam", "", params)
	if err != nil {
		return nil, err
	}
	shazamResponse, err := unmarshalJSONResponse(body)
	if err != nil {
		return nil, err
	}
	return shazamResponse, nil
}
