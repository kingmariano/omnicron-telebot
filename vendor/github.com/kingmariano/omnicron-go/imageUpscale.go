/*
This image upscale  function uses different and high efficient image Upscale AI models on Replicate. Please note that cost varies as you some of this model.
*/
package omnicron

import (
	"context"
	"os"
)

// These model is described low because the model doesnt support provide some additional functionalities for handling image upscaling effectively.
type ReplicateLowImageUpscaleGenerationModel string

// These model is described high because the model  supports the functionalities to provide for effective image scaling.
type ReplicateHighImageUpscaleGenerationModel string

const (

	// Model on Replicate: https://replicate.com/nightmareai/real-esrgan
	RealErsganModel ReplicateLowImageUpscaleGenerationModel = "nightmareai/real-esrgan"

	// Model on Replicate: https://replicate.com/philz1337x/clarity-upscaler
	ClarityUpscalerModel ReplicateHighImageUpscaleGenerationModel = "philz1337x/clarity-upscaler"
)

type LowImageUpscaleGenerationParams struct {
	Image       *os.File `form:"image"`
	Scale       *float64 `form:"scale,omitempty"`
	FaceEnhance *bool    `form:"face_enhance,omitempty"`
}

type HighImageUpscaleGenerationParams struct {
	Image                 *os.File `form:"image"`
	Prompt                *string  `form:"prompt,omitempty"`
	NegativePrompt        *string  `form:"negative_prompt,omitempty"`
	ScaleFactor           *float64 `form:"scale_factor,omitempty"`
	Dynamic               *float64 `form:"dynamic,omitempty"`
	Creativity            *float64 `form:"creativity,omitempty"`
	Resemblance           *float64 `form:"resemblance,omitempty"`
	TilingWidth           *int     `form:"tiling_width,omitempty"`
	TilingHeight          *int     `form:"tiling_height,omitempty"`
	SdModel               *string  `form:"sd_model,omitempty"`
	Scheduler             *string  `form:"scheduler,omitempty"`
	NumInferenceSteps     *int     `form:"num_inference_steps,omitempty"`
	Seed                  *int     `form:"seed,omitempty"`
	Downscaling           *bool    `form:"downscaling,omitempty"`
	DownscalingResolution *int     `form:"downscaling_resolution,omitempty"`
	Sharpen               *float64 `form:"sharpen,omitempty"`
	OutputFormat          *string  `form:"output_format,omitempty"`
}

// LowImageUpscaleGenerationModelAndParams groups the low imageupscale generation model with its parameters.
type LowImageUpscaleGenerationModelAndParams struct {
	Model      ReplicateLowImageUpscaleGenerationModel
	Parameters LowImageUpscaleGenerationParams
}

// HighImageUpscaleGenerationModelAndParams groups the high imageupscale generation model with its parameters.
type HighImageUpscaleGenerationModelAndParams struct {
	Model      ReplicateHighImageUpscaleGenerationModel
	Parameters HighImageUpscaleGenerationParams
}

func (c *Client) LowImageUpscaling(ctx context.Context, req LowImageUpscaleGenerationModelAndParams) (*GabsContainer, error) {
	if req.Model == "" {
		return nil, ErrModelNotFound
	}
	body, err := c.newFormWithFilePostRequest(ctx, "/replicate/imageupscale", string(req.Model), req.Parameters)
	if err != nil {
		return nil, err
	}
	lowImageUpscaleResponse, err := unmarshalJSONResponse(body)
	if err != nil {
		return nil, err
	}
	return lowImageUpscaleResponse, nil
}

func (c *Client) HighImageUpscaling(ctx context.Context, req HighImageUpscaleGenerationModelAndParams) (*GabsContainer, error) {
	if req.Model == "" {
		return nil, ErrModelNotFound
	}
	body, err := c.newFormWithFilePostRequest(ctx, "/replicate/imageupscale", string(req.Model), req.Parameters)
	if err != nil {
		return nil, err
	}
	highImageUpscaleResponse, err := unmarshalJSONResponse(body)
	if err != nil {
		return nil, err
	}
	return highImageUpscaleResponse, nil
}
