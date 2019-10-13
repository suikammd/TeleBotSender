package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"time"
)

type send func(to tb.Recipient, what *tb.Message, options ...interface{}) (*tb.Message, error)
type sendAlbum func(to tb.Recipient, a tb.Album, options ...interface{}) ([]tb.Message, error)

func initBot(botToken string) (bot *tb.Bot, err error) {
	return tb.NewBot(tb.Settings{
		Token:  botToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
}

func startBot(botToken string, personalId int, chats ...*tb.Chat) {
	bot, err := initBot(botToken)
	if (err != nil) {
		log.Fatal(err)
		return
	}

	var photos []tb.InputMedia
	var startTime = time.Time{}
	var endTime = time.Time{}
	goroutine := make(map[string]bool)

	bot.Handle(tb.OnPhoto, func(m *tb.Message) {
		albumId := m.AlbumID
		if albumId == "" {
			sendNormalChat(bot.Forward, m, chats...)
			return
		}
		if albumId != "" {
			if startTime.Equal(time.Time{}) {
				startTime = time.Now().Local()
				endTime = startTime.Add(time.Second * 8)
			}
			now := time.Now().Local()
			if (!endTime.Equal(time.Time{}) && now.Before(endTime)) {
				photos = append(photos, m.Photo)
			}
			if _, ok := goroutine[albumId]; !ok {
				go func() {
					select {
					case <-time.After(time.Second * 8):
						album := tb.Album{}
						album = append(album, photos...)
						sendAlbumChat(bot.SendAlbum, album, chats...)

						startTime = time.Time{}
						endTime = time.Time{}
						return
					}
				}()
				goroutine[albumId] = true
			}
		}
	})
	bot.Handle(tb.OnVideo, func(m *tb.Message) {
		sendNormalChat(bot.Forward, m, chats...)
	})
	bot.Handle(tb.OnText, func(m *tb.Message) {
		if m.Sender.ID == personalId {
			for _, chat := range chats {
				bot.Send(chat, m.Text)
				return
			}
		}
		sendNormalChat(bot.Forward, m, chats...)
	})
	bot.Handle(tb.OnVideoNote, func(m *tb.Message) {
		sendNormalChat(bot.Forward, m, chats...)
	})
	bot.Handle(tb.OnDocument, func(m *tb.Message) {
		sendNormalChat(bot.Forward, m, chats...)
	})
	bot.Handle(tb.OnSticker, func(m *tb.Message) {
		sendNormalChat(bot.Forward, m, chats...)
	})

	bot.Start()
}

func sendNormalChat(f send, m *tb.Message, chats ...*tb.Chat) {
	for _, chat := range chats {
		f(chat, m)
	}
}

func sendAlbumChat(f sendAlbum, a tb.Album, chats ...*tb.Chat)  {
	for _, chat := range chats {
		f(chat, a)
	}
}
