package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

  "strings"
  "fmt"
  "log"
  "errors"
)

func handleCallback(user DBUser, query *tgbotapi.CallbackQuery) {
  var err error;

  if query.Message == nil {
    log.Printf("[handleCallback] Query without message: %v", query);
    //TODO: Handle
    return;
  }

  edit := tgbotapi.NewEditMessageTextAndMarkup(query.From.ID, query.Message.MessageID, query.Message.Text, *query.Message.ReplyMarkup);

  switch parts := strings.Split(query.Data, ";"); parts[0] {
  case "todo":
    err = todoCallback(user, query, parts, &edit);

  default:
    err = errors.New(fmt.Sprintf("Invalid query data, possible threat: %v", query.Data));
  }

  if err != nil {
    log.Printf("[handleCallback] Error: %v", err);
    //TODO: Handle
    return;
  }

  _, err = Bot.Send(edit);
  if err != nil {
    log.Printf("[handleCallback] Error sending message: %v", err);
  }
}
