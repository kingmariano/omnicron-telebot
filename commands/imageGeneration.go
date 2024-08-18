package commands

import (
	"fmt"
	"log"

	"github.com/looplab/fsm"

	"github.com/kingmariano/omnicron-telebot/config"
	tele "gopkg.in/telebot.v3"
)

type ImageGenerationSession struct {
	Model  string
	Prompt string
	FSM    *fsm.FSM
}

var ImageGenButtons = []MarkUp{
	{Text: "sdxl Lightning 4step", Unique: "001", Data: "bytedance/sdxl-lightning-4step"},
	{Text: "realvisxl v2.0", Unique: "002", Data: "lucataco/realvisxl-v2.0"},
	{Text: "playground v2.5 1024px aesthetic", Unique: "003", Data: "playgroundai/playground-v2.5-1024px-aesthetic"},
	{Text: "dreamshaper xl turbo", Unique: "004", Data: "lucataco/dreamshaper-xl-turbo"},
	{Text: "astra", Unique: "005", Data: "lorenzomarines/astra"},
}

func ImageGeneration(b *config.BotConfig) tele.HandlerFunc {
	return func(c tele.Context) error {
		// Generate inline buttons
		selector := SendMarkupSelector(ImageGenButtons, &tele.ReplyMarkup{})
		// Send the message with the inline keyboard
		err := c.Send("Select AI model\ntype *cancel* to cancel the request", selector)
		if err != nil {
			return err
		}
		log.Println("groups created")
		newBot := c.Bot()
		callback := c.Callback()
		newBot.Handle(tele.OnCallback, func(d tele.Context) error {
			log.Println("callback noticed")
			nextPrompt := fmt.Sprintf("You selected %s model\nEnter text prompt to generate stunning and unique images. (ex: *A peaceful sunset over a calm ocean, with vibrant colors reflecting in the water.", callback.Data)
			return c.Send(nextPrompt)
		})
		log.Println("botgroup exiting..")
		isClosed, err := b.Bot.Close()
		if err != nil {
			return err
		}
		log.Printf("bot has %v", isClosed)
		return err
	}
}

// first thing is idle
// func newSession() *ImageGenerationSession {
// 	f := fsm.NewFSM(
// 		"idle",
// 		[]fsm.EventDesc{
// 			{Name: "start", Src: []string{"idle"}, Dst: "selectModel"},
// 			{Name: "model", Src: []string{"selectModel"}, Dst: "enterPrompt"},
// 			{Name: "Prompt", Src: []string{"enterPrompt"}, Dst: "generateImage"},
// 			{Name: "cancel", Src: []string{"idle"}, Dst: "cancel"},
// 		},
// 		fsm.Callbacks{
// 			"enter_start": func(_ context.Context, e *fsm.Event) {
// 				fmt.Println("Preparing to start the image generation process...")
// 			},
// 			"leave_idle": func(_ context.Context, e *fsm.Event) {
// 				fmt.Println("Leaving idle state...")
// 			},
// 			"enter_selectModel": func(_ context.Context, e *fsm.Event) {
// 				fmt.Println("Entering selectModel state. Please choose an AI model.")
// 			},
// 			"before_model": func(_ context.Context, e *fsm.Event) {
// 				fmt.Println("Model selection in progress...")
// 			},
// 			"leave_selectModel": func(_ context.Context, e *fsm.Event) {
// 				fmt.Println("Model selected, moving to prompt entry.")
// 			},
// 			"enter_enterPrompt": func(_ context.Context, e *fsm.Event) {
// 				fmt.Println("Entering enterPrompt state. Please enter your image prompt.")
// 			},
// 			"before_Prompt": func(_ context.Context, e *fsm.Event) {
// 				fmt.Println("Prompt entry in progress...")
// 			},
// 			"leave_enterPrompt": func(_ context.Context, e *fsm.Event) {
// 				fmt.Println("Prompt entered, generating the image...")
// 			},
// 			"enter_generateImage": func(_ context.Context, e *fsm.Event) {
// 				fmt.Println("Image generation started.")
// 			},
// 			"before_cancel": func(_ context.Context, e *fsm.Event) {
// 				fmt.Println("Cancelling the process and returning to idle state...")
// 			},
// 			"enter_cancel": func(_ context.Context, e *fsm.Event) {
// 				fmt.Println("Process cancelled.")
// 			},
// 		},
// 	)
// 	newSession := &ImageGenerationSession{
// 		FSM: f,
// 	}
// 	return newSession
// }
