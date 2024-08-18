package middleware

import (
	"context"
	"time"
    "github.com/kingmariano/omnicron-telebot/config"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)


func  HandleUserLimitReached(b *config.BotConfig) tele.MiddlewareFunc{
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			ctx, cancel := context.WithTimeout(b.Context, time.Minute)
		   defer cancel()
		   username := c.Sender().Username
			user, err := b.DB.GetUserByUsername(ctx, username)
			if err != nil{
				b.Logger.Warn("Error getting user", zap.String("error", err.Error()))
				return c.Send("An error occurred, couldn't fetch user. Please register by using the /start and try again.")
			}
			if !user.IsSubscribed.Bool && user.Points.Int32==0{
				b.Logger.Info("Error user limit has been reached")
                   return c.Send("you have reached your limit, buy premium with the /invoice to continue enjoying the service. ")
			}
			return next(c)
		}
	}
}
