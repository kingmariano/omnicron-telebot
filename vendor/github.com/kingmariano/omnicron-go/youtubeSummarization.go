package omnicron

import "context"

type YoutubeSummarizationParams struct {
	URL string `json:"url"`
}

func (c *Client) YoutubeSummarization(ctx context.Context, params *YoutubeSummarizationParams) (*GabsContainer, error) {
	if params.URL == "" {
		return nil, ErrVideoDownloadNoURLProvided
	}
	body, err := c.newJSONPostRequest(ctx, "/youtubesummarization", "", params)
	if err != nil {
		return nil, err
	}
	youtubeSummarizationResponse, err := unmarshalJSONResponse(body)
	if err != nil {
		return nil, err
	}
	return youtubeSummarizationResponse, nil

}
