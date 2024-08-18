package config

import (
	"context"
    tele "gopkg.in/telebot.v3"
	"go.uber.org/zap"
	"github.com/kingmariano/omnicron-go"
	"github.com/kingmariano/omnicron-telebot/internal/database"
)

type BotConfig struct {
	DB            *database.Queries
	Bot           *tele.Bot
	Context       context.Context
	BotToken      string
	ProviderToken string
	MyAPIKey      string
	DBURL         string
	Logger        *zap.Logger
	ChatHistory   map[int64][]omnicron.Message // map to store chat history with user ID as the key
}