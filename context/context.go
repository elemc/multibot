package context

import (
	"io/ioutil"
	"os"

	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

const telegramMaximumMessageSize = 4096

// Options is a type for store all application options
type Options struct {
	AppName   string
	APIKey    string
	PgSQLDSN  string
	LogLevel  string
	Debug     bool
	PluginDir string
}

// MultiBotContext is a struct with methods for interact bot with plugins
type MultiBotContext struct {
	db      *pg.DB
	bot     *tgbotapi.BotAPI
	options *Options
}

// InitContext initialize context and return it pointer
func InitContext(db *pg.DB, bot *tgbotapi.BotAPI, options *Options) *MultiBotContext {
	mbc := &MultiBotContext{
		db:      db,
		bot:     bot,
		options: options,
	}
	return mbc
}

// SendMessage send message from bot to chat with ID == chatID and text
// if replyID != 0 message send as reply
func (ctx *MultiBotContext) SendMessage(chatID int64, text string, replyID int) {
	var (
		err error
	)

	if len(text) > telegramMaximumMessageSize {
		log.Debugf("Message to big, size %d, send as file", len(text))
		ctx.SendMessage(chatID, "* Сообщение слишком большое. Текст будет отправлен в виде файла! *", replyID)
		filename := ctx.saveTextToFileAndGetName(text)
		ctx.SendFile(chatID, replyID, filename, "text/plain")
		return
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	if replyID != 0 {
		msg.ReplyToMessageID = replyID
	}
	if _, err = ctx.bot.Send(msg); err != nil {
		// oops, try to send as plain text
		log.Warnf("oops, unable to send markdown message [%s]: %s. Try to send as plain text.", text, err)
		msg.ParseMode = ""
		if _, err = ctx.bot.Send(msg); err != nil {
			log.Errorf("Unable to send message to %d with text [%s] and reply [%d]: %s", chatID, text, replyID, err)
			return
		}
	}
}

// SendFile send file with file name as filename from bot to chat with ID == chatID
// if replyID != 0 message send as reply
func (ctx *MultiBotContext) SendFile(chatID int64, replyID int, filename string, mime string) {
	var (
		err error
	)
	msg := tgbotapi.NewDocumentUpload(chatID, filename)
	msg.MimeType = mime

	if replyID <= 0 {
		msg.ReplyToMessageID = replyID
	}

	if _, err = ctx.bot.Send(msg); err != nil {
		log.Errorf("Unable to send file: %s", err)
	}
}

func (ctx *MultiBotContext) saveTextToFileAndGetName(msgText string) string {
	var (
		err error
		f   *os.File
	)
	if f, err = ioutil.TempFile("", ctx.options.AppName); err != nil {
		log.Errorf("Unable to open temporary file: %s", err)
		return ""
	}
	defer f.Close()
	if _, err = f.WriteString(msgText); err != nil {
		log.Errorf("Unable to write text to file: %s", err)
		return ""
	}
	return f.Name()
}
