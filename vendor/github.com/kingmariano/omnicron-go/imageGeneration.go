/*
This image generation function uses different and high efficient image generation AI models on Replicate. Please note that cost varies as you some of this model.
*/
package omnicron

import (
	"context"
	"os"
)

// ReplicateLowImageGenerationModel defines the low image generation models on Replicate.
// These models do not support image-to-image processing.
type ReplicateLowImageGenerationModel string

// ReplicateHighImageGenerationModel defines the high image generation models on Replicate.
// These models support image-to-image processing.
type ReplicateHighImageGenerationModel string

const (
	// Model on Replicate: https://replicate.com/bytedance/sdxl-lightning-4step
	SDXLLightning4stepModel ReplicateLowImageGenerationModel = "bytedance/sdxl-lightning-4step"
	// Model on Replicate: https://replicate.com/lucataco/realvisxl-v2.0
	RealvisxlV20Model ReplicateHighImageGenerationModel = "lucataco/realvisxl-v2.0"
	// Model on Replicate: https://replicate.com/playgroundai/playground-v2.5-1024px-aesthetic
	PlaygroundV251024pxAestheticModel ReplicateHighImageGenerationModel = "playgroundai/playground-v2.5-1024px-aesthetic"
	// Model on Replicate: https://replicate.com/lucataco/dreamshaper-xl-turbo
	DreamshaperXlTurboModel ReplicateLowImageGenerationModel = "lucataco/dreamshaper-xl-turbo"
	// Model on Replicate: https://replicate.com/lorenzomarines/astra
	AstraModel ReplicateHighImageGenerationModel = "lorenzomarines/astra"
)

// LowImageGenerationParams defines the parameters for low image generation models.
type LowImageGenerationParams struct {
	Prompt            string   `json:"prompt"`
	Width             *int     `json:"width,omitempty"`
	Height            *int     `json:"height,omitempty"`
	Scheduler         *string  `json:"scheduler,omitempty"`
	NumOutputs        *int     `json:"num_outputs,omitempty"`
	GuidanceScale     *float64 `json:"guidance_scale,omitempty"`
	NegativePrompt    *string  `json:"negative_prompt,omitempty"`
	NumInferenceSteps *int     `json:"num_inference_steps,omitempty"`
}

// LowImageGenerationModelAndParams groups the low image generation model with its parameters.
type LowImageGenerationModelAndParams struct {
	Model      ReplicateLowImageGenerationModel
	Parameters *LowImageGenerationParams
}

// HighImageGenerationParams defines the parameters for high image generation models.
type HighImageGenerationParams struct {
	Prompt            string   `form:"prompt"`
	Width             *int     `form:"width,omitempty"`
	Height            *int     `form:"height,omitempty"`
	Scheduler         *string  `form:"scheduler,omitempty"`
	NumOutputs        *int     `form:"num_outputs,omitempty"`
	GuidanceScale     *float64 `form:"guidance_scale,omitempty"`
	NegativePrompt    *string  `form:"negative_prompt,omitempty"`
	NumInferenceSteps *int     `form:"num_inference_steps,omitempty"`
	LoraScale         *float64 `form:"lora_scale,omitempty"`
	Image             *os.File `form:"image,omitempty"`
	Mask              *os.File `form:"mask,omitempty"`
	PromptStrength    *float64 `form:"prompt_strength,omitempty"`
	ApplyWatermark    *bool    `form:"apply_watermark,omitempty"`
	Seed              *int     `form:"seed,omitempty"`
}

// HighImageGenerationModelAndParams groups the high image generation model with its parameters.
type HighImageGenerationModelAndParams struct {
	Model      ReplicateHighImageGenerationModel
	Parameters HighImageGenerationParams
}

// LowImageGeneration handles the low image generation request for models that do not support image-to-image processing.
func (c *Client) LowImageGeneration(ctx context.Context, req LowImageGenerationModelAndParams) (*GabsContainer, error) {
	if req.Model == "" {
		return nil, ErrModelNotFound
	}
	body, err := c.newJSONPostRequest(ctx, "/replicate/imagegeneration", string(req.Model), req.Parameters)
	if err != nil {
		return nil, err
	}
	lowImageGenResponse, err := unmarshalJSONResponse(body)
	if err != nil {
		return nil, err
	}
	return lowImageGenResponse, nil
}

// HighImageGeneration handles the high image generation request for models that support image-to-image processing.
func (c *Client) HighImageGeneration(ctx context.Context, req HighImageGenerationModelAndParams) (*GabsContainer, error) {
	if req.Model == "" {
		return nil, ErrModelNotFound
	}
	if req.Parameters.Prompt == "" {
		return nil, ErrPromptMissing
	}

	body, err := c.newFormWithFilePostRequest(ctx, "/replicate/imagegeneration", string(req.Model), req.Parameters)
	if err != nil {
		return nil, err
	}
	highImageGenResponse, err := unmarshalJSONResponse(body)
	if err != nil {
		return nil, err
	}
	return highImageGenResponse, nil
}
