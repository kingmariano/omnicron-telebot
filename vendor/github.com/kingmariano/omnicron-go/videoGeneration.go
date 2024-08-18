/*
This video  Generation function uses the zeroscope-v2-xl model on replicate. Please note that cost varies as you this model.
*/
package omnicron

import (
	"context"
	"os"
)

// These model is described high because it is the only model provided for video generation
type ReplicateHighVideoGenerationModel string

const (
	// Model on Replicate: https://replicate.com/anotherjesse/zeroscope-v2-xl
	ZeroScopeV2XLModel ReplicateHighVideoGenerationModel = "anotherjesse/zeroscope-v2-xl" // This model is very powerful and costly. It's recommended for generating high-quality videos.
)

type HighVideoGenerationParams struct {
	Prompt            string   `json:"prompt"`
	NegativePrompt    *string  `json:"negative_prompt,omitempty"`
	InitVideo         *os.File `json:"init_video,omitempty"`
	InitWeight        *float64 `json:"init_weight,omitempty"`
	NumFrames         *int     `json:"num_frames,omitempty"`
	NumInferenceSteps *int     `json:"num_inferences_steps,omitempty"`
	Width             *int     `json:"width,omitempty"`
	Height            *int     `json:"height,omitempty"`
	GuidanceScale     *float64 `json:"guidance_scale,omitempty"`
	FPS               *int     `json:"fps,omitempty"`
	VideoModel        *string  `json:"video_model,omitempty"`
	BatchSize         *int     `json:"batch_size,omitempty"`
	RemoveWatermark   *bool    `json:"remove_watermark,omitempty"`
}

type HighVideoGenerationModelAndParams struct {
	Model      ReplicateHighVideoGenerationModel
	Parameters HighVideoGenerationParams
}

func (c *Client) VideoGeneration(ctx context.Context, req HighVideoGenerationModelAndParams) (*GabsContainer, error) {
	if req.Model == "" {
		return nil, ErrModelNotFound
	}
	if req.Parameters.Prompt == "" {
		return nil, ErrPromptMissing
	}

	body, err := c.newFormWithFilePostRequest(ctx, "/replicate/videogeneration", string(req.Model), req.Parameters)
	if err != nil {
		return nil, err
	}
	videoGenResponse, err := unmarshalJSONResponse(body)
	if err != nil {
		return nil, err
	}
	return videoGenResponse, nil
}
