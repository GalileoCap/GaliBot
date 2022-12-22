package main;

import (
  tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"log"
);

var Bot *tgbotapi.BotAPI;

func telegramInit(apiToken string) {
  var err error;
  Bot, err = tgbotapi.NewBotAPI(apiToken);
  if err != nil {
    log.Fatalf("[telegramInit] %v", err);
  }

  //Bot.Debug = Config.Test; //TODO: ?
  log.Printf("[telegramInit] Running as %s", Bot.Self.UserName);
}

func receiveUpdates() {
  u := tgbotapi.NewUpdate(0); //TODO: Last update ID
  u.Timeout = 60; //TODO: What does this do?

  updates := Bot.GetUpdatesChan(u);
  for update := range updates {
    log.Printf("[receiveUpdates] New update userID=%v, chatID=%v", update.SentFrom().ID, update.FromChat().ID);
  }
}
