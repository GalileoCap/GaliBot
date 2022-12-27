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

  registerCommands();
}

func receiveUpdates() {
  u := tgbotapi.NewUpdate(0); //TODO: Last update ID
  u.Timeout = 60; //TODO: What does this do?

  updates := Bot.GetUpdatesChan(u);
  for update := range updates {
    log.Printf("[receiveUpdates] New update uid=%v", update.SentFrom().ID);
    user, err := dbGetUser(update.SentFrom());
    if err != nil {
      log.Printf("[receiveUpdates] Error getting user: %v", err);
    }

    if user.Permissions == "block" { //A: Ignore them
      continue;
    }

    if update.Message != nil {
      if update.Message.IsCommand() {
        handleCommand(user, update.Message);
      } else {
        //TODO: Mode handler
      }
    } else if update.CallbackQuery != nil {
      //TODO: Handle query
    } else {
      log.Printf("[receiveUpdates] Unhandled update: %+v", update) //TODO: Print which type rather the entire thing
    }

    err = dbSaveUser(user);
    if err != nil {
      log.Printf("[receiveUpdates] Error saving user: %v", err);
      continue;
    }
  }
}

func newReply(user *User, msg *tgbotapi.Message) tgbotapi.MessageConfig {
  reply := tgbotapi.NewMessage(user.ID, ""); //TODO: Default text
  reply.ReplyToMessageID = msg.MessageID; 
  return reply;
}

func sendMessage(msg tgbotapi.MessageConfig) error {
  _, err := Bot.Send(msg);
  return err;
}
