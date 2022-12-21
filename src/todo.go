package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

  "fmt"
	"errors"
  //"log"
)

func todoCMD(user DBUser, msg *tgbotapi.Message, reply *tgbotapi.MessageConfig) error {
  reply.Text = "Here's your TODO hub";
  reply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
    tgbotapi.NewInlineKeyboardRow(
      tgbotapi.NewInlineKeyboardButtonData("Add", "todo;add"),
    ),
    tgbotapi.NewInlineKeyboardRow(
      tgbotapi.NewInlineKeyboardButtonData("List All", "todo;list"),
      tgbotapi.NewInlineKeyboardButtonData("Random", "todo;random"),
    ),
  );

  return nil;
}

func todoCallback(user DBUser, query *tgbotapi.CallbackQuery, parts []string, edit *tgbotapi.EditMessageTextConfig) error {
  if len(parts) < 2 {
    return errors.New(fmt.Sprintf("[todoCallback] Invalid query data, possible threat: %v", query.Data));
  }

  var err error = nil;
  switch parts[1] {
  case "add":
    err = todoAdd(user, query, parts, edit);

  default:
    err = errors.New(fmt.Sprintf("[todoCallback] Invalid query data, possible threat: %v", query.Data));
  }

  return err;
}

func todoAdd(user DBUser, query *tgbotapi.CallbackQuery, parts []string, edit *tgbotapi.EditMessageTextConfig) error {
  var err error = nil;

  switch len(parts) {
  case 2:
    *edit.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
      tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Urgent", "todo;add;urgent"),
        tgbotapi.NewInlineKeyboardButtonData("Medium", "todo;add;medium"),
        tgbotapi.NewInlineKeyboardButtonData("Low", "todo;add;low"),
      ),
      tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Long", "todo;add;long"),
        tgbotapi.NewInlineKeyboardButtonData("Medium", "todo;add;medium"),
        tgbotapi.NewInlineKeyboardButtonData("Short", "todo;add;short"),
      ),
      tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Tag", "todo;add;tag"),
      ),
    );

  default:
    err = errors.New("TODO: Longer parts");
  }

  return err;
}
