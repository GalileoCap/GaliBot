package main

import (
  tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

  "os"
  "encoding/json"
  "log"
  //"fmt"
)

type Credentials struct {
  Token string
};

var Bot *tgbotapi.BotAPI;

func getCredentials(fpath string) Credentials {
  data, err := os.ReadFile(fpath); //TODO: Config path
  if err != nil {
    log.Fatalf("[getCredentials] Error opening \"%v\": %v", fpath, err);
  }

  var credentials Credentials;
  err = json.Unmarshal(data, &credentials);
  if err != nil {
    log.Fatalf("[getCredentials] Error decoding \"%v\": %v", fpath, err);
  }

  log.Printf("[getCredentials] Successful %v", credentials);

  return credentials;
}

func main() {
  credentials := getCredentials(".credentials.json"); //TODO: Configure path
  Bot, err := tgbotapi.NewBotAPI(credentials.Token);
  if err != nil {
    log.Panicf("[main] Error NewBotAPI: %v", err);
  }

  //Bot.Debug = true;
  log.Printf(Bot.Self.UserName);

  u := tgbotapi.NewUpdate(0); //TODO: Last update +1
  u.Timeout = 60; //TODO: What?

  updates := Bot.GetUpdatesChan(u);

  for update := range updates {
    if update.Message != nil { // If we got a message
      log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

      msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
      msg.ReplyToMessageID = update.Message.MessageID

      Bot.Send(msg)
    }
  }
}
