package telegramApi

import (
	"fmt"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/oleksiy-os/porto-events/internal/model"
	log "github.com/sirupsen/logrus"
	"strings"
)

type (
	Bot struct {
		bot    *tgbot.BotAPI
		config Telegram
	}

	Telegram struct {
		ApiToken    string `toml:"bot_api_token"`
		ChannelId   string `toml:"channel_id"`
		ChannelName string `toml:"channel_name"`
	}
)

func (t *Bot) Publish(e *model.Event) error {
	msg := `<b><a href="%s">%s</a></b> &#10;%s &#10;📍 <a href="%s">%s</a> &#10;🗓 %s &#10;🕒 %s &#10;%s`
	e.Description = truncateString(e, &msg)
	msg = fmt.Sprintf(msg, e.Url, e.Title, e.Description, e.LocationMap, e.Place, e.DateText, e.Time, e.Days)

	photo := tgbot.FileURL(e.Image)
	cnf := tgbot.NewPhotoToChannel(t.config.ChannelId, photo)
	cnf.ParseMode = "HTML"
	cnf.DisableNotification = true
	cnf.Caption = msg
	if _, err := t.bot.Send(cnf); err != nil {
		log.Errorln("error send message", err)
		return err
	}

	return nil
}

// truncateString description to telegram limit 1024 characters for cation(text message) with photo
//
// https://core.telegram.org/bots/api#inputmediaphoto
func truncateString(ev *model.Event, msg *string) string {
	limit := 1024
	d := ev.Description

	fieldsLen := len(*msg) +
		len(ev.Url) +
		len(ev.Title) +
		len(ev.LocationMap) +
		len(ev.Place) +
		len(ev.DateText) +
		len(ev.Time) +
		len(ev.Days)

	if fieldsLen+len(d) <= limit {
		return d // no need changes, all fields len less than limit
	}

	limit -= fieldsLen

	cutToLastDot := strings.LastIndex(d[:limit], ".") + 1

	return d[:cutToLastDot]
}

func New(config Telegram) *Bot {
	bot, err := tgbot.NewBotAPI(config.ApiToken)
	if err != nil {
		log.Fatal(err)
	}

	// Set this to true to log all interactions with telegram servers
	bot.Debug = true

	return &Bot{
		bot:    bot,
		config: config,
	}
}
