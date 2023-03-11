package main

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func TestIsRelevantMessage(t *testing.T) {
	// Define a test message
	message := tgbotapi.Message{
		Text: "This is a test message containing a keyword and a price of 50$.",
		Chat: &tgbotapi.Chat{
			ID:   456,
			Type: "group",
		},
	}

	// Define a test semantic cluster
	cluster := SemanticCluster{
		KeywordGroups: []KeywordGroup{
			{
				Keywords: []string{"test"},
			},
			{
				Keywords: []string{"keyword"},
			},
		},
		NegativeWords: []string{"spam"},
		MinPrice:      10,
		MaxPrice:      100,
	}

	// Test that the message is relevant to the semantic cluster
	if !isRelevantMessage(cluster, &message) {
		t.Errorf("Expected message to be relevant to cluster, but it was not.")
	}

	// Modify the test message to include a negative word
	message.Text += " This message is spam."
	if isRelevantMessage(cluster, &message) {
		t.Errorf("Expected message to be irrelevant due to negative word, but it was not.")
	}

	// Modify the test message to include a price outside the specified range
	message.Text = "This is a test message containing a keyword and a price of 500$."
	if isRelevantMessage(cluster, &message) {
		t.Errorf("Expected message to be irrelevant due to price outside range, but it was not.")
	}

	// Modify the test message to be from a non-group chat
	message.Chat.Type = "private"
	if isRelevantMessage(cluster, &message) {
		t.Errorf("Expected message to be irrelevant due to non-group chat, but it was not.")
	}
}
