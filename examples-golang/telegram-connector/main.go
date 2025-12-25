package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	vastar "github.com/fullstack-aidev/vastar-wf-connector-sdk-bin/sdk-golang"
	ipc "github.com/fullstack-aidev/vastar-wf-connector-sdk-bin/sdk-golang/protocol"
)

////////////////////////////////////////////////////////////////////////////////
// CONFIG
////////////////////////////////////////////////////////////////////////////////

const (
	TelegramBaseURL = "https://api.telegram.org"
	RequestTimeout  = 15000 // 15 detik (safe)
)

////////////////////////////////////////////////////////////////////////////////
// PAYLOAD
////////////////////////////////////////////////////////////////////////////////

type telegramMessageRequest struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

////////////////////////////////////////////////////////////////////////////////
// CONNECTOR
////////////////////////////////////////////////////////////////////////////////

type TelegramConnector struct {
	client *vastar.RuntimeClient
	token  string
	chatID string
	logger *log.Logger
}

func NewTelegramConnector(token, chatID string) (*TelegramConnector, error) {
	if token == "" || chatID == "" {
		return nil, fmt.Errorf("telegram token or chat id missing")
	}

	client, err := vastar.NewRuntimeClient()
	if err != nil {
		return nil, fmt.Errorf("vastar runtime unreachable: %w", err)
	}

	return &TelegramConnector{
		client: client,
		token:  token,
		chatID: chatID,
		logger: log.New(os.Stdout, "[TELEGRAM] ", log.LstdFlags),
	}, nil
}

func (c *TelegramConnector) Close() error {
	c.logger.Println("closing IPC connection")
	return c.client.Close()
}

////////////////////////////////////////////////////////////////////////////////
// SEND MESSAGE (PRODUCTION SAFE)
////////////////////////////////////////////////////////////////////////////////

func (c *TelegramConnector) SendMessage(message string) error {
	payload := telegramMessageRequest{
		ChatID: c.chatID,
		Text:   message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", TelegramBaseURL, c.token)

	req := vastar.POST(url).
		WithHeader("Content-Type", "application/json").
		WithBody(body).
		WithTimeout(RequestTimeout)

	start := time.Now()
	resp, err := c.client.ExecuteHTTP(req)
	if err != nil {
		return fmt.Errorf("IPC error: %w", err)
	}

	if resp.ErrorClass != ipc.ErrorClassSuccess {
		return fmt.Errorf("runtime error class: %s", resp.ErrorClass.String())
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("telegram API error %d: %s", resp.StatusCode, string(resp.Body))
	}

	c.logger.Printf("message sent successfully (latency: %v)", time.Since(start))
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// MAIN
////////////////////////////////////////////////////////////////////////////////

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("VASTAR TELEGRAM NOTIFICATION CONNECTOR (PRODUCTION)")
	fmt.Println(strings.Repeat("=", 80))

	connector, err := NewTelegramConnector(token, chatID)
	if err != nil {
		log.Fatal(err)
	}
	defer connector.Close()

	msg := fmt.Sprintf(
		"🚀 *Enterprise Status*\nTime: %s\nStatus: ONLINE\nEngine: Vastar IPC",
		time.Now().Format(time.RFC1123),
	)

	if err := connector.SendMessage(msg); err != nil {
		log.Fatal("send failed:", err)
	}

	fmt.Println("✅ Telegram notification delivered")
	fmt.Println(strings.Repeat("=", 80))
}
