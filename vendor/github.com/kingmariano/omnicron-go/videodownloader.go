package omnicron

import (
	"context"
)

type VideoDownloadParams struct {
	// url could be youtube, vimeo, facebook etc. see [!sites supported](https://github.com/iawia002/lux?tab=readme-ov-file#supported-sites)
	URL string `json:"url"`
	// resolution could be 1080p, 720p, 480p, 240p depending on the youtube. leave blank If you are not sure.
	Resolution string `json:"resolution,omitempty"`
}

// this download video function downloads the video url from any of the supported sites uploads to cloudinary and returns the direct download url
func (c *Client) DownloadVideo(ctx context.Context, req *VideoDownloadParams) (*GabsContainer, error) {
	if req.URL == "" {
		return nil, ErrVideoDownloadNoURLProvided
	}
	body, err := c.newJSONPostRequest(ctx, "/downloadvideo", "", req)
	if err != nil {
		return nil, err
	}
	videoDownloadResponse, err := unmarshalJSONResponse(body)
	if err != nil {
		return nil, err
	}
	return videoDownloadResponse, nil
}
