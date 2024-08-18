package omnicron

import (
	"context"
	"os"
)

// model is ranked from lowest to highest based on their strength and abilites to perform Speech To Text (STT) generation. Check out the models on replicate.

type ReplicateLowSTTModel string

type ReplicateHighSTTModel string

const (
	// Model on Replicate: https://replicate.com/openai/whisper
	WhisperModel ReplicateLowSTTModel = "openai/whisper"
	// Model on Replicate: https://replicate.com/turian/insanely-fast-whisper-with-video
	InsanelyFastWhisperWithVideoModel ReplicateHighSTTModel = "turian/insanely-fast-whisper-with-video"
)

type LowSTTParams struct {
	Audio                   *os.File `form:"audio"`
	Transcription           *string  `form:"transcription,omitempty"`
	Temperature             *float64 `form:"temperature,omitempty"`
	Translate               *bool    `form:"translate,omitempty"`
	InitialPrompt           *string  `form:"initial_prompt,omitempty"`
	ConditionOnPreviousText *bool    `form:"condition_on_previous_text,omitempty"`
}
type HighSTTParams struct {
	AudioFile *os.File `form:"audio"`
	URL       *string  `form:"url,omitempty"`
	Task      *string  `form:"task,omitempty"`
	BatchSize *int     `form:"batch_size,omitempty"`
	Timestamp *string  `form:"timestamp,omitempty"`
}

type LowSTTModelAndParams struct {
	Model      ReplicateLowSTTModel
	Parameters LowSTTParams
}

type HighSTTModelAndParams struct {
	Model      ReplicateHighSTTModel
	Parameters HighSTTParams
}

// speech to text generation function
func (c *Client) LowSTTGeneration(ctx context.Context, req LowSTTModelAndParams) (*GabsContainer, error) {
	if req.Model == "" {
		return nil, ErrModelNotFound
	}
	body, err := c.newFormWithFilePostRequest(ctx, "/replicate/stt", string(req.Model), req.Parameters)
	if err != nil {
		return nil, err
	}
	lowSTTGenResponse, err := unmarshalJSONResponse(body)
	if err != nil {
		return nil, err
	}
	return lowSTTGenResponse, nil
}
func (c *Client) HighSTTGeneration(ctx context.Context, req HighSTTModelAndParams) (*GabsContainer, error) {
	if req.Model == "" {
		return nil, ErrModelNotFound
	}
	body, err := c.newFormWithFilePostRequest(ctx, "/replicate/stt", string(req.Model), req.Parameters)
	if err != nil {
		return nil, err
	}
	highSTTGenResponse, err := unmarshalJSONResponse(body)
	if err != nil {
		return nil, err
	}
	return highSTTGenResponse, nil
}
