package omnicron

import (
	"context"
	"os"
)

type DocGPTParams struct {
	File   *os.File `form:"file"`
	Prompt string   `form:"prompt"`
}

func (c *Client) DocGPT(ctx context.Context, req DocGPTParams) (*GabsContainer, error) {
	if req.File == nil {
		return nil, ErrNoFileProvided
	}
	if req.Prompt == "" {
		return nil, ErrPromptMissing
	}

	body, err := c.newFormWithFilePostRequest(ctx, "/docgpt", "", req)
	if err != nil {
		return nil, err
	}
	docGPTResponse, err := unmarshalJSONResponse(body)
	if err != nil {
		return nil, err
	}
	return docGPTResponse, nil
}
