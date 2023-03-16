package main

import (
	"fmt"
	"strconv"

	"github.com/NicoNex/echotron/v3"
)

// GetExpencesReplyMarkup converts expences to inline keyboard markup
func GetExpencesReplyMarkup(expences []Expence) echotron.InlineKeyboardMarkup {
	rows := make([][]echotron.InlineKeyboardButton, 0, len(expences))
	for _, e := range expences {
		row := expenceToInlineButtonRow(e)
		rows = append(rows, row)
	}
	return echotron.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

// expenceToInlineButtonRow converts single expence to inline buttons row
func expenceToInlineButtonRow(e Expence) []echotron.InlineKeyboardButton {
	return []echotron.InlineKeyboardButton{
		{Text: e.Category, CallbackData: "1"},
		{Text: strconv.Itoa(e.Amount), CallbackData: "2"},
		{Text: "❌", CallbackData: fmt.Sprintf("/del %v", e.Id)},
	}
}
