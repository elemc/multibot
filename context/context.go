package context

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	log "github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

const telegramMaximumMessageSize = 4096

// Options is a type for store all application options
type Options struct {
	AppName         string
	APIKey          string
	PgSQLDSN        string
	LogLevel        string
	Debug           bool
	PluginDir       string
	PluginsSettings map[string]map[string]interface{}
}

// MultiBotContext is a struct with methods for interact bot with plugins
type MultiBotContext struct {
	db      *pg.DB
	bot     *tgbotapi.BotAPI
	options *Options
	log     *log.Logger
}

// InitContext initialize context and return it pointer
func InitContext(db *pg.DB, bot *tgbotapi.BotAPI, options *Options, l *log.Logger) *MultiBotContext {
	mbc := &MultiBotContext{
		db:      db,
		bot:     bot,
		options: options,
		log:     l,
	}
	return mbc
}

// SendMessageText send message from bot to chat with ID == chatID and text
// if replyID != 0 message send as reply
func (ctx *MultiBotContext) SendMessageText(chatID int64, text string, replyID int, replyMarkup interface{}) {
	ctx.sendMessage(chatID, text, replyID, "", replyMarkup)
}

// SendMessageMarkdown send message from bot to chat with ID == chatID and text
// if replyID != 0 message send as reply
func (ctx *MultiBotContext) SendMessageMarkdown(chatID int64, text string, replyID int, replyMarkup interface{}) {
	ctx.sendMessage(chatID, text, replyID, "Markdown", replyMarkup)
}

func (ctx *MultiBotContext) sendMessage(chatID int64, text string, replyID int, parseMode string, replyMarkup interface{}) {
	var (
		err error
	)

	if len(text) > telegramMaximumMessageSize {
		log.Debugf("Message to big, size %d, send as file", len(text))
		ctx.SendMessageMarkdown(chatID, "* Сообщение слишком большое. Текст будет отправлен в виде файла! *", replyID, replyMarkup)
		filename := ctx.saveTextToFileAndGetName(text)
		ctx.SendFile(chatID, replyID, filename, "text/plain")
		return
	}

	msg := tgbotapi.NewMessage(chatID, text)
	if replyMarkup != nil {
		msg.ReplyMarkup = replyMarkup
	}
	msg.ParseMode = parseMode
	if replyID != 0 {
		msg.ReplyToMessageID = replyID
	}
	if _, err = ctx.bot.Send(msg); err != nil && parseMode != "" {
		// oops, try to send as plain text
		log.Warnf("oops, unable to send markdown message [%s]: %s. Try to send as plain text.", text, err)
		msg.ParseMode = ""
		if _, err = ctx.bot.Send(msg); err != nil {
			log.Errorf("Unable to send message to %d with text [%s] and reply [%d]: %s", chatID, text, replyID, err)
			return
		}
	} else if err != nil {
		log.Errorf("Unable to send plain message [%s]: %s", text, err)
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

// DBCreateTable create table in database by struct
func (ctx *MultiBotContext) DBCreateTable(data interface{}) (err error) {
	err = ctx.db.CreateTable(data, &orm.CreateTableOptions{IfNotExists: true})
	return
}

// GetDB return database pointer
func (ctx *MultiBotContext) GetDB() *pg.DB {
	return ctx.db
}

// Log return main log pointer
func (ctx *MultiBotContext) Log() *log.Logger {
	return ctx.log
}

// GetOptions return options pointer
func (ctx *MultiBotContext) GetOptions(pluginName string) map[string]interface{} {
	return ctx.options.PluginsSettings[pluginName]
}

// DBInsert insert data to database
func (ctx *MultiBotContext) DBInsert(data interface{}) (err error) {
	err = ctx.db.Insert(data)
	return
}

// GetFile function download file from telegram and store in 'filename'
func (ctx *MultiBotContext) GetFile(fileID, dir string) (filename string, err error) {
	var (
		botf tgbotapi.File
		resp *http.Response
		file *os.File
	)

	fc := tgbotapi.FileConfig{FileID: fileID}
	if botf, err = ctx.bot.GetFile(fc); err != nil {
		log.Errorf("Unable to get file FileID [%s]: %s", fileID, err)
		return
	}

	filename = botf.FilePath
	fullPath := filepath.Join(dir, filename)

	// check directory
	path := filepath.Dir(fullPath)
	if err = os.MkdirAll(path, 0755); err != nil {
		log.Errorf("Unable to make directories for FileID [%s]: %s", fileID, err)
		return
	}

	link := botf.Link(ctx.options.APIKey)
	if resp, err = http.Get(link); err != nil {
		log.Errorf("Unable to get file from [%s]: %s", link, err)
		if resp.Body != nil {
			resp.Body.Close()
		}
		return
	}
	defer resp.Body.Close()

	if file, err = os.Create(fullPath); err != nil {
		log.Errorf("Unable to create file [%s]: %s", fullPath, err)
		return
	}
	defer file.Close()
	if _, err = io.Copy(file, resp.Body); err != nil {
		log.Errorf("Unable to copy data from link %s to file %s: %s", link, fullPath, err)
		return
	}

	log.Debugf("File downloaded for FileID [%s] to %s", fileID, fullPath)
	return
}
