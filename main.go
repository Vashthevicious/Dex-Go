package main

import (
	"github.com/fortinj1354/Dex-Go/metrics"
	"github.com/fortinj1354/Dex-Go/models"
	"github.com/fortinj1354/Dex-Go/settings"
	"github.com/fortinj1354/Dex-Go/telegram"
	"math/rand"
	"time"
)

func main() {
	settings.LoadSettings()
	models.MakeDB(settings.GetDBAddr())
	metrics.StartUp()
	rand.Seed(time.Now().UnixNano())

	//Any channel commands
	telegram.Register(`^\/ping`, 0, telegram.TestCmd)
	telegram.Register(`^\/mods`, 0, telegram.SummonMods)

	//Main channel commands
	telegram.Register(`^\/help`, settings.GetChannelID(), telegram.MainChannelHelp)
	telegram.Register(`^\/ask .+`, settings.GetChannelID(), telegram.MarkovTalkAbout)
	telegram.Register(`.* c/d$`, settings.GetChannelID(), telegram.YesOrNo)
	telegram.RegisterRegex(`\/roll (?P<dice>\d+)d(?P<sides>\d+).*`, settings.GetChannelID(), telegram.DiceRoll)
	telegram.Register(`.*`, settings.GetChannelID(), telegram.HandleUsers)
	telegram.Register(`.*`, settings.GetChannelID(), telegram.HandleMarkov)

	//Control channel commands
	telegram.Register(`^\/help`, settings.GetControlID(), telegram.ControlChannelHelp)
	telegram.Register(`^\/info @.+`, settings.GetControlID(), telegram.FindUserByUsername)
	telegram.Register(`^\/info \d+`, settings.GetControlID(), telegram.FindUserByUserID)
	telegram.Register(`^\/warn @.+ .+`, settings.GetControlID(), telegram.WarnUserByUsername)
	telegram.Register(`^\/warn \d+ .+`, settings.GetControlID(), telegram.WarnUserByID)
	telegram.Register(`^\/find .+`, settings.GetControlID(), telegram.LookupAlias)
	telegram.Register(`^\/status`, settings.GetControlID(), telegram.GetBotStatus)
	telegram.Register(`^\/count`, settings.GetControlID(), telegram.MarkovCount)

	//Callbacks
	telegram.RegisterCallback(`^\/togglemods \d+`, telegram.ToggleMods)
	telegram.RegisterCallback(`^\/getwarnings \d+`, telegram.DisplayWarnings)
	telegram.RegisterCallback(`^\/getaliases \d+`, telegram.DisplayAliases)
	telegram.RegisterCallback(`^\/ban \d+`, telegram.PreBan)
	telegram.RegisterCallback(`^\/callbackinfo \d+`, telegram.CallbackInfo)
	telegram.RegisterCallback(`^\/preconfirmban \d+`, telegram.PreConfirmBan)
	telegram.RegisterCallback(`^\/confirmban \d+`, telegram.ConfirmBan)
	telegram.RegisterCallback(`^\/togglemarkov \d+`, telegram.ToggleMarkov)

	telegram.InitBot(settings.GetBotToken())
}
