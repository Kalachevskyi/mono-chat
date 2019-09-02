// Copyright © 2019 Volodymyr Kalachevskyi <v.kalachevskyi@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package handlers is an interface adapters of application
package handlers

import (
	"io"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

// CsvUC - represents a usecase interface for processing business logic of a CSV report
type CsvUC interface {
	Validate(name string) error
	GetFile(url string) (io.ReadCloser, error)
	Parse(chatID int64, fileName string, r io.Reader) (io.Reader, error)
}

// NewFileReport - builds "FileReport" internal handler
func NewFileReport(csvUC CsvUC, botWrapper *BotWrapper) *FileReport {
	return &FileReport{csvUC: csvUC, BotWrapper: botWrapper}
}

// FileReport - represents an internal handler for processing a CSV report
type FileReport struct {
	csvUC CsvUC
	*BotWrapper
}

// Handle - process the CSV report, send the result to the user
func (f *FileReport) Handle(u tg.Update) {
	if err := f.csvUC.Validate(u.Message.Document.FileName); err != nil {
		f.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}

	fileTG, err := f.bot.GetFile(tg.FileConfig{FileID: u.Message.Document.FileID})
	if err != nil {
		f.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}

	file, err := f.csvUC.GetFile(fileTG.Link(f.bot.Token))
	if err != nil {
		f.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}

	fileResp, err := f.csvUC.Parse(u.Message.Chat.ID, u.Message.Document.FileName, file)
	if err != nil {
		f.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}

	name := u.Message.Document.FileName
	reader := tg.FileReader{
		Name:   name,
		Reader: fileResp,
		Size:   -1,
	}
	msg := tg.NewDocumentUpload(u.Message.Chat.ID, reader)
	f.sendMSG(msg)

	close(file, f.log)
}
