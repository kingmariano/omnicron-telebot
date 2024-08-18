package omnicron

import "context"

type G4FRequest struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model,omitempty"`
	Stream   bool      `json:"stream,omitempty"`
	Proxy    string    `json:"proxy,omitempty"`
	Timeout  int       `json:"timeout,omitempty"`
	Shuffle  bool      `json:"shuffle,omitempty"`
	ImageURL string    `json:"image_url,omitempty"`
}

// this uses the g4f library by xtekky: https://github.com/xtekky/gpt4free
func (c *Client) GPT4Free(ctx context.Context, req *G4FRequest) (*GabsContainer, error) {
	if len(req.Messages) == 0 {
		return nil, ErrGroqChatCompletionNoMessage
	}
	body, err := c.newJSONPostRequest(ctx, "/gpt4free", "", req)
	if err != nil {
		return nil, err
	}
	g4fResponse, err := unmarshalJSONResponse(body)
	if err != nil {
		return nil, err
	}
	return g4fResponse, nil
}
