package omnicron

import "errors"

var ErrNoQueryParameter = errors.New("no query parameter specified")
var ErrModelNotFound = errors.New("model not found")
var ErrGroqChatCompletionNoMessage = errors.New("message field is required")
var ErrPromptMissing = errors.New("prompt field is required")
var ErrNoFileProvided = errors.New("no file provided")
var ErrVideoDownloadNoURLProvided = errors.New("video download no URL provided")
var ErrConvertToMP3NoURLOrFile = errors.New("please provide either a valid URL or a file")
var ErrSongNotProvided = errors.New("please provide the song input name")

// ErrorResponse represents an error response from the API.
type ErrorResponse struct {
	Error string `json:"error"`
}
