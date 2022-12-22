package main;

/*
 * TODO:
 *  - For now (during initial testing) the database is kept in memory. Will be using SQL later
 */

import (
  tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

  "log"
);

type ModeT int;
const (
  Normal ModeT = iota
);

type User struct {
  ID int64;
  FirstName string;
  LastName string `json:",omitempty"`;
  UserName string `json:",omitempty"`;

  Permissions string; //U: admin/allow/block
  Mode string;
};
var Users map[int64]User;

func dbInit() {
  Users = make(map[int64]User);
}

func dbGetUser(requestUser *tgbotapi.User) User {
  user, present := Users[requestUser.ID];

  if !present { //A: Create it
    user = User{
      ID: requestUser.ID,
      FirstName: requestUser.FirstName,
      LastName: requestUser.LastName,
      UserName: requestUser.UserName,

      Mode: "",
    };

    //A: Set permissions
    if contains(Config.Admin, user.ID) { //A: Is an admin
      user.Permissions = "admin";
    } else if contains(Config.Block, user.ID) { //A: Is blocked
      user.Permissions = "block";
    } else { //A: Default
      user.Permissions = "allow";
    }
    //TODO: Simplify

    log.Printf("[dbGetUser] New user: %+v", user);
    Users[requestUser.ID] = user; //A: Save it
  }

  return user;
}

func dbSetUser(user *User) {
  Users[user.ID] = *user;
}
