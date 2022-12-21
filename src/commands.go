package main

import (
  tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

  "fmt"
	"log"
  //"errors"
)

type Command struct {
  Description string;
  Admin bool; //U: Admin required to use this command

  Function func(DBUser, *tgbotapi.Message, *tgbotapi.MessageConfig) error;
};
var Commands map[string]Command = map[string]Command{
  "ping": {Description: "ping", Function: pingCMD},
  "ip": {Description: "ip", Admin: true, Function: ipCMD},
};

func registerCommands() error { //U: Registers all known commands
  commands := make([]tgbotapi.BotCommand, 0, len(Commands));
  for name, command := range Commands {
    commands = append(commands, tgbotapi.BotCommand{Command: name, Description: command.Description});
  }
  config := tgbotapi.NewSetMyCommands(commands...);
  _, err := Bot.Request(config);
  return err;

  //TODO: Scope
  //TODO: Simplify
}

func pingCMD(user DBUser, msg *tgbotapi.Message, reply *tgbotapi.MessageConfig) error {
  reply.Text = "pong";
  return nil;
}

func ipCMD(user DBUser, msg *tgbotapi.Message, reply *tgbotapi.MessageConfig) error {
  change, err := getIP();
  if err == nil {
    reply.Text = fmt.Sprintf("IP (%v): %v", change, CurrIP);
  }
  return err;
}

func handleCommand(user DBUser, msg *tgbotapi.Message) {
  var err error;

  reply := tgbotapi.NewMessage(user.ChatId, "");
  reply.ReplyToMessageID = msg.MessageID;

  cmd, prs := Commands[msg.Command()];
  if !prs {
    reply.Text = fmt.Sprintf("Unknown command: /%v, try asking for /help", msg.Command());
    goto SEND;
  }

  if cmd.Admin && user.Permissions != "admin" {
    reply.Text = fmt.Sprintf("You don't have the correct permissions for the command: /%v", msg.Command());
    goto SEND;
  }

  err = cmd.Function(user, msg, &reply);
  if err != nil { //A: Reset the reply
    log.Printf("[handleCommand] Error: %v", err);
    reply = tgbotapi.NewMessage(user.ChatId, "There was an error"); 
    reply.ReplyToMessageID = msg.MessageID;

    //TODO: Identifier
  }

SEND:
  _, err = Bot.Send(reply);
  if err != nil {
    log.Printf("[handleCommand] Error sending message: %v", err);
  }
}
