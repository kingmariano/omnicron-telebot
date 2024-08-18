package commands

import (
	"context"
	"fmt"
	"time"
   "github.com/google/uuid"
   "github.com/kingmariano/omnicron-telebot/config"
	"github.com/kingmariano/omnicron-telebot/internal/database"
   "go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func Start(b *config.BotConfig) tele.HandlerFunc{
   return func(c tele.Context) error{
      ctx, cancel := context.WithTimeout(b.Context, time.Minute)
		defer cancel()
      username := c.Sender().Username
      userID := c.Sender().ID
      user, err := b.DB.CreateUser(ctx, database.CreateUserParams{
      ID:             uuid.New(),
      UserName:       username,
      TelegramID: int32(userID),
      })
    if err != nil {
      b.Logger.Error("an error occurred while creating a new user: User already exists", zap.String("error", err.Error()))
      return c.Reply(fmt.Sprintf("Hello %s, How can i help you?", username))
    }
	
    return c.Reply(initMessage(user.UserName))
   }
}
func initMessage(username string) string{
	return fmt.Sprintf("Hello %s, I am Omnibot, your friendly and versatile AI Bot. I can perform a wide range of operations, which includes image/video generation, audio/video transcription and more. You can start by chatting with me like a  friend or use the command (ex: /command_name) to perform other tasks . Have fun chatting! 😊", username)
}