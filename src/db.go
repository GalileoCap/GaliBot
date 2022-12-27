package main

/*
 * TODO:
 *  - For now (during initial testing) the database is kept in memory. Will be using SQL later
 */

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/muesli/cache2go"

  "time"
	"errors"
	"fmt"
	"log"
);

type ModeT int;
const (
  Normal ModeT = iota
);

type User struct {
  // This block is stored in the db
  ID int64;
  FirstName string;
  LastName string `json:",omitempty"`;
  UserName string `json:",omitempty"`;
  Permissions string; //U: admin/allow/block

  // This block only exists on the cache
  Mode string;
};
var _Users map[int64]*User; //U: Temporary until I have persistent storage //TODO: Remove
var Users *cache2go.CacheTable;

func dbInit() {
  _Users = make(map[int64]*User);
  Users = cache2go.Cache("users");
  Users.SetDataLoader(dbUsersDataLoader);
  Users.SetAboutToDeleteItemCallback(func (item *cache2go.CacheItem) { log.Printf("Delete: %v", item.Data().(*User).ID); });
}

func dbUsersDataLoader(key interface{}, args ...interface{}) *cache2go.CacheItem {
  err := args[1].(*error);

  user, present := _Users[key.(int64)]; //TODO: Load from the database
  if !present {
    *err = errors.New(fmt.Sprintf("[dbUsersDataLoader] User not present: %v", key.(int64)));
    return nil;
  }

  *err = nil;
  return cache2go.NewCacheItem(user.ID, time.Duration(Config.CacheLifespan) * time.Second, user);
}

func dbGetUser(requestUser *tgbotapi.User) (*User, error) {
  var err error;
  item, _ := Users.Value(requestUser.ID, requestUser, &err); //NOTE: The error returned is useless
  if err == nil { //A: They are cache'd
    return item.Data().(*User), nil; 
  }

  //TODO: Check err == "Not in db"

  //A: They're a new user
  user, err := dbNewUser(requestUser);
  if err != nil {
    return nil, err;
  }

  return user, nil;
}

func dbSaveUser(user *User) error {
  Users.Add(user.ID, time.Duration(Config.CacheLifespan) * time.Second, user);
  _Users[user.ID] = user; //TODO: Persistent
  return nil;
}

func dbNewUser(requestUser *tgbotapi.User) (*User, error) {
  user := &User{
    ID: requestUser.ID,
    FirstName: requestUser.FirstName,
    LastName: requestUser.LastName,
    UserName: requestUser.UserName,
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

  err := dbSaveUser(user);
  if err != nil {
    return nil, err;
  }

  log.Printf("[dbNewUser] New user: %+v", user);
  return user, nil;
}
