package main

import (
  tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

  "strings"
  "strconv"
  "time"
  "fmt"
  "log"
)

type FOODDOUser struct {
  ID int64 //U: Telegram's UserID of the entry's owner

  Breakfast int16
  Lunch int16
  Merienda int16
  Dinner int16
}

type FOODDOEntry struct {
  ID int64 //U: Entry's unique ID
  UserID int64 //U: Telegram's UserID of the entry's owner

  Date time.Time
  Type int // 0: Breakfast, 1: Lunch, 2: Merienda, 3: Dinner, 4: Extra
  Description string
  Skipped bool

  Meat bool
  Veggies bool
  Fruit bool

  Message *tgbotapi.Message
}

var FOODDOTypes = []string{ "Breakfast", "Lunch", "Merienda", "Dinner", "Extra" } //TODO: Const

func fooddoEntry(entry *FOODDOEntry) string {
  date := entry.Date.Format("Mon 02/01/2006, 03:04")
  if entry.Skipped {
    return fmt.Sprintf("- Date: %v\n.\t- Type: %v.\n\t- Skipped.", date, FOODDOTypes[entry.Type])
  } else {
    tags := []string{}
    if entry.Meat {
      tags = append(tags, "Meat")
    }
    if entry.Veggies {
      tags = append(tags, "Veggies")
    }
    if entry.Fruit {
      tags = append(tags, "Fruit")
    }
    return fmt.Sprintf("- Date: %v\n\t- Type: %v.\n\t- Tags:%v\n\t- Description: %v.", date, entry.Type, tags, entry.Description)
  }
}

func fooddoAddEntry(entry *FOODDOEntry) error {
  result, err := DB.Exec("INSERT INTO fooddo_entries (UserID, Date, Type, Description, Skipped, Meat, Veggies, Fruit) VALUES (?,?,?,?,?,?,?,?);", entry.UserID, entry.Date.String(), entry.Type, entry.Description, entry.Skipped, entry.Meat, entry.Veggies, entry.Fruit);
  if err != nil {
    return err
  }
  
  id, err := result.LastInsertId()
  if err != nil {
    return err
  }

  entry.ID = id
  return nil
}

func fooddoRoutine() {
  firstDelay := time.Now().Minute() % 30
  log.Printf("[fooddoRoutine] First delay for %v minutes", firstDelay)
  time.Sleep(time.Duration(firstDelay) * time.Minute) //A: Wait until first half-hour
  for {
    now := time.Now().Hour() * 100 + time.Now().Minute() //A: To military time //TODO: Round down to 00/30
    rows, err := DB.Query("SELECT * FROM fooddo_users WHERE Breakfast = ? OR Lunch = ? OR Merienda = ? OR Dinner = ?", now, now, now, now)
    if err != nil {
      log.Printf("[fooddoRoutine] Error with query: %v", err)
    }
    for rows.Next() {
      var user FOODDOUser
      if err := rows.Scan(&user.ID, &user.Breakfast, &user.Lunch, &user.Merienda, &user.Dinner); err != nil {
        log.Printf("[fooddoRoutine] Error scanning: %v", err)
      }
      log.Print(user)
    }

    //TODO: Send message asking for an entry to all registered users with this time

    rows.Close()
    time.Sleep(30 * time.Minute)
  }
}

func fooddoCMD(user *User, msg *tgbotapi.Message, reply *tgbotapi.MessageConfig) error {
  reply.Text = "This is your FoodDO HUB" //TODO: Better description
  reply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
    tgbotapi.NewInlineKeyboardRow(
      tgbotapi.NewInlineKeyboardButtonData("New Entry", "fooddoNew;"),
    ),
    tgbotapi.NewInlineKeyboardRow(
      tgbotapi.NewInlineKeyboardButtonData("List All", "fooddoList;"),
      tgbotapi.NewInlineKeyboardButtonData("See My Data", "fooddoData;"),
    ),
    tgbotapi.NewInlineKeyboardRow(
      tgbotapi.NewInlineKeyboardButtonData("Edit Entry", "fooddoEdit;"),
      tgbotapi.NewInlineKeyboardButtonData("Config", "fooddoConfig;"),
    ),
  )
  return nil
}

func fooddoNewCB(user *User, query string, msg *tgbotapi.Message, edit *tgbotapi.EditMessageTextConfig) error {
  if user.FOODDOEntry == nil {
    user.FOODDOEntry = &FOODDOEntry{
      UserID: user.ID,
      Date: time.Now(),

      Message: msg,
    }
    //TODO: Start message
    return fooddoNewCB(user, "", msg, edit)
  }

  var action string

  switch parts := strings.Split(query, ";"); parts[0] {
  default:
    return ErrCBInvalidQuery

  case "":
    return fooddoNewCB(user, "type", msg, edit)
    
  case "type":
    if len(parts) == 1 { //A: Hasn't chosen yet
      action = "Please choose the type of meal."
      *edit.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
          tgbotapi.NewInlineKeyboardButtonData("Breakfast", "fooddoNew;type;0"),
          tgbotapi.NewInlineKeyboardButtonData("Lunch", "fooddoNew;type;1"),
        ),
        tgbotapi.NewInlineKeyboardRow(
          tgbotapi.NewInlineKeyboardButtonData("Merienda", "fooddoNew;type;2"),
          tgbotapi.NewInlineKeyboardButtonData("Dinner", "fooddoNew;type;3"),
        ),
        tgbotapi.NewInlineKeyboardRow(
          tgbotapi.NewInlineKeyboardButtonData("Extra", "fooddoNew;type;4"),
        ),
        tgbotapi.NewInlineKeyboardRow(
          tgbotapi.NewInlineKeyboardButtonData("Cancel", "fooddoNew;cancel"),
          tgbotapi.NewInlineKeyboardButtonData("Back", "fooddoNew;hub"),
        ),
      )
    } else { // len(parts) > 1
      t, err := strconv.Atoi(parts[1])
      if err != nil {
        //TODO: Prompt again
        return err
      }
      user.FOODDOEntry.Type = t
      return fooddoNewCB(user, "description", msg, edit)
    }

  case "description":
    enterMode(user, "fooddoNewDescribe")

    action = "Please describe your meal, or press skip if you skipped this meal."
    *edit.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
      tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Skip", "fooddoNew;skip"),
      ),
      tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Cancel", "fooddoNew;cancel"),
        tgbotapi.NewInlineKeyboardButtonData("Back", "fooddoNew;type"),
      ),
    )

  case "skip":
    user.FOODDOEntry.Skipped = true
    return fooddoNewCB(user, "edit", msg, edit)

  case "tags":
    if len(parts) == 2 {
      switch parts[1] {
      case "meat": user.FOODDOEntry.Meat = !user.FOODDOEntry.Meat
      case "veggies": user.FOODDOEntry.Veggies = !user.FOODDOEntry.Veggies
      case "fruit": user.FOODDOEntry.Fruit = !user.FOODDOEntry.Fruit
      default: return ErrCBInvalidQuery
      }
    }

    action = "Please choose the tags that apply to your meal and then move to the next step."
    *edit.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
      tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Meat (%v)", user.FOODDOEntry.Meat), "fooddoNew;tags;meat"),
        tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Veggies (%v)", user.FOODDOEntry.Veggies), "fooddoNew;tags;veggies"),
        tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Fruit (%v)", user.FOODDOEntry.Fruit), "fooddoNew;tags;fruit"),
      ),
      tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Next", "fooddoNew;edit"),
      ),
      tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Cancel", "fooddoNew;cancel"),
        tgbotapi.NewInlineKeyboardButtonData("Back", "fooddoNew;description"),
      ),
    )
  
  case "edit":
    action = fmt.Sprintf("Do you wish to edit your entry?\nCurrent values:\n%v", fooddoEntry(user.FOODDOEntry))
    *edit.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
      tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Edit", "fooddoNew;"),
        tgbotapi.NewInlineKeyboardButtonData("Confirm", "fooddoNew;confirm"),
      ),
      tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Cancel", "fooddoNew;cancel"),
        tgbotapi.NewInlineKeyboardButtonData("Back", "fooddoNew;describe"), //TODO: Back depending on skip/tags
      ),
    )
  
  case "confirm":
    err := fooddoAddEntry(user.FOODDOEntry)
    if err != nil {
      return err
    }

    edit.Text = fmt.Sprintf("Successfully created new FOODDO entry!\nWith values:\n%v\n\nYou can use the unique ID %v to edit/delete it later.", fooddoEntry(user.FOODDOEntry), user.FOODDOEntry.ID)
    *edit.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
      tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Edit Now", fmt.Sprintf("fooddoEdit;%v", user.FOODDOEntry.ID)),
      ),
      tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Back to hub", "fooddoNew;hub"),
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

  return nil
}

func fooddoNewDescribeMode(user *User, msg *tgbotapi.Message) {
  cancelMode(user)
  user.FOODDOEntry.Description = msg.Text

  edit := newEdit(user, user.FOODDOEntry.Message)
  err := fooddoNewCB(user, "tags", user.FOODDOEntry.Message, &edit)
  if err != nil {
    //TODO: Warn
    return
  }
  sendEdit(edit)
}
