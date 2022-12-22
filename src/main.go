package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
);

type ConfigT struct {
  Token string
  TestToken string `json=",omitempty"`
  Test bool
};
var Config ConfigT;

func parseConfig(path string, test bool) {
  data, err := os.ReadFile(path);
  if err != nil {
    log.Fatalf("[readToken] %v", err);
  }

  if err = json.Unmarshal(data, &Config); err != nil {
    log.Fatalf("[readToken] %v", err);
  }

  Config.Test = test;
}

func main() {
  //A: Register command flags
  configPath := flag.String("configPath", "config.json", "Path to the config file");
  apiToken := flag.String("token", "", "Your bot's API token (precedence over config and test)");
  test := flag.Bool("test", false, "Run in test mode (requires TestToken in config or --apiToken, precedence over config)");
  flag.Parse();

  parseConfig(*configPath, *test);

  if *apiToken == "" { //A: Make sure to have the API token
    if Config.Test {
      *apiToken = Config.TestToken;
    } else {
      *apiToken = Config.Token;
    }
  }

  telegramInit(*apiToken);

  receiveUpdates();
}
