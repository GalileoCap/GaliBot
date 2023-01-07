package main

import (
	tggenbot "github.com/GalileoCap/tgGenBot/api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

  "strings"
  "strconv"
  "time"
  "fmt"
)

func newCB(user *tggenbot.User, data string, msg *tgbotapi.Message) error {
  info, present := user.PluginInfo["Test"]
  if !present {
    user.PluginInfo["Test"] = &Entry{
      UserID: user.ID,
      Date: time.Now(),

      Message: msg,
    }
    //TODO: Start message
    return newCB(user, "", msg)
  }
  edit := tgbotapi.NewEditMessageTextAndMarkup(msg.Chat.ID, msg.MessageID, "", *msg.ReplyMarkup)
  entry := info.(*Entry)
  var action string

  switch parts := strings.Split(data, ";"); parts[0] {
  default:
    return tggenbot.ErrCBInvalidQuery

  case "":
    return newCB(user, "type", msg)
    
  case "type":
    if len(parts) == 1 { //A: Hasn't chosen yet
      action = "Please choose the type of meal."
      *edit.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
          tgbotapi.NewInlineKeyboardButtonData("Breakfast", "FoodDO;new;type;0"),
          tgbotapi.NewInlineKeyboardButtonData("Lunch", "FoodDO;new;type;1"),
        ),
        tgbotapi.NewInlineKeyboardRow(
          tgbotapi.NewInlineKeyboardButtonData("Merienda", "FoodDO;new;type;2"),
          tgbotapi.NewInlineKeyboardButtonData("Dinner", "FoodDO;new;type;3"),
        ),
        tgbotapi.NewInlineKeyboardRow(
          tgbotapi.NewInlineKeyboardButtonData("Extra", "FoodDO;new;type;4"),
        ),
        tgbotapi.NewInlineKeyboardRow(
          tgbotapi.NewInlineKeyboardButtonData("Cancel", "FoodDO;new;cancel"),
          tgbotapi.NewInlineKeyboardButtonData("Back", "FoodDO;new;hub"),
        ),
      )
    } else { // len(parts) > 1
      t, err := strconv.Atoi(parts[1])
      if err != nil {
        //TODO: Prompt again
        return err
      }
      entry.Type = t
      return newCB(user, "description", msg)
    }

  case "description":
    user.NextMessage = newDescribe

    action = "Please describe your meal, or press skip if you skipped this meal."
    *edit.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
      tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Skip", "FoodDO;new;skip"),
      ),
      tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Cancel", "FoodDO;new;cancel"),
        tgbotapi.NewInlineKeyboardButtonData("Back", "FoodDO;new;type"),
      ),
    )

  case "skip":
    entry.Skipped = true
    return newCB(user, "edit", msg)

  case "tags":
    if len(parts) == 2 {
      switch parts[1] {
      case "meat": entry.Meat = !entry.Meat
      case "veggies": entry.Veggies = !entry.Veggies
      case "fruit": entry.Fruit = !entry.Fruit
      default: return tggenbot.ErrCBInvalidQuery
      }
    }

    action = "Please choose the tags that apply to your meal and then move to the next step."
    *edit.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
      tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Meat (%v)", entry.Meat), "FoodDO;new;tags;meat"),
        tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Veggies (%v)", entry.Veggies), "FoodDO;new;tags;veggies"),
        tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Fruit (%v)", entry.Fruit), "FoodDO;new;tags;fruit"),
      ),
      tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Next", "FoodDO;new;edit"),
      ),
      tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Cancel", "FoodDO;new;cancel"),
        tgbotapi.NewInlineKeyboardButtonData("Back", "FoodDO;new;description"),
      ),
    )
  
  case "edit":
    action = fmt.Sprintf("Do you wish to edit your entry?\nCurrent values:\n%v", fmtEntry(entry))
    *edit.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
      tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Edit", "FoodDO;new;"),
        tgbotapi.NewInlineKeyboardButtonData("Confirm", "FoodDO;new;confirm"),
      ),
      tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Cancel", "FoodDO;new;cancel"),
        tgbotapi.NewInlineKeyboardButtonData("Back", "FoodDO;new;describe"), //TODO: Back depending on skip/tags
      ),
    )
  
  case "confirm":
    err := addEntry(entry)
    if err != nil {
      return err
    }

    edit.Text = fmt.Sprintf("Successfully created new FOODDO entry!\nWith values:\n%v\n\nYou can use the unique ID %v to edit/delete it later.", fmtEntry(entry), entry.ID)
    *edit.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
      tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Edit Now", fmt.Sprintf("fooddoEdit;%v", entry.ID)),
      ),
      tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Back to hub", "FoodDO;new;hub"),
      ),
    )
    return nil

  case "hub":
    //TODO: Go back to the hub
  
  case "cancel":
    //TODO: Cancel
    //TODO: Also auto-cancel after inactivity
  }

  edit.Text = action

  _, err := tggenbot.SendMessage(edit)
  return err
}

func newDescribe(user *tggenbot.User, msg *tgbotapi.Message) error {
  info, _ := user.PluginInfo["Test"]
  entry := info.(*Entry)
  entry.Description = msg.Text
  return newCB(user, "tags", entry.Message)
}
