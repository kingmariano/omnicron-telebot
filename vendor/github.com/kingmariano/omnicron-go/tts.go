package omnicron

import (
	"context"
	"os"
)

// model is ranked from lowest to highest based on their strength and abilites to perform Text To Speech (TTS) generation. Check out the models on replicate.
type ReplicateLowTTSModel string

type ReplicateMediumTTSModel string

type ReplicateHighTTSModel string

const (
	// Model on Replicate: https://replicate.com/lucataco/xtts-v2
	XTTSV2Model ReplicateLowTTSModel = "lucataco/xtts-v2"
	// Model on Replicate: https://replicate.com/zsxkib/realistic-voice-cloning
	RealisticVoiceCloningModel ReplicateMediumTTSModel = "zsxkib/realistic-voice-cloning"
	// Model on Replicate: https://replicate.com/chenxwh/openvoice
	OpenVoiceModel ReplicateHighTTSModel = "chenxwh/openvoice"
)

type LowTTSParams struct {
	Text         string   `form:"text"`
	Speaker      *os.File `form:"speaker"`
	Language     *string  `form:"language,omitempty"`
	CleanupVoice *bool    `form:"cleanup_voice,omitempty"`
}

type MediumTTSParams struct {
	SongInput                 *os.File `form:"song_input"`
	RvcModel                  *string  `form:"rvc_model,omitempty"`
	CustomRvcModelDownloadURL *string  `form:"custom_rvc_model_download_url,omitempty"`
	PitchChange               *string  `form:"pitch_change,omitempty"`
	IndexRate                 *float64 `form:"index_rate,omitempty"`
	FilterRaidus              *int     `form:"filter_raidus,omitempty"`
	RmsMixRate                *float64 `form:"rms_mix_rate,omitempty"`
	PitchDetectionAlgorithm   *string  `form:"pitch_detection_algorithm,omitempty"`
	CrepeHopLength            *int     `form:"crepe_hop_length,omitempty"`
	Protect                   *float64 `form:"protect,omitempty"`
	MainVocalsVolumeChange    *float64 `form:"main_vocals_volume_change,omitempty"`
	BackupVocalsVolumeChange  *float64 `form:"backup_vocals_volume_change,omitempty"`
	InstrumentalVolumeChange  *float64 `form:"instrumental_volume_change,omitempty"`
	PitchChangeAll            *float64 `form:"pitch_change_all,omitempty"`
	ReverbSize                *float64 `form:"reverb_size,omitempty"`
	ReverbWetness             *float64 `form:"reverb_wetness,omitempty"`
	ReverbDryness             *float64 `form:"reverb_dryness,omitempty"`
	ReverbDamping             *float64 `form:"reverb_damping,omitempty"`
	OutputFormat              *string  `form:"output_format,omitempty"`
}
type HighTTSParams struct {
	Audio    *os.File `form:"audio"`
	Text     string   `form:"text"`
	Language *string  `form:"language,omitempty"`
	Speed    *float64 `form:"speed,omitempty"`
}

type LowTTSModelAndParams struct {
	Model      ReplicateLowTTSModel
	Parameters LowTTSParams
}

type MediumTTSModelAndParams struct {
	Model      ReplicateMediumTTSModel
	Parameters MediumTTSParams
}

type HighTTSModelAndParams struct {
	Model      ReplicateHighTTSModel
	Parameters HighTTSParams
}

func (c *Client) LowTTSGeneration(ctx context.Context, req LowTTSModelAndParams) (*GabsContainer, error) {
	if req.Model == "" {
		return nil, ErrModelNotFound
	}
	body, err := c.newFormWithFilePostRequest(ctx, "/replicate/tts", string(req.Model), req.Parameters)
	if err != nil {
		return nil, err
	}
	lowTTSGenResponse, err := unmarshalJSONResponse(body)
	if err != nil {
		return nil, err
	}
	return lowTTSGenResponse, nil
}

func (c *Client) MediumTTSGeneration(ctx context.Context, req MediumTTSModelAndParams) (*GabsContainer, error) {
	if req.Model == "" {
		return nil, ErrModelNotFound
	}
	body, err := c.newFormWithFilePostRequest(ctx, "/replicate/tts", string(req.Model), req.Parameters)
	if err != nil {
		return nil, err
	}
	mediumTTSGenResponse, err := unmarshalJSONResponse(body)
	if err != nil {
		return nil, err
	}
	return mediumTTSGenResponse, nil
}

func (c *Client) HighTTSGeneration(ctx context.Context, req HighTTSModelAndParams) (*GabsContainer, error) {
	if req.Model == "" {
		return nil, ErrModelNotFound
	}
	body, err := c.newFormWithFilePostRequest(ctx, "/replicate/tts", string(req.Model), req.Parameters)
	if err != nil {
		return nil, err
	}
	highTTSGenResponse, err := unmarshalJSONResponse(body)
	if err != nil {
		return nil, err
	}
	return highTTSGenResponse, nil
}
