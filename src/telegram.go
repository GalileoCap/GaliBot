package main

import (
  tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

  "fmt"
	"log"
  "time"
)

func ipUpdater() {
  ticker := time.NewTicker(24 * time.Hour); //TODO: Configure time
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

func command(chatID int64, Message *tgbotapi.Message) {
  msg := tgbotapi.NewMessage(chatID, "");
  msg.ReplyToMessageID = Message.MessageID;

  user, err := getUser(chatID);
  if err != nil {
    log.Fatalf("[command] getUser error: %v", err);
    //TODO: Handle
  }

  if user.Permissions != "block" {
    switch Message.Command() {
      case "ping": msg.Text = "pong";
      case "ip":
        if user.Permissions != "admin" {
          msg.Text = fmt.Sprintf("Unknown command: /%v", Message.Command()); //TODO: Repeated code
          break;
        }
        if _, err := getIP(); err != nil {
          msg.Text = "Error getting IP";
        } else {
          msg.Text = fmt.Sprintf("IP: %v", CurrIP);
        }
      default: msg.Text = fmt.Sprintf("Unknown command: /%v, try asking for /help", Message.Command());
    }
  } else {
    msg.Text = "You're not in the allowlist, please ask your local admin to add you";
  }

  if _, err := Bot.Send(msg); err != nil {
    log.Printf("[command] Error sending message: %v", err);
  }
}

func listenForMessages() {
  u := tgbotapi.NewUpdate(0); //TODO: Last update +1
  u.Timeout = 60; //TODO: What?

  updates := Bot.GetUpdatesChan(u);

  for update := range updates {
    chatID := update.FromChat().ID; message := update.Message; //A: Rename
    log.Printf("[listenForMessages] Received update from: %v", chatID);

    if message != nil && message.IsCommand() {
      command(chatID, message);
    } //TODO: Non-message updates
  }
}
