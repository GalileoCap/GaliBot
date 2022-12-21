package main

import (
  tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"

	"os"
	"io/ioutil" //TODO: Use os
	"net/http"
	"encoding/json"
	"log"
)

var DB *sql.DB;
var Bot *tgbotapi.BotAPI;

const MyChatID int64 = 1129477471; //TODO: Database
var CurrIP string;

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

func main() {
  //A: Init database
  _DB, err := sql.Open("mysql", "root:root@tcp(database:3306)/galibot");
  DB = _DB; //A: Rename to make it global
  if err != nil {
    log.Fatal(err);
  }
  defer DB.Close();

  //A: Init bot
  credentials := getCredentials(".credentials.json"); //TODO: Configure path
  _Bot, err := tgbotapi.NewBotAPI(credentials.Token);
  Bot = _Bot; //A: Rename to make it global
  if err != nil {
    log.Fatalf("[main] Error NewBotAPI: %v", err);
  }
  //Bot.Debug = true;
  log.Printf("[main] Running as %v", Bot.Self.UserName);

  if err = registerCommands(); err != nil {
    log.Fatalf("[main] Error on registerCommands: %v", err);
  }

  go ipUpdater();
  listenForMessages();
}
