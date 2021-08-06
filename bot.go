package main

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/NicoNex/echotron/v3"
)

var (
	addCommandPattern = regexp.MustCompile(`/add ([[:alpha:]а-яА-ЯёЁ]+) (\d+)`)
	delCommandPattern = regexp.MustCompile(`/del (\d+)`)
)

const wellcomeMessage = `
This bot will help you store and visualize your expences!

Bot commands:
*/add <category> <amount>* will add new expence (example: ` + "`/add taxi 200)`" + `
*/get* will show your month expences (you can delete them in place)
*/chart* will show you chart of your month expences
`

type bot struct {
	chatID int64
	echotron.API
}

func newBot(chatID int64) echotron.Bot {
	return &bot{
		chatID,
		echotron.NewAPI(token),
	}
}

func (b *bot) Update(update *echotron.Update) {
	switch {
	case update.Message != nil:
		b.handleMessage(update.Message)
	case update.CallbackQuery != nil:
		b.handleCallbackQuery(update.CallbackQuery)
	}
}

// handleMessage handles default message
func (b *bot) handleMessage(m *echotron.Message) {
	messageText := m.Text
	userID := m.From.ID
	switch {
	case messageText == "/start":
		b.sendWelcome()
	case messageText == "/get":
		b.sendExpences(userID)
	case messageText == "/chart":
		b.sendChart(userID)
	case addCommandPattern.MatchString(messageText):
		b.addExpence(messageText, userID)
	}
}

// sendWelcome handles /start command and sends
// welcome message with bot instructions
func (b *bot) sendWelcome() {
	messageOpts := echotron.MessageOptions{ParseMode: echotron.Markdown}
	b.SendMessage(wellcomeMessage, b.chatID, &messageOpts)
}

// sendExpences handles /get command and sends user his own expences
func (b *bot) sendExpences(userID int64) {
	expences, err := expenceService.GetUserExpences(userID)
	if err != nil {
		b.SendMessage("Couldn't get your expences", b.chatID, nil)
	} else if len(expences) == 0 {
		b.SendMessage("You haven't any expences yet", b.chatID, nil)
	} else {
		replyMarkup := GetExpencesReplyMarkup(expences)
		messageOptions := getReplyMarkupMessageOptions(replyMarkup)
		b.SendMessage("Your current month expences:", b.chatID, &messageOptions)
	}
}

// getReplyMarkupMessageOptions wraps reply markup in echotron.MessageOptions
// that is required to send messages with reply markup
func getReplyMarkupMessageOptions(rm echotron.ReplyMarkup) echotron.MessageOptions {
	return echotron.MessageOptions{
		BaseOptions: echotron.BaseOptions{ReplyMarkup: rm},
	}
}

// sendChart handles /chart command and sends user his own expences chart
func (b *bot) sendChart(userID int64) {
	expences, err := expenceService.GetUserExpencesGroupedByCategory(userID)
	if err != nil {
		b.SendMessage("Couldn't make chart", b.chatID, nil)
	} else if len(expences) == 0 {
		b.SendMessage("You haven't any expences yet", b.chatID, nil)
	} else {
		buf, err := RenderExpencesGroupedByCategory(expences)
		if err != nil {
			b.SendMessage("Error while rendering expences chart", b.chatID, nil)
		} else {
			f := echotron.NewInputFileBytes("chart.png", buf.Bytes())
			b.SendPhoto(f, b.chatID, nil)
		}
	}
}

// addExpence handles /add expence command and adds new expence with given params
func (b *bot) addExpence(messageText string, userID int64) {
	category, amount := parseAddExpenceCommand(messageText)
	category = strings.ToLower(category)

	if _, err := expenceService.InsertExpence(userID, category, amount); err != nil {
		b.SendMessage("Error while adding expence", b.chatID, nil)
	} else {
		b.SendMessage("Expence was added", b.chatID, nil)
	}
}

// parseAddExpenceCommand parses /add expence command
func parseAddExpenceCommand(commandText string) (category string, amount int) {
	submatch := addCommandPattern.FindStringSubmatch(commandText)
	category = submatch[1]
	amountI64, _ := strconv.ParseInt(submatch[2], 10, 64)
	amount = int(amountI64)

	return
}

// handleCallbackQuery handles callback queries
func (b *bot) handleCallbackQuery(c *echotron.CallbackQuery) {
	if delCommandPattern.MatchString(c.Data) {
		b.delExpence(c)
	}
}

// delExpence handles /del callback query and deletes expence with passed expence ID
func (b *bot) delExpence(c *echotron.CallbackQuery) {
	expenceID := parseDelCommand(c.Data)
	err := expenceService.DeleteExpence(expenceID)

	if err != nil {
		b.SendMessage("Error while deleting expence", b.chatID, nil)
	} else {
		b.updateExpencesMarkup(c.Message, c.From.ID)
		b.SendMessage("Expence was deleted", b.chatID, nil)
	}
}

// parseDelCommand parses /del expence command
func parseDelCommand(commandText string) (expenceID int64) {
	submatch := delCommandPattern.FindStringSubmatch(commandText)
	expenceID, _ = strconv.ParseInt(submatch[1], 10, 64)

	return
}

// updateExpencesMarkup updates expences markup with new data
// due to deleting one of the expences
func (b *bot) updateExpencesMarkup(m *echotron.Message, userID int64) {
	expences, err := expenceService.GetUserExpences(userID)
	if err != nil {
		return
	}
	replyMarkup := GetExpencesReplyMarkup(expences)
	opts := echotron.MessageReplyMarkup{ReplyMarkup: replyMarkup}
	message := echotron.NewMessageID(b.chatID, m.ID)
	b.EditMessageReplyMarkup(message, &opts)
}
