package main

import (
	tggenbot "github.com/GalileoCap/tgGenBot/api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"net/http"
	"io"
	"time"
	"fmt"
  "log"
)

var Export = tggenbot.Plugin{
  Name: "IP",
  Version: "0.1.0",
  Init: Init,

  Commands: map[string]tggenbot.Command{
    "ping": {
      Description: "Check if the bot is up",
      Admin: false,
      Tester: false,
      Skip: false,
      Function: cmdPing,
    },
    "ip": {
      Description: "Gets the bot's public IP (and registers you to get other updates)",
      Admin: true,
      Tester: false,
      Skip: true,
      Function: cmdIP,
    },
  },
}

var currIP string;
var users = []int64{1129477471}

func Init() error {
  ticker := time.NewTicker(time.Second)
  go func() {
    for {
      <-ticker.C
      changed, err := getIP()
      if err != nil {
        log.Print(err)
        continue
      }
      if !changed {
        continue
      }

      for _, userID := range users {
        msg := tgbotapi.NewMessage(
          userID,
          fmt.Sprintf("New IP: %v", currIP),
        )
        _, err = tggenbot.SendMessage(msg)
      }
    }
  }()

  return nil
}

func addUser(userID int64) bool {
  for _, uid := range users {
    if userID == uid {
      return false
    }
  }

  users = append(users, userID)
  return true
}

func getIP() (bool, error) {
  resp, err := http.Get("https://api.ipify.org?format=text");
  if err != nil {
    return false, err;
  }
  defer resp.Body.Close();

  newIP_b, err := io.ReadAll(resp.Body);
  if err != nil {
    return false, err;
  }

  newIP, prevIP := string(newIP_b), currIP;
  currIP = newIP;

  return currIP != prevIP, nil;
}

func cmdIP(user *tggenbot.User, msg *tgbotapi.Message) error {
  addUser(user.ID) //TODO: Announce

  changed, err := getIP()
  if err != nil {
    return err
  }

  reply := tgbotapi.NewMessage(
    msg.Chat.ID,
    fmt.Sprintf("IP (%v): %v", changed, currIP),
  )
  reply.ReplyToMessageID = msg.MessageID
  _, err = tggenbot.SendMessage(reply)
  return err
}

func cmdPing(user *tggenbot.User, msg *tgbotapi.Message) error {
  reply := tgbotapi.NewMessage(msg.Chat.ID, "pong")
  reply.ReplyToMessageID = msg.MessageID
  _, err := tggenbot.SendMessage(reply)
  return err
}
