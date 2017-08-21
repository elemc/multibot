package main

import (
	"fmt"
	"multibot/context"
	"path/filepath"

	"github.com/labstack/gommon/log"

	"gopkg.in/telegram-bot-api.v4"
)

var (
	ctx                *context.MultiBotContext
	options            map[string]interface{}
	filesDir           string
	secretPhraseLength int
)

// InitPlugin initialize plugin if it needed
func InitPlugin(mbc *context.MultiBotContext) error {
	ctx = mbc
	options = ctx.GetOptions(GetName())
	if st, ok := options["files_dir"]; ok && st != nil {
		filesDir = st.(string)
	}
	if st, ok := options["secret_length"]; ok && st != nil {
		secretPhraseLength = st.(int)
	}
	if err := ctx.DBCreateTable(&File{}); err != nil {
		ctx.Log().Errorf("Unable to create table for files: %s", err)
		return err
	}
	return nil
}

// GetName function returns plugin name
func GetName() string {
	return "filer"
}

// GetDescription function returns plugin description
func GetDescription() string {
	return "Plugin for upload files and transfer secret phrase to another people for download uploaded files"
}

// GetCommands return plugin commands for bot
func GetCommands() []string {
	return []string{}
}

// UpdateHandler function call for each update
func UpdateHandler(update tgbotapi.Update) (err error) {
	var f *File
	if update.Message.Document != nil {
		f = newFile(update.Message.Document.FileID, update.Message.Document.FileName, update.Message.Document.MimeType)
	} else if update.Message.Photo != nil {
		var max tgbotapi.PhotoSize
		for _, photo := range *update.Message.Photo {
			if photo.Height > max.Height {
				max = photo
			}
		}
		f = newFile(max.FileID, max.FileID, "image/jpeg")
	} else if update.Message.Text != "" {
		go searchSecretPhrase(update.Message)
		return
	}
	if err = f.Upload(); err != nil {
		return
	}
	ctx.SendMessageMarkdown(update.Message.Chat.ID, fmt.Sprintf("Секрет: %s", f.SecretPhrase), 0, nil)
	return nil
}

// RunCommand handler start if bot get one of commands
func RunCommand(command string, update tgbotapi.Update) (err error) {
	return
}

// StartCommand handler start if bot get one command 'start'
func StartCommand(update tgbotapi.Update) (err error) {
	msg := fmt.Sprintf(`Тебя приветствует плагин "Файлер"
Отправь боту файл, в ответ получишь секретное слово, запиши его и передай тому, кому адресован файл, по этому секретному слову файл можно будет получить.
Приятного пользования "Файлером"!`)
	ctx.SendMessageMarkdown(update.Message.Chat.ID, msg, 0, nil)
	return
}

func searchSecretPhrase(msg *tgbotapi.Message) {
	var (
		files []*File
		err   error
	)

	if files, err = getFiles(); err != nil {
		log.Errorf("Unable to get all files: %s", err)
		return
	}

	for _, f := range files {
		if f.SecretPhrase == msg.Text {
			ctx.SendFile(msg.Chat.ID, msg.MessageID, filepath.Join(filesDir, f.FileName), f.MimeType)
		}
	}
}
