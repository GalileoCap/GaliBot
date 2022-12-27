package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	//"fmt"
  "log"
	//"errors"
)

type Mode struct {
  Function func(*User, *tgbotapi.Message);
};
var Modes map[string]Mode = map[string]Mode{ //U: Add all possible modes here
  "": {Function: modeDefault},
};

func handleMode(user *User, msg *tgbotapi.Message) {
  mode, present := Modes[user.Mode];
  if !present {
    log.Printf("[handleMode] Mode not present: %v", user.Mode);
    cancelMode(user); //A: Reset the user's mode
    //TODO: Error message and identifier
  }

  mode.Function(user, msg);
}

func enterMode(user *User, mode string) {
  user.Mode = mode;
}

func cancelMode(user *User) {
  user.Mode = "";
}

func modeDefault(user *User, msg *tgbotapi.Message) {
  //A: Ignore regular messages
}
