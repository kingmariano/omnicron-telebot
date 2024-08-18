package commands

import (
	"context"
	"github.com/kingmariano/omnicron-go"
	"github.com/kingmariano/omnicron-telebot/config"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
	"time"
)

var chatResponse string
var systemPrompt = `You are Omnibot, a versatile Telegram bot with a wide range of capabilities. Your primary mode is engaging in natural and friendly conversation. You can help with various tasks such as chat completions, image generation, music downloading, and more.

You will only mention specific commands when they are relevant to the conversation or when the user explicitly asks for them. The following commands are available if the user requests specific functionalities:

    /generate_video - Generate Quality Videos
    /transcribe_audio - Transcribe various audio formats to written text
    /text2speech - Generate quality audio output from text
    /download_video_url - Download videos from websites like YouTube, Vimeo, Instagram, and more.
    /download_song - Easily download any music to your device.
    /shazam - Identify any song using the Shazam algorithm.
    /generate_music - Generate quality music by entering a prompt.
    /convert2mp3 - Convert any audio or video format to MP3
    /upscale_image - Enhance image resolution and quality effortlessly.
    /youtube_summarization - Summarize YouTube videos quickly and easily.
    /image2text - Convert images to text effortlessly using advanced OCR technology.
    /generate_image - Generate stunning images.

When users say "hello" or engage in general conversation, respond naturally and avoid mentioning specific commands unless the user prompts for or hints at a specific function.`

func Chat(b *config.BotConfig) tele.HandlerFunc {

	return func(c tele.Context) error {
		ctx, cancel := context.WithTimeout(b.Context, 3*time.Minute)
		defer cancel()

		omniClient := omnicron.NewClient(b.MyAPIKey)
		c.Notify(tele.Typing)
		receivedMessage := c.Text()
		userID := c.Sender().ID
		// Initialize chat history for the user if not already present
		if _, ok := b.ChatHistory[userID]; !ok {
			b.ChatHistory[userID] = []omnicron.Message{
				{
					Role:    "system",
					Content: systemPrompt,
				},
			}
		}

		// Append the received message to the chat history
		b.ChatHistory[userID] = append(b.ChatHistory[userID], omnicron.Message{
			Role:    "user",
			Content: receivedMessage,
		})

		// Attempt to get a response from the GroqChatCompletion endpoint
		grokResponse, err := omniClient.GroqChatCompletion(ctx, &omnicron.GroqChatCompletionParams{
			Messages: b.ChatHistory[userID],
			Model:    "mixtral-8x7b-32768",
		})
		if err != nil {
			b.Logger.Error("an error occurred while using the Groq endpoint:", zap.String("error", err.Error()))

			// Use the GPT4Free endpoint as a backup
			g4fResponse, err := omniClient.GPT4Free(ctx, &omnicron.G4FRequest{
				Messages: b.ChatHistory[userID],
			})
			if err != nil {
				b.Logger.Error("an error occurred while using the GPT4Free endpoint:", zap.String("error", err.Error()))
				return c.Reply("Sorry, an error occurred while processing your request.")
			}

			chatResponse = g4fResponse.Path("response").Data().(string)
		} else {
			chatResponse = grokResponse.Path("choices.0.message.content").Data().(string)
		}
		// Append the assistant's response to the chat history
		b.ChatHistory[userID] = append(b.ChatHistory[userID], omnicron.Message{
			Role:    "assistant",
			Content: chatResponse,
		})
		return c.Reply(chatResponse)
	}
}
