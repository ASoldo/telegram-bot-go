// go run main.go --bot-token="6015443726:AAG8DVXj4o9yLPw5aE4hF93iYPOVQRUylaI" --chat-id=6229440871 --clusters-json=./clusters.json
package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type KeywordGroup struct {
	Keywords []string `json:"keywords"`
}

type SemanticCluster struct {
	ChatID        int64          `json:"chat_id"`
	KeywordGroups []KeywordGroup `json:"keyword_groups"`
	NegativeWords []string       `json:"negative_words"`
	MinPrice      int            `json:"min_price"`
	MaxPrice      int            `json:"max_price"`
}

func main() {
	botToken := flag.String("bot-token", "", "Telegram Bot API token")
	chatID := flag.Int64("chat-id", 0, "Telegram chat ID to forward messages to")
	clustersJSONFile := flag.String(
		"clusters-json",
		"",
		"path to JSON file containing semantic clusters",
	)
	flag.Parse()

	if *botToken == "" {
		log.Fatal("missing bot token argument")
	}

	if *chatID == 0 {
		log.Fatal("missing chat ID argument")
	}

	clustersJSON, err := os.ReadFile(*clustersJSONFile)
	if err != nil {
		log.Fatal(err)
	}

	clusters := make([]SemanticCluster, 0)

	err = json.Unmarshal(clustersJSON, &clusters)
	if err != nil {
		log.Fatal(err)
	}

	if len(clusters) == 0 {
		log.Fatal("no semantic clusters defined")
	}

	bot, err := tgbotapi.NewBotAPI(*botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		for _, cluster := range clusters {
			if isRelevantMessage(cluster, update.Message) {
				forwardMessage(bot, cluster, update.Message, *chatID)
			}
		}
	}
}

func isRelevantMessage(cluster SemanticCluster, message *tgbotapi.Message) bool {
	log.Printf("Message from %s chat with type %s", message.Chat.Title, message.Chat.Type)
	// ignore messages from non-group or non-channel chats
	if message.Chat.Type != "group" && message.Chat.Type != "channel" {
		return false
	}
	for _, keywordGroup := range cluster.KeywordGroups {
		containsKeywordGroup := false
		for _, keyword := range keywordGroup.Keywords {
			if strings.Contains(message.Text, keyword) {
				containsKeywordGroup = true
				break
			}
		}
		if !containsKeywordGroup {
			return false
		}
	}

	for _, negativeWord := range cluster.NegativeWords {
		if strings.Contains(message.Text, negativeWord) {
			return false
		}
	}

	priceRegex := regexp.MustCompile(`(\d+)\$`)
	matches := priceRegex.FindStringSubmatch(message.Text)
	if len(matches) > 1 {
		price, _ := strconv.Atoi(matches[1])
		if price < cluster.MinPrice || price > cluster.MaxPrice {
			return false
		}
	}

	return true
}

func forwardMessage(
	bot *tgbotapi.BotAPI,
	cluster SemanticCluster,
	message *tgbotapi.Message,
	chatID int64,
) {
	targetChatID := chatID

	if message.Chat.Type == "channel" {
		targetChatID = cluster.ChatID
	}

	msg := tgbotapi.NewForward(targetChatID, message.Chat.ID, message.MessageID)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error forwarding message: %v", err)
	}
}
