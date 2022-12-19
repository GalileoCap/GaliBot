package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"

	"os"
	"io/ioutil" //TODO: Use os
	"net/http"
	"encoding/json"
  "fmt"
	"log"
	"time"
)

var MyChatID int64 = 1129477471;
var CurrIP string;

type PermLvl int;
const (
  Admin PermLvl = 0;
  Allow PermLvl = 1;
  Block PermLvl = 2;
)
var Permissions = map[int64]PermLvl{
  MyChatID: Admin,
  5629879871: Admin,
};

type Credentials struct {
  Token string
};

func getCredentials(fpath string) Credentials {
  data, err := os.ReadFile(fpath); //TODO: Config path
  if err != nil {
    log.Fatalf("[getCredentials] Error opening \"%v\": %v", fpath, err);
  }

  var credentials Credentials;
  err = json.Unmarshal(data, &credentials);
  if err != nil {
    log.Fatalf("[getCredentials] Error decoding \"%v\": %v", fpath, err);
  }

  return credentials;
}

func getIP() (bool, error) {
  resp, err := http.Get("https://api.ipify.org?format=text");
  if err != nil {
    log.Printf("[ipUpdater] Error getting IP: %v", err);
    return false, err;
  }
  defer resp.Body.Close();
  
  newIP_b, err := ioutil.ReadAll(resp.Body);
  if err != nil {
    log.Printf("[ipUpdater] Error reading response body: %v", err);
    return false, err;
  }
  newIP := string(newIP_b);
  prevIP := CurrIP;
  CurrIP = newIP;

  return prevIP != newIP, nil;
} 

func ipUpdater(Bot *tgbotapi.BotAPI) {
  ticker := time.NewTicker(24 * time.Hour); //TODO: Configure time
  for range ticker.C {
    changed, err := getIP();
    if err != nil {
      log.Printf("[ipUpdater] Error getting IP: %v", err);
      //TODO: Announce error
      continue;
    }

    if changed {
      Bot.Send(tgbotapi.NewMessage(MyChatID, fmt.Sprintf("[ipUpdater] newIP: %v", CurrIP)));
    }
  }
}

func command(chatID int64, Message *tgbotapi.Message, Bot *tgbotapi.BotAPI) {
  msg := tgbotapi.NewMessage(chatID, "");
  msg.ReplyToMessageID = Message.MessageID;

  if permLvl, ok := Permissions[chatID]; ok && permLvl != Block {
    switch Message.Command() {
      case "ping": msg.Text = "pong";
      case "ip":
        if permLvl != Admin {
          msg.Text = fmt.Sprintf("Unknown command: /%v", Message.Command()); //TODO: Repeated code
          break;
        }
        if _, err := getIP(); err != nil {
          msg.Text = "Error getting IP";
        } else {
          msg.Text = fmt.Sprintf("IP: %v", CurrIP);
        }
      default: msg.Text = fmt.Sprintf("Unknown command: /%v, try asking for /help", Message.Command());
    }
  } else {
    msg.Text = "You're not in the allowlist, please ask your local admin to add you";
  }

  if _, err := Bot.Send(msg); err != nil {
    log.Printf("[command] Error sending message: %v", err);
  }
}

func listenForMessages(Bot *tgbotapi.BotAPI) {
  u := tgbotapi.NewUpdate(0); //TODO: Last update +1
  u.Timeout = 60; //TODO: What?

  updates := Bot.GetUpdatesChan(u);

  for update := range updates {
    chatID := update.FromChat().ID; message := update.Message; //A: Rename
    log.Printf("[listenForMessages] Received update from: %v", chatID);

    if message != nil && message.IsCommand() {
      command(chatID, message, Bot);
    } //TODO: Non-message updates
  }
}

type User struct {
  id int64;
  name string;
}

func main() {
  db, err := sql.Open("mysql", "root:root@tcp(database:3306)/galibot");
  if err != nil {
    log.Fatal(err);
  }
  defer db.Close();

  rows, err := db.Query("SELECT * FROM users");
  if err != nil {
    log.Fatal(err);
  }
  defer rows.Close();

  var users []User;
  for rows.Next() {
    var user User;
    if err = rows.Scan(&user.id, &user.name); err != nil {
      log.Fatal(err);
    }
    users = append(users, user);
  }
  if err = rows.Err(); err != nil {
    log.Fatal(err);
  }

  log.Print("USERS", users);

  return;

  credentials := getCredentials(".credentials.json"); //TODO: Configure path
  Bot, err := tgbotapi.NewBotAPI(credentials.Token);
  if err != nil {
    log.Panicf("[main] Error NewBotAPI: %v", err);
  }
  //Bot.Debug = true;
  log.Printf("[main] Running as %v", Bot.Self.UserName);

  go ipUpdater(Bot);
  listenForMessages(Bot);
}

//TODO: Split into multiple files
