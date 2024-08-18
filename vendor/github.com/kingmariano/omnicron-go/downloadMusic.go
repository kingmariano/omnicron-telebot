package omnicron

import (
	"context"
)

type MusicRequest struct {
	Song string `json:"song"`
}

// the downloadMusic function takes a song as input, downloads the song and return the direct cloudinary url. something to note: use the search music function to get the song before using it as input. Do not use any song name directly to avoid inaccuracy.
func (c *Client) DownloadMusic(ctx context.Context, req *MusicRequest) (*GabsContainer, error) {
	if req.Song == "" {
		return nil, ErrSongNotProvided
	}
	body, err := c.newJSONPostRequest(ctx, "/downloadmusic", "", req)
	if err != nil {
		return nil, err
	}
	musicDownloadResponse, err := unmarshalJSONResponse(body)
	if err != nil {
		return nil, err
	}
	return musicDownloadResponse, nil
}
