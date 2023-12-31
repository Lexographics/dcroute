package utils

import "fmt"

func GetEmoji(emojiName string, emojiID string) string {
	return fmt.Sprintf("<:%s:%s>", emojiName, emojiID)
}

func GetReaction(emojiName string, emojiID string) string {
	return fmt.Sprintf("%s:%s", emojiName, emojiID)
}