package telegram

import (
	"bytes"
	"github.com/fortinj1354/Dex-Go/metrics"
	"github.com/fortinj1354/Dex-Go/models"
	"github.com/fortinj1354/Dex-Go/settings"
	"gopkg.in/telegram-bot-api.v4"
	"math/rand"
	"strconv"
)

//Help response for main channel
func MainChannelHelp(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	newMsg := tgbotapi.NewMessage(settings.GetChannelID(), "/mods - Call the chat moderators\n/ask <word> - Ask the bot about a word\n<anything> c/d - Ask the bot the hard questions in life\n/roll <dice>d<sides> - Get a dice roll")
	bot.Send(newMsg)
}

//Help response for control channel
func ControlChannelHelp(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	newMsg := tgbotapi.NewMessage(settings.GetControlID(), "/info <@username OR TelegramID> - Get information about a username\n/warn <@username OR TelegramID> <warning message> - Record a warning for a user\n/find <display name> - Find a user by their display name\n/status - Get bot status information\n/count - Get the number of chains in the markov database")
	bot.Send(newMsg)
}

//Handles the /mods command
func SummonMods(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	user := models.ChatUserFromTGID(upd.Message.From.ID, upd.Message.From.UserName)
	if user.PingAllowed {
		newFwd := tgbotapi.NewForward(settings.GetControlID(), upd.Message.Chat.ID, upd.Message.MessageID)
		bot.Send(newFwd)
		newMsg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Summoning mods!")
		bot.Send(newMsg)
	} else {
		newMsg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Sorry, you are banned from /mods")
		bot.Send(newMsg)
	}
}

//Returns uptime of the bot
func GetBotStatus(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetControlID() {
		newMess := tgbotapi.NewMessage(settings.GetControlID(), "Time since startup: "+metrics.TimeSinceStart().String())
		bot.Send(newMess)
	}
}

//c/d questions
func YesOrNo(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetChannelID() {
		if rand.Intn(2) == 0 {
			bot.Send(tgbotapi.NewMessage(settings.GetChannelID(), "c"))
		} else {
			bot.Send(tgbotapi.NewMessage(settings.GetChannelID(), "d"))
		}
	}
}

func DiceRoll(upd tgbotapi.Update, bot *tgbotapi.BotAPI, regexMatch map[string]string) {
	if upd.Message.Chat.ID == settings.GetChannelID() {
		dice, err := strconv.Atoi(regexMatch["dice"])
		if err != nil {
			panic(err)
		}
		sides, err := strconv.Atoi(regexMatch["sides"])
		if err != nil {
			panic(err)
		}

		if dice > 0 && sides > 0 {
			var outMessage tgbotapi.MessageConfig
			if dice > 10000 {
				outMessage = tgbotapi.NewMessage(settings.GetChannelID(), "Too many dice")
			} else if dice > 50 || sides > 100000 {
				diceSum := 0
				for i := 0; i < dice; i++ {
					diceSum += rand.Intn(sides) + 1
				}
				outMessage = tgbotapi.NewMessage(settings.GetChannelID(), "Total: "+strconv.Itoa(diceSum))
			} else if dice == 1 {
				outMessage = tgbotapi.NewMessage(settings.GetChannelID(), "Roll: "+strconv.Itoa(rand.Intn(sides)+1))
			} else {
				var rolls bytes.Buffer
				diceSum := 0
				for i := 0; i < dice; i++ {
					roll := rand.Intn(sides) + 1
					diceSum += roll
					if i != 0 {
						rolls.WriteString(", ")
					}
					rolls.WriteString(strconv.Itoa(roll))
				}
				outMessage = tgbotapi.NewMessage(settings.GetChannelID(), "Total: "+strconv.Itoa(diceSum)+" Rolls: "+rolls.String())
			}
			bot.Send(outMessage)
		}
	}
}
