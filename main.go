package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const notAuthorizedMessage = "Not authorized: Use ChatID:BotToken"

func alertHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	if len(s) != 2 {
		http.Error(w, notAuthorizedMessage, http.StatusUnauthorized)
		return
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		http.Error(w, notAuthorizedMessage, http.StatusUnauthorized)
		return
	}

	chatID := pair[0]
	token := pair[1]

	fmt.Printf("%v : %v\n", chatID, token)

	defer r.Body.Close()
	var alert AlertBody
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&alert); err != nil {
		http.Error(w, "Invalid Alert Body: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = handleAlert(chatID, token, &alert)
	if err != nil {
		http.Error(w, "Failed to handle alert", http.StatusBadRequest)
		zap.L().Error("handleAlert fail", zap.Error(err))
		return
	}

	w.Write([]byte("Ok!\n"))
}

func maskBotToken(tok string) string {
	split := strings.SplitN(tok, ":", 2)
	return split[0] + ":*****"
}

func handleAlert(chatID, botToken string, alert *AlertBody) error {
	zap.L().Info(
		"handling alert",
		zap.String("chat_id", chatID),
		zap.String("token", maskBotToken(botToken)),
		zap.String("alert_title", alert.Title))

	err := NewBot(BotConfig{
		Token: botToken,
	}).SendMessage(&SendMessageRequest{
		ChatID:    chatID,
		Text:      buildMessage(alert),
		ParseMode: "html",
	})

	zap.L().Info("result", zap.Error(err))

	return nil
}

func initLogger() {
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.CallerKey = ""
	logger, _ := config.Build()
	zap.ReplaceGlobals(logger)
}

func main() {
	initLogger()

	var listenAddr string

	flag.StringVar(&TelegramAPIURL, "tg-api-url", TelegramAPIURL, "Format string for Telegram API URL")
	flag.StringVar(&listenAddr, "listen", ":80", "Listen on [addr]:port")

	flag.Parse()

	zap.L().Info("starting server",
		zap.String("listen", listenAddr),
		zap.String("tg_api_url", TelegramAPIURL))
	http.HandleFunc("/alert", alertHandler)

	http.ListenAndServe(listenAddr, nil)
}
