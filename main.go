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

var Bot *tgbotapi.BotAPI;

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

func ipUpdater(Bot *tgbotapi.BotAPI) {
  ticker := time.NewTicker(24 * time.Hour); //TODO: Configure time
  var ip string;
  for range ticker.C {
    resp, err := http.Get("https://api.ipify.org?format=text");
    if err != nil {
      log.Printf("[ipUpdater] Error getting IP: %v", err);
      //TODO: Handle
    }
    
    newIP_b, err := ioutil.ReadAll(resp.Body);
    if err != nil {
      log.Printf("[ipUpdater] Error reading response body: %v", err);
      //TODO: Handle
    }
    newIP := string(newIP_b);
    if ip != newIP {
      ip = newIP;
      log.Printf("[ipUpdater] IP changed to: %v", ip);
      Bot.Send(tgbotapi.NewMessage(1129477471, fmt.Sprintf("[ipUpdater] newIP: %v", ip))); //TODO: Configure chatID
    }

    resp.Body.Close();
  }
  //TODO: Simplify
  //TODO: Callable
}

func listenForMessages(Bot *tgbotapi.BotAPI) {
  u := tgbotapi.NewUpdate(0); //TODO: Last update +1
  u.Timeout = 60; //TODO: What?

  updates := Bot.GetUpdatesChan(u);

  for update := range updates {
    if update.Message != nil { //TODO: 
      log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

      msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
      msg.ReplyToMessageID = update.Message.MessageID

      Bot.Send(msg)
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
