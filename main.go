package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"os"
	"io/ioutil" //TODO: Use os
	"net/http"
	"encoding/json"
  "fmt"
	"log"
	"time"
)

type Credentials struct {
  Token string
};

var myChatID int64 = 1129477471;
var currIP string;

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

func checkPermissions(chatID int64) bool {
  return chatID == myChatID;
  //TODO: Allowlist
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
  prevIP := currIP;
  currIP = newIP;

  return prevIP != newIP, nil;
} 

func ipUpdater(Bot *tgbotapi.BotAPI) {
  //ticker := time.NewTicker(24 * time.Hour); //TODO: Configure time
  ticker := time.NewTicker(time.Second); //TODO: Configure time
  for range ticker.C {
    changed, err := getIP();
    if err != nil {
      log.Printf("[ipUpdater] Error getting IP: %v", err);
      //TODO: Announce error
      continue;
    }

    if changed {
      Bot.Send(tgbotapi.NewMessage(myChatID, fmt.Sprintf("[ipUpdater] newIP: %v", currIP)));
    }
  }
}

func listenForMessages(Bot *tgbotapi.BotAPI) {
  u := tgbotapi.NewUpdate(0); //TODO: Last update +1
  u.Timeout = 60; //TODO: What?

  updates := Bot.GetUpdatesChan(u);

  for update := range updates {
    chatID := update.FromChat().ID; message := update.Message; //A: Rename
    log.Printf("[listenForMessages] Received update from: %v", chatID);

    if message == nil { //A: Ignore non-Message updates //TODO
      continue;
    }
    if !message.IsCommand() { //A: Ignore non-Command messages
      continue;
    }

    msg := tgbotapi.NewMessage(chatID, "");
    msg.ReplyToMessageID = message.MessageID;

    if checkPermissions(chatID) {
      switch update.Message.Command() {
        case "ping": msg.Text = "pong";
        case "ip":
          //TODO: Higher permissions
          if _, err := getIP(); err != nil {
            msg.Text = "Error getting IP";
          } else {
            msg.Text = fmt.Sprintf("IP: %v", currIP);
          }
        default: msg.Text = fmt.Sprintf("Unknown command: /%v", message.Command());
      }
    } else {
      msg.Text = "You're not in the allowlist";
    }

    if _, err := Bot.Send(msg); err != nil {
      log.Printf("[listenForMessages] Error sending message: %v", err);
    }
  }
}

func main() {
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
