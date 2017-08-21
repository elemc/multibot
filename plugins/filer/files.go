package main

import (
	"bytes"
	"math/rand"
	"time"

	"github.com/labstack/gommon/log"
)

// File is a struct for store secret phrases for uploaded files
type File struct {
	FileID       string
	FileName     string
	SecretPhrase string
	DocumentName string
	MimeType     string
}

func newFile(id string, docname string, mime string) (f *File) {
	f = &File{
		FileID:       id,
		DocumentName: docname,
		MimeType:     mime,
	}
	f.generateNewSecretPhrase()
	ctx.Log().Debugf("New file: %#v", f)
	return
}

// Upload function upload file and store secret phrase in database
func (f *File) Upload() (err error) {
	if f.FileName, err = ctx.GetFile(f.FileID, filesDir); err != nil {
		log.Errorf("Unable to get file from telegram with ID %s: %s", f.FileID, err)
		return
	}
	if err = ctx.GetDB().Insert(f); err != nil {
		log.Errorf("Unable to store file in database: %s", err)
		return
	}

	return
}

func (f *File) generateNewSecretPhrase() {
	var buf bytes.Buffer
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < secretPhraseLength; i++ {
		b := byte(r.Intn(93)) + 33
		buf.WriteByte(b)
	}
	f.SecretPhrase = buf.String()
}

func getFiles() (files []*File, err error) {
	err = ctx.GetDB().Model(&files).Select()
	return
}
