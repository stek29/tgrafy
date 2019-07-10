package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

var (
	TelegramAPIURL = "https://api.telegram.org/bot%s/%s"
)

type Bot struct {
	token  string
	client *http.Client
}

type BotConfig struct {
	Token  string
	Client *http.Client
}

func NewBot(cfg BotConfig) *Bot {
	client := cfg.Client
	if client == nil {
		client = http.DefaultClient
	}
	return &Bot{
		token:  cfg.Token,
		client: client,
	}
}

type SendMessageRequest struct {
	ChatID                string `json:"chat_id"`
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode,omitempty"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview,omitempty"`
	DisableNotification   bool   `json:"disable_notification,omitempty"`
	ReplyToMessageID      int    `json:"reply_to_message_id,omitempty"`
}

func (b *Bot) doRequest(method string, payload interface{}) error {
	url := fmt.Sprintf(TelegramAPIURL, b.token, method)

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	zap.L().Debug("response", zap.ByteString("resp", respBody))

	return nil
}

func (b *Bot) SendMessage(r *SendMessageRequest) error {
	return b.doRequest("sendMessage", r)
}
