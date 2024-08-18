package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/kingmariano/omnicron-go"
	"github.com/kingmariano/omnicron-telebot/commands"
	"github.com/kingmariano/omnicron-telebot/config"
	"github.com/kingmariano/omnicron-telebot/internal/database"
	"github.com/kingmariano/omnicron-telebot/middleware"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

// loadenv loads environment variables from a .env file.
func loadenv() (string, string, string, string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		return "", "", "", "", errors.New("unable to get bot token")
	}
	myAPIKey := os.Getenv("MY_API_KEY")
	if myAPIKey == "" {
		return token, "", "", "", errors.New("unable to get API Key")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return token, myAPIKey, "", "", errors.New("unable to get database URL")
	}
	providerToken := os.Getenv("PROVIDER_TOKEN")
	if providerToken == "" {
		return token, myAPIKey, dbURL, "", errors.New("unable to get provider token")
	}
	return token, myAPIKey, dbURL, providerToken, nil
}

// main is the entry point of the application.
func main() {
	token, myAPIKey, dbURL, providerToken, err := loadenv()
	logger := GetLogger()
	if err != nil {
		logger.Fatal("failed to load environment variables", zap.String("error", err.Error()))
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.String("error", err.Error()))
	}
	defer conn.Close()

	queries := database.New(conn)

	bot, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
		Verbose: true,
		ParseMode: tele.ModeMarkdown,
	})
	if err != nil {
		logger.Fatal("failed to create bot", zap.String("error", err.Error()))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	botConfig := &config.BotConfig{
		Bot:           bot,
		DB:            queries,
		Context:       ctx,
		BotToken:      token,
		MyAPIKey:      myAPIKey,
		DBURL:         dbURL,
		ProviderToken: providerToken,
		Logger:        &logger,
		ChatHistory:   make(map[int64][]omnicron.Message), // Initialize the chat history map
	}
    botConfig.Bot.Use(middleware.HandleUserLimitReached(botConfig))
	// Register command handlers
	registerHandlers(botConfig)

	// Channel to listen for termination signals (graceful shutdown)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		logger.Info("shutting down bot")
		cancel()       // Cancel the context to stop background operations
		bot.Stop()     // Stop the bot gracefully
	}()

	logger.Info("starting bot")
	botConfig.Bot.Start()
}

// registerHandlers registers all the command handlers for the bot.
func registerHandlers(botConfig *config.BotConfig) {
	botConfig.Bot.Handle("/start", commands.Start(botConfig))
	botConfig.Bot.Handle("/invoice", commands.Invoice(botConfig))
	botConfig.Bot.Handle("/points", commands.GetPoints(botConfig))
	botConfig.Bot.Handle("/generate_image", commands.ImageGeneration(botConfig))
	botConfig.Bot.Handle(tele.OnText, commands.Chat(botConfig))
	botConfig.Bot.Handle(tele.OnCheckout, commands.Checkout(botConfig))
	botConfig.Bot.Handle(tele.OnPayment, commands.Payment(botConfig))
}
