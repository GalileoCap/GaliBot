package main

import (
	tggenbot "github.com/GalileoCap/tgGenBot/api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

  "time"
  "log"
)

var Export = tggenbot.Plugin{
  Name: "FoodDO",
  Version: "0.1.0",
  Init: Init,

  Commands: map[string]tggenbot.Command{
    "fooddo": {
      Description: "Show FoodDO hub",
      Admin: false,
      Tester: false,
      Skip: false,
      Function: cmdFooddo,
    },
    "foodo": {
      Description: "Show FoodDO hub (typo)",
      Admin: false,
      Tester: false,
      Skip: true,
      Function: cmdFooddo,
    },
  },

  Callbacks: map[string]tggenbot.Callback{
    "new": {
      Function: newCB,
    },
  },
}

func Init() error {
  err := InitDB()
  if err != nil {
    return err
  }

  go routine()

	return nil
}

func cmdFooddo(user *tggenbot.User, msg *tgbotapi.Message) error {
  reply := tgbotapi.NewMessage(msg.Chat.ID, "pong")
  reply.ReplyToMessageID = msg.MessageID

  reply.Text = "This is your FoodDO HUB" //TODO: Better description
  reply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
    tgbotapi.NewInlineKeyboardRow(
      tgbotapi.NewInlineKeyboardButtonData("New Entry", "FoodDO;new"),
    ),
    tgbotapi.NewInlineKeyboardRow(
      tgbotapi.NewInlineKeyboardButtonData("List All", "FoodDO;list"),
      tgbotapi.NewInlineKeyboardButtonData("See My Data", "FoodDO;data"),
    ),
    tgbotapi.NewInlineKeyboardRow(
      tgbotapi.NewInlineKeyboardButtonData("Edit Entry", "FoodDO;edit"),
      tgbotapi.NewInlineKeyboardButtonData("Config", "FoodDO;config"),
    ),
  )

  _, err := tggenbot.SendMessage(reply)
  return err
}

func routine() {
  firstDelay := time.Now().Minute() % 30
  log.Printf("[fooddoRoutine] First delay for %v minutes", firstDelay)
  time.Sleep(time.Duration(firstDelay) * time.Minute) //A: Wait until first half-hour
  for {
    now := time.Now().Hour() * 100 + time.Now().Minute() //A: To military time //TODO: Round down to 00/30
    rows, err := db.Query("SELECT * FROM fooddo_users WHERE Breakfast = ? OR Lunch = ? OR Merienda = ? OR Dinner = ?", now, now, now, now)
    if err != nil {
      log.Printf("[fooddoRoutine] Error with query: %v", err)
    }
    for rows.Next() {
      var user User
      if err := rows.Scan(&user.ID, &user.Breakfast, &user.Lunch, &user.Merienda, &user.Dinner); err != nil {
        log.Printf("[fooddoRoutine] Error scanning: %v", err)
      }
      log.Print(user)
    }

    //TODO: Send message asking for an entry to all registered users with this time

    rows.Close()
    time.Sleep(30 * time.Minute) //TODO: User ticker
  }
}
