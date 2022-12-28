package main;

import (
  tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

  "fmt"
  "log"
  //"errors"
)

type Command struct {
  Description string;
  Admin bool; //U: Admin required to use this command

  Function func(*User, *tgbotapi.Message, *tgbotapi.MessageConfig) error;
};
var Commands map[string]Command = map[string]Command{ //U: Add commands to be registered here
  "test": {Description: "Whatever I might be testing right now", Admin: true, Function: cmdTest},
  "ing": {Description: "Whatever I might be testing right now", Admin: true, Function: cmdIng},

  "ping": {Description: "Ping me", Function: cmdPing},
  "cancel": {Description: "Leave the current mode", Function: cmdCancel},

  "ip": {Description: "...", Admin: true, Function: cmdIP},
};

func registerCommands() { //U: Registers all known commands
  commands := make([]tgbotapi.BotCommand, 0, len(Commands));
  for name, command := range Commands {
    commands = append(commands, tgbotapi.BotCommand{Command: name, Description: command.Description});
  }

  _, err := Bot.Request(tgbotapi.NewSetMyCommands(commands...));
  if err != nil {
    log.Fatalf("[registerCommands] %v", err); //TODO: Is this fatal? Maybe send message to admin
  }

  //TODO: Scope
  //TODO: Simplify
}

func handleCommand(user *User, msg *tgbotapi.Message) {
  var err error;
  reply := newReply(user, msg); 

  cmd, present := Commands[msg.Command()];
  if !present {
    reply.Text = fmt.Sprintf("Unknown command /%v, try asking for /help", msg.Command());
    goto SEND;
  }

  if cmd.Admin && user.Permissions != 0 {
    reply.Text = fmt.Sprintf("You don't have the correct permissions for the command /%v", msg.Command());
    goto SEND;
  }

  err = cmd.Function(user, msg, &reply);
  if err != nil { //A: Reset the reply
    log.Printf("[handleCommand] Error on handle: %v", err);

    reply = newReply(user, msg); //A: Reset the reply
    reply.Text = fmt.Sprintf("There was an error handling the command /%v.\nPlease retry in a bit, or ask your local admin for help", msg.Command());

    //TODO: Error identifier
  }

SEND:
  err = sendMessage(reply);
  if err != nil {
    log.Printf("[handleCommand] Error on send: %v", err);
  }

  //TODO: Should all commands cancelMode?
}

func cmdTest(user *User, msg *tgbotapi.Message, reply *tgbotapi.MessageConfig) error {
  return nil;
}

func cmdIng(user *User, msg *tgbotapi.Message, reply *tgbotapi.MessageConfig) error {
  return nil;
}

func cmdPing(user *User, msg *tgbotapi.Message, reply *tgbotapi.MessageConfig) error {
  reply.Text = "pong";
  return nil;
}

func cmdCancel(user *User, msg *tgbotapi.Message, reply *tgbotapi.MessageConfig) error {
  reply.Text = fmt.Sprintf("Exited mode: %v", user.Mode);
  cancelMode(user);
  return nil;
}

func cmdIP(user *User, msg *tgbotapi.Message, reply *tgbotapi.MessageConfig) error {
  changed, err := updateIP();
  if err == nil {
    reply.Text = fmt.Sprintf("IP (%v): %v", changed, CurrIP);
  }
  return err;
}
