package main

import (
  tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

  "fmt"
	"log"
  "time"
)

func ipUpdater() {
  ticker := time.NewTicker(time.Hour); //TODO: Configure time
  for range ticker.C {
    changed, err := getIP();
    if err != nil {
      log.Printf("[ipUpdater] Error getting IP: %v", err);
      //TODO: Announce error
      continue;
    }

    if changed {
      Bot.Send(tgbotapi.NewMessage(MyChatID, fmt.Sprintf("[ipUpdater] newIP: %v", CurrIP)));
    }
  }
}

func listenForMessages() {
  u := tgbotapi.NewUpdate(0); //TODO: Last update +1
  u.Timeout = 60; //TODO: What?

  updates := Bot.GetUpdatesChan(u);

  for update := range updates {
    chatID := update.FromChat().ID; msg := update.Message //A: Rename
    log.Printf("[listenForMessages] Received update from: %v", chatID);

    user, err := getUser(chatID);
    if err != nil {
      log.Printf("[listenForMessages] getUser error: %v", err);
      //TODO: Handle
      continue;
    }

    if user.Permissions == "block" {
      //TODO: Warn
      continue;
    }

    if msg != nil && msg.IsCommand() {
      handleCommand(user, msg);
    }

    //TODO: Non-message updates
  }
}
