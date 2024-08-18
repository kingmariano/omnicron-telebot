package omnicron

import (
	"context"
	"os"
)

// convert to url params takes a valid url or file and converts to mp3.
type ConvertToMP3Params struct {
	URL  string   `form:"url,omitempty"`
	File *os.File `form:"file,omitempty"`
}

// convert to mp3 function accepts the request parameters converts either the URL string or audio file to mp3 and then returns the result clouding direct url.
func (c *Client) ConvertToMP3(ctx context.Context, req ConvertToMP3Params) (*GabsContainer, error) {
	if req.URL == "" && req.File == nil {
		return nil, ErrConvertToMP3NoURLOrFile
	}
	body, err := c.newFormWithFilePostRequest(ctx, "/convert2mp3", "", req)
	if err != nil {
		return nil, err
	}
	convertToMP3Response, err := unmarshalJSONResponse(body)
	if err != nil {
		return nil, err
	}
	return convertToMP3Response, nil
}
