package telegram

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"regexp"
)

var commands []*BotCommand
var regexCommands []*RegexCommand
var callbacks []*BotCommand

type BotCommand struct {
	MatchCmd   *regexp.Regexp
	Chan       int64
	HandleFunc func(tgbotapi.Update, *tgbotapi.BotAPI)
}

type RegexCommand struct {
	MatchCmd   *regexp.Regexp
	Chan       int64
	HandleFunc func(tgbotapi.Update, *tgbotapi.BotAPI, map[string]string)
}

func TestCmd(updateIn tgbotapi.Update, botIn *tgbotapi.BotAPI) {
	c := tgbotapi.NewMessage(updateIn.Message.Chat.ID, "Butts")
	botIn.Send(c)
}

func Register(regexIn string, chanIn int64, handleFunc func(tgbotapi.Update, *tgbotapi.BotAPI)) {
	if commands == nil {
		commands = []*BotCommand{}
	}
	newCommand := &BotCommand{}
	newCommand.MatchCmd = regexp.MustCompile(regexIn)
	newCommand.Chan = chanIn
	newCommand.HandleFunc = handleFunc
	commands = append(commands, newCommand)
}
func RegisterRegex(regexIn string, chanIn int64, handleFunc func(tgbotapi.Update, *tgbotapi.BotAPI, map[string]string)) {
	if regexCommands == nil {
		regexCommands = []*RegexCommand{}
	}
	newCommand := &RegexCommand{}
	newCommand.MatchCmd = regexp.MustCompile(regexIn)
	newCommand.Chan = chanIn
	newCommand.HandleFunc = handleFunc
	regexCommands = append(regexCommands, newCommand)
}
func RegisterCallback(regexIn string, handleFunc func(tgbotapi.Update, *tgbotapi.BotAPI)) {
	if callbacks == nil {
		callbacks = []*BotCommand{}
	}
	newCallback := &BotCommand{}
	newCallback.MatchCmd = regexp.MustCompile(regexIn)
	newCallback.HandleFunc = handleFunc
	callbacks = append(callbacks, newCallback)
}
func ProcessMessage(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	message := upd.Message.Text
	if message == "" {
		return
	}
	for _, cmd := range commands {
		if cmd.MatchCmd.Match([]byte(message)) {
			if cmd.Chan == upd.Message.Chat.ID || cmd.Chan == 0 {
				go cmd.HandleFunc(upd, bot)
			}
		}
	}
	for _, cmd := range regexCommands {
		if cmd.MatchCmd.Match([]byte(message)) {
			if cmd.Chan == upd.Message.Chat.ID || cmd.Chan == 0 {
				go cmd.HandleFunc(upd, bot, resolveRegexCaptureGroups(cmd.MatchCmd, message))
			}
		}
	}
}
func resolveRegexCaptureGroups(regex *regexp.Regexp, text string) map[string]string {
	match := regex.FindStringSubmatch(text)
	result := make(map[string]string)
	for i, name := range regex.SubexpNames() {
		if i != 0 {
			result[name] = match[i]
		}
	}
	return result
}
func ProcessCallback(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	data := upd.CallbackQuery.Data
	if data == "" {
		return
	}
	for _, cback := range callbacks {
		if cback.MatchCmd.Match([]byte(data)) {
			go cback.HandleFunc(upd, bot)
		}
	}
}
func InitBot(botToken string) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}
	bot.Debug = false
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}
	for update := range updates {
		if update.Message != nil {
			outLog := fmt.Sprintf("Message: %s %s>%s", update.Message.From.FirstName, update.Message.From.LastName, update.Message.Text)
			fmt.Println(outLog)
			//fmt.Println(update.Message.Chat.ID)
			ProcessMessage(update, bot)
		} else if update.CallbackQuery != nil {
			fmt.Println("Callback handled: " + update.CallbackQuery.Data)
			ProcessCallback(update, bot)
			//config := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			//fmt.Println(update.CallbackQuery.Data)
			//bot.AnswerCallbackQuery(config)
		} else {
			continue
		}
	}
}
