package commands

import (
	"context"
	"fmt"
	"time"
	"github.com/kingmariano/omnicron-telebot/config"
	tele "gopkg.in/telebot.v3"
)
func GetPoints(b *config.BotConfig) tele.HandlerFunc{
   return func(c tele.Context) error{
	    ctx, cancel := context.WithTimeout(b.Context, time.Minute)
		defer cancel()
		username := c.Sender().Username
		user, err := b.DB.GetUserByUsername(ctx, username)
		if err != nil{
			return err
		}
		points := user.Points.Int32
		return c.Reply(fmt.Sprintf("You have %d points", points))
   }
}