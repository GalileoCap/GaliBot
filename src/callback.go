package main;

import (
  tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

  //"fmt"
  "log"
  "strings"
  //"errors"
)

type Callback struct {
  Function func(*User, string, *tgbotapi.Message, *tgbotapi.EditMessageTextConfig) error;
};
var Callbacks map[string]Callback = map[string]Callback{ //U: Add all possible callbacks here
  "todoNew": {Function: cbTodoNew},
};

func handleCallback(user *User, query *tgbotapi.CallbackQuery) {
  var err error;
  edit := newEdit(user, query.Message);

  data := strings.SplitN(query.Data, ";", 2);
  if len(data) != 2 {
    log.Printf("[handleCallback] Invalid callback: %v", query.Data);
    //TODO: Send message with error identifier
    return;
  }

  cb, present := Callbacks[data[0]];
  if !present {
    log.Printf("[handleCallback] Invalid callback: %v", query.Data);
    //TODO: Send message with error identifier
    return;
  }

  err = cb.Function(user, data[1], query.Message, &edit);
  if err != nil { //A: Reset the reply
    log.Printf("[handleCallback] Error on handle: %v", err);
    //TODO: Send message with error identifier
    return;
  }

  err = sendEdit(edit);
  if err != nil {
    log.Printf("[handleCallback] Error on send: %v", err);
  }
}
