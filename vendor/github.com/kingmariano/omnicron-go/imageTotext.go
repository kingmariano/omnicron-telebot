package omnicron

import (
	"context"
	"os"
)

type ImageToTextParams struct {
	File *os.File `form:"file"`
}

func (c *Client) ImageToText(ctx context.Context, params ImageToTextParams) (*GabsContainer, error) {
	if params.File == nil {
		return nil, ErrNoFileProvided
	}
	body, err := c.newFormWithFilePostRequest(ctx, "/image2text", "", params)
	if err != nil {
		return nil, err
	}
	imageToTextResponse, err := unmarshalJSONResponse(body)
	if err != nil {
		return nil, err
	}
	return imageToTextResponse, nil
}
