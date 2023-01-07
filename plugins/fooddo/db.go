package main

import (
  "database/sql"
  _ "github.com/mattn/go-sqlite3"

	//tggenbot "github.com/GalileoCap/tgGenBot/api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

  "time"
)

type User struct {
  ID int64 //U: Telegram's UserID of the entry's owner

  Breakfast int16
  Lunch int16
  Merienda int16
  Dinner int16
}

type Entry struct {
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

var db *sql.DB

func InitDB() error {
  //TODO
  return nil
}

func addEntry(entry *Entry) error {
  result, err := db.Exec("INSERT INTO entries (UserID, Date, Type, Description, Skipped, Meat, Veggies, Fruit) VALUES (?,?,?,?,?,?,?,?);", entry.UserID, entry.Date.String(), entry.Type, entry.Description, entry.Skipped, entry.Meat, entry.Veggies, entry.Fruit);
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
