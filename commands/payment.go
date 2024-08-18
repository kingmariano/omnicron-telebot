package commands

import (
	"context"
	"database/sql"
	"math"
	"time"
   "github.com/kingmariano/omnicron-telebot/config"
	"github.com/google/uuid"
	"github.com/kingmariano/omnicron-telebot/internal/database"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func Invoice(b *config.BotConfig) tele.HandlerFunc {

	return func(c tele.Context) error {
		amount := 30000 // Amount in smallest units of currency

		price := tele.Price{
			Label:  "One-time Subscription",
			Amount: amount,
		}

		invoice := &tele.Invoice{
			Title:       "One-time Payment",
			Description: "Get a lifetime subscription to our service.",
			Payload:     c.Sender().Username, // Identify the user
			Token:       b.ProviderToken,
			Currency:    "USD",
			Prices:      []tele.Price{price},
			Start:       "start",
			Photo: &tele.Photo{
				File:    tele.FromDisk("omnicron.jpg"),
				Caption: "One time payment to access all the premium service omnibot offers.",
			},
			NeedName:            true,
			NeedPhoneNumber:     true,
			NeedEmail:           true,
			NeedShippingAddress: true,
		}
		return c.Send(invoice)
	}
}
func Checkout(b *config.BotConfig) tele.HandlerFunc {
	return func(c tele.Context) error {
		return c.Accept()
	}
}
func Payment(b *config.BotConfig) tele.HandlerFunc {
	return func(c tele.Context) error {
		ctx, cancel := context.WithTimeout(b.Context, time.Minute)
		defer cancel()

		payment := c.Message().Payment
		username := payment.Payload
		telegramChargeID := payment.TelegramChargeID
		providerChargeID := payment.ProviderChargeID
		user, err := b.DB.GetUserByUsername(ctx, username)
		if err != nil {
			b.Logger.Error("an error occurred while getting user by username:", zap.String("error", err.Error()))
			return c.Send("an error occurred")
		}
		//update sender subscription status
		err = b.DB.UpdateUserSubscriptionStatus(ctx, database.UpdateUserSubscriptionStatusParams{
			Points: sql.NullInt32{
				Int32: math.MaxInt32,
				Valid: true,
			},
			TelegramID: user.TelegramID,
		})
		if err != nil {
			b.Logger.Error("an error occurred while updating subscription status:", zap.String("error", err.Error()))
			return c.Send("an error occurred while updating subscription status 😢")
		}
		// add sender to the list of subscribers

		b.DB.AddUserToSubscription(ctx, database.AddUserToSubscriptionParams{
			ID:               uuid.New(),
			UserName:         user.UserName,
			TelegramID:       user.TelegramID,
			TelegramChargeID: telegramChargeID,
			ProviderChargeID: providerChargeID,
		})
		return c.Send("Thank you for your payment! Your subscription is now active.")
	}
}
