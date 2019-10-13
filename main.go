package main

import (
	"flag"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {

	var (
		botToken   string
		groupId    int64
		personalId int
	)

	flag.StringVar(&botToken, "b", "", "bot token")
	flag.Int64Var(&groupId, "g", int64(0), "group id")
	flag.IntVar(&personalId, "p", 0, "personal id")
	flag.Parse()

	// send to chats
	chat := &tb.Chat{
		ID:   groupId,
		Type: tb.ChatGroup,
	}

	println(botToken, groupId, personalId)
	startBot(botToken, personalId, chat)
}
